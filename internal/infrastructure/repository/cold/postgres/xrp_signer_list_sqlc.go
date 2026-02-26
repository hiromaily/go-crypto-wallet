package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	domainXRP "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/xrp"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/postgres/sqlcgen"
)

// XRPSignerListRepositorySqlc is repository for xrp_signer_list table using sqlc
type XRPSignerListRepositorySqlc struct {
	queries *sqlcgen.Queries
}

// NewXRPSignerListRepositorySqlc returns XRPSignerListRepositorySqlc object
func NewXRPSignerListRepositorySqlc(dbConn *sql.DB) *XRPSignerListRepositorySqlc {
	return &XRPSignerListRepositorySqlc{
		queries: sqlcgen.New(dbConn),
	}
}

// convertToXRPSignerList converts sqlcgen.XrpSignerList to domain.XRPSignerList entity.
func convertToXRPSignerList(sqlcList *sqlcgen.XrpSignerList) *domainXRP.XRPSignerList {
	list := &domainXRP.XRPSignerList{
		ID:           sqlcList.ID,
		AccountID:    sqlcList.AccountID,
		SignerQuorum: uint32(sqlcList.SignerQuorum),
		IsActive:     sqlcList.IsActive,
	}

	if sqlcList.SetTxHash.Valid {
		list.SetTxHash = &sqlcList.SetTxHash.String
	}
	if sqlcList.CreatedAt.Valid {
		list.CreatedAt = &sqlcList.CreatedAt.Time
	}
	if sqlcList.UpdatedAt.Valid {
		list.UpdatedAt = &sqlcList.UpdatedAt.Time
	}

	return list
}

// GetByAccountID returns the active signer list for an account
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

	return convertToXRPSignerList(&sqlcList), nil
}

// GetByID returns a signer list by its ID
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

	return convertToXRPSignerList(&sqlcList), nil
}

// GetHistoryByAccountID returns all signer lists (active and inactive) for an account
func (r *XRPSignerListRepositorySqlc) GetHistoryByAccountID(
	ctx context.Context, accountID string,
) ([]*domainXRP.XRPSignerList, error) {
	sqlcLists, err := r.queries.GetXRPSignerListHistoryByAccountID(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetXRPSignerListHistoryByAccountID(): %w", err)
	}

	result := make([]*domainXRP.XRPSignerList, 0, len(sqlcLists))
	for i := range sqlcLists {
		result = append(result, convertToXRPSignerList(&sqlcLists[i]))
	}

	return result, nil
}

// Insert creates a new signer list record and returns the inserted ID
func (r *XRPSignerListRepositorySqlc) Insert(
	ctx context.Context, signerList *domainXRP.XRPSignerList,
) (int64, error) {
	params := sqlcgen.InsertXRPSignerListParams{
		AccountID:    signerList.AccountID,
		SignerQuorum: int32(signerList.SignerQuorum),
		IsActive:     signerList.IsActive,
	}

	if signerList.SetTxHash != nil {
		params.SetTxHash = sql.NullString{String: *signerList.SetTxHash, Valid: true}
	}
	if signerList.CreatedAt != nil {
		params.CreatedAt = sql.NullTime{Time: *signerList.CreatedAt, Valid: true}
	} else {
		params.CreatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	}

	// Postgres RETURNING id returns the id directly
	id, err := r.queries.InsertXRPSignerList(ctx, params)
	if err != nil {
		return 0, fmt.Errorf("failed to call InsertXRPSignerList(): %w", err)
	}

	return id, nil
}

// UpdateStatus updates the active status of a signer list
func (r *XRPSignerListRepositorySqlc) UpdateStatus(
	ctx context.Context, id int64, isActive bool,
) error {
	params := sqlcgen.UpdateXRPSignerListStatusParams{
		IsActive:  isActive,
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		ID:        id,
	}

	_, err := r.queries.UpdateXRPSignerListStatus(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to call UpdateXRPSignerListStatus(): %w", err)
	}

	return nil
}

// UpdateTxHash updates the SignerListSet transaction hash
func (r *XRPSignerListRepositorySqlc) UpdateTxHash(
	ctx context.Context, id int64, txHash string,
) error {
	params := sqlcgen.UpdateXRPSignerListTxHashParams{
		SetTxHash: sql.NullString{String: txHash, Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		ID:        id,
	}

	_, err := r.queries.UpdateXRPSignerListTxHash(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to call UpdateXRPSignerListTxHash(): %w", err)
	}

	return nil
}

// DeactivateByAccountID deactivates all signer lists for an account
func (r *XRPSignerListRepositorySqlc) DeactivateByAccountID(
	ctx context.Context, accountID string,
) error {
	params := sqlcgen.DeactivateXRPSignerListByAccountIDParams{
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		AccountID: accountID,
	}

	_, err := r.queries.DeactivateXRPSignerListByAccountID(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to call DeactivateXRPSignerListByAccountID(): %w", err)
	}

	return nil
}
