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
	MultisigQuorum    uint32  // 0 = single-sig (default); >1 = multisig JSON-format file
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

// XRP Multi-signature and Regular Key Use Cases

// SetRegularKeyUseCase creates SetRegularKey transactions for XRP accounts.
// This allows an account to authorize a regular key for signing transactions.
type SetRegularKeyUseCase interface {
	Execute(ctx context.Context, input SetRegularKeyInput) (SetRegularKeyOutput, error)
}

// SetSignerListUseCase creates SignerListSet transactions for XRP multi-sig accounts.
// This configures a signer list enabling multi-signature transactions.
type SetSignerListUseCase interface {
	Execute(ctx context.Context, input SetSignerListInput) (SetSignerListOutput, error)
}

// CreateMultisigTxUseCase creates a pending multi-sig transaction awaiting signatures.
type CreateMultisigTxUseCase interface {
	Execute(ctx context.Context, input CreateMultisigTxInput) (CreateMultisigTxOutput, error)
}

// AddMultisigSignatureUseCase adds a signature to a pending multi-sig transaction.
type AddMultisigSignatureUseCase interface {
	Execute(ctx context.Context, input AddMultisigSignatureInput) (AddMultisigSignatureOutput, error)
}

// SubmitMultisigTxUseCase submits a ready multi-sig transaction to the network.
type SubmitMultisigTxUseCase interface {
	Execute(ctx context.Context, input SubmitMultisigTxInput) (SubmitMultisigTxOutput, error)
}

// SetRegularKeyInput represents input for setting a regular key
type SetRegularKeyInput struct {
	AccountAddress    string // XRP account to set regular key for
	RegularKeyAddress string // Regular key address to authorize (empty to remove)
}

// SetRegularKeyOutput represents output from setting a regular key
type SetRegularKeyOutput struct {
	FileName  string // Generated transaction file
	TxJSON    string // Transaction JSON
	AccountID string // Account that was modified
}

// SignerEntry represents a signer for multi-sig configuration
type SignerEntry struct {
	Account string
	Weight  uint32
}

// SetSignerListInput represents input for setting a signer list
type SetSignerListInput struct {
	AccountAddress    string        // XRP account to configure multi-sig for
	SenderAccountType string        // Account type label (e.g. "payment") for key lookup by keygen signer
	SignerQuorum      uint32        // Minimum weight required for signatures
	SignerEntries     []SignerEntry // List of authorized signers with weights
}

// SetSignerListOutput represents output from setting a signer list
type SetSignerListOutput struct {
	FileName     string // Generated transaction file
	TxJSON       string // Transaction JSON
	SignerListID int64  // Database ID of the created signer list
}

// CreateMultisigTxInput represents input for creating a pending multi-sig transaction
type CreateMultisigTxInput struct {
	AccountAddress  string  // XRP account sending the transaction
	ReceiverAddress string  // Destination address
	Amount          float64 // Amount in XRP
	TxType          string  // Transaction type (e.g., "Payment")
}

// CreateMultisigTxOutput represents output from creating a pending multi-sig transaction
type CreateMultisigTxOutput struct {
	TxUUID         string // Unique identifier for tracking
	PendingID      int64  // Database ID of the pending transaction
	UnsignedTxJSON string // JSON of the unsigned transaction
	RequiredQuorum uint32 // Quorum needed to sign
}

// AddMultisigSignatureInput represents input for adding a signature
type AddMultisigSignatureInput struct {
	TxUUID        string // Transaction UUID to sign
	SignerAccount string // Account of the signer
	SignedTxBlob  string // Signed transaction blob from this signer
}

// AddMultisigSignatureOutput represents output from adding a signature
type AddMultisigSignatureOutput struct {
	CurrentWeight  uint32 // Total weight after adding signature
	RequiredQuorum uint32 // Quorum needed
	IsReady        bool   // True if quorum is met
	CombinedTxBlob string // Combined tx blob (only set when ready)
}

// SubmitMultisigTxInput represents input for submitting a multi-sig transaction
type SubmitMultisigTxInput struct {
	TxUUID string // Transaction UUID to submit
}

// SubmitMultisigTxOutput represents output from submitting a multi-sig transaction
type SubmitMultisigTxOutput struct {
	TxHash   string // Submitted transaction hash
	IsQueued bool   // True if queued for validation
}

// ETH Safe Multisig Use Cases

// CreateETHMultisigTransactionUseCase creates an unsigned Safe multisig transaction proposal file.
// No database record is created; the output is a JSON file written to disk.
type CreateETHMultisigTransactionUseCase interface {
	Execute(ctx context.Context, input CreateETHMultisigTransactionInput) (CreateETHMultisigTransactionOutput, error)
}

// SendETHMultisigTransactionUseCase submits a fully signed Safe multisig transaction to Ethereum.
type SendETHMultisigTransactionUseCase interface {
	Execute(ctx context.Context, input SendETHMultisigTransactionInput) (SendETHMultisigTransactionOutput, error)
}

// ETHSafeInfoUseCase retrieves on-chain Safe contract state.
type ETHSafeInfoUseCase interface {
	Execute(ctx context.Context, input ETHSafeInfoInput) (ETHSafeInfoOutput, error)
}

// CreateETHMultisigTransactionInput is the input for the ETH multisig proposal creation use case.
type CreateETHMultisigTransactionInput struct {
	SafeAddress string  // EIP-55 checksummed Safe proxy address
	To          string  // Recipient address (EIP-55 checksum)
	AmountEther float64 // Amount in Ether (converted to Wei internally)
	Threshold   int     // Required signer count (m in m-of-n)
	ActionType  string  // "deposit", "payment", "transfer"
}

// CreateETHMultisigTransactionOutput is the output of the ETH multisig proposal creation use case.
type CreateETHMultisigTransactionOutput struct {
	FilePath string // Absolute path of the written unsigned multisig JSON file
	UUID     string // Generated proposal UUID
}

// SendETHMultisigTransactionInput is the input for the ETH multisig submission use case.
type SendETHMultisigTransactionInput struct {
	FilePath string // Path to the fully signed multisig JSON file
}

// SendETHMultisigTransactionOutput is the output of the ETH multisig submission use case.
type SendETHMultisigTransactionOutput struct {
	TxHash string // Submitted Ethereum transaction hash
}

// ETHSafeInfoInput is the input for the ETH Safe info use case.
type ETHSafeInfoInput struct {
	SafeAddress string // EIP-55 checksummed Safe proxy address
}

// ETHSafeInfoOutput is the output of the ETH Safe info use case.
type ETHSafeInfoOutput struct {
	Owners     []string // EIP-55 checksummed owner addresses
	Threshold  uint64   // Current signature threshold
	Nonce      uint64   // Current Safe nonce
	BalanceWei string   // Safe balance in Wei (decimal string)
}
