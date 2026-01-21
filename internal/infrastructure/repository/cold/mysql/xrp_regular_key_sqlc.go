package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	domainXrp "github.com/hiromaily/go-crypto-wallet/internal/domain/xrp"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/mysql/sqlcgen"
)

// XRPRegularKeyRepositorySqlc is repository for xrp_regular_key table using sqlc
type XRPRegularKeyRepositorySqlc struct {
	queries *sqlcgen.Queries
}

// NewXRPRegularKeyRepositorySqlc returns XRPRegularKeyRepositorySqlc object
func NewXRPRegularKeyRepositorySqlc(dbConn *sql.DB) *XRPRegularKeyRepositorySqlc {
	return &XRPRegularKeyRepositorySqlc{
		queries: sqlcgen.New(dbConn),
	}
}

// convertToXRPRegularKey converts sqlcgen.XrpRegularKey to domain.XRPRegularKey entity.
func convertToXRPRegularKey(sqlcKey *sqlcgen.XrpRegularKey) *domainXrp.XRPRegularKey {
	key := &domainXrp.XRPRegularKey{
		ID:                sqlcKey.ID,
		AccountID:         sqlcKey.AccountID,
		RegularKeyAddress: sqlcKey.RegularKeyAddress,
		PublicKey:         sqlcKey.PublicKey,
		PublicKeyHex:      sqlcKey.PublicKeyHex,
		IsActive:          sqlcKey.IsActive,
	}

	if sqlcKey.SetTxHash.Valid {
		key.SetTxHash = &sqlcKey.SetTxHash.String
	}
	if sqlcKey.CreatedAt.Valid {
		key.CreatedAt = &sqlcKey.CreatedAt.Time
	}
	if sqlcKey.RotatedAt.Valid {
		key.RotatedAt = &sqlcKey.RotatedAt.Time
	}

	return key
}

// GetByAccountID returns the active regular key for an account
func (r *XRPRegularKeyRepositorySqlc) GetByAccountID(
	ctx context.Context, accountID string,
) (*domainXrp.XRPRegularKey, error) {
	sqlcKey, err := r.queries.GetXRPRegularKeyByAccountID(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No active regular key found
		}
		return nil, fmt.Errorf("failed to call GetXRPRegularKeyByAccountID(): %w", err)
	}

	return convertToXRPRegularKey(&sqlcKey), nil
}

// GetAllByAccountID returns all regular keys (active and inactive) for an account
func (r *XRPRegularKeyRepositorySqlc) GetAllByAccountID(
	ctx context.Context, accountID string,
) ([]*domainXrp.XRPRegularKey, error) {
	sqlcKeys, err := r.queries.GetXRPRegularKeysByAccountID(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetXRPRegularKeysByAccountID(): %w", err)
	}

	result := make([]*domainXrp.XRPRegularKey, 0, len(sqlcKeys))
	for i := range sqlcKeys {
		result = append(result, convertToXRPRegularKey(&sqlcKeys[i]))
	}

	return result, nil
}

// GetActiveKeys returns all currently active regular keys
func (r *XRPRegularKeyRepositorySqlc) GetActiveKeys(
	ctx context.Context,
) ([]*domainXrp.XRPRegularKey, error) {
	sqlcKeys, err := r.queries.GetActiveXRPRegularKeys(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetActiveXRPRegularKeys(): %w", err)
	}

	result := make([]*domainXrp.XRPRegularKey, 0, len(sqlcKeys))
	for i := range sqlcKeys {
		result = append(result, convertToXRPRegularKey(&sqlcKeys[i]))
	}

	return result, nil
}

// GetByAddress returns a regular key by its address
func (r *XRPRegularKeyRepositorySqlc) GetByAddress(
	ctx context.Context, regularKeyAddress string,
) (*domainXrp.XRPRegularKey, error) {
	sqlcKey, err := r.queries.GetXRPRegularKeyByAddress(ctx, regularKeyAddress)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Regular key not found
		}
		return nil, fmt.Errorf("failed to call GetXRPRegularKeyByAddress(): %w", err)
	}

	return convertToXRPRegularKey(&sqlcKey), nil
}

// Insert creates a new regular key record
func (r *XRPRegularKeyRepositorySqlc) Insert(
	ctx context.Context, key *domainXrp.XRPRegularKey,
) (int64, error) {
	params := sqlcgen.InsertXRPRegularKeyParams{
		AccountID:         key.AccountID,
		RegularKeyAddress: key.RegularKeyAddress,
		PublicKey:         key.PublicKey,
		PublicKeyHex:      key.PublicKeyHex,
		IsActive:          key.IsActive,
	}

	if key.SetTxHash != nil {
		params.SetTxHash = sql.NullString{String: *key.SetTxHash, Valid: true}
	}
	if key.CreatedAt != nil {
		params.CreatedAt = sql.NullTime{Time: *key.CreatedAt, Valid: true}
	} else {
		params.CreatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	}

	result, err := r.queries.InsertXRPRegularKey(ctx, params)
	if err != nil {
		return 0, fmt.Errorf("failed to call InsertXRPRegularKey(): %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get LastInsertId(): %w", err)
	}

	return id, nil
}

// UpdateStatus updates the active status and rotated_at timestamp
func (r *XRPRegularKeyRepositorySqlc) UpdateStatus(
	ctx context.Context, id int64, isActive bool,
) error {
	params := sqlcgen.UpdateXRPRegularKeyStatusParams{
		IsActive: isActive,
		ID:       id,
	}

	if !isActive {
		params.RotatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	}

	_, err := r.queries.UpdateXRPRegularKeyStatus(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to call UpdateXRPRegularKeyStatus(): %w", err)
	}

	return nil
}

// DeactivateByAccountID deactivates all regular keys for an account
func (r *XRPRegularKeyRepositorySqlc) DeactivateByAccountID(
	ctx context.Context, accountID string,
) error {
	_, err := r.queries.DeactivateXRPRegularKeyByAccountID(ctx, sqlcgen.DeactivateXRPRegularKeyByAccountIDParams{
		RotatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		AccountID: accountID,
	})
	if err != nil {
		return fmt.Errorf("failed to call DeactivateXRPRegularKeyByAccountID(): %w", err)
	}

	return nil
}

// UpdateTxHash updates the SetRegularKey transaction hash
func (r *XRPRegularKeyRepositorySqlc) UpdateTxHash(
	ctx context.Context, id int64, txHash string,
) error {
	_, err := r.queries.UpdateXRPRegularKeyTxHash(ctx, sqlcgen.UpdateXRPRegularKeyTxHashParams{
		SetTxHash: sql.NullString{String: txHash, Valid: true},
		ID:        id,
	})
	if err != nil {
		return fmt.Errorf("failed to call UpdateXRPRegularKeyTxHash(): %w", err)
	}

	return nil
}
