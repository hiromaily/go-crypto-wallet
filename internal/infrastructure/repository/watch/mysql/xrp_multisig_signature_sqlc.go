package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	domainXrp "github.com/hiromaily/go-crypto-wallet/internal/domain/xrp"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/mysql/sqlcgen"
)

// XRPMultisigSignatureRepositorySqlc is repository for xrp_multisig_signature table using sqlc
type XRPMultisigSignatureRepositorySqlc struct {
	queries *sqlcgen.Queries
}

// NewXRPMultisigSignatureRepositorySqlc returns XRPMultisigSignatureRepositorySqlc object
func NewXRPMultisigSignatureRepositorySqlc(dbConn *sql.DB) *XRPMultisigSignatureRepositorySqlc {
	return &XRPMultisigSignatureRepositorySqlc{
		queries: sqlcgen.New(dbConn),
	}
}

// convertToXRPMultisigSignature converts sqlcgen.XrpMultisigSignature to domain.XRPMultisigSignature entity.
func convertToXRPMultisigSignature(sqlcSig *sqlcgen.XrpMultisigSignature) *domainXrp.XRPMultisigSignature {
	sig := &domainXrp.XRPMultisigSignature{
		ID:                sqlcSig.ID,
		PendingMultisigID: sqlcSig.PendingMultisigID,
		SignerAccount:     sqlcSig.SignerAccount,
		SignedTxBlob:      sqlcSig.SignedTxBlob,
		SignerWeight:      sqlcSig.SignerWeight,
	}

	if sqlcSig.SignedAt.Valid {
		sig.SignedAt = &sqlcSig.SignedAt.Time
	}

	return sig
}

// GetByID returns a signature by its ID
func (r *XRPMultisigSignatureRepositorySqlc) GetByID(
	ctx context.Context, id int64,
) (*domainXrp.XRPMultisigSignature, error) {
	sqlcSig, err := r.queries.GetXRPMultisigSignatureByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, fmt.Errorf("failed to call GetXRPMultisigSignatureByID(): %w", err)
	}

	return convertToXRPMultisigSignature(&sqlcSig), nil
}

// GetByPendingID returns all signatures for a pending multi-sig transaction
func (r *XRPMultisigSignatureRepositorySqlc) GetByPendingID(
	ctx context.Context, pendingMultisigID int64,
) ([]*domainXrp.XRPMultisigSignature, error) {
	sqlcSigs, err := r.queries.GetXRPMultisigSignaturesByPendingID(ctx, pendingMultisigID)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetXRPMultisigSignaturesByPendingID(): %w", err)
	}

	result := make([]*domainXrp.XRPMultisigSignature, 0, len(sqlcSigs))
	for i := range sqlcSigs {
		result = append(result, convertToXRPMultisigSignature(&sqlcSigs[i]))
	}

	return result, nil
}

// GetByPendingAndSigner returns a signature by pending ID and signer account
func (r *XRPMultisigSignatureRepositorySqlc) GetByPendingAndSigner(
	ctx context.Context, pendingMultisigID int64, signerAccount string,
) (*domainXrp.XRPMultisigSignature, error) {
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

	return convertToXRPMultisigSignature(&sqlcSig), nil
}

// GetSignatureCount returns the number of signatures for a pending transaction
func (r *XRPMultisigSignatureRepositorySqlc) GetSignatureCount(
	ctx context.Context, pendingMultisigID int64,
) (int64, error) {
	count, err := r.queries.GetSignatureCountByPendingID(ctx, pendingMultisigID)
	if err != nil {
		return 0, fmt.Errorf("failed to call GetSignatureCountByPendingID(): %w", err)
	}

	return count, nil
}

// GetTotalWeight returns the total weight of signatures for a pending transaction
func (r *XRPMultisigSignatureRepositorySqlc) GetTotalWeight(
	ctx context.Context, pendingMultisigID int64,
) (uint32, error) {
	result, err := r.queries.GetTotalSignedWeightByPendingID(ctx, pendingMultisigID)
	if err != nil {
		return 0, fmt.Errorf("failed to call GetTotalSignedWeightByPendingID(): %w", err)
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

// Insert creates a new signature and returns the inserted ID
func (r *XRPMultisigSignatureRepositorySqlc) Insert(
	ctx context.Context, signature *domainXrp.XRPMultisigSignature,
) (int64, error) {
	params := sqlcgen.InsertXRPMultisigSignatureParams{
		PendingMultisigID: signature.PendingMultisigID,
		SignerAccount:     signature.SignerAccount,
		SignedTxBlob:      signature.SignedTxBlob,
		SignerWeight:      signature.SignerWeight,
	}

	if signature.SignedAt != nil {
		params.SignedAt = sql.NullTime{Time: *signature.SignedAt, Valid: true}
	} else {
		params.SignedAt = sql.NullTime{Time: time.Now(), Valid: true}
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

// DeleteByPendingID deletes all signatures for a pending multi-sig transaction
func (r *XRPMultisigSignatureRepositorySqlc) DeleteByPendingID(
	ctx context.Context, pendingMultisigID int64,
) error {
	err := r.queries.DeleteXRPMultisigSignaturesByPendingID(ctx, pendingMultisigID)
	if err != nil {
		return fmt.Errorf("failed to call DeleteXRPMultisigSignaturesByPendingID(): %w", err)
	}

	return nil
}
