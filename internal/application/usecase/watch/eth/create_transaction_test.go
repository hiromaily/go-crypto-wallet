package eth_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	dtoeth "github.com/hiromaily/go-crypto-wallet/internal/application/dto/eth"
	apieth "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/eth"
	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	watchusecaseeth "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch/eth"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainETH "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/eth"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	ethapiamocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/eth/mocks"
	repomocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/watch/mocks"
	storagemocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/transaction/mocks"
)

// ethCreateTxTestDependencies holds all mock dependencies for testing
type ethCreateTxTestDependencies struct {
	ethClient    *ethapiamocks.MockERC20er
	addrRepo     *repomocks.MockAddressRepositorier
	txRepo       *repomocks.MockTxRepositorier
	txDetailRepo *repomocks.MockETHDetailTXRepositorier
	payReqRepo   *repomocks.MockPaymentRequestRepositorier
	txFileRepo   *storagemocks.MockTransactionFileRepositorier
}

func newEthCreateTxTestDeps(t *testing.T) *ethCreateTxTestDependencies {
	t.Helper()
	return &ethCreateTxTestDependencies{
		ethClient:    ethapiamocks.NewMockERC20er(t),
		addrRepo:     repomocks.NewMockAddressRepositorier(t),
		txRepo:       repomocks.NewMockTxRepositorier(t),
		txDetailRepo: repomocks.NewMockETHDetailTXRepositorier(t),
		payReqRepo:   repomocks.NewMockPaymentRequestRepositorier(t),
		txFileRepo:   storagemocks.NewMockTransactionFileRepositorier(t),
	}
}

func newEthCreateTxUseCase(deps *ethCreateTxTestDependencies) watchusecase.CreateTransactionUseCase {
	return watchusecaseeth.NewCreateTransactionUseCase(
		deps.ethClient,
		nil, // dbConn — nil is safe for tests that do not reach DB code
		deps.addrRepo,
		deps.txRepo,
		deps.txDetailRepo,
		deps.payReqRepo,
		deps.txFileRepo,
		domainAccount.AccountTypeDeposit,
		domainAccount.AccountTypePayment,
	)
}

// ---- Helper values -----

const (
	testETHFromAddr = "0xAbcdef1234567890AbCdEf1234567890AbCdEf12"
	testETHToAddr   = "0xDeadBeef1234567890DeadBeef1234567890abcd"
	testETHRawTx    = "0x02c0" // minimal EIP-2718 encoded tx for testing
	testETHTxPath   = "deposit_1_unsigned_0_1534744535097796209.json"
)

func testETHAddress(addr string, accountType domainAccount.AccountType) *domainAddress.Address {
	return &domainAddress.Address{
		WalletAddress: addr,
		AccountType:   accountType,
	}
}

// ---- Tests ----

func TestETHCreateTransactionUseCase_Execute_InvalidActionType(t *testing.T) {
	t.Parallel()
	deps := newEthCreateTxTestDeps(t)
	useCase := newEthCreateTxUseCase(deps)

	output, err := useCase.Execute(context.Background(), watchusecase.CreateTransactionInput{
		ActionType: "invalid",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid action type")
	assert.Empty(t, output.FileName)
}

func TestETHCreateTransactionUseCase_Execute_DepositNoAddresses(t *testing.T) {
	t.Parallel()
	deps := newEthCreateTxTestDeps(t)
	useCase := newEthCreateTxUseCase(deps)

	// Return empty address list — no transactions to create
	deps.addrRepo.EXPECT().
		GetAll(domainAccount.AccountTypeClient).
		Return(nil, nil)

	output, err := useCase.Execute(context.Background(), watchusecase.CreateTransactionInput{
		ActionType: string(domainTx.ActionTypeDeposit),
	})

	require.NoError(t, err)
	assert.Empty(t, output.FileName)
}

func TestETHCreateTransactionUseCase_Execute_DepositAllZeroBalance(t *testing.T) {
	t.Parallel()
	deps := newEthCreateTxTestDeps(t)
	useCase := newEthCreateTxUseCase(deps)

	addrs := []*domainAddress.Address{
		testETHAddress(testETHFromAddr, domainAccount.AccountTypeClient),
	}
	deps.addrRepo.EXPECT().
		GetAll(domainAccount.AccountTypeClient).
		Return(addrs, nil)
	// Balance is zero → address is skipped
	deps.ethClient.EXPECT().
		GetBalance(context.Background(), testETHFromAddr, domainETH.QuantityTagLatest).
		Return(new(big.Int), nil)

	output, err := useCase.Execute(context.Background(), watchusecase.CreateTransactionInput{
		ActionType: string(domainTx.ActionTypeDeposit),
	})

	require.NoError(t, err)
	assert.Empty(t, output.FileName)
}

func TestETHCreateTransactionUseCase_Execute_TransferSenderEqualsReceiver(t *testing.T) {
	t.Parallel()
	deps := newEthCreateTxTestDeps(t)
	useCase := newEthCreateTxUseCase(deps)

	output, err := useCase.Execute(context.Background(), watchusecase.CreateTransactionInput{
		ActionType:      string(domainTx.ActionTypeTransfer),
		SenderAccount:   domainAccount.AccountTypeDeposit,
		ReceiverAccount: domainAccount.AccountTypeDeposit,
		Amount:          1.0,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "sender and receiver is same")
	assert.Empty(t, output.FileName)
}

func TestETHCreateTransactionUseCase_Execute_TransferInvalidReceiver(t *testing.T) {
	t.Parallel()
	deps := newEthCreateTxTestDeps(t)
	useCase := newEthCreateTxUseCase(deps)

	output, err := useCase.Execute(context.Background(), watchusecase.CreateTransactionInput{
		ActionType:      string(domainTx.ActionTypeTransfer),
		SenderAccount:   domainAccount.AccountTypeDeposit,
		ReceiverAccount: domainAccount.AccountTypeClient, // not allowed as receiver
		Amount:          1.0,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "client, authorization account is not allowed as receiver")
	assert.Empty(t, output.FileName)
}

func TestETHCreateTransactionUseCase_Execute_TransferZeroBalance(t *testing.T) {
	t.Parallel()
	deps := newEthCreateTxTestDeps(t)
	useCase := newEthCreateTxUseCase(deps)

	deps.addrRepo.EXPECT().
		GetOneUnAllocated(domainAccount.AccountTypeDeposit).
		Return(testETHAddress(testETHFromAddr, domainAccount.AccountTypeDeposit), nil)
	deps.ethClient.EXPECT().
		GetBalance(context.Background(), testETHFromAddr, domainETH.QuantityTagLatest).
		Return(new(big.Int), nil) // zero balance

	output, err := useCase.Execute(context.Background(), watchusecase.CreateTransactionInput{
		ActionType:      string(domainTx.ActionTypeTransfer),
		SenderAccount:   domainAccount.AccountTypeDeposit,
		ReceiverAccount: domainAccount.AccountTypePayment,
		Amount:          1.0,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "sender has no balance")
	assert.Empty(t, output.FileName)
}

func TestETHCreateTransactionUseCase_Execute_DepositEIP1559Path(t *testing.T) {
	t.Parallel()
	t.Skip("Skipping until DB transaction is mocked - requires comprehensive DB mock setup")

	deps := newEthCreateTxTestDeps(t)
	useCase := newEthCreateTxUseCase(deps)

	addrs := []*domainAddress.Address{
		testETHAddress(testETHFromAddr, domainAccount.AccountTypeClient),
	}
	depositAddr := testETHAddress(testETHToAddr, domainAccount.AccountTypeDeposit)

	deps.addrRepo.EXPECT().
		GetAll(domainAccount.AccountTypeClient).
		Return(addrs, nil)
	deps.ethClient.EXPECT().
		GetBalance(context.Background(), testETHFromAddr, domainETH.QuantityTagLatest).
		Return(new(big.Int).SetUint64(1_000_000_000_000_000_000), nil) // 1 ETH
	deps.addrRepo.EXPECT().
		GetOneUnAllocated(domainAccount.AccountTypeDeposit).
		Return(depositAddr, nil)

	// EIP-1559 path: SupportsEIP1559 returns true → CreateRawTransactionEIP1559 is called
	deps.ethClient.EXPECT().
		SupportsEIP1559(context.Background()).
		Return(true)
	deps.ethClient.EXPECT().
		CreateRawTransactionEIP1559(
			context.Background(), testETHFromAddr, testETHToAddr, uint64(0), 0,
		).
		Return(
			&domainETH.RawTx{
				UUID:  "test-uuid",
				From:  testETHFromAddr,
				To:    testETHToAddr,
				Value: *new(big.Int).SetUint64(999_000_000_000_000_000),
				Nonce: 0,
				TxHex: testETHRawTx,
			},
			&apieth.TxCreateParams{
				UUID:                 "test-uuid",
				FromAddress:          testETHFromAddr,
				ToAddress:            testETHToAddr,
				Amount:               999_000_000_000_000_000,
				Fee:                  1_000_000_000_000_000,
				GasLimit:             21000,
				Nonce:                0,
				EthTxType:            2,
				ChainID:              31337,
				MaxFeePerGas:         20_000_000_000,
				MaxPriorityFeePerGas: 2_000_000_000,
			},
			nil,
		)

	// DB mocks would be needed here (InsertUnsignedTx, InsertBulk, etc.)
	// txFileRepo expects WriteETHJSONFile with EIP-1559 data
	deps.txFileRepo.EXPECT().
		CreateFilePath(domainTx.ActionTypeDeposit, domainTx.TxTypeUnsigned, int64(1), 0).
		Return(testETHTxPath)
	deps.txFileRepo.EXPECT().
		WriteETHJSONFile(testETHTxPath, &dtoeth.ETHTransactionFile{
			Version:              1,
			TxType:               string(domainTx.TxTypeUnsigned),
			EthTxType:            2,
			ChainID:              31337,
			Nonce:                0,
			From:                 testETHFromAddr,
			To:                   testETHToAddr,
			Value:                "999000000000000000",
			Gas:                  21000,
			MaxFeePerGas:         "20000000000",
			MaxPriorityFeePerGas: "2000000000",
			RawTxHex:             testETHRawTx,
		}).
		Return(testETHTxPath, nil)

	output, err := useCase.Execute(context.Background(), watchusecase.CreateTransactionInput{
		ActionType: string(domainTx.ActionTypeDeposit),
	})

	require.NoError(t, err)
	assert.Equal(t, testETHTxPath, output.FileName)
}

func TestETHCreateTransactionUseCase_Execute_DepositLegacyPath(t *testing.T) {
	t.Parallel()
	t.Skip("Skipping until DB transaction is mocked - requires comprehensive DB mock setup")

	deps := newEthCreateTxTestDeps(t)
	useCase := newEthCreateTxUseCase(deps)

	addrs := []*domainAddress.Address{
		testETHAddress(testETHFromAddr, domainAccount.AccountTypeClient),
	}
	depositAddr := testETHAddress(testETHToAddr, domainAccount.AccountTypeDeposit)

	deps.addrRepo.EXPECT().
		GetAll(domainAccount.AccountTypeClient).
		Return(addrs, nil)
	deps.ethClient.EXPECT().
		GetBalance(context.Background(), testETHFromAddr, domainETH.QuantityTagLatest).
		Return(new(big.Int).SetUint64(1_000_000_000_000_000_000), nil) // 1 ETH
	deps.addrRepo.EXPECT().
		GetOneUnAllocated(domainAccount.AccountTypeDeposit).
		Return(depositAddr, nil)

	// Legacy path: SupportsEIP1559 returns false → CreateRawTransaction is called
	deps.ethClient.EXPECT().
		SupportsEIP1559(context.Background()).
		Return(false)
	deps.ethClient.EXPECT().
		CreateRawTransaction(
			context.Background(), testETHFromAddr, testETHToAddr, uint64(0), 0,
		).
		Return(
			&domainETH.RawTx{
				UUID:  "test-uuid",
				From:  testETHFromAddr,
				To:    testETHToAddr,
				Value: *new(big.Int).SetUint64(999_000_000_000_000_000),
				Nonce: 0,
				TxHex: testETHRawTx,
			},
			&apieth.TxCreateParams{
				UUID:        "test-uuid",
				FromAddress: testETHFromAddr,
				ToAddress:   testETHToAddr,
				Amount:      999_000_000_000_000_000,
				Fee:         1_000_000_000_000_000,
				GasLimit:    21000,
				Nonce:       0,
				EthTxType:   0,
				GasPrice:    20_000_000_000,
			},
			nil,
		)

	// txFileRepo expects WriteETHJSONFile with legacy data
	deps.txFileRepo.EXPECT().
		CreateFilePath(domainTx.ActionTypeDeposit, domainTx.TxTypeUnsigned, int64(1), 0).
		Return(testETHTxPath)
	deps.txFileRepo.EXPECT().
		WriteETHJSONFile(testETHTxPath, &dtoeth.ETHTransactionFile{
			Version:   1,
			TxType:    string(domainTx.TxTypeUnsigned),
			EthTxType: 0,
			Nonce:     0,
			From:      testETHFromAddr,
			To:        testETHToAddr,
			Value:     "999000000000000000",
			Gas:       21000,
			GasPrice:  "20000000000",
			RawTxHex:  testETHRawTx,
		}).
		Return(testETHTxPath, nil)

	output, err := useCase.Execute(context.Background(), watchusecase.CreateTransactionInput{
		ActionType: string(domainTx.ActionTypeDeposit),
	})

	require.NoError(t, err)
	assert.Equal(t, testETHTxPath, output.FileName)
}
