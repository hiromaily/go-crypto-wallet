package xrp

import (
	"encoding/json"
	"fmt"
	"time"
)

// TransactionFileVersion is the current version of the transaction file format.
const TransactionFileVersion = 1

// TransactionFile represents the JSON structure for XRP transaction files.
// This format aligns with BTC's PSBT concept for offline signing support.
type TransactionFile struct {
	// Version is the file format version for future compatibility
	Version int `json:"version"`

	// Chain identifies the blockchain (always "xrp")
	Chain string `json:"chain"`

	// Network identifies the network (mainnet, testnet, devnet)
	Network string `json:"network"`

	// CreatedAt is the timestamp when the file was created
	CreatedAt time.Time `json:"created_at"`

	// ActionType is the transaction action (deposit, payment, transfer)
	ActionType string `json:"action_type"`

	// SenderAccountType is the account type of the sender (client, deposit, payment, etc.)
	SenderAccountType string `json:"sender_account_type"`

	// Transactions is the list of transactions in this file
	Transactions []TransactionEntry `json:"transactions"`
}

// TransactionEntry represents a single transaction in the file.
type TransactionEntry struct {
	// UUID is the unique identifier for tracking this transaction
	UUID string `json:"uuid"`

	// UnsignedTx is the unsigned transaction JSON (present for unsigned transactions)
	UnsignedTx *TxInput `json:"unsigned_tx,omitempty"`

	// SignedTxID is the transaction ID after signing
	SignedTxID string `json:"signed_tx_id,omitempty"`

	// SignedTxBlob is the signed transaction blob (hex)
	SignedTxBlob string `json:"signed_tx_blob,omitempty"`

	// SenderAccount is the sender's XRP address
	SenderAccount string `json:"sender_account"`

	// SenderAccountType is the account type (client, deposit, payment, etc.)
	SenderAccountType string `json:"sender_account_type"`

	// ReceiverAccount is the receiver's XRP address
	ReceiverAccount string `json:"receiver_account"`

	// ReceiverAccountType is the receiver's account type
	ReceiverAccountType string `json:"receiver_account_type"`

	// Amount is the transaction amount in drops
	Amount string `json:"amount"`

	// SignatureCount is the number of signatures applied
	SignatureCount int `json:"signature_count"`

	// RequiredSignatures is the number of signatures required (1 for single-sig, more for multi-sig)
	RequiredSignatures int `json:"required_signatures"`
}

// NewTransactionFile creates a new transaction file with default values.
func NewTransactionFile(network, actionType, senderAccountType string) *TransactionFile {
	return &TransactionFile{
		Version:           TransactionFileVersion,
		Chain:             "xrp",
		Network:           network,
		CreatedAt:         time.Now().UTC(),
		ActionType:        actionType,
		SenderAccountType: senderAccountType,
		Transactions:      make([]TransactionEntry, 0),
	}
}

// AddTransaction adds a new unsigned transaction entry to the file.
func (f *TransactionFile) AddTransaction(
	uuid string,
	unsignedTx *TxInput,
	senderAccount, senderAccountType string,
	receiverAccount, receiverAccountType string,
	amount string,
	requiredSignatures int,
) {
	entry := TransactionEntry{
		UUID:                uuid,
		UnsignedTx:          unsignedTx,
		SenderAccount:       senderAccount,
		SenderAccountType:   senderAccountType,
		ReceiverAccount:     receiverAccount,
		ReceiverAccountType: receiverAccountType,
		Amount:              amount,
		SignatureCount:      0,
		RequiredSignatures:  requiredSignatures,
	}
	f.Transactions = append(f.Transactions, entry)
}

// ToJSON serializes the transaction file to JSON.
func (f *TransactionFile) ToJSON() ([]byte, error) {
	data, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transaction file: %w", err)
	}
	return data, nil
}

// TransactionFileFromJSON deserializes a transaction file from JSON.
func TransactionFileFromJSON(data []byte) (*TransactionFile, error) {
	var file TransactionFile
	if err := json.Unmarshal(data, &file); err != nil {
		return nil, fmt.Errorf("failed to unmarshal transaction file: %w", err)
	}
	return &file, nil
}

// IsUnsigned returns true if the transaction file contains unsigned transactions.
func (f *TransactionFile) IsUnsigned() bool {
	for _, tx := range f.Transactions {
		if tx.UnsignedTx != nil && tx.SignedTxBlob == "" {
			return true
		}
	}
	return false
}

// IsSigned returns true if all transactions in the file are signed.
func (f *TransactionFile) IsSigned() bool {
	for _, tx := range f.Transactions {
		if tx.SignedTxBlob == "" {
			return false
		}
	}
	return len(f.Transactions) > 0
}

// IsComplete returns true if all transactions have the required signatures.
func (f *TransactionFile) IsComplete() bool {
	for _, tx := range f.Transactions {
		if tx.SignatureCount < tx.RequiredSignatures {
			return false
		}
	}
	return len(f.Transactions) > 0
}

// SignedTransactionFile represents a file containing signed transactions.
// This is used after the signing process is complete.
type SignedTransactionFile struct {
	// Version is the file format version
	Version int `json:"version"`

	// Chain identifies the blockchain
	Chain string `json:"chain"`

	// Network identifies the network
	Network string `json:"network"`

	// CreatedAt is when the original unsigned file was created
	CreatedAt time.Time `json:"created_at"`

	// SignedAt is when the file was signed
	SignedAt time.Time `json:"signed_at"`

	// ActionType is the transaction action
	ActionType string `json:"action_type"`

	// SenderAccountType is the sender's account type
	SenderAccountType string `json:"sender_account_type"`

	// Transactions is the list of signed transactions
	Transactions []SignedTransactionEntry `json:"transactions"`
}

// SignedTransactionEntry represents a signed transaction entry.
type SignedTransactionEntry struct {
	// UUID is the unique identifier for tracking
	UUID string `json:"uuid"`

	// SignedTxID is the transaction hash/ID
	SignedTxID string `json:"signed_tx_id"`

	// SignedTxBlob is the signed transaction blob (hex)
	SignedTxBlob string `json:"signed_tx_blob"`

	// SenderAccount is the sender's XRP address
	SenderAccount string `json:"sender_account"`

	// SenderAccountType is the sender's account type
	SenderAccountType string `json:"sender_account_type"`

	// ReceiverAccount is the receiver's XRP address
	ReceiverAccount string `json:"receiver_account"`

	// ReceiverAccountType is the receiver's account type
	ReceiverAccountType string `json:"receiver_account_type"`

	// Amount is the transaction amount in drops
	Amount string `json:"amount"`

	// LastLedgerSequence is the max ledger version for this tx
	LastLedgerSequence uint64 `json:"last_ledger_sequence"`

	// SignatureCount is the number of signatures applied
	SignatureCount int `json:"signature_count"`

	// RequiredSignatures is the number of signatures required
	RequiredSignatures int `json:"required_signatures"`
}

// NewSignedTransactionFile creates a signed transaction file from an unsigned one.
func NewSignedTransactionFile(unsignedFile *TransactionFile) *SignedTransactionFile {
	return &SignedTransactionFile{
		Version:           unsignedFile.Version,
		Chain:             unsignedFile.Chain,
		Network:           unsignedFile.Network,
		CreatedAt:         unsignedFile.CreatedAt,
		SignedAt:          time.Now().UTC(),
		ActionType:        unsignedFile.ActionType,
		SenderAccountType: unsignedFile.SenderAccountType,
		Transactions:      make([]SignedTransactionEntry, 0),
	}
}

// AddSignedTransaction adds a signed transaction entry.
func (f *SignedTransactionFile) AddSignedTransaction(
	uuid, signedTxID, signedTxBlob string,
	senderAccount, senderAccountType string,
	receiverAccount, receiverAccountType string,
	amount string,
	lastLedgerSequence uint64,
	signatureCount, requiredSignatures int,
) {
	entry := SignedTransactionEntry{
		UUID:                uuid,
		SignedTxID:          signedTxID,
		SignedTxBlob:        signedTxBlob,
		SenderAccount:       senderAccount,
		SenderAccountType:   senderAccountType,
		ReceiverAccount:     receiverAccount,
		ReceiverAccountType: receiverAccountType,
		Amount:              amount,
		LastLedgerSequence:  lastLedgerSequence,
		SignatureCount:      signatureCount,
		RequiredSignatures:  requiredSignatures,
	}
	f.Transactions = append(f.Transactions, entry)
}

// ToJSON serializes the signed transaction file to JSON.
func (f *SignedTransactionFile) ToJSON() ([]byte, error) {
	data, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal signed transaction file: %w", err)
	}
	return data, nil
}

// SignedTransactionFileFromJSON deserializes a signed transaction file from JSON.
func SignedTransactionFileFromJSON(data []byte) (*SignedTransactionFile, error) {
	var file SignedTransactionFile
	if err := json.Unmarshal(data, &file); err != nil {
		return nil, fmt.Errorf("failed to unmarshal signed transaction file: %w", err)
	}
	return &file, nil
}
