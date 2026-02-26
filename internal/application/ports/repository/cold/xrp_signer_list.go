package cold

import (
	"context"

	domainXRP "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/xrp"
)

// XRPSignerListRepositorier is the repository interface for XRP multi-signature signer list management.
// A signer list enables multi-signature transactions where multiple parties must sign
// to authorize a transaction.
// Reference: https://xrpl.org/docs/concepts/accounts/multi-signing
type XRPSignerListRepositorier interface {
	// GetByAccountID returns the active signer list for an account
	GetByAccountID(ctx context.Context, accountID string) (*domainXRP.XRPSignerList, error)
	// GetByID returns a signer list by its ID
	GetByID(ctx context.Context, id int64) (*domainXRP.XRPSignerList, error)
	// GetHistoryByAccountID returns all signer lists (active and inactive) for an account
	GetHistoryByAccountID(ctx context.Context, accountID string) ([]*domainXRP.XRPSignerList, error)
	// Insert creates a new signer list record and returns the inserted ID
	Insert(ctx context.Context, signerList *domainXRP.XRPSignerList) (int64, error)
	// UpdateStatus updates the active status of a signer list
	UpdateStatus(ctx context.Context, id int64, isActive bool) error
	// UpdateTxHash updates the SignerListSet transaction hash
	UpdateTxHash(ctx context.Context, id int64, txHash string) error
	// DeactivateByAccountID deactivates all signer lists for an account
	DeactivateByAccountID(ctx context.Context, accountID string) error
}

// XRPSignerEntryRepositorier is the repository interface for XRP signer entries.
// Each entry represents an authorized signer with a weight that contributes to the quorum.
type XRPSignerEntryRepositorier interface {
	// GetByListID returns all signer entries for a signer list
	GetByListID(ctx context.Context, signerListID int64) ([]*domainXRP.XRPSignerEntry, error)
	// GetByID returns a signer entry by its ID
	GetByID(ctx context.Context, id int64) (*domainXRP.XRPSignerEntry, error)
	// GetByListAndAccount returns a signer entry by list ID and signer account
	GetByListAndAccount(
		ctx context.Context, signerListID int64, signerAccount string,
	) (*domainXRP.XRPSignerEntry, error)
	// GetTotalWeight returns the total weight of all signers for a list
	GetTotalWeight(ctx context.Context, signerListID int64) (uint32, error)
	// Insert creates a new signer entry and returns the inserted ID
	Insert(ctx context.Context, entry *domainXRP.XRPSignerEntry) (int64, error)
	// InsertBulk creates multiple signer entries for a list
	InsertBulk(ctx context.Context, entries []*domainXRP.XRPSignerEntry) error
	// DeleteByListID deletes all signer entries for a signer list
	DeleteByListID(ctx context.Context, signerListID int64) error
}
