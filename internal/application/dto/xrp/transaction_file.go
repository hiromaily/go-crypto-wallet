package xrp

import (
	"errors"
	"regexp"
)

// XRPTransactionFile represents a JSON transaction file structure for XRP transactions.
// This format supports offline signing and multi-signature workflows aligned with Bitcoin PSBT patterns.
type XRPTransactionFile struct {
	Version      string                `json:"version"`      // Semantic versioning (e.g., "1.0.0")
	Chain        string                `json:"chain"`        // Always "XRP"
	Network      string                `json:"network"`      // "mainnet" | "testnet"
	CreatedAt    string                `json:"createdAt"`    // ISO 8601 timestamp
	Transactions []XRPTransactionEntry `json:"transactions"` // Transaction entries
}

// XRPTransactionEntry represents a single transaction entry in the transaction file.
type XRPTransactionEntry struct {
	UUID               string  `json:"uuid"`               // UUIDv7 for traceability
	UnsignedData       TxInput `json:"unsignedData"`       // Full TxInput structure
	SenderAccount      string  `json:"senderAccount"`      // XRP classic address
	SenderAccountType  string  `json:"senderAccountType"`  // Domain account type
	SignatureCount     int     `json:"signatureCount"`     // Current signatures
	RequiredSignatures int     `json:"requiredSignatures"` // Quorum threshold
	SignedBlob         *string `json:"signedBlob"`         // Hex-encoded signed tx (null if unsigned)
	IsComplete         bool    `json:"isComplete"`         // Completion flag
}

var (
	// semVerPattern matches semantic versioning pattern (Major.Minor.Patch)
	semVerPattern = regexp.MustCompile(`^\d+\.\d+\.\d+$`)

	// ErrInvalidVersion is returned when version format is invalid
	ErrInvalidVersion = errors.New("version must follow semantic versioning pattern")
	// ErrInvalidChain is returned when chain is not XRP
	ErrInvalidChain = errors.New("chain must be XRP")
	// ErrInvalidNetwork is returned when network is not mainnet or testnet
	ErrInvalidNetwork = errors.New("network must be mainnet or testnet")
	// ErrEmptyTransactions is returned when transactions array is empty
	ErrEmptyTransactions = errors.New("transactions array cannot be empty")
	// ErrInvalidSignatureCount is returned when signature count is invalid
	ErrInvalidSignatureCount = errors.New("signatureCount must be >= 0")
	// ErrSignatureCountExceeds is returned when signature count exceeds required
	ErrSignatureCountExceeds = errors.New("signatureCount cannot exceed requiredSignatures")
	// ErrSignedBlobNotNull is returned when signedBlob is not null with zero signatures
	ErrSignedBlobNotNull = errors.New("signedBlob must be null when signatureCount is 0")
)

// Validate validates the XRPTransactionFile structure and enforces all invariants.
func (f *XRPTransactionFile) Validate() error {
	// Validate version format (semantic versioning: Major.Minor.Patch)
	if !semVerPattern.MatchString(f.Version) {
		return ErrInvalidVersion
	}

	// Validate chain
	if f.Chain != "XRP" {
		return ErrInvalidChain
	}

	// Validate network
	if f.Network != "mainnet" && f.Network != "testnet" {
		return ErrInvalidNetwork
	}

	// Validate transactions array is not empty
	if len(f.Transactions) == 0 {
		return ErrEmptyTransactions
	}

	// Validate each transaction entry
	for _, tx := range f.Transactions {
		if err := tx.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// Validate validates the XRPTransactionEntry structure and enforces all invariants.
func (e *XRPTransactionEntry) Validate() error {
	// Validate signature count is non-negative
	if e.SignatureCount < 0 {
		return ErrInvalidSignatureCount
	}

	// Validate signature count does not exceed required
	if e.SignatureCount > e.RequiredSignatures {
		return ErrSignatureCountExceeds
	}

	// Validate signedBlob is null when signatureCount is 0
	if e.SignatureCount == 0 && e.SignedBlob != nil {
		return ErrSignedBlobNotNull
	}

	return nil
}
