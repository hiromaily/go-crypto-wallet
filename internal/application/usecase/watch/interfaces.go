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

// ImportDescriptorUseCase imports descriptors into watch wallet
type ImportDescriptorUseCase interface {
	Import(ctx context.Context, input ImportDescriptorInput) (ImportDescriptorOutput, error)
}

// CreatePaymentRequestUseCase creates payment requests
type CreatePaymentRequestUseCase interface {
	Execute(ctx context.Context, input CreatePaymentRequestInput) error
}

// AggregateMuSig2SignaturesUseCase aggregates MuSig2 partial signatures (BTC only)
type AggregateMuSig2SignaturesUseCase interface {
	Execute(ctx context.Context, input AggregateMuSig2SignaturesInput) (AggregateMuSig2SignaturesOutput, error)
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

// ImportDescriptorInput contains import parameters
type ImportDescriptorInput struct {
	FilePath     string
	AccountType  domainAccount.AccountType
	StartIndex   uint32
	Count        uint32
	ValidateOnly bool
}

// ImportDescriptorOutput contains import results
type ImportDescriptorOutput struct {
	DescriptorsImported int
	AddressesGenerated  int
	Errors              []string
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

// PartialSignatureData represents a partial signature from a signer
type PartialSignatureData struct {
	SignerID        string
	Signature       [32]byte // Partial signature scalar (S)
	NonceCommitment [33]byte // Public nonce commitment (R), compressed
}

// AggregateMuSig2SignaturesInput represents input for aggregating MuSig2 signatures
type AggregateMuSig2SignaturesInput struct {
	PSBTBase64          string
	PartialSignatures   []PartialSignatureData
	SignerPublicKeys    [][33]byte // Public keys of all participating signers
	CombinedNonce       [66]byte   // Combined public nonce from all signers
	AggregatedPublicKey [33]byte   // Compressed aggregated public key
	MessageHash         [32]byte
}

// AggregateMuSig2SignaturesOutput represents output from aggregating MuSig2 signatures
type AggregateMuSig2SignaturesOutput struct {
	FinalPSBT  string
	FinalTxHex string
	IsComplete bool
	TxID       string
}
