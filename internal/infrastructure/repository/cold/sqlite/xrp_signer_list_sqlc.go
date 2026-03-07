package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	domainXRP "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/xrp"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/sqlite/sqlcgen"
)

const sqliteTimestampFormat = "2006-01-02 15:04:05"

// XRPSignerListRepositorySqlc is repository for xrp_signer_list table using sqlc (SQLite driver).
//
// DeactivateByAccountID followed by Insert is non-atomic — this is an accepted CLI limitation.
// A crash between the two operations leaves the account with no active signer list,
// which can be repaired by re-running the set-signer-list command.
type XRPSignerListRepositorySqlc struct {
	queries *sqlcgen.Queries
}

// NewXRPSignerListRepositorySqlc returns a new XRPSignerListRepositorySqlc.
func NewXRPSignerListRepositorySqlc(dbConn *sql.DB) *XRPSignerListRepositorySqlc {
	return &XRPSignerListRepositorySqlc{
		queries: sqlcgen.New(dbConn),
	}
}

// convertToXRPSignerList converts a sqlcgen.XrpSignerList (SQLite row) to a domain entity.
// SQLite stores timestamps as TEXT and booleans as INTEGER (0/1).
func convertToXRPSignerListFromSQLite(sqlcList *sqlcgen.XrpSignerList) *domainXRP.XRPSignerList {
	list := &domainXRP.XRPSignerList{
		ID:           sqlcList.ID,
		AccountID:    sqlcList.AccountID,
		SignerQuorum: uint32(sqlcList.SignerQuorum),
		IsActive:     sqlcList.IsActive != 0,
	}

	if sqlcList.SetTxHash.Valid {
		list.SetTxHash = &sqlcList.SetTxHash.String
	}
	if sqlcList.CreatedAt.Valid {
		t, err := time.Parse(sqliteTimestampFormat, sqlcList.CreatedAt.String)
		if err == nil {
			list.CreatedAt = &t
		}
	}
	if sqlcList.UpdatedAt.Valid {
		t, err := time.Parse(sqliteTimestampFormat, sqlcList.UpdatedAt.String)
		if err == nil {
			list.UpdatedAt = &t
		}
	}

	return list
}

// GetByAccountID returns the active signer list for an account.
func (r *XRPSignerListRepositorySqlc) GetByAccountID(
	ctx context.Context, accountID string,
) (*domainXRP.XRPSignerList, error) {
	sqlcList, err := r.queries.GetXRPSignerListByAccountID(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No active signer list found
		}
		return nil, fmt.Errorf("failed to call GetXRPSignerListByAccountID(): %w", err)
	}

	return convertToXRPSignerListFromSQLite(&sqlcList), nil
}

// GetByID returns a signer list by its ID.
func (r *XRPSignerListRepositorySqlc) GetByID(
	ctx context.Context, id int64,
) (*domainXRP.XRPSignerList, error) {
	sqlcList, err := r.queries.GetXRPSignerListByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Signer list not found
		}
		return nil, fmt.Errorf("failed to call GetXRPSignerListByID(): %w", err)
	}

	return convertToXRPSignerListFromSQLite(&sqlcList), nil
}

// GetHistoryByAccountID returns all signer lists (active and inactive) for an account.
func (r *XRPSignerListRepositorySqlc) GetHistoryByAccountID(
	ctx context.Context, accountID string,
) ([]*domainXRP.XRPSignerList, error) {
	sqlcLists, err := r.queries.GetXRPSignerListHistoryByAccountID(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetXRPSignerListHistoryByAccountID(): %w", err)
	}

	result := make([]*domainXRP.XRPSignerList, 0, len(sqlcLists))
	for i := range sqlcLists {
		result = append(result, convertToXRPSignerListFromSQLite(&sqlcLists[i]))
	}

	return result, nil
}

// Insert creates a new signer list record and returns the inserted ID.
func (r *XRPSignerListRepositorySqlc) Insert(
	ctx context.Context, signerList *domainXRP.XRPSignerList,
) (int64, error) {
	createdAt := time.Now()
	if signerList.CreatedAt != nil {
		createdAt = *signerList.CreatedAt
	}

	params := sqlcgen.InsertXRPSignerListParams{
		AccountID:    signerList.AccountID,
		SignerQuorum: int64(signerList.SignerQuorum),
		IsActive:     boolToInt64(signerList.IsActive),
		CreatedAt:    sql.NullString{String: createdAt.Format(sqliteTimestampFormat), Valid: true},
	}

	if signerList.SetTxHash != nil {
		params.SetTxHash = sql.NullString{String: *signerList.SetTxHash, Valid: true}
	}

	result, err := r.queries.InsertXRPSignerList(ctx, params)
	if err != nil {
		return 0, fmt.Errorf("failed to call InsertXRPSignerList(): %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get LastInsertId(): %w", err)
	}

	return id, nil
}

// UpdateStatus updates the active status of a signer list.
func (r *XRPSignerListRepositorySqlc) UpdateStatus(
	ctx context.Context, id int64, isActive bool,
) error {
	params := sqlcgen.UpdateXRPSignerListStatusParams{
		IsActive:  boolToInt64(isActive),
		UpdatedAt: sql.NullString{String: time.Now().Format(sqliteTimestampFormat), Valid: true},
		ID:        id,
	}

	_, err := r.queries.UpdateXRPSignerListStatus(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to call UpdateXRPSignerListStatus(): %w", err)
	}

	return nil
}

// UpdateTxHash updates the SignerListSet transaction hash.
func (r *XRPSignerListRepositorySqlc) UpdateTxHash(
	ctx context.Context, id int64, txHash string,
) error {
	params := sqlcgen.UpdateXRPSignerListTxHashParams{
		SetTxHash: sql.NullString{String: txHash, Valid: true},
		UpdatedAt: sql.NullString{String: time.Now().Format(sqliteTimestampFormat), Valid: true},
		ID:        id,
	}

	_, err := r.queries.UpdateXRPSignerListTxHash(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to call UpdateXRPSignerListTxHash(): %w", err)
	}

	return nil
}

// DeactivateByAccountID deactivates all active signer lists for an account.
//
// Note: This operation is not atomic with the subsequent Insert call in SetSignerListUseCase.
// A crash between DeactivateByAccountID and Insert leaves the account with no active signer list.
// This is an accepted CLI limitation — re-running set-signer-list will restore the list.
func (r *XRPSignerListRepositorySqlc) DeactivateByAccountID(
	ctx context.Context, accountID string,
) error {
	params := sqlcgen.DeactivateXRPSignerListByAccountIDParams{
		UpdatedAt: sql.NullString{String: time.Now().Format(sqliteTimestampFormat), Valid: true},
		AccountID: accountID,
	}

	_, err := r.queries.DeactivateXRPSignerListByAccountID(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to call DeactivateXRPSignerListByAccountID(): %w", err)
	}

	return nil
}
