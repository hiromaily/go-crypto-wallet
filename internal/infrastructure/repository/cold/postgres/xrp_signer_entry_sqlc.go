package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	domainXrp "github.com/hiromaily/go-crypto-wallet/internal/domain/xrp"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/postgres/sqlcgen"
)

// XRPSignerEntryRepositorySqlc is repository for xrp_signer_entry table using sqlc
type XRPSignerEntryRepositorySqlc struct {
	queries *sqlcgen.Queries
}

// NewXRPSignerEntryRepositorySqlc returns XRPSignerEntryRepositorySqlc object
func NewXRPSignerEntryRepositorySqlc(dbConn *sql.DB) *XRPSignerEntryRepositorySqlc {
	return &XRPSignerEntryRepositorySqlc{
		queries: sqlcgen.New(dbConn),
	}
}

// convertToXRPSignerEntry converts sqlcgen.XrpSignerEntry to domain.XRPSignerEntry entity.
func convertToXRPSignerEntry(sqlcEntry *sqlcgen.XrpSignerEntry) *domainXrp.XRPSignerEntry {
	entry := &domainXrp.XRPSignerEntry{
		ID:            sqlcEntry.ID,
		SignerListID:  sqlcEntry.SignerListID,
		SignerAccount: sqlcEntry.SignerAccount,
		SignerWeight:  uint32(sqlcEntry.SignerWeight),
	}

	if sqlcEntry.CreatedAt.Valid {
		entry.CreatedAt = &sqlcEntry.CreatedAt.Time
	}

	return entry
}

// GetByListID returns all signer entries for a signer list
func (r *XRPSignerEntryRepositorySqlc) GetByListID(
	ctx context.Context, signerListID int64,
) ([]*domainXrp.XRPSignerEntry, error) {
	sqlcEntries, err := r.queries.GetXRPSignerEntriesByListID(ctx, signerListID)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetXRPSignerEntriesByListID(): %w", err)
	}

	result := make([]*domainXrp.XRPSignerEntry, 0, len(sqlcEntries))
	for i := range sqlcEntries {
		result = append(result, convertToXRPSignerEntry(&sqlcEntries[i]))
	}

	return result, nil
}

// GetByID returns a signer entry by its ID
func (r *XRPSignerEntryRepositorySqlc) GetByID(
	ctx context.Context, id int64,
) (*domainXrp.XRPSignerEntry, error) {
	sqlcEntry, err := r.queries.GetXRPSignerEntryByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Signer entry not found
		}
		return nil, fmt.Errorf("failed to call GetXRPSignerEntryByID(): %w", err)
	}

	return convertToXRPSignerEntry(&sqlcEntry), nil
}

// GetByListAndAccount returns a signer entry by list ID and signer account
func (r *XRPSignerEntryRepositorySqlc) GetByListAndAccount(
	ctx context.Context, signerListID int64, signerAccount string,
) (*domainXrp.XRPSignerEntry, error) {
	params := sqlcgen.GetXRPSignerEntryByListAndAccountParams{
		SignerListID:  signerListID,
		SignerAccount: signerAccount,
	}

	sqlcEntry, err := r.queries.GetXRPSignerEntryByListAndAccount(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Signer entry not found
		}
		return nil, fmt.Errorf("failed to call GetXRPSignerEntryByListAndAccount(): %w", err)
	}

	return convertToXRPSignerEntry(&sqlcEntry), nil
}

// GetTotalWeight returns the total weight of all signers for a list
func (r *XRPSignerEntryRepositorySqlc) GetTotalWeight(
	ctx context.Context, signerListID int64,
) (uint32, error) {
	result, err := r.queries.GetTotalWeightByListID(ctx, signerListID)
	if err != nil {
		return 0, fmt.Errorf("failed to call GetTotalWeightByListID(): %w", err)
	}

	// The result is interface{} due to COALESCE, we need to handle different types
	switch v := result.(type) {
	case int64:
		return uint32(v), nil
	case float64:
		return uint32(v), nil
	default:
		return 0, nil
	}
}

// Insert creates a new signer entry and returns the inserted ID
func (r *XRPSignerEntryRepositorySqlc) Insert(
	ctx context.Context, entry *domainXrp.XRPSignerEntry,
) (int64, error) {
	params := sqlcgen.InsertXRPSignerEntryParams{
		SignerListID:  entry.SignerListID,
		SignerAccount: entry.SignerAccount,
		SignerWeight:  int32(entry.SignerWeight),
	}

	if entry.CreatedAt != nil {
		params.CreatedAt = sql.NullTime{Time: *entry.CreatedAt, Valid: true}
	} else {
		params.CreatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	}

	// Postgres RETURNING id returns the id directly
	id, err := r.queries.InsertXRPSignerEntry(ctx, params)
	if err != nil {
		return 0, fmt.Errorf("failed to call InsertXRPSignerEntry(): %w", err)
	}

	return id, nil
}

// InsertBulk creates multiple signer entries for a list
func (r *XRPSignerEntryRepositorySqlc) InsertBulk(
	ctx context.Context, entries []*domainXrp.XRPSignerEntry,
) error {
	for _, entry := range entries {
		_, err := r.Insert(ctx, entry)
		if err != nil {
			return fmt.Errorf("failed to insert signer entry for account %s: %w", entry.SignerAccount, err)
		}
	}

	return nil
}

// DeleteByListID deletes all signer entries for a signer list
func (r *XRPSignerEntryRepositorySqlc) DeleteByListID(
	ctx context.Context, signerListID int64,
) error {
	err := r.queries.DeleteXRPSignerEntriesByListID(ctx, signerListID)
	if err != nil {
		return fmt.Errorf("failed to call DeleteXRPSignerEntriesByListID(): %w", err)
	}

	return nil
}
