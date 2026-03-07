package xrp_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	dtoxrp "github.com/hiromaily/go-crypto-wallet/internal/application/dto/xrp"
	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	"github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch/xrp"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	xrpapiamocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/xrp/mocks"
	repomocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/watch/mocks"
	storagemocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/transaction/mocks"
	dbtxmocks "github.com/hiromaily/go-crypto-wallet/pkg/db/tx/mocks"
	pkguuid "github.com/hiromaily/go-crypto-wallet/pkg/uuid"
)

// testDependencies holds all mock dependencies for testing
type testDependencies struct {
	accountInfo  *xrpapiamocks.MockAccountInfoProvider
	txPreparer   *xrpapiamocks.MockTransactionPreparer
	addrRepo     *repomocks.MockAddressRepositorier
	txRepo       *repomocks.MockTxRepositorier
	txDetailRepo *repomocks.MockXRPDetailTXRepositorier
	payReqRepo   *repomocks.MockPaymentRequestRepositorier
	txFileRepo   *storagemocks.MockTransactionFileRepositorier
	uuidHandler  pkguuid.UUIDHandler
	unitOfWork   *dbtxmocks.MockUnitOfWork
	network      string
}

// newTestDependencies creates all mock dependencies
func newTestDependencies(t *testing.T) *testDependencies {
	t.Helper()
	return &testDependencies{
		accountInfo:  xrpapiamocks.NewMockAccountInfoProvider(t),
		txPreparer:   xrpapiamocks.NewMockTransactionPreparer(t),
		addrRepo:     repomocks.NewMockAddressRepositorier(t),
		txRepo:       repomocks.NewMockTxRepositorier(t),
		txDetailRepo: repomocks.NewMockXRPDetailTXRepositorier(t),
		payReqRepo:   repomocks.NewMockPaymentRequestRepositorier(t),
		txFileRepo:   storagemocks.NewMockTransactionFileRepositorier(t),
		uuidHandler:  pkguuid.NewGoogleUUIDHandler(),
		network:      "testnet",
	}
}

// createUseCase creates a new CreateTransactionUseCase with the test dependencies
func createUseCase(deps *testDependencies) watchusecase.CreateTransactionUseCase {
	return xrp.NewCreateTransactionUseCase(
		deps.accountInfo,
		deps.txPreparer,
		deps.unitOfWork, // nil is safe for tests that don't reach DB operations
		deps.uuidHandler,
		deps.addrRepo,
		deps.txRepo,
		deps.txDetailRepo,
		deps.payReqRepo,
		deps.txFileRepo,
		domainAccount.AccountTypeDeposit,
		domainAccount.AccountTypePayment,
		deps.network,
	)
}

func TestCreateTransactionUseCase_Execute_InvalidActionType(t *testing.T) {
	t.Parallel()
	deps := newTestDependencies(t)
	useCase := createUseCase(deps)

	input := watchusecase.CreateTransactionInput{
		ActionType: "invalid",
	}

	output, err := useCase.Execute(context.Background(), input)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid action type")
	assert.Empty(t, output.FileName)
}

func TestCreateTransactionUseCase_Execute_TransferWithInsufficientBalance(t *testing.T) {
	t.Parallel()
	deps := newTestDependencies(t)
	useCase := createUseCase(deps)

	senderAddr := &domainAddress.Address{
		ID:            1,
		WalletAddress: "rSenderAddress123",
	}

	// Setup mocks
	deps.addrRepo.EXPECT().GetOneUnAllocated(domainAccount.AccountTypePayment).Return(senderAddr, nil)
	deps.accountInfo.EXPECT().GetBalance(mock.Anything, "rSenderAddress123").Return(10.0, nil)

	input := watchusecase.CreateTransactionInput{
		ActionType:      domainTx.ActionTypeTransfer.String(),
		SenderAccount:   domainAccount.AccountTypePayment,
		ReceiverAccount: domainAccount.AccountTypeDeposit,
		Amount:          100.0, // More than balance
	}

	output, err := useCase.Execute(context.Background(), input)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "sender balance")
	assert.Contains(t, err.Error(), "insufficient")
	assert.Empty(t, output.FileName)
}

func TestCreateTransactionUseCase_Execute_TransferWithValidBalance(t *testing.T) {
	t.Parallel()
	t.Skip("Skipping until DB transaction is mocked - requires comprehensive DB mock setup")
	// This test requires:
	// 1. Mock database connection with transaction support
	// 2. Mock all repository operations within the transaction
	// 3. Mock file write operations
	// This is complex and should be tested in integration tests instead
}

func TestCreateTransactionUseCase_Execute_GetBalanceError(t *testing.T) {
	t.Parallel()
	deps := newTestDependencies(t)
	useCase := createUseCase(deps)

	senderAddr := &domainAddress.Address{
		ID:            1,
		WalletAddress: "rSenderAddress123",
	}

	// Setup mocks
	deps.addrRepo.EXPECT().GetOneUnAllocated(domainAccount.AccountTypePayment).Return(senderAddr, nil)
	deps.accountInfo.EXPECT().GetBalance(mock.Anything, "rSenderAddress123").
		Return(0.0, errors.New("network error"))

	input := watchusecase.CreateTransactionInput{
		ActionType:      domainTx.ActionTypeTransfer.String(),
		SenderAccount:   domainAccount.AccountTypePayment,
		ReceiverAccount: domainAccount.AccountTypeDeposit,
		Amount:          10.0,
	}

	output, err := useCase.Execute(context.Background(), input)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get balance")
	assert.Contains(t, err.Error(), "network error")
	assert.Empty(t, output.FileName)
}

func TestCreateTransactionUseCase_Execute_GetAddressError(t *testing.T) {
	t.Parallel()
	deps := newTestDependencies(t)
	useCase := createUseCase(deps)

	// Setup mocks
	deps.addrRepo.EXPECT().GetOneUnAllocated(domainAccount.AccountTypePayment).
		Return((*domainAddress.Address)(nil), errors.New("database error"))

	input := watchusecase.CreateTransactionInput{
		ActionType:      domainTx.ActionTypeTransfer.String(),
		SenderAccount:   domainAccount.AccountTypePayment,
		ReceiverAccount: domainAccount.AccountTypeDeposit,
		Amount:          10.0,
	}

	output, err := useCase.Execute(context.Background(), input)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get sender address")
	assert.Contains(t, err.Error(), "database error")
	assert.Empty(t, output.FileName)
}

func TestCreateTransactionUseCase_Execute_NilSenderAddress(t *testing.T) {
	t.Parallel()
	deps := newTestDependencies(t)
	useCase := createUseCase(deps)

	// Setup mocks - repository returns nil address without error (no unallocated address exists)
	deps.addrRepo.EXPECT().GetOneUnAllocated(domainAccount.AccountTypePayment).
		Return((*domainAddress.Address)(nil), nil)

	input := watchusecase.CreateTransactionInput{
		ActionType:      domainTx.ActionTypeTransfer.String(),
		SenderAccount:   domainAccount.AccountTypePayment,
		ReceiverAccount: domainAccount.AccountTypeDeposit,
		Amount:          10.0,
	}

	output, err := useCase.Execute(context.Background(), input)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "no unallocated address found")
	assert.Contains(t, err.Error(), "sender account")
	assert.Empty(t, output.FileName)
}

func TestCreateTransactionUseCase_Execute_NilReceiverAddress(t *testing.T) {
	t.Parallel()
	deps := newTestDependencies(t)
	useCase := createUseCase(deps)

	senderAddr := &domainAddress.Address{
		ID:            1,
		WalletAddress: "rSenderAddress123",
	}

	// Setup mocks
	deps.addrRepo.EXPECT().GetOneUnAllocated(domainAccount.AccountTypePayment).Return(senderAddr, nil)
	deps.accountInfo.EXPECT().GetBalance(mock.Anything, "rSenderAddress123").Return(100.0, nil)
	// Receiver address returns nil without error
	deps.addrRepo.EXPECT().GetOneUnAllocated(domainAccount.AccountTypeDeposit).
		Return((*domainAddress.Address)(nil), nil)

	input := watchusecase.CreateTransactionInput{
		ActionType:      domainTx.ActionTypeTransfer.String(),
		SenderAccount:   domainAccount.AccountTypePayment,
		ReceiverAccount: domainAccount.AccountTypeDeposit,
		Amount:          10.0,
	}

	output, err := useCase.Execute(context.Background(), input)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "no unallocated address found")
	assert.Contains(t, err.Error(), "receiver account")
	assert.Empty(t, output.FileName)
}

func TestCreateTransactionUseCase_Dependencies(t *testing.T) {
	t.Parallel()
	t.Run("uses segregated interfaces instead of full XRPer interface", func(t *testing.T) {
		t.Parallel()
		deps := newTestDependencies(t)

		// Create use case - should accept AccountInfoProvider and TransactionPreparer separately
		useCase := xrp.NewCreateTransactionUseCase(
			deps.accountInfo, // AccountInfoProvider (not full XRPer)
			deps.txPreparer,  // TransactionPreparer (not full XRPer)
			nil,
			deps.uuidHandler,
			deps.addrRepo,
			deps.txRepo,
			deps.txDetailRepo,
			deps.payReqRepo,
			deps.txFileRepo,
			domainAccount.AccountTypeDeposit,
			domainAccount.AccountTypePayment,
			"testnet",
		)

		assert.NotNil(t, useCase, "use case should be created with segregated interfaces")
	})
}

func TestCreateTransactionUseCase_Execute_TransferMultisig_WritesJSONFile(t *testing.T) {
	t.Parallel()

	deps := newTestDependencies(t)
	deps.unitOfWork = dbtxmocks.NewMockUnitOfWork(t)

	senderAddr := &domainAddress.Address{ID: 1, WalletAddress: "rSenderAddr123"}
	receiverAddr := &domainAddress.Address{ID: 2, WalletAddress: "rReceiverAddr456"}
	txJSON := &dtoxrp.TxInput{
		TransactionType:    "Payment",
		Amount:             "10000000",
		Fee:                "12",
		Sequence:           100,
		LastLedgerSequence: 200,
	}
	mockTx := dbtxmocks.NewMockTransaction(t)

	deps.addrRepo.EXPECT().GetOneUnAllocated(domainAccount.AccountTypePayment).Return(senderAddr, nil)
	deps.accountInfo.EXPECT().GetBalance(mock.Anything, "rSenderAddr123").Return(50.0, nil)
	deps.addrRepo.EXPECT().GetOneUnAllocated(domainAccount.AccountTypeDeposit).Return(receiverAddr, nil)
	deps.txPreparer.EXPECT().
		CreateRawTransaction(mock.Anything, "rSenderAddr123", "rReceiverAddr456", 10.0, mock.Anything).
		Return(txJSON, `{"TransactionType":"Payment"}`, nil)

	deps.unitOfWork.EXPECT().Begin(mock.Anything).Return(mockTx, nil)
	deps.txRepo.EXPECT().WithTransaction(mockTx).Return(deps.txRepo, nil)
	deps.txDetailRepo.EXPECT().WithTransaction(mockTx).Return(deps.txDetailRepo, nil)
	deps.payReqRepo.EXPECT().WithTransaction(mockTx).Return(deps.payReqRepo, nil)
	deps.txRepo.EXPECT().InsertUnsignedTx(domainTx.ActionTypeTransfer).Return(int64(42), nil)
	deps.txDetailRepo.EXPECT().InsertBulk(mock.Anything).Return(nil)
	mockTx.EXPECT().Commit().Return(nil)

	deps.txFileRepo.EXPECT().
		CreateFilePath(domainTx.ActionTypeTransfer, domainTx.TxTypeUnsigned, int64(42), 0).
		Return("/tmp/multisig_transfer_42_0.json")
	deps.txFileRepo.EXPECT().
		WriteXRPJSONFile("/tmp/multisig_transfer_42_0.json", mock.MatchedBy(func(f *dtoxrp.XRPTransactionFile) bool {
			return f.Chain == "XRP" &&
				f.Network == "testnet" &&
				f.Version == "1.0.0" &&
				len(f.Transactions) == 1 &&
				f.Transactions[0].RequiredSignatures == 2 &&
				f.Transactions[0].SignatureCount == 0 &&
				!f.Transactions[0].IsComplete &&
				f.Transactions[0].SignedBlob == nil
		})).
		Return("/tmp/multisig_transfer_42_0.json", nil)

	useCase := createUseCase(deps)
	input := watchusecase.CreateTransactionInput{
		ActionType:      domainTx.ActionTypeTransfer.String(),
		SenderAccount:   domainAccount.AccountTypePayment,
		ReceiverAccount: domainAccount.AccountTypeDeposit,
		Amount:          10.0,
		MultisigQuorum:  2,
	}

	output, err := useCase.Execute(context.Background(), input)

	require.NoError(t, err)
	assert.Equal(t, "/tmp/multisig_transfer_42_0.json", output.FileName)
}

func TestCreateTransactionUseCase_Execute_TransferSingleSig_WritesTextFile(t *testing.T) {
	t.Parallel()

	deps := newTestDependencies(t)
	deps.unitOfWork = dbtxmocks.NewMockUnitOfWork(t)

	senderAddr := &domainAddress.Address{ID: 1, WalletAddress: "rSenderAddr123"}
	receiverAddr := &domainAddress.Address{ID: 2, WalletAddress: "rReceiverAddr456"}
	txJSON := &dtoxrp.TxInput{
		TransactionType:    "Payment",
		Amount:             "10000000",
		Fee:                "12",
		Sequence:           100,
		LastLedgerSequence: 200,
	}
	mockTx := dbtxmocks.NewMockTransaction(t)

	deps.addrRepo.EXPECT().GetOneUnAllocated(domainAccount.AccountTypePayment).Return(senderAddr, nil)
	deps.accountInfo.EXPECT().GetBalance(mock.Anything, "rSenderAddr123").Return(50.0, nil)
	deps.addrRepo.EXPECT().GetOneUnAllocated(domainAccount.AccountTypeDeposit).Return(receiverAddr, nil)
	deps.txPreparer.EXPECT().
		CreateRawTransaction(mock.Anything, "rSenderAddr123", "rReceiverAddr456", 10.0, mock.Anything).
		Return(txJSON, `{"TransactionType":"Payment"}`, nil)

	deps.unitOfWork.EXPECT().Begin(mock.Anything).Return(mockTx, nil)
	deps.txRepo.EXPECT().WithTransaction(mockTx).Return(deps.txRepo, nil)
	deps.txDetailRepo.EXPECT().WithTransaction(mockTx).Return(deps.txDetailRepo, nil)
	deps.payReqRepo.EXPECT().WithTransaction(mockTx).Return(deps.payReqRepo, nil)
	deps.txRepo.EXPECT().InsertUnsignedTx(domainTx.ActionTypeTransfer).Return(int64(42), nil)
	deps.txDetailRepo.EXPECT().InsertBulk(mock.Anything).Return(nil)
	mockTx.EXPECT().Commit().Return(nil)

	deps.txFileRepo.EXPECT().
		CreateFilePath(domainTx.ActionTypeTransfer, domainTx.TxTypeUnsigned, int64(42), 0).
		Return("/tmp/transfer_42_0.json")
	// Single-sig path: WriteFileSlice is called, NOT WriteXRPJSONFile
	deps.txFileRepo.EXPECT().
		WriteFileSlice("/tmp/transfer_42_0.json", mock.Anything).
		Return("/tmp/transfer_42_0.json", nil)

	useCase := createUseCase(deps)
	input := watchusecase.CreateTransactionInput{
		ActionType:      domainTx.ActionTypeTransfer.String(),
		SenderAccount:   domainAccount.AccountTypePayment,
		ReceiverAccount: domainAccount.AccountTypeDeposit,
		Amount:          10.0,
		MultisigQuorum:  0, // Single-sig path
	}

	output, err := useCase.Execute(context.Background(), input)

	require.NoError(t, err)
	assert.Equal(t, "/tmp/transfer_42_0.json", output.FileName)
}
