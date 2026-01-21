package ripple

// Instructions represents transaction instructions for XRP operations.
type Instructions struct {
	Fee                    string
	MaxFee                 string
	MaxLedgerVersion       uint64
	MaxLedgerVersionOffset uint64
	Sequence               uint64
	SignersCount           uint64
}

// TxInput represents an XRP transaction input.
type TxInput struct {
	TransactionType    string
	Account            string
	Amount             string
	Destination        string
	Fee                string
	Flags              uint64
	LastLedgerSequence uint64
	Sequence           uint64
	SigningPubKey      string
	TxnSignature       string
	Hash               string
}

// SentTx represents the result of a sent XRP transaction.
type SentTx struct {
	ResultCode          string
	ResultMessage       string
	EngineResult        string
	EngineResultCode    int
	EngineResultMessage string
	TxBlob              string
	TxJSON              TxInput
}

// TxInfo represents detailed information about an XRP transaction.
type TxInfo struct {
	Type          string
	Address       string
	Sequence      int
	ID            string
	Specification TxSpecification
	Outcome       TxOutcome
}

// TxSpecification is part of TxInfo.
type TxSpecification struct {
	Source      TxSpecSource
	Destination TxSpecDestination
}

// TxSpecSource is part of TxInfo.
type TxSpecSource struct {
	Address   string
	MaxAmount TxAmount
}

// TxAmount is part of TxInfo.
type TxAmount struct {
	Currency string
	Value    string
}

// TxTotalPrice is part of TxInfo.
type TxTotalPrice struct {
	Currency     string
	Counterparty string
	Value        string
}

// TxSpecDestination is part of TxInfo.
type TxSpecDestination struct {
	Address string
}

// TxOutcome is part of TxInfo.
type TxOutcome struct {
	Result           string
	Timestamp        string
	Fee              string
	BalanceChanges   map[string][]TxAmount
	OrderbookChanges map[string][]TxOrderbookChange
	LedgerVersion    int
	IndexInLedger    int
	DeliveredAmount  TxAmount
}

// TxOrderbookChange is part of TxInfo.
type TxOrderbookChange struct {
	Direction         string
	Quantity          TxAmount
	TotalPrice        TxTotalPrice
	Sequence          int
	Status            string
	MakerExchangeRate string
}

// SignerListEntry represents a single signer in a SignerList.
// Reference: https://xrpl.org/docs/references/protocol/ledger-data/ledger-entry-types/signerlist
type SignerListEntry struct {
	Account      string
	SignerWeight uint32
}

// SignerListSetTxInput represents a SignerListSet transaction input.
// Reference: https://xrpl.org/docs/references/protocol/transactions/types/signerlistset
type SignerListSetTxInput struct {
	TransactionType    string
	Account            string
	SignerQuorum       uint32
	SignerEntries      []SignerListEntry
	Fee                string
	Flags              uint64
	LastLedgerSequence uint64
	Sequence           uint64
	SigningPubKey      string
	TxnSignature       string
	Hash               string
}

// IssuedCurrencyAmount represents a token amount with currency and issuer.
// Used for TrustSet and other token-related transactions.
// Reference: https://xrpl.org/docs/references/protocol/data-types/currency-formats#issued-currency-amount-format
type IssuedCurrencyAmount struct {
	Currency string
	Issuer   string
	Value    string
}

// TrustSetTxInput represents a TrustSet transaction input.
// Reference: https://xrpl.org/docs/references/protocol/transactions/types/trustset
type TrustSetTxInput struct {
	TransactionType    string
	Account            string
	LimitAmount        IssuedCurrencyAmount
	QualityIn          uint32
	QualityOut         uint32
	Fee                string
	Flags              uint64
	LastLedgerSequence uint64
	Sequence           uint64
	SigningPubKey      string
	TxnSignature       string
	Hash               string
}

// EscrowCreateTxInput represents an EscrowCreate transaction input.
// Reference: https://xrpl.org/docs/references/protocol/transactions/types/escrowcreate
type EscrowCreateTxInput struct {
	TransactionType    string
	Account            string
	Amount             string
	Destination        string
	CancelAfter        uint32
	FinishAfter        uint32
	Condition          string
	DestinationTag     uint32
	Fee                string
	Flags              uint64
	LastLedgerSequence uint64
	Sequence           uint64
	SigningPubKey      string
	TxnSignature       string
	Hash               string
}

// EscrowFinishTxInput represents an EscrowFinish transaction input.
// Reference: https://xrpl.org/docs/references/protocol/transactions/types/escrowfinish
type EscrowFinishTxInput struct {
	TransactionType    string
	Account            string
	Owner              string
	OfferSequence      uint32
	Condition          string
	Fulfillment        string
	Fee                string
	Flags              uint64
	LastLedgerSequence uint64
	Sequence           uint64
	SigningPubKey      string
	TxnSignature       string
	Hash               string
}

// EscrowCancelTxInput represents an EscrowCancel transaction input.
// Reference: https://xrpl.org/docs/references/protocol/transactions/types/escrowcancel
type EscrowCancelTxInput struct {
	TransactionType    string
	Account            string
	Owner              string
	OfferSequence      uint32
	Fee                string
	Flags              uint64
	LastLedgerSequence uint64
	Sequence           uint64
	SigningPubKey      string
	TxnSignature       string
	Hash               string
}

// PaymentChannelCreateTxInput represents a PaymentChannelCreate transaction input.
// Reference: https://xrpl.org/docs/references/protocol/transactions/types/paymentchannelcreate
type PaymentChannelCreateTxInput struct {
	TransactionType    string
	Account            string
	Amount             string
	Destination        string
	SettleDelay        uint32
	PublicKey          string
	CancelAfter        uint32
	DestinationTag     uint32
	SourceTag          uint32
	Fee                string
	Flags              uint64
	LastLedgerSequence uint64
	Sequence           uint64
	SigningPubKey      string
	TxnSignature       string
	Hash               string
}

// PaymentChannelFundTxInput represents a PaymentChannelFund transaction input.
// Reference: https://xrpl.org/docs/references/protocol/transactions/types/paymentchannelfund
type PaymentChannelFundTxInput struct {
	TransactionType    string
	Account            string
	Channel            string
	Amount             string
	Expiration         uint32
	Fee                string
	Flags              uint64
	LastLedgerSequence uint64
	Sequence           uint64
	SigningPubKey      string
	TxnSignature       string
	Hash               string
}

// PaymentChannelClaimTxInput represents a PaymentChannelClaim transaction input.
// Reference: https://xrpl.org/docs/references/protocol/transactions/types/paymentchannelclaim
type PaymentChannelClaimTxInput struct {
	TransactionType    string
	Account            string
	Channel            string
	Balance            string
	Amount             string
	Signature          string
	PublicKey          string
	Fee                string
	Flags              uint64
	LastLedgerSequence uint64
	Sequence           uint64
	SigningPubKey      string
	TxnSignature       string
	Hash               string
}

// XRPKeyType represents the type of XRP cryptographic key.
type XRPKeyType string

const (
	XRPKeyTypeSecp256k1 XRPKeyType = "secp256k1"
	XRPKeyTypeEd25519   XRPKeyType = "ed25519"
)
