package xrp

import (
	"errors"
	"time"
)

// MultisigStatus represents the status of a pending multi-signature transaction.
type MultisigStatus string

const (
	// MultisigStatusPending indicates the transaction is awaiting signatures.
	MultisigStatusPending MultisigStatus = "pending"
	// MultisigStatusReady indicates enough signatures collected to meet quorum.
	MultisigStatusReady MultisigStatus = "ready"
	// MultisigStatusSubmitted indicates the combined transaction has been submitted.
	MultisigStatusSubmitted MultisigStatus = "submitted"
	// MultisigStatusConfirmed indicates the transaction was validated on the ledger.
	MultisigStatusConfirmed MultisigStatus = "confirmed"
	// MultisigStatusFailed indicates the transaction failed during submission or validation.
	MultisigStatusFailed MultisigStatus = "failed"
	// MultisigStatusExpired indicates the transaction expired before enough signatures.
	MultisigStatusExpired MultisigStatus = "expired"
)

// XRPPendingMultisig represents a multi-signature transaction awaiting signatures.
type XRPPendingMultisig struct {
	ID              int64
	TxUUID          string         // Unique identifier for the transaction
	AccountID       string         // XRP account address (r...) that will sign this transaction
	UnsignedTxJSON  string         // JSON representation of the unsigned transaction
	XRPTxType       string         // XRP transaction type (Payment, SignerListSet, etc.)
	RequiredQuorum  uint32         // Quorum required from signer list
	CurrentWeight   uint32         // Sum of collected signature weights
	Status          MultisigStatus // Current status of the pending transaction
	CombinedTxBlob  *string        // Combined signed transaction blob (when ready)
	SubmittedTxHash *string        // Transaction hash after submission
	ExpiresAt       *time.Time     // When this pending transaction expires
	CreatedAt       *time.Time     // When the pending transaction was created
	UpdatedAt       *time.Time     // When the pending transaction was last updated
}

// NewXRPPendingMultisig creates a new XRPPendingMultisig entity.
func NewXRPPendingMultisig(
	txUUID string,
	accountID string,
	unsignedTxJSON string,
	xrpTxType string,
	requiredQuorum uint32,
	expiresAt *time.Time,
) (*XRPPendingMultisig, error) {
	if txUUID == "" {
		return nil, errors.New("transaction UUID cannot be empty")
	}
	if accountID == "" {
		return nil, errors.New("account ID cannot be empty")
	}
	if unsignedTxJSON == "" {
		return nil, errors.New("unsigned transaction JSON cannot be empty")
	}
	if xrpTxType == "" {
		return nil, errors.New("XRP transaction type cannot be empty")
	}
	if requiredQuorum == 0 {
		return nil, errors.New("required quorum must be greater than 0")
	}

	now := time.Now()
	return &XRPPendingMultisig{
		TxUUID:         txUUID,
		AccountID:      accountID,
		UnsignedTxJSON: unsignedTxJSON,
		XRPTxType:      xrpTxType,
		RequiredQuorum: requiredQuorum,
		CurrentWeight:  0,
		Status:         MultisigStatusPending,
		ExpiresAt:      expiresAt,
		CreatedAt:      &now,
	}, nil
}

// AddWeight adds signature weight and checks if quorum is met.
// Returns true if quorum is now met.
func (p *XRPPendingMultisig) AddWeight(weight uint32) bool {
	p.CurrentWeight += weight
	now := time.Now()
	p.UpdatedAt = &now

	if p.CurrentWeight >= p.RequiredQuorum {
		p.Status = MultisigStatusReady
		return true
	}
	return false
}

// SetReady marks the transaction as ready for submission.
func (p *XRPPendingMultisig) SetReady(combinedTxBlob string) {
	p.Status = MultisigStatusReady
	p.CombinedTxBlob = &combinedTxBlob
	now := time.Now()
	p.UpdatedAt = &now
}

// SetSubmitted marks the transaction as submitted.
func (p *XRPPendingMultisig) SetSubmitted(txHash string) {
	p.Status = MultisigStatusSubmitted
	p.SubmittedTxHash = &txHash
	now := time.Now()
	p.UpdatedAt = &now
}

// SetConfirmed marks the transaction as confirmed on the ledger.
func (p *XRPPendingMultisig) SetConfirmed() {
	p.Status = MultisigStatusConfirmed
	now := time.Now()
	p.UpdatedAt = &now
}

// SetFailed marks the transaction as failed.
func (p *XRPPendingMultisig) SetFailed() {
	p.Status = MultisigStatusFailed
	now := time.Now()
	p.UpdatedAt = &now
}

// SetExpired marks the transaction as expired.
func (p *XRPPendingMultisig) SetExpired() {
	p.Status = MultisigStatusExpired
	now := time.Now()
	p.UpdatedAt = &now
}

// IsExpired checks if the pending transaction has expired.
func (p *XRPPendingMultisig) IsExpired() bool {
	if p.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*p.ExpiresAt)
}

// CanAddSignature checks if the transaction can accept more signatures.
func (p *XRPPendingMultisig) CanAddSignature() bool {
	return p.Status == MultisigStatusPending && !p.IsExpired()
}

// XRPMultisigSignature represents a signature from one signer for a multi-sig transaction.
type XRPMultisigSignature struct {
	ID                int64
	PendingMultisigID int64      // Reference to XRPPendingMultisig.ID
	SignerAccount     string     // XRP address of the signer (r...)
	SignedTxBlob      string     // Signed transaction blob from this signer
	SignerWeight      uint32     // Weight of this signature
	SignedAt          *time.Time // When the signature was added
}

// NewXRPMultisigSignature creates a new XRPMultisigSignature entity.
func NewXRPMultisigSignature(
	pendingMultisigID int64,
	signerAccount string,
	signedTxBlob string,
	signerWeight uint32,
) (*XRPMultisigSignature, error) {
	if pendingMultisigID <= 0 {
		return nil, errors.New("pending multisig ID must be positive")
	}
	if signerAccount == "" {
		return nil, errors.New("signer account cannot be empty")
	}
	if signedTxBlob == "" {
		return nil, errors.New("signed transaction blob cannot be empty")
	}
	if signerWeight == 0 {
		return nil, errors.New("signer weight must be greater than 0")
	}

	now := time.Now()
	return &XRPMultisigSignature{
		PendingMultisigID: pendingMultisigID,
		SignerAccount:     signerAccount,
		SignedTxBlob:      signedTxBlob,
		SignerWeight:      signerWeight,
		SignedAt:          &now,
	}, nil
}
