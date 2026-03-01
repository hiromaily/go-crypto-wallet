package btc_test

import (
	"context"
	"errors"
	"testing"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	dtobtc "github.com/hiromaily/go-crypto-wallet/internal/application/dto/btc"
	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	"github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch/btc"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
	bitcoinmocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/btc/mocks"
	repomocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/watch/mocks"
	storagemocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/transaction/mocks"
)

// testDependencies holds all mock dependencies for testing
type testDependencies struct {
	btcClient    *bitcoinmocks.MockBitcoiner
	addrRepo     *repomocks.MockAddressRepositorier
	txRepo       *repomocks.MockBTCTxRepositorier
	txInputRepo  *repomocks.MockTxInputRepositorier
	txOutputRepo *repomocks.MockTxOutputRepositorier
	payReqRepo   *repomocks.MockPaymentRequestRepositorier
	txFileRepo   *storagemocks.MockTransactionFileRepositorier
}

// newTestDependencies creates all mock dependencies
func newTestDependencies(t *testing.T) *testDependencies {
	t.Helper()
	return &testDependencies{
		btcClient:    bitcoinmocks.NewMockBitcoiner(t),
		addrRepo:     repomocks.NewMockAddressRepositorier(t),
		txRepo:       repomocks.NewMockBTCTxRepositorier(t),
		txInputRepo:  repomocks.NewMockTxInputRepositorier(t),
		txOutputRepo: repomocks.NewMockTxOutputRepositorier(t),
		payReqRepo:   repomocks.NewMockPaymentRequestRepositorier(t),
		txFileRepo:   storagemocks.NewMockTransactionFileRepositorier(t),
	}
}

// createUseCase creates a new CreateTransactionUseCase with the test dependencies
// Note: dbConn is set to nil because current tests focus on validation and error paths
// that don't reach the database transaction code. For tests that cover successful
// transaction creation and persistence, a mock database connection (e.g., sqlmock)
// would be required. See test comments below for more details.
func createUseCase(deps *testDependencies) watchusecase.CreateTransactionUseCase {
	return btc.NewCreateTransactionUseCase(
		deps.btcClient,
		nil, // dbConn - nil is safe for current tests that don't reach DB operations
		deps.addrRepo,
		deps.txRepo,
		deps.txInputRepo,
		deps.txOutputRepo,
		deps.payReqRepo,
		deps.txFileRepo,
		domainAccount.AccountTypeDeposit,
		domainAccount.AccountTypePayment,
		domainWallet.WalletTypeWatchOnly,
	)
}

// TestNewCreateTransactionUseCase tests the constructor
func TestNewCreateTransactionUseCase(t *testing.T) {
	t.Run("creates use case successfully with nil dependencies", func(t *testing.T) {
		useCase := btc.NewCreateTransactionUseCase(
			nil, // btcClient
			nil, // dbConn
			nil, // addrRepo
			nil, // txRepo
			nil, // txInputRepo
			nil, // txOutputRepo
			nil, // payReqRepo
			nil, // txFileRepo
			domainAccount.AccountTypeDeposit,
			domainAccount.AccountTypePayment,
			domainWallet.WalletTypeWatchOnly,
		)

		assert.NotNil(t, useCase, "use case should not be nil")
	})

	t.Run("returns correct interface type", func(t *testing.T) {
		useCase := btc.NewCreateTransactionUseCase(
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			domainAccount.AccountTypeDeposit,
			domainAccount.AccountTypePayment,
			domainWallet.WalletTypeWatchOnly,
		)

		// Verify it implements the interface
		assert.Implements(t, (*watchusecase.CreateTransactionUseCase)(nil), useCase)
	})

	t.Run("creates use case with mock dependencies", func(t *testing.T) {
		deps := newTestDependencies(t)
		useCase := createUseCase(deps)

		assert.NotNil(t, useCase, "use case should not be nil")
		assert.Implements(t, (*watchusecase.CreateTransactionUseCase)(nil), useCase)
	})
}

// TestExecute_InvalidActionType tests Execute with invalid action type
func TestExecute_InvalidActionType(t *testing.T) {
	deps := newTestDependencies(t)
	useCase := createUseCase(deps)

	input := watchusecase.CreateTransactionInput{
		ActionType: "invalid_action",
	}

	output, err := useCase.Execute(context.Background(), input)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid action type")
	assert.Empty(t, output.TransactionHex)
	assert.Empty(t, output.FileName)
}

// TestExecute_DepositTransaction tests deposit transaction creation
func TestExecute_DepositTransaction(t *testing.T) {
	t.Run("returns empty when no unspent transactions", func(t *testing.T) {
		deps := newTestDependencies(t)
		useCase := createUseCase(deps)

		// Setup mocks
		deps.btcClient.EXPECT().ConfirmationBlock().Return(uint64(6))
		deps.btcClient.EXPECT().ListUnspentByAccount(
			domainAccount.AccountTypeClient,
			uint64(6),
		).Return([]dtobtc.UnspentOutput{}, nil)
		deps.btcClient.EXPECT().GetUnspentListAddrs(
			[]dtobtc.UnspentOutput{},
			domainAccount.AccountTypeClient,
		).Return([]string{})

		input := watchusecase.CreateTransactionInput{
			ActionType:    domainTx.ActionTypeDeposit.String(),
			AdjustmentFee: 0.0,
		}

		output, err := useCase.Execute(context.Background(), input)

		require.NoError(t, err)
		assert.Empty(t, output.TransactionHex)
		assert.Empty(t, output.FileName)
	})

	t.Run("returns error when ListUnspentByAccount fails", func(t *testing.T) {
		deps := newTestDependencies(t)
		useCase := createUseCase(deps)

		// Setup mocks
		deps.btcClient.EXPECT().ConfirmationBlock().Return(uint64(6))
		deps.btcClient.EXPECT().ListUnspentByAccount(
			domainAccount.AccountTypeClient,
			uint64(6),
		).Return(nil, errors.New("RPC error"))

		input := watchusecase.CreateTransactionInput{
			ActionType: domainTx.ActionTypeDeposit.String(),
		}

		output, err := useCase.Execute(context.Background(), input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create transaction")
		assert.Empty(t, output.TransactionHex)
	})
}

// TestExecute_PaymentTransaction tests payment transaction creation
func TestExecute_PaymentTransaction(t *testing.T) {
	t.Run("returns empty when no payment requests", func(t *testing.T) {
		deps := newTestDependencies(t)
		useCase := createUseCase(deps)

		// Setup mocks for payment transaction - no payment requests
		deps.payReqRepo.EXPECT().GetAll().Return(nil, nil)

		input := watchusecase.CreateTransactionInput{
			ActionType: domainTx.ActionTypePayment.String(),
		}

		output, err := useCase.Execute(context.Background(), input)

		require.NoError(t, err)
		assert.Empty(t, output.TransactionHex)
		assert.Empty(t, output.FileName)
	})

	t.Run("returns error when GetAll payment requests fails", func(t *testing.T) {
		deps := newTestDependencies(t)
		useCase := createUseCase(deps)

		deps.payReqRepo.EXPECT().GetAll().Return(nil, errors.New("database error"))

		input := watchusecase.CreateTransactionInput{
			ActionType: domainTx.ActionTypePayment.String(),
		}

		output, err := useCase.Execute(context.Background(), input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create transaction")
		assert.Empty(t, output.TransactionHex)
	})
}

// TestExecute_TransferTransaction tests transfer transaction creation
func TestExecute_TransferTransaction(t *testing.T) {
	t.Run("returns error when sender and receiver are same", func(t *testing.T) {
		deps := newTestDependencies(t)
		useCase := createUseCase(deps)

		input := watchusecase.CreateTransactionInput{
			ActionType:      domainTx.ActionTypeTransfer.String(),
			SenderAccount:   domainAccount.AccountTypeDeposit,
			ReceiverAccount: domainAccount.AccountTypeDeposit, // Same as sender
			Amount:          0.001,
		}

		output, err := useCase.Execute(context.Background(), input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sender and receiver is same")
		assert.Empty(t, output.TransactionHex)
	})

	t.Run("returns error when receiver is client account", func(t *testing.T) {
		deps := newTestDependencies(t)
		useCase := createUseCase(deps)

		input := watchusecase.CreateTransactionInput{
			ActionType:      domainTx.ActionTypeTransfer.String(),
			SenderAccount:   domainAccount.AccountTypeDeposit,
			ReceiverAccount: domainAccount.AccountTypeClient, // Not allowed
			Amount:          0.001,
		}

		output, err := useCase.Execute(context.Background(), input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid receiver account")
		assert.Empty(t, output.TransactionHex)
	})

	t.Run("returns error when receiver is authorization account", func(t *testing.T) {
		deps := newTestDependencies(t)
		useCase := createUseCase(deps)

		input := watchusecase.CreateTransactionInput{
			ActionType:      domainTx.ActionTypeTransfer.String(),
			SenderAccount:   domainAccount.AccountTypeDeposit,
			ReceiverAccount: domainAccount.AccountTypeAuthorization, // Not allowed
			Amount:          0.001,
		}

		output, err := useCase.Execute(context.Background(), input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid receiver account")
		assert.Empty(t, output.TransactionHex)
	})

	t.Run("returns error when GetBalanceByAccount fails", func(t *testing.T) {
		deps := newTestDependencies(t)
		useCase := createUseCase(deps)

		deps.btcClient.EXPECT().ConfirmationBlock().Return(uint64(6))
		deps.btcClient.EXPECT().GetBalanceByAccount(
			domainAccount.AccountTypeDeposit,
			uint64(6),
		).Return(btcutil.Amount(0), errors.New("RPC error"))

		input := watchusecase.CreateTransactionInput{
			ActionType:      domainTx.ActionTypeTransfer.String(),
			SenderAccount:   domainAccount.AccountTypeDeposit,
			ReceiverAccount: domainAccount.AccountTypePayment,
			Amount:          0.01,
		}

		output, err := useCase.Execute(context.Background(), input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create transaction")
		assert.Empty(t, output.TransactionHex)
	})

	t.Run("returns error when sender balance is insufficient", func(t *testing.T) {
		deps := newTestDependencies(t)
		useCase := createUseCase(deps)

		// Setup mocks
		deps.btcClient.EXPECT().ConfirmationBlock().Return(uint64(6))
		deps.btcClient.EXPECT().GetBalanceByAccount(
			domainAccount.AccountTypeDeposit,
			uint64(6),
		).Return(btcutil.Amount(500000), nil) // Less than required

		input := watchusecase.CreateTransactionInput{
			ActionType:      domainTx.ActionTypeTransfer.String(),
			SenderAccount:   domainAccount.AccountTypeDeposit,
			ReceiverAccount: domainAccount.AccountTypePayment,
			Amount:          0.01, // 0.01 BTC required, but only 0.005 BTC available
		}

		output, err := useCase.Execute(context.Background(), input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "balance is insufficient")
		assert.Empty(t, output.TransactionHex)
	})
}

// TestExecute_UnsupportedActionType tests unsupported action types
func TestExecute_UnsupportedActionType(t *testing.T) {
	deps := newTestDependencies(t)
	useCase := createUseCase(deps)

	// Note: This tests the action type validation
	input := watchusecase.CreateTransactionInput{
		ActionType: "unknown",
	}

	output, err := useCase.Execute(context.Background(), input)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid action type")
	assert.Empty(t, output.TransactionHex)
}

// TestMockExpectations verifies that all mock expectations work correctly
func TestMockExpectations(t *testing.T) {
	t.Run("verifies mock expectations are properly set and verified", func(t *testing.T) {
		deps := newTestDependencies(t)
		useCase := createUseCase(deps)

		// Setup minimal mocks
		deps.btcClient.EXPECT().ConfirmationBlock().Return(uint64(6))
		deps.btcClient.EXPECT().ListUnspentByAccount(
			domainAccount.AccountTypeClient,
			uint64(6),
		).Return([]dtobtc.UnspentOutput{}, nil)
		deps.btcClient.EXPECT().GetUnspentListAddrs(
			[]dtobtc.UnspentOutput{},
			domainAccount.AccountTypeClient,
		).Return([]string{})

		input := watchusecase.CreateTransactionInput{
			ActionType: domainTx.ActionTypeDeposit.String(),
		}

		_, err := useCase.Execute(context.Background(), input)
		require.NoError(t, err)

		// Mock expectations are automatically verified by testify when test ends
		// AssertExpectations is called automatically by mockery-generated mocks
	})
}

// Test Coverage Notes:
//
// Current test coverage focuses on:
// - Constructor validation and interface compliance
// - Error handling paths (validation errors, RPC errors, insufficient balance)
// - Input validation (invalid action types, invalid accounts)
// - Early return paths (no unspent transactions, no payment requests)
//
// Note on dbConn being nil:
// The dbConn parameter is set to nil in these tests because the current test cases
// focus on validation and error paths that return before reaching the database
// transaction code (insertTransaction method). The dbConn is only used in the
// success path when a transaction is successfully created and needs to be persisted.
// For tests that cover the successful transaction creation and persistence path,
// a mock database connection (e.g., using sqlmock) would be required.
//
// Success path testing (future work):
// To test the full success path including database persistence, the following
// would be needed:
// 1. Mock database connection using sqlmock (github.com/DATA-DOG/go-sqlmock)
// 2. Mock Bitcoin client with realistic address decoding (btcutil.DecodeAddress)
// 3. Mock file repository for PSBT file writing
// 4. Complete mock setup for all repository operations (InsertUnsignedTx,
//    InsertBulk for inputs/outputs, UpdatePaymentID, etc.)
//
// These comprehensive tests would be better suited for integration tests rather
// than unit tests, as they require extensive mocking of external systems.
//
// PSBT-specific functionality testing:
// The generatePSBTFile method creates PSBT files instead of CSV files.
// Integration tests should verify:
// - PSBT creation for deposit transactions (client → deposit)
// - PSBT creation for payment transactions (payment → user addresses)
// - PSBT creation for transfer transactions (account → account)
// - PSBT includes all required metadata (amounts, scriptPubKeys, redeem scripts)
// - PSBT files follow naming convention: {actionType}_{txID}_{txType}_{signedCount}_{timestamp}.psbt
// - PSBT files are base64-encoded and BIP174 compliant
// - Generated PSBT files can be read by Keygen wallet (offline signing)
