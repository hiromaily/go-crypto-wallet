package watch

import (
	"context"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
)

// CreateTransactionUseCase creates unsigned transactions
type CreateTransactionUseCase interface {
	Execute(ctx context.Context, input CreateTransactionInput) (CreateTransactionOutput, error)
}

// MonitorTransactionUseCase monitors transaction status and balances
type MonitorTransactionUseCase interface {
	UpdateTxStatus(ctx context.Context) error
	MonitorBalance(ctx context.Context, input MonitorBalanceInput) error
}

// SendTransactionUseCase sends signed transactions to the network
type SendTransactionUseCase interface {
	Execute(ctx context.Context, input SendTransactionInput) (SendTransactionOutput, error)
}

// ImportAddressUseCase imports addresses from files
type ImportAddressUseCase interface {
	Execute(ctx context.Context, input ImportAddressInput) error
}

// CreatePaymentRequestUseCase creates payment requests
type CreatePaymentRequestUseCase interface {
	Execute(ctx context.Context, input CreatePaymentRequestInput) error
}

// Input/Output DTOs

// CreateTransactionInput represents input for creating a transaction
type CreateTransactionInput struct {
	ActionType        string // "deposit", "payment", "transfer"
	SenderAccount     domainAccount.AccountType
	ReceiverAccount   domainAccount.AccountType
	Amount            float64
	AdjustmentFee     float64
	PaymentRequestIDs []int64 // For payment transactions
}

// CreateTransactionOutput represents output from creating a transaction
type CreateTransactionOutput struct {
	TransactionHex string
	FileName       string
}

// MonitorBalanceInput represents input for monitoring balance
type MonitorBalanceInput struct {
	ConfirmationNum uint64
}

// SendTransactionInput represents input for sending a transaction
type SendTransactionInput struct {
	FilePath string
}

// SendTransactionOutput represents output from sending a transaction
type SendTransactionOutput struct {
	TxID string
}

// ImportAddressInput represents input for importing addresses
type ImportAddressInput struct {
	FileName string
	Rescan   bool
}

// CreatePaymentRequestInput represents input for creating payment requests
type CreatePaymentRequestInput struct {
	AmountList []float64
}
