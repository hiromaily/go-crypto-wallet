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
	repomocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/watch/mocks"
	storagemocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/transaction/mocks"
	pkguuid "github.com/hiromaily/go-crypto-wallet/pkg/uuid"
)

// MockAccountInfoProvider is a mock implementation of AccountInfoProvider for testing
type MockAccountInfoProvider struct {
	mock.Mock
}

func (m *MockAccountInfoProvider) GetAccountInfo(
	ctx context.Context,
	address string,
) (*dtoxrp.ResponseGetAccountInfo, error) {
	args := m.Called(ctx, address)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dtoxrp.ResponseGetAccountInfo), args.Error(1)
}

func (m *MockAccountInfoProvider) GetBalance(ctx context.Context, addr string) (float64, error) {
	args := m.Called(ctx, addr)
	return args.Get(0).(float64), args.Error(1)
}

// MockTransactionPreparer is a mock implementation of TransactionPreparer for testing
type MockTransactionPreparer struct {
	mock.Mock
}

func (m *MockTransactionPreparer) CreateRawTransaction(
	ctx context.Context,
	senderAccount, receiverAccount string,
	amount float64,
	instructions *dtoxrp.Instructions,
) (*dtoxrp.TxInput, string, error) {
	args := m.Called(ctx, senderAccount, receiverAccount, amount, instructions)
	if args.Get(0) == nil {
		return nil, "", args.Error(2)
	}
	return args.Get(0).(*dtoxrp.TxInput), args.String(1), args.Error(2)
}

// testDependencies holds all mock dependencies for testing
type testDependencies struct {
	accountInfo  *MockAccountInfoProvider
	txPreparer   *MockTransactionPreparer
	addrRepo     *repomocks.MockAddressRepositorier
	txRepo       *repomocks.MockTxRepositorier
	txDetailRepo *repomocks.MockXRPDetailTXRepositorier
	payReqRepo   *repomocks.MockPaymentRequestRepositorier
	txFileRepo   *storagemocks.MockTransactionFileRepositorier
	uuidHandler  pkguuid.UUIDHandler
}

// newTestDependencies creates all mock dependencies
func newTestDependencies(t *testing.T) *testDependencies {
	t.Helper()
	return &testDependencies{
		accountInfo:  &MockAccountInfoProvider{},
		txPreparer:   &MockTransactionPreparer{},
		addrRepo:     repomocks.NewMockAddressRepositorier(t),
		txRepo:       repomocks.NewMockTxRepositorier(t),
		txDetailRepo: repomocks.NewMockXRPDetailTXRepositorier(t),
		payReqRepo:   repomocks.NewMockPaymentRequestRepositorier(t),
		txFileRepo:   storagemocks.NewMockTransactionFileRepositorier(t),
		uuidHandler:  pkguuid.NewGoogleUUIDHandler(),
	}
}

// createUseCase creates a new CreateTransactionUseCase with the test dependencies
func createUseCase(deps *testDependencies) watchusecase.CreateTransactionUseCase {
	return xrp.NewCreateTransactionUseCase(
		deps.accountInfo,
		deps.txPreparer,
		nil, // dbConn - nil is safe for tests that don't reach DB operations
		deps.uuidHandler,
		deps.addrRepo,
		deps.txRepo,
		deps.txDetailRepo,
		deps.payReqRepo,
		deps.txFileRepo,
		domainAccount.AccountTypeDeposit,
		domainAccount.AccountTypePayment,
	)
}

func TestCreateTransactionUseCase_Execute_InvalidActionType(t *testing.T) {
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
	deps := newTestDependencies(t)
	useCase := createUseCase(deps)

	senderAddr := &domainAddress.Address{
		ID:            1,
		WalletAddress: "rSenderAddress123",
	}

	// Setup mocks
	deps.addrRepo.On("GetOneUnAllocated", domainAccount.AccountTypePayment).Return(senderAddr, nil)
	deps.accountInfo.On("GetBalance", mock.Anything, "rSenderAddress123").Return(10.0, nil)

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
	deps.addrRepo.AssertExpectations(t)
	deps.accountInfo.AssertExpectations(t)
}

func TestCreateTransactionUseCase_Execute_TransferWithValidBalance(t *testing.T) {
	t.Skip("Skipping until DB transaction is mocked - requires comprehensive DB mock setup")
	// This test requires:
	// 1. Mock database connection with transaction support
	// 2. Mock all repository operations within the transaction
	// 3. Mock file write operations
	// This is complex and should be tested in integration tests instead
}

func TestCreateTransactionUseCase_Execute_GetBalanceError(t *testing.T) {
	deps := newTestDependencies(t)
	useCase := createUseCase(deps)

	senderAddr := &domainAddress.Address{
		ID:            1,
		WalletAddress: "rSenderAddress123",
	}

	// Setup mocks
	deps.addrRepo.On("GetOneUnAllocated", domainAccount.AccountTypePayment).Return(senderAddr, nil)
	deps.accountInfo.On("GetBalance", mock.Anything, "rSenderAddress123").
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
	deps.addrRepo.AssertExpectations(t)
	deps.accountInfo.AssertExpectations(t)
}

func TestCreateTransactionUseCase_Execute_GetAddressError(t *testing.T) {
	deps := newTestDependencies(t)
	useCase := createUseCase(deps)

	// Setup mocks
	deps.addrRepo.On("GetOneUnAllocated", domainAccount.AccountTypePayment).
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
	deps.addrRepo.AssertExpectations(t)
}

func TestCreateTransactionUseCase_Execute_NilSenderAddress(t *testing.T) {
	deps := newTestDependencies(t)
	useCase := createUseCase(deps)

	// Setup mocks - repository returns nil address without error (no unallocated address exists)
	deps.addrRepo.On("GetOneUnAllocated", domainAccount.AccountTypePayment).
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
	deps.addrRepo.AssertExpectations(t)
}

func TestCreateTransactionUseCase_Execute_NilReceiverAddress(t *testing.T) {
	deps := newTestDependencies(t)
	useCase := createUseCase(deps)

	senderAddr := &domainAddress.Address{
		ID:            1,
		WalletAddress: "rSenderAddress123",
	}

	// Setup mocks
	deps.addrRepo.On("GetOneUnAllocated", domainAccount.AccountTypePayment).Return(senderAddr, nil)
	deps.accountInfo.On("GetBalance", mock.Anything, "rSenderAddress123").Return(100.0, nil)
	// Receiver address returns nil without error
	deps.addrRepo.On("GetOneUnAllocated", domainAccount.AccountTypeDeposit).
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
	deps.addrRepo.AssertExpectations(t)
	deps.accountInfo.AssertExpectations(t)
}

func TestCreateTransactionUseCase_Dependencies(t *testing.T) {
	t.Run("uses segregated interfaces instead of full XRPer interface", func(t *testing.T) {
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
		)

		assert.NotNil(t, useCase, "use case should be created with segregated interfaces")
	})
}
