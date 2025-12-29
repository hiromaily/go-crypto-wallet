package multisig

import (
	"context"
	"time"
)

// NonceRepository defines the interface for persistent nonce storage.
// This interface follows the Repository pattern and is part of the domain layer.
//
// SECURITY: Proper nonce management is CRITICAL for MuSig2 security.
// Nonce reuse will leak the private key!
type NonceRepository interface {
	// SaveNonce stores a nonce commitment for a specific signer and transaction.
	// Returns an error if a nonce already exists for this signer-transaction pair (enforces uniqueness).
	SaveNonce(ctx context.Context, signerID, transactionID string, publicNonce [66]byte) error

	// GetNonce retrieves a nonce commitment for a specific signer and transaction.
	// Returns ErrNonceNotFound if the nonce does not exist.
	GetNonce(ctx context.Context, signerID, transactionID string) (*NonceCommitment, error)

	// GetAllNoncesForTransaction retrieves all nonce commitments for a transaction.
	// Used in Round 1 to check if all participants have submitted their nonces.
	GetAllNoncesForTransaction(ctx context.Context, transactionID string) ([]*NonceCommitment, error)

	// GetUnusedNoncesForTransaction retrieves all unused nonce commitments for a transaction.
	// Returns only nonces that have not been marked as used.
	GetUnusedNoncesForTransaction(ctx context.Context, transactionID string) ([]*NonceCommitment, error)

	// MarkNonceUsed marks a nonce as used after successful signing.
	// This prevents accidental reuse of the nonce.
	MarkNonceUsed(ctx context.Context, signerID, transactionID string) error

	// DeleteNoncesForTransaction removes all nonces associated with a transaction.
	// Used for cleanup after transaction completion or cancellation.
	DeleteNoncesForTransaction(ctx context.Context, transactionID string) error

	// CleanupOldUnusedNonces removes unused nonces older than the specified time.
	// Used for periodic cleanup of stale nonces (e.g., from failed transactions).
	CleanupOldUnusedNonces(ctx context.Context, olderThan time.Time) error
}

// Sentinel errors for nonce repository operations.
var (
	// ErrNonceNotFound indicates that a requested nonce does not exist.
	ErrNonceNotFound = NewError("nonce not found")

	// ErrNonceDuplicate indicates that a nonce already exists for the signer-transaction pair.
	ErrNonceDuplicate = NewError("nonce already exists for this signer and transaction")
)

// NewError creates a domain error.
func NewError(msg string) error {
	return &domainError{msg: msg}
}

type domainError struct {
	msg string
}

func (e *domainError) Error() string {
	return e.msg
}
