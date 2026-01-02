package xrp

import (
	"errors"
	"time"

	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
)

// XrpDetailTx represents an XRP transaction detail in the domain layer
type XrpDetailTx struct {
	ID                    int64
	TxID                  int64
	UUID                  string
	CurrentTxType         domainTx.TxType
	SenderAccount         string
	SenderAddress         string
	ReceiverAccount       string
	ReceiverAddress       string
	Amount                string
	XrpTxType             string
	Fee                   string
	Flags                 uint64
	LastLedgerSequence    uint64
	Sequence              uint64
	SigningPubkey         string
	TxnSignature          string
	Hash                  string
	EarliestLedgerVersion uint64
	SignedTxID            string
	TxBlob                string
	SentUpdatedAt         *time.Time
}

// NewXrpDetailTx creates a new XrpDetailTx entity
func NewXrpDetailTx(
	txID int64,
	uuid string,
	currentTxType domainTx.TxType,
	senderAccount string,
	senderAddress string,
	receiverAccount string,
	receiverAddress string,
	amount string,
	xrpTxType string,
	fee string,
	flags uint64,
	lastLedgerSequence uint64,
	sequence uint64,
) (*XrpDetailTx, error) {
	if uuid == "" {
		return nil, errors.New("uuid cannot be empty")
	}
	if senderAddress == "" {
		return nil, errors.New("sender address cannot be empty")
	}
	if receiverAddress == "" {
		return nil, errors.New("receiver address cannot be empty")
	}
	if amount == "" {
		return nil, errors.New("amount cannot be empty")
	}

	return &XrpDetailTx{
		TxID:               txID,
		UUID:               uuid,
		CurrentTxType:      currentTxType,
		SenderAccount:      senderAccount,
		SenderAddress:      senderAddress,
		ReceiverAccount:    receiverAccount,
		ReceiverAddress:    receiverAddress,
		Amount:             amount,
		XrpTxType:          xrpTxType,
		Fee:                fee,
		Flags:              flags,
		LastLedgerSequence: lastLedgerSequence,
		Sequence:           sequence,
	}, nil
}

// SetSigningData sets signing-related data
func (x *XrpDetailTx) SetSigningData(signingPubkey, txnSignature, hash string) {
	x.SigningPubkey = signingPubkey
	x.TxnSignature = txnSignature
	x.Hash = hash
}

// SetSentData sets sent transaction data
func (x *XrpDetailTx) SetSentData(signedTxID, txBlob string, earliestLedgerVersion uint64) {
	x.SignedTxID = signedTxID
	x.TxBlob = txBlob
	x.EarliestLedgerVersion = earliestLedgerVersion
	now := time.Now()
	x.SentUpdatedAt = &now
}

// UpdateTxType updates the transaction type
func (x *XrpDetailTx) UpdateTxType(txType domainTx.TxType) {
	x.CurrentTxType = txType
}
