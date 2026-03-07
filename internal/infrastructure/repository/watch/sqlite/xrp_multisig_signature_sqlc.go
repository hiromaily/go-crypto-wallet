package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	domainXRP "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/xrp"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/sqlite/sqlcgen"
)

// XRPMultisigSignatureRepositorySqlc is repository for xrp_multisig_signature table using sqlc (SQLite).
type XRPMultisigSignatureRepositorySqlc struct {
	queries *sqlcgen.Queries
}

// NewXRPMultisigSignatureRepositorySqlc returns a new XRPMultisigSignatureRepositorySqlc.
func NewXRPMultisigSignatureRepositorySqlc(dbConn *sql.DB) *XRPMultisigSignatureRepositorySqlc {
	return &XRPMultisigSignatureRepositorySqlc{
		queries: sqlcgen.New(dbConn),
	}
}

// convertToXRPMultisigSignatureFromSQLite converts a sqlcgen.XrpMultisigSignature (SQLite row) to a domain entity.
// SQLite stores timestamps as TEXT and weights as INTEGER.
func convertToXRPMultisigSignatureFromSQLite(
	sqlcSig *sqlcgen.XrpMultisigSignature,
) *domainXRP.XRPMultisigSignature {
	sig := &domainXRP.XRPMultisigSignature{
		ID:                sqlcSig.ID,
		PendingMultisigID: sqlcSig.PendingMultisigID,
		SignerAccount:     sqlcSig.SignerAccount,
		SignedTxBlob:      sqlcSig.SignedTxBlob,
		SignerWeight:      uint32(sqlcSig.SignerWeight),
	}

	if sqlcSig.SignedAt.Valid {
		t, err := time.Parse(sqliteTimestampFormat, sqlcSig.SignedAt.String)
		if err == nil {
			sig.SignedAt = &t
		}
	}

	return sig
}

// GetByID returns a signature by its ID.
func (r *XRPMultisigSignatureRepositorySqlc) GetByID(
	ctx context.Context, id int64,
) (*domainXRP.XRPMultisigSignature, error) {
	sqlcSig, err := r.queries.GetXRPMultisigSignatureByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, fmt.Errorf("failed to call GetXRPMultisigSignatureByID(): %w", err)
	}

	return convertToXRPMultisigSignatureFromSQLite(&sqlcSig), nil
}

// GetByPendingID returns all signatures for a pending multi-sig transaction.
func (r *XRPMultisigSignatureRepositorySqlc) GetByPendingID(
	ctx context.Context, pendingMultisigID int64,
) ([]*domainXRP.XRPMultisigSignature, error) {
	sqlcSigs, err := r.queries.GetXRPMultisigSignaturesByPendingID(ctx, pendingMultisigID)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetXRPMultisigSignaturesByPendingID(): %w", err)
	}

	result := make([]*domainXRP.XRPMultisigSignature, 0, len(sqlcSigs))
	for i := range sqlcSigs {
		result = append(result, convertToXRPMultisigSignatureFromSQLite(&sqlcSigs[i]))
	}

	return result, nil
}

// GetByPendingAndSigner returns a signature by pending ID and signer account.
func (r *XRPMultisigSignatureRepositorySqlc) GetByPendingAndSigner(
	ctx context.Context, pendingMultisigID int64, signerAccount string,
) (*domainXRP.XRPMultisigSignature, error) {
	params := sqlcgen.GetXRPMultisigSignatureByPendingAndSignerParams{
		PendingMultisigID: pendingMultisigID,
		SignerAccount:     signerAccount,
	}

	sqlcSig, err := r.queries.GetXRPMultisigSignatureByPendingAndSigner(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, fmt.Errorf("failed to call GetXRPMultisigSignatureByPendingAndSigner(): %w", err)
	}

	return convertToXRPMultisigSignatureFromSQLite(&sqlcSig), nil
}

// GetSignatureCount returns the number of signatures for a pending transaction.
func (r *XRPMultisigSignatureRepositorySqlc) GetSignatureCount(
	ctx context.Context, pendingMultisigID int64,
) (int64, error) {
	count, err := r.queries.GetSignatureCountByPendingID(ctx, pendingMultisigID)
	if err != nil {
		return 0, fmt.Errorf("failed to call GetSignatureCountByPendingID(): %w", err)
	}

	return count, nil
}

// GetTotalWeight returns the total weight of signatures for a pending transaction.
func (r *XRPMultisigSignatureRepositorySqlc) GetTotalWeight(
	ctx context.Context, pendingMultisigID int64,
) (uint32, error) {
	result, err := r.queries.GetTotalSignedWeightByPendingID(ctx, pendingMultisigID)
	if err != nil {
		return 0, fmt.Errorf("failed to call GetTotalSignedWeightByPendingID(): %w", err)
	}

	// SQLite COALESCE(SUM(...), 0) may return different numeric types.
	switch v := result.(type) {
	case int64:
		return uint32(v), nil
	case float64:
		return uint32(v), nil
	default:
		return 0, nil
	}
}

// Insert creates a new signature and returns the inserted ID.
func (r *XRPMultisigSignatureRepositorySqlc) Insert(
	ctx context.Context, signature *domainXRP.XRPMultisigSignature,
) (int64, error) {
	signedAt := time.Now()
	if signature.SignedAt != nil {
		signedAt = *signature.SignedAt
	}

	params := sqlcgen.InsertXRPMultisigSignatureParams{
		PendingMultisigID: signature.PendingMultisigID,
		SignerAccount:     signature.SignerAccount,
		SignedTxBlob:      signature.SignedTxBlob,
		SignerWeight:      int64(signature.SignerWeight),
		SignedAt:          sql.NullString{String: signedAt.Format(sqliteTimestampFormat), Valid: true},
	}

	result, err := r.queries.InsertXRPMultisigSignature(ctx, params)
	if err != nil {
		return 0, fmt.Errorf("failed to call InsertXRPMultisigSignature(): %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get LastInsertId(): %w", err)
	}

	return id, nil
}

// DeleteByPendingID deletes all signatures for a pending multi-sig transaction.
func (r *XRPMultisigSignatureRepositorySqlc) DeleteByPendingID(
	ctx context.Context, pendingMultisigID int64,
) error {
	err := r.queries.DeleteXRPMultisigSignaturesByPendingID(ctx, pendingMultisigID)
	if err != nil {
		return fmt.Errorf("failed to call DeleteXRPMultisigSignaturesByPendingID(): %w", err)
	}

	return nil
}
