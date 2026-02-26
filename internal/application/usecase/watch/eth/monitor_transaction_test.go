package eth_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	watchusecaseeth "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch/eth"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainETH "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/eth"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	ethapiamocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/eth/mocks"
	repomocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/watch/mocks"
)

// ethMonitorTxTestDependencies holds all mock dependencies for monitor transaction testing.
type ethMonitorTxTestDependencies struct {
	ethClient    *ethapiamocks.MockEtherTxMonitor
	addrRepo     *repomocks.MockAddressRepositorier
	txDetailRepo *repomocks.MockETHDetailTXRepositorier
}

func newEthMonitorTxTestDeps(t *testing.T) *ethMonitorTxTestDependencies {
	t.Helper()
	return &ethMonitorTxTestDependencies{
		ethClient:    ethapiamocks.NewMockEtherTxMonitor(t),
		addrRepo:     repomocks.NewMockAddressRepositorier(t),
		txDetailRepo: repomocks.NewMockETHDetailTXRepositorier(t),
	}
}

func newEthMonitorTxUseCase(deps *ethMonitorTxTestDependencies) watchusecase.MonitorTransactionUseCase {
	return watchusecaseeth.NewMonitorTransactionUseCase(
		deps.ethClient,
		deps.addrRepo,
		deps.txDetailRepo,
		testConfirmNum,
	)
}

const (
	testMonitorSentHash = "0xaabbccdd1234567890aabbccdd1234567890aabbccdd1234567890aabbccdd12"
	testConfirmNum      = uint64(6)
)

func testSuccessReceipt() *domainETH.TransactionReceipt {
	return &domainETH.TransactionReceipt{
		TransactionHash: testMonitorSentHash,
		BlockNumber:     100,
		From:            testETHFromAddr,
		To:              testETHToAddr,
		Status:          1, // success
	}
}

func testFailedReceipt() *domainETH.TransactionReceipt {
	return &domainETH.TransactionReceipt{
		TransactionHash: testMonitorSentHash,
		BlockNumber:     100,
		From:            testETHFromAddr,
		To:              testETHToAddr,
		Status:          0, // failure/revert
		RevertReason:    "execution reverted: insufficient balance",
	}
}

// --- UpdateTxStatus tests ---

func TestETHMonitorTransactionUseCase_UpdateTxStatus_GetSentHashTxError(t *testing.T) {
	t.Parallel()
	deps := newEthMonitorTxTestDeps(t)
	useCase := newEthMonitorTxUseCase(deps)

	deps.txDetailRepo.EXPECT().
		GetSentHashTx(domainTx.TxTypeSent).
		Return(nil, errors.New("db connection error"))

	err := useCase.UpdateTxStatus(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "fail to call updateStatusTxTypeSent()")
}

func TestETHMonitorTransactionUseCase_UpdateTxStatus_NoSentTransactions(t *testing.T) {
	t.Parallel()
	deps := newEthMonitorTxTestDeps(t)
	useCase := newEthMonitorTxUseCase(deps)

	deps.txDetailRepo.EXPECT().
		GetSentHashTx(domainTx.TxTypeSent).
		Return([]string{}, nil)

	err := useCase.UpdateTxStatus(context.Background())
	require.NoError(t, err)
}

// Task 9.1: Transaction not yet mined (receipt is nil) — skip without error.
func TestETHMonitorTransactionUseCase_UpdateTxStatus_TxNotYetMined(t *testing.T) {
	t.Parallel()
	deps := newEthMonitorTxTestDeps(t)
	useCase := newEthMonitorTxUseCase(deps)

	deps.txDetailRepo.EXPECT().
		GetSentHashTx(domainTx.TxTypeSent).
		Return([]string{testMonitorSentHash}, nil)
	deps.ethClient.EXPECT().
		GetTxReceipt(context.Background(), testMonitorSentHash).
		Return(nil, nil) // nil, nil means not yet mined

	err := useCase.UpdateTxStatus(context.Background())
	require.NoError(t, err)
}

// Task 9.2: Transaction fails on-chain (status=0) — transition to TxTypeCancel and free sender address.
func TestETHMonitorTransactionUseCase_UpdateTxStatus_FailedTransaction(t *testing.T) {
	t.Parallel()
	deps := newEthMonitorTxTestDeps(t)
	useCase := newEthMonitorTxUseCase(deps)

	deps.txDetailRepo.EXPECT().
		GetSentHashTx(domainTx.TxTypeSent).
		Return([]string{testMonitorSentHash}, nil)
	deps.ethClient.EXPECT().
		GetTxReceipt(context.Background(), testMonitorSentHash).
		Return(testFailedReceipt(), nil)
	deps.txDetailRepo.EXPECT().
		UpdateTxTypeBySentHashTx(domainTx.TxTypeCancel, testMonitorSentHash).
		Return(int64(1), nil)
	deps.addrRepo.EXPECT().
		UpdateIsAllocated(false, testETHFromAddr).
		Return(int64(1), nil)

	err := useCase.UpdateTxStatus(context.Background())
	require.NoError(t, err)
}

// Task 9.1: Successful transaction but insufficient confirmations — no status change.
func TestETHMonitorTransactionUseCase_UpdateTxStatus_InsufficientConfirmations(t *testing.T) {
	t.Parallel()
	deps := newEthMonitorTxTestDeps(t)
	useCase := newEthMonitorTxUseCase(deps)

	deps.txDetailRepo.EXPECT().
		GetSentHashTx(domainTx.TxTypeSent).
		Return([]string{testMonitorSentHash}, nil)
	deps.ethClient.EXPECT().
		GetTxReceipt(context.Background(), testMonitorSentHash).
		Return(testSuccessReceipt(), nil)
	deps.ethClient.EXPECT().
		GetConfirmation(context.Background(), testMonitorSentHash).
		Return(uint64(3), nil) // 3 < 6 threshold

	err := useCase.UpdateTxStatus(context.Background())
	require.NoError(t, err)
}

// Task 9.1: Confirmed successful transaction — transition to TxTypeDone and update is_allocated.
func TestETHMonitorTransactionUseCase_UpdateTxStatus_ConfirmedSuccess(t *testing.T) {
	t.Parallel()
	deps := newEthMonitorTxTestDeps(t)
	useCase := newEthMonitorTxUseCase(deps)

	deps.txDetailRepo.EXPECT().
		GetSentHashTx(domainTx.TxTypeSent).
		Return([]string{testMonitorSentHash}, nil)
	deps.ethClient.EXPECT().
		GetTxReceipt(context.Background(), testMonitorSentHash).
		Return(testSuccessReceipt(), nil)
	deps.ethClient.EXPECT().
		GetConfirmation(context.Background(), testMonitorSentHash).
		Return(uint64(8), nil) // 8 >= 6 threshold
	deps.txDetailRepo.EXPECT().
		UpdateTxTypeBySentHashTx(domainTx.TxTypeDone, testMonitorSentHash).
		Return(int64(1), nil)
	deps.addrRepo.EXPECT().
		UpdateIsAllocated(false, testETHFromAddr).
		Return(int64(1), nil)

	err := useCase.UpdateTxStatus(context.Background())
	require.NoError(t, err)
}

// Task 9.2: Node connectivity failure with cancelled context — retries respect ctx.Done().
func TestETHMonitorTransactionUseCase_UpdateTxStatus_ConnectivityFailure(t *testing.T) {
	t.Parallel()
	deps := newEthMonitorTxTestDeps(t)
	useCase := newEthMonitorTxUseCase(deps)

	connectivityErr := errors.New("connection refused")

	deps.txDetailRepo.EXPECT().
		GetSentHashTx(domainTx.TxTypeSent).
		Return([]string{testMonitorSentHash}, nil)
	// With a pre-cancelled context, the backoff select returns ctx.Done immediately after 1st failure
	deps.ethClient.EXPECT().
		GetTxReceipt(mock.Anything, testMonitorSentHash).
		Return(nil, connectivityErr).
		Times(1)

	// Cancel context immediately so backoff select returns ctx.Err instead of sleeping
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// UpdateTxStatus warns on per-tx failure but does not propagate it
	err := useCase.UpdateTxStatus(ctx)
	require.NoError(t, err) // per-tx errors are logged, not returned
}

// UpdateTxTypeBySentHashTx failure on success path — error is propagated.
func TestETHMonitorTransactionUseCase_UpdateTxStatus_UpdateTxTypeDoneError(t *testing.T) {
	t.Parallel()
	deps := newEthMonitorTxTestDeps(t)
	useCase := newEthMonitorTxUseCase(deps)

	deps.txDetailRepo.EXPECT().
		GetSentHashTx(domainTx.TxTypeSent).
		Return([]string{testMonitorSentHash}, nil)
	deps.ethClient.EXPECT().
		GetTxReceipt(context.Background(), testMonitorSentHash).
		Return(testSuccessReceipt(), nil)
	deps.ethClient.EXPECT().
		GetConfirmation(context.Background(), testMonitorSentHash).
		Return(uint64(8), nil)
	deps.txDetailRepo.EXPECT().
		UpdateTxTypeBySentHashTx(domainTx.TxTypeDone, testMonitorSentHash).
		Return(int64(0), errors.New("db write error"))

	// per-tx error is logged at WARN, not propagated by UpdateTxStatus
	err := useCase.UpdateTxStatus(context.Background())
	require.NoError(t, err)
}

// UpdateIsAllocated failure on success path — error is propagated.
func TestETHMonitorTransactionUseCase_UpdateTxStatus_UpdateIsAllocatedError(t *testing.T) {
	t.Parallel()
	deps := newEthMonitorTxTestDeps(t)
	useCase := newEthMonitorTxUseCase(deps)

	deps.txDetailRepo.EXPECT().
		GetSentHashTx(domainTx.TxTypeSent).
		Return([]string{testMonitorSentHash}, nil)
	deps.ethClient.EXPECT().
		GetTxReceipt(context.Background(), testMonitorSentHash).
		Return(testSuccessReceipt(), nil)
	deps.ethClient.EXPECT().
		GetConfirmation(context.Background(), testMonitorSentHash).
		Return(uint64(8), nil)
	deps.txDetailRepo.EXPECT().
		UpdateTxTypeBySentHashTx(domainTx.TxTypeDone, testMonitorSentHash).
		Return(int64(1), nil)
	deps.addrRepo.EXPECT().
		UpdateIsAllocated(false, testETHFromAddr).
		Return(int64(0), errors.New("db write error"))

	// per-tx error is logged at WARN, not propagated by UpdateTxStatus
	err := useCase.UpdateTxStatus(context.Background())
	require.NoError(t, err)
}

// Task 9.1: GetConfirmation fails — warn and continue.
func TestETHMonitorTransactionUseCase_UpdateTxStatus_GetConfirmationError(t *testing.T) {
	t.Parallel()
	deps := newEthMonitorTxTestDeps(t)
	useCase := newEthMonitorTxUseCase(deps)

	deps.txDetailRepo.EXPECT().
		GetSentHashTx(domainTx.TxTypeSent).
		Return([]string{testMonitorSentHash}, nil)
	deps.ethClient.EXPECT().
		GetTxReceipt(context.Background(), testMonitorSentHash).
		Return(testSuccessReceipt(), nil)
	deps.ethClient.EXPECT().
		GetConfirmation(context.Background(), testMonitorSentHash).
		Return(uint64(0), errors.New("node not responding"))

	err := useCase.UpdateTxStatus(context.Background())
	// per-tx error is logged at WARN, not propagated
	require.NoError(t, err)
}

// --- MonitorBalance tests ---

func TestETHMonitorTransactionUseCase_MonitorBalance_GetAllAddressError(t *testing.T) {
	t.Parallel()
	deps := newEthMonitorTxTestDeps(t)
	useCase := newEthMonitorTxUseCase(deps)

	deps.addrRepo.EXPECT().
		GetAllAddress(domainAccount.AccountTypeClient).
		Return(nil, errors.New("db error"))

	err := useCase.MonitorBalance(context.Background(), watchusecase.MonitorBalanceInput{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "fail to call addrRepo.GetAllAddress()")
}

func TestETHMonitorTransactionUseCase_MonitorBalance_Success(t *testing.T) {
	t.Parallel()
	deps := newEthMonitorTxTestDeps(t)
	useCase := newEthMonitorTxUseCase(deps)

	addrs := []string{testETHFromAddr}
	balance := big.NewInt(1000000000000000000) // 1 ETH

	for _, acnt := range []domainAccount.AccountType{
		domainAccount.AccountTypeClient,
		domainAccount.AccountTypeDeposit,
		domainAccount.AccountTypePayment,
		domainAccount.AccountTypeStored,
	} {
		deps.addrRepo.EXPECT().
			GetAllAddress(acnt).
			Return(addrs, nil)
		deps.ethClient.EXPECT().
			GetTotalBalance(context.Background(), addrs).
			Return(balance, []domainETH.UserAmount{{Address: testETHFromAddr, Amount: uint64(balance.Int64())}})
	}

	err := useCase.MonitorBalance(context.Background(), watchusecase.MonitorBalanceInput{})
	require.NoError(t, err)
}
