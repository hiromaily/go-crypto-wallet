package watch

import (
	"context"

	domainXRP "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/xrp"
)

// XRPPendingMultisigRepositorier is the repository interface for XRP pending multi-signature transactions.
// Pending multi-sig transactions wait for enough signatures to meet the quorum before submission.
type XRPPendingMultisigRepositorier interface {
	// GetByID returns a pending multi-sig transaction by its ID
	GetByID(ctx context.Context, id int64) (*domainXRP.XRPPendingMultisig, error)
	// GetByUUID returns a pending multi-sig transaction by its UUID
	GetByUUID(ctx context.Context, txUUID string) (*domainXRP.XRPPendingMultisig, error)
	// GetByAccountIDAndStatus returns pending multi-sig transactions for an account with a specific status
	GetByAccountIDAndStatus(
		ctx context.Context, accountID string, status domainXRP.MultisigStatus,
	) ([]*domainXRP.XRPPendingMultisig, error)
	// GetByStatus returns all pending multi-sig transactions with a specific status
	GetByStatus(ctx context.Context, status domainXRP.MultisigStatus) ([]*domainXRP.XRPPendingMultisig, error)
	// GetExpired returns pending transactions that have expired
	GetExpired(ctx context.Context) ([]*domainXRP.XRPPendingMultisig, error)
	// Insert creates a new pending multi-sig transaction and returns the inserted ID
	Insert(ctx context.Context, pending *domainXRP.XRPPendingMultisig) (int64, error)
	// UpdateWeight updates the current weight of collected signatures
	UpdateWeight(ctx context.Context, id int64, currentWeight uint32) error
	// UpdateStatus updates the status of a pending transaction
	UpdateStatus(ctx context.Context, id int64, status domainXRP.MultisigStatus) error
	// SetReady marks the transaction as ready with the combined tx blob
	SetReady(ctx context.Context, id int64, combinedTxBlob string) error
	// SetSubmitted marks the transaction as submitted with the tx hash
	SetSubmitted(ctx context.Context, id int64, submittedTxHash string) error
	// SetConfirmed marks the transaction as confirmed on the ledger
	SetConfirmed(ctx context.Context, id int64) error
	// SetFailed marks the transaction as failed
	SetFailed(ctx context.Context, id int64) error
	// ExpireAll marks all expired pending transactions as expired
	ExpireAll(ctx context.Context) (int64, error)
}

// XRPMultisigSignatureRepositorier is the repository interface for multi-sig signatures.
// Each signature represents one signer's contribution to a multi-sig transaction.
type XRPMultisigSignatureRepositorier interface {
	// GetByID returns a signature by its ID
	GetByID(ctx context.Context, id int64) (*domainXRP.XRPMultisigSignature, error)
	// GetByPendingID returns all signatures for a pending multi-sig transaction
	GetByPendingID(ctx context.Context, pendingMultisigID int64) ([]*domainXRP.XRPMultisigSignature, error)
	// GetByPendingAndSigner returns a signature by pending ID and signer account
	GetByPendingAndSigner(
		ctx context.Context, pendingMultisigID int64, signerAccount string,
	) (*domainXRP.XRPMultisigSignature, error)
	// GetSignatureCount returns the number of signatures for a pending transaction
	GetSignatureCount(ctx context.Context, pendingMultisigID int64) (int64, error)
	// GetTotalWeight returns the total weight of signatures for a pending transaction
	GetTotalWeight(ctx context.Context, pendingMultisigID int64) (uint32, error)
	// Insert creates a new signature and returns the inserted ID
	Insert(ctx context.Context, signature *domainXRP.XRPMultisigSignature) (int64, error)
	// DeleteByPendingID deletes all signatures for a pending multi-sig transaction
	DeleteByPendingID(ctx context.Context, pendingMultisigID int64) error
}
