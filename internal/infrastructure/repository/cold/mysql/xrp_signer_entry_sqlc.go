package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	domainXRP "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/xrp"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/mysql/sqlcgen"
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
func convertToXRPSignerEntry(sqlcEntry *sqlcgen.XrpSignerEntry) *domainXRP.XRPSignerEntry {
	entry := &domainXRP.XRPSignerEntry{
		ID:            sqlcEntry.ID,
		SignerListID:  sqlcEntry.SignerListID,
		SignerAccount: sqlcEntry.SignerAccount,
		SignerWeight:  sqlcEntry.SignerWeight,
	}

	if sqlcEntry.CreatedAt.Valid {
		entry.CreatedAt = &sqlcEntry.CreatedAt.Time
	}

	return entry
}

// GetByListID returns all signer entries for a signer list
func (r *XRPSignerEntryRepositorySqlc) GetByListID(
	ctx context.Context, signerListID int64,
) ([]*domainXRP.XRPSignerEntry, error) {
	sqlcEntries, err := r.queries.GetXRPSignerEntriesByListID(ctx, signerListID)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetXRPSignerEntriesByListID(): %w", err)
	}

	result := make([]*domainXRP.XRPSignerEntry, 0, len(sqlcEntries))
	for i := range sqlcEntries {
		result = append(result, convertToXRPSignerEntry(&sqlcEntries[i]))
	}

	return result, nil
}

// GetByID returns a signer entry by its ID
func (r *XRPSignerEntryRepositorySqlc) GetByID(
	ctx context.Context, id int64,
) (*domainXRP.XRPSignerEntry, error) {
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
) (*domainXRP.XRPSignerEntry, error) {
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
	case []byte:
		// MySQL might return as []byte
		var total int64
		for _, b := range v {
			total = total*10 + int64(b-'0')
		}
		return uint32(total), nil
	default:
		return 0, nil
	}
}

// Insert creates a new signer entry and returns the inserted ID
func (r *XRPSignerEntryRepositorySqlc) Insert(
	ctx context.Context, entry *domainXRP.XRPSignerEntry,
) (int64, error) {
	params := sqlcgen.InsertXRPSignerEntryParams{
		SignerListID:  entry.SignerListID,
		SignerAccount: entry.SignerAccount,
		SignerWeight:  entry.SignerWeight,
	}

	if entry.CreatedAt != nil {
		params.CreatedAt = sql.NullTime{Time: *entry.CreatedAt, Valid: true}
	} else {
		params.CreatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	}

	result, err := r.queries.InsertXRPSignerEntry(ctx, params)
	if err != nil {
		return 0, fmt.Errorf("failed to call InsertXRPSignerEntry(): %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get LastInsertId(): %w", err)
	}

	return id, nil
}

// InsertBulk creates multiple signer entries for a list
func (r *XRPSignerEntryRepositorySqlc) InsertBulk(
	ctx context.Context, entries []*domainXRP.XRPSignerEntry,
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
