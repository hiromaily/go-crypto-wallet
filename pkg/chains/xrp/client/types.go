// Package client provides a WebSocket client for XRPL node communication.
// It wraps the xrpscan/xrpl-go library and implements the small interfaces
// defined in internal/application/ports/api/xrp.
package client

// AccountInfo represents account information retrieved from the XRP Ledger.
// Returned by Client.GetAccountInfo.
type AccountInfo struct {
	Sequence                                  uint64
	XrpBalance                                string
	OwnerCount                                uint64
	PreviousAffectingTransactionID            string
	PreviousAffectingTransactionLedgerVersion uint64
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

// SentTx represents the result of a submitted XRP transaction.
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

// TxSpecDestination is part of TxInfo.
type TxSpecDestination struct {
	Address string
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
