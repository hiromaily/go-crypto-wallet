package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	domainXRP "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/xrp"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/sqlite/sqlcgen"
)

// XRPPendingMultisigRepositorySqlc is repository for xrp_pending_multisig table using sqlc (SQLite).
type XRPPendingMultisigRepositorySqlc struct {
	queries *sqlcgen.Queries
}

// NewXRPPendingMultisigRepositorySqlc returns a new XRPPendingMultisigRepositorySqlc.
func NewXRPPendingMultisigRepositorySqlc(dbConn *sql.DB) *XRPPendingMultisigRepositorySqlc {
	return &XRPPendingMultisigRepositorySqlc{
		queries: sqlcgen.New(dbConn),
	}
}

// convertToXRPPendingMultisigFromSQLite converts a sqlcgen.XrpPendingMultisig (SQLite row) to a domain entity.
// SQLite stores timestamps as TEXT and integers instead of native time/enum types.
func convertToXRPPendingMultisigFromSQLite(
	sqlcPending *sqlcgen.XrpPendingMultisig,
) *domainXRP.XRPPendingMultisig {
	pending := &domainXRP.XRPPendingMultisig{
		ID:             sqlcPending.ID,
		TxUUID:         sqlcPending.TxUuid,
		AccountID:      sqlcPending.AccountID,
		UnsignedTxJSON: sqlcPending.UnsignedTxJson,
		XRPTxType:      sqlcPending.XrpTxType,
		RequiredQuorum: uint32(sqlcPending.RequiredQuorum),
		CurrentWeight:  uint32(sqlcPending.CurrentWeight),
		Status:         domainXRP.MultisigStatus(sqlcPending.Status),
	}

	if sqlcPending.CombinedTxBlob.Valid {
		pending.CombinedTxBlob = &sqlcPending.CombinedTxBlob.String
	}
	if sqlcPending.SubmittedTxHash.Valid {
		pending.SubmittedTxHash = &sqlcPending.SubmittedTxHash.String
	}
	if sqlcPending.ExpiresAt.Valid {
		t, err := time.Parse(sqliteTimestampFormat, sqlcPending.ExpiresAt.String)
		if err == nil {
			pending.ExpiresAt = &t
		}
	}
	if sqlcPending.CreatedAt.Valid {
		t, err := time.Parse(sqliteTimestampFormat, sqlcPending.CreatedAt.String)
		if err == nil {
			pending.CreatedAt = &t
		}
	}
	if sqlcPending.UpdatedAt.Valid {
		t, err := time.Parse(sqliteTimestampFormat, sqlcPending.UpdatedAt.String)
		if err == nil {
			pending.UpdatedAt = &t
		}
	}

	return pending
}

// GetByID returns a pending multi-sig transaction by its ID.
func (r *XRPPendingMultisigRepositorySqlc) GetByID(
	ctx context.Context, id int64,
) (*domainXRP.XRPPendingMultisig, error) {
	sqlcPending, err := r.queries.GetXRPPendingMultisigByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, fmt.Errorf("failed to call GetXRPPendingMultisigByID(): %w", err)
	}

	return convertToXRPPendingMultisigFromSQLite(&sqlcPending), nil
}

// GetByUUID returns a pending multi-sig transaction by its UUID.
func (r *XRPPendingMultisigRepositorySqlc) GetByUUID(
	ctx context.Context, txUUID string,
) (*domainXRP.XRPPendingMultisig, error) {
	sqlcPending, err := r.queries.GetXRPPendingMultisigByUUID(ctx, txUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, fmt.Errorf("failed to call GetXRPPendingMultisigByUUID(): %w", err)
	}

	return convertToXRPPendingMultisigFromSQLite(&sqlcPending), nil
}

// GetByAccountIDAndStatus returns pending multi-sig transactions for an account with a specific status.
func (r *XRPPendingMultisigRepositorySqlc) GetByAccountIDAndStatus(
	ctx context.Context, accountID string, status domainXRP.MultisigStatus,
) ([]*domainXRP.XRPPendingMultisig, error) {
	params := sqlcgen.GetXRPPendingMultisigsByAccountIDParams{
		AccountID: accountID,
		Status:    string(status),
	}

	sqlcPendings, err := r.queries.GetXRPPendingMultisigsByAccountID(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetXRPPendingMultisigsByAccountID(): %w", err)
	}

	result := make([]*domainXRP.XRPPendingMultisig, 0, len(sqlcPendings))
	for i := range sqlcPendings {
		result = append(result, convertToXRPPendingMultisigFromSQLite(&sqlcPendings[i]))
	}

	return result, nil
}

// GetByStatus returns all pending multi-sig transactions with a specific status.
func (r *XRPPendingMultisigRepositorySqlc) GetByStatus(
	ctx context.Context, status domainXRP.MultisigStatus,
) ([]*domainXRP.XRPPendingMultisig, error) {
	sqlcPendings, err := r.queries.GetXRPPendingMultisigsByStatus(ctx, string(status))
	if err != nil {
		return nil, fmt.Errorf("failed to call GetXRPPendingMultisigsByStatus(): %w", err)
	}

	result := make([]*domainXRP.XRPPendingMultisig, 0, len(sqlcPendings))
	for i := range sqlcPendings {
		result = append(result, convertToXRPPendingMultisigFromSQLite(&sqlcPendings[i]))
	}

	return result, nil
}

// GetExpired returns pending transactions that have expired.
func (r *XRPPendingMultisigRepositorySqlc) GetExpired(
	ctx context.Context,
) ([]*domainXRP.XRPPendingMultisig, error) {
	nowStr := sql.NullString{String: time.Now().Format(sqliteTimestampFormat), Valid: true}

	sqlcPendings, err := r.queries.GetExpiredXRPPendingMultisigs(ctx, nowStr)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetExpiredXRPPendingMultisigs(): %w", err)
	}

	result := make([]*domainXRP.XRPPendingMultisig, 0, len(sqlcPendings))
	for i := range sqlcPendings {
		result = append(result, convertToXRPPendingMultisigFromSQLite(&sqlcPendings[i]))
	}

	return result, nil
}

// Insert creates a new pending multi-sig transaction and returns the inserted ID.
func (r *XRPPendingMultisigRepositorySqlc) Insert(
	ctx context.Context, pending *domainXRP.XRPPendingMultisig,
) (int64, error) {
	createdAt := time.Now()
	if pending.CreatedAt != nil {
		createdAt = *pending.CreatedAt
	}

	params := sqlcgen.InsertXRPPendingMultisigParams{
		TxUuid:         pending.TxUUID,
		AccountID:      pending.AccountID,
		UnsignedTxJson: pending.UnsignedTxJSON,
		XrpTxType:      pending.XRPTxType,
		RequiredQuorum: int64(pending.RequiredQuorum),
		CurrentWeight:  int64(pending.CurrentWeight),
		Status:         string(pending.Status),
		CreatedAt:      sql.NullString{String: createdAt.Format(sqliteTimestampFormat), Valid: true},
	}

	if pending.ExpiresAt != nil {
		params.ExpiresAt = sql.NullString{
			String: pending.ExpiresAt.Format(sqliteTimestampFormat),
			Valid:  true,
		}
	}

	result, err := r.queries.InsertXRPPendingMultisig(ctx, params)
	if err != nil {
		return 0, fmt.Errorf("failed to call InsertXRPPendingMultisig(): %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get LastInsertId(): %w", err)
	}

	return id, nil
}

// UpdateWeight updates the current weight of collected signatures.
func (r *XRPPendingMultisigRepositorySqlc) UpdateWeight(
	ctx context.Context, id int64, currentWeight uint32,
) error {
	params := sqlcgen.UpdateXRPPendingMultisigWeightParams{
		CurrentWeight: int64(currentWeight),
		UpdatedAt:     sql.NullString{String: time.Now().Format(sqliteTimestampFormat), Valid: true},
		ID:            id,
	}

	_, err := r.queries.UpdateXRPPendingMultisigWeight(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to call UpdateXRPPendingMultisigWeight(): %w", err)
	}

	return nil
}

// UpdateStatus updates the status of a pending transaction.
func (r *XRPPendingMultisigRepositorySqlc) UpdateStatus(
	ctx context.Context, id int64, status domainXRP.MultisigStatus,
) error {
	params := sqlcgen.UpdateXRPPendingMultisigStatusParams{
		Status:    string(status),
		UpdatedAt: sql.NullString{String: time.Now().Format(sqliteTimestampFormat), Valid: true},
		ID:        id,
	}

	_, err := r.queries.UpdateXRPPendingMultisigStatus(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to call UpdateXRPPendingMultisigStatus(): %w", err)
	}

	return nil
}

// SetReady marks the transaction as ready with the combined tx blob.
func (r *XRPPendingMultisigRepositorySqlc) SetReady(
	ctx context.Context, id int64, combinedTxBlob string,
) error {
	params := sqlcgen.UpdateXRPPendingMultisigReadyParams{
		CombinedTxBlob: sql.NullString{String: combinedTxBlob, Valid: true},
		UpdatedAt:      sql.NullString{String: time.Now().Format(sqliteTimestampFormat), Valid: true},
		ID:             id,
	}

	_, err := r.queries.UpdateXRPPendingMultisigReady(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to call UpdateXRPPendingMultisigReady(): %w", err)
	}

	return nil
}

// SetSubmitted marks the transaction as submitted with the tx hash.
func (r *XRPPendingMultisigRepositorySqlc) SetSubmitted(
	ctx context.Context, id int64, submittedTxHash string,
) error {
	params := sqlcgen.UpdateXRPPendingMultisigSubmittedParams{
		SubmittedTxHash: sql.NullString{String: submittedTxHash, Valid: true},
		UpdatedAt:       sql.NullString{String: time.Now().Format(sqliteTimestampFormat), Valid: true},
		ID:              id,
	}

	_, err := r.queries.UpdateXRPPendingMultisigSubmitted(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to call UpdateXRPPendingMultisigSubmitted(): %w", err)
	}

	return nil
}

// SetConfirmed marks the transaction as confirmed on the ledger.
func (r *XRPPendingMultisigRepositorySqlc) SetConfirmed(
	ctx context.Context, id int64,
) error {
	params := sqlcgen.UpdateXRPPendingMultisigConfirmedParams{
		UpdatedAt: sql.NullString{String: time.Now().Format(sqliteTimestampFormat), Valid: true},
		ID:        id,
	}

	_, err := r.queries.UpdateXRPPendingMultisigConfirmed(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to call UpdateXRPPendingMultisigConfirmed(): %w", err)
	}

	return nil
}

// SetFailed marks the transaction as failed.
func (r *XRPPendingMultisigRepositorySqlc) SetFailed(
	ctx context.Context, id int64,
) error {
	params := sqlcgen.UpdateXRPPendingMultisigFailedParams{
		UpdatedAt: sql.NullString{String: time.Now().Format(sqliteTimestampFormat), Valid: true},
		ID:        id,
	}

	_, err := r.queries.UpdateXRPPendingMultisigFailed(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to call UpdateXRPPendingMultisigFailed(): %w", err)
	}

	return nil
}

// ExpireAll marks all expired pending transactions as expired.
func (r *XRPPendingMultisigRepositorySqlc) ExpireAll(
	ctx context.Context,
) (int64, error) {
	nowStr := time.Now().Format(sqliteTimestampFormat)
	params := sqlcgen.ExpireXRPPendingMultisigsParams{
		UpdatedAt: sql.NullString{String: nowStr, Valid: true},
		ExpiresAt: sql.NullString{String: nowStr, Valid: true},
	}

	result, err := r.queries.ExpireXRPPendingMultisigs(ctx, params)
	if err != nil {
		return 0, fmt.Errorf("failed to call ExpireXRPPendingMultisigs(): %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get RowsAffected(): %w", err)
	}

	return rowsAffected, nil
}
