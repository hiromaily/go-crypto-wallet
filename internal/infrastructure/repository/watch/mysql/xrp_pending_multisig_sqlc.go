package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	domainXRP "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/xrp"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/mysql/sqlcgen"
)

// XRPPendingMultisigRepositorySqlc is repository for xrp_pending_multisig table using sqlc
type XRPPendingMultisigRepositorySqlc struct {
	queries *sqlcgen.Queries
}

// NewXRPPendingMultisigRepositorySqlc returns XRPPendingMultisigRepositorySqlc object
func NewXRPPendingMultisigRepositorySqlc(dbConn *sql.DB) *XRPPendingMultisigRepositorySqlc {
	return &XRPPendingMultisigRepositorySqlc{
		queries: sqlcgen.New(dbConn),
	}
}

// convertToXRPPendingMultisig converts sqlcgen.XrpPendingMultisig to domain.XRPPendingMultisig entity.
func convertToXRPPendingMultisig(sqlcPending *sqlcgen.XrpPendingMultisig) *domainXRP.XRPPendingMultisig {
	pending := &domainXRP.XRPPendingMultisig{
		ID:             sqlcPending.ID,
		TxUUID:         sqlcPending.TxUuid,
		AccountID:      sqlcPending.AccountID,
		UnsignedTxJSON: sqlcPending.UnsignedTxJson,
		XRPTxType:      sqlcPending.XrpTxType,
		RequiredQuorum: sqlcPending.RequiredQuorum,
		CurrentWeight:  sqlcPending.CurrentWeight,
		Status:         domainXRP.MultisigStatus(sqlcPending.Status),
	}

	if sqlcPending.CombinedTxBlob.Valid {
		pending.CombinedTxBlob = &sqlcPending.CombinedTxBlob.String
	}
	if sqlcPending.SubmittedTxHash.Valid {
		pending.SubmittedTxHash = &sqlcPending.SubmittedTxHash.String
	}
	if sqlcPending.ExpiresAt.Valid {
		pending.ExpiresAt = &sqlcPending.ExpiresAt.Time
	}
	if sqlcPending.CreatedAt.Valid {
		pending.CreatedAt = &sqlcPending.CreatedAt.Time
	}
	if sqlcPending.UpdatedAt.Valid {
		pending.UpdatedAt = &sqlcPending.UpdatedAt.Time
	}

	return pending
}

// convertDomainStatusToSqlc converts domain MultisigStatus to sqlcgen XrpPendingMultisigStatus
func convertDomainStatusToSqlc(status domainXRP.MultisigStatus) sqlcgen.XrpPendingMultisigStatus {
	switch status {
	case domainXRP.MultisigStatusPending:
		return sqlcgen.XrpPendingMultisigStatusPending
	case domainXRP.MultisigStatusReady:
		return sqlcgen.XrpPendingMultisigStatusReady
	case domainXRP.MultisigStatusSubmitted:
		return sqlcgen.XrpPendingMultisigStatusSubmitted
	case domainXRP.MultisigStatusConfirmed:
		return sqlcgen.XrpPendingMultisigStatusConfirmed
	case domainXRP.MultisigStatusFailed:
		return sqlcgen.XrpPendingMultisigStatusFailed
	case domainXRP.MultisigStatusExpired:
		return sqlcgen.XrpPendingMultisigStatusExpired
	default:
		return sqlcgen.XrpPendingMultisigStatusPending
	}
}

// GetByID returns a pending multi-sig transaction by its ID
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

	return convertToXRPPendingMultisig(&sqlcPending), nil
}

// GetByUUID returns a pending multi-sig transaction by its UUID
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

	return convertToXRPPendingMultisig(&sqlcPending), nil
}

// GetByAccountIDAndStatus returns pending multi-sig transactions for an account with a specific status
func (r *XRPPendingMultisigRepositorySqlc) GetByAccountIDAndStatus(
	ctx context.Context, accountID string, status domainXRP.MultisigStatus,
) ([]*domainXRP.XRPPendingMultisig, error) {
	params := sqlcgen.GetXRPPendingMultisigsByAccountIDParams{
		AccountID: accountID,
		Status:    convertDomainStatusToSqlc(status),
	}

	sqlcPendings, err := r.queries.GetXRPPendingMultisigsByAccountID(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetXRPPendingMultisigsByAccountID(): %w", err)
	}

	result := make([]*domainXRP.XRPPendingMultisig, 0, len(sqlcPendings))
	for i := range sqlcPendings {
		result = append(result, convertToXRPPendingMultisig(&sqlcPendings[i]))
	}

	return result, nil
}

// GetByStatus returns all pending multi-sig transactions with a specific status
func (r *XRPPendingMultisigRepositorySqlc) GetByStatus(
	ctx context.Context, status domainXRP.MultisigStatus,
) ([]*domainXRP.XRPPendingMultisig, error) {
	sqlcStatus := convertDomainStatusToSqlc(status)

	sqlcPendings, err := r.queries.GetXRPPendingMultisigsByStatus(ctx, sqlcStatus)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetXRPPendingMultisigsByStatus(): %w", err)
	}

	result := make([]*domainXRP.XRPPendingMultisig, 0, len(sqlcPendings))
	for i := range sqlcPendings {
		result = append(result, convertToXRPPendingMultisig(&sqlcPendings[i]))
	}

	return result, nil
}

// GetExpired returns pending transactions that have expired
func (r *XRPPendingMultisigRepositorySqlc) GetExpired(
	ctx context.Context,
) ([]*domainXRP.XRPPendingMultisig, error) {
	sqlcPendings, err := r.queries.GetExpiredXRPPendingMultisigs(ctx, sql.NullTime{Time: time.Now(), Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to call GetExpiredXRPPendingMultisigs(): %w", err)
	}

	result := make([]*domainXRP.XRPPendingMultisig, 0, len(sqlcPendings))
	for i := range sqlcPendings {
		result = append(result, convertToXRPPendingMultisig(&sqlcPendings[i]))
	}

	return result, nil
}

// Insert creates a new pending multi-sig transaction and returns the inserted ID
func (r *XRPPendingMultisigRepositorySqlc) Insert(
	ctx context.Context, pending *domainXRP.XRPPendingMultisig,
) (int64, error) {
	params := sqlcgen.InsertXRPPendingMultisigParams{
		TxUuid:         pending.TxUUID,
		AccountID:      pending.AccountID,
		UnsignedTxJson: pending.UnsignedTxJSON,
		XrpTxType:      pending.XRPTxType,
		RequiredQuorum: pending.RequiredQuorum,
		CurrentWeight:  pending.CurrentWeight,
		Status:         convertDomainStatusToSqlc(pending.Status),
	}

	if pending.ExpiresAt != nil {
		params.ExpiresAt = sql.NullTime{Time: *pending.ExpiresAt, Valid: true}
	}
	if pending.CreatedAt != nil {
		params.CreatedAt = sql.NullTime{Time: *pending.CreatedAt, Valid: true}
	} else {
		params.CreatedAt = sql.NullTime{Time: time.Now(), Valid: true}
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

// UpdateWeight updates the current weight of collected signatures
func (r *XRPPendingMultisigRepositorySqlc) UpdateWeight(
	ctx context.Context, id int64, currentWeight uint32,
) error {
	params := sqlcgen.UpdateXRPPendingMultisigWeightParams{
		CurrentWeight: currentWeight,
		UpdatedAt:     sql.NullTime{Time: time.Now(), Valid: true},
		ID:            id,
	}

	_, err := r.queries.UpdateXRPPendingMultisigWeight(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to call UpdateXRPPendingMultisigWeight(): %w", err)
	}

	return nil
}

// UpdateStatus updates the status of a pending transaction
func (r *XRPPendingMultisigRepositorySqlc) UpdateStatus(
	ctx context.Context, id int64, status domainXRP.MultisigStatus,
) error {
	params := sqlcgen.UpdateXRPPendingMultisigStatusParams{
		Status:    convertDomainStatusToSqlc(status),
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		ID:        id,
	}

	_, err := r.queries.UpdateXRPPendingMultisigStatus(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to call UpdateXRPPendingMultisigStatus(): %w", err)
	}

	return nil
}

// SetReady marks the transaction as ready with the combined tx blob
func (r *XRPPendingMultisigRepositorySqlc) SetReady(
	ctx context.Context, id int64, combinedTxBlob string,
) error {
	params := sqlcgen.UpdateXRPPendingMultisigReadyParams{
		CombinedTxBlob: sql.NullString{String: combinedTxBlob, Valid: true},
		UpdatedAt:      sql.NullTime{Time: time.Now(), Valid: true},
		ID:             id,
	}

	_, err := r.queries.UpdateXRPPendingMultisigReady(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to call UpdateXRPPendingMultisigReady(): %w", err)
	}

	return nil
}

// SetSubmitted marks the transaction as submitted with the tx hash
func (r *XRPPendingMultisigRepositorySqlc) SetSubmitted(
	ctx context.Context, id int64, submittedTxHash string,
) error {
	params := sqlcgen.UpdateXRPPendingMultisigSubmittedParams{
		SubmittedTxHash: sql.NullString{String: submittedTxHash, Valid: true},
		UpdatedAt:       sql.NullTime{Time: time.Now(), Valid: true},
		ID:              id,
	}

	_, err := r.queries.UpdateXRPPendingMultisigSubmitted(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to call UpdateXRPPendingMultisigSubmitted(): %w", err)
	}

	return nil
}

// SetConfirmed marks the transaction as confirmed on the ledger
func (r *XRPPendingMultisigRepositorySqlc) SetConfirmed(
	ctx context.Context, id int64,
) error {
	params := sqlcgen.UpdateXRPPendingMultisigConfirmedParams{
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		ID:        id,
	}

	_, err := r.queries.UpdateXRPPendingMultisigConfirmed(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to call UpdateXRPPendingMultisigConfirmed(): %w", err)
	}

	return nil
}

// SetFailed marks the transaction as failed
func (r *XRPPendingMultisigRepositorySqlc) SetFailed(
	ctx context.Context, id int64,
) error {
	params := sqlcgen.UpdateXRPPendingMultisigFailedParams{
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		ID:        id,
	}

	_, err := r.queries.UpdateXRPPendingMultisigFailed(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to call UpdateXRPPendingMultisigFailed(): %w", err)
	}

	return nil
}

// ExpireAll marks all expired pending transactions as expired
func (r *XRPPendingMultisigRepositorySqlc) ExpireAll(
	ctx context.Context,
) (int64, error) {
	now := time.Now()
	params := sqlcgen.ExpireXRPPendingMultisigsParams{
		UpdatedAt: sql.NullTime{Time: now, Valid: true},
		ExpiresAt: sql.NullTime{Time: now, Valid: true},
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
