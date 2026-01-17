package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"modernc.org/sqlite"
	sqlitelib "modernc.org/sqlite/lib"

	"github.com/hiromaily/go-crypto-wallet/internal/domain/multisig"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/sqlite/sqlcgen"
)

// NonceRepositorySqlc implements NonceRepository using SQLC-generated code.
// This repository is used by cold wallets (keygen and sign) for secure nonce storage.
type NonceRepositorySqlc struct {
	queries *sqlcgen.Queries
}

// NewNonceRepositorySqlc creates a new SQLC-based nonce repocold.
func NewNonceRepositorySqlc(dbConn *sql.DB) *NonceRepositorySqlc {
	return &NonceRepositorySqlc{
		queries: sqlcgen.New(dbConn),
	}
}

// SaveNonce stores a nonce commitment for a specific signer and transaction.
func (r *NonceRepositorySqlc) SaveNonce(
	ctx context.Context,
	signerID, transactionID string,
	publicNonce [66]byte,
) error {
	_, err := r.queries.SaveNonce(ctx, sqlcgen.SaveNonceParams{
		SignerID:      signerID,
		TransactionID: transactionID,
		PublicNonce:   publicNonce[:],
		IsUsed:        boolToInt64(false),
		CreatedAt:     sql.NullString{String: time.Now().Format("2006-01-02 15:04:05"), Valid: true},
	})
	if err != nil {
		// Check for duplicate key error
		// SQLite error code: SQLITE_CONSTRAINT_PRIMARYKEY=1555 or SQLITE_CONSTRAINT_UNIQUE=2067
		var sqliteErr *sqlite.Error
		if errors.As(err, &sqliteErr) {
			code := sqliteErr.Code()
			if code == sqlitelib.SQLITE_CONSTRAINT ||
				code == sqlitelib.SQLITE_CONSTRAINT_UNIQUE ||
				code == sqlitelib.SQLITE_CONSTRAINT_PRIMARYKEY {
				return multisig.ErrNonceDuplicate
			}
		}
		return fmt.Errorf("failed to save nonce: %w", err)
	}

	return nil
}

// GetNonce retrieves a nonce commitment for a specific signer and transaction.
func (r *NonceRepositorySqlc) GetNonce(
	ctx context.Context,
	signerID, transactionID string,
) (*multisig.NonceCommitment, error) {
	nonce, err := r.queries.GetNonceBySignerAndTx(ctx, sqlcgen.GetNonceBySignerAndTxParams{
		SignerID:      signerID,
		TransactionID: transactionID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, multisig.ErrNonceNotFound
		}
		return nil, fmt.Errorf("failed to get nonce: %w", err)
	}

	return convertToNonceCommitment(&nonce)
}

// GetAllNoncesForTransaction retrieves all nonce commitments for a transaction.
func (r *NonceRepositorySqlc) GetAllNoncesForTransaction(
	ctx context.Context,
	transactionID string,
) ([]*multisig.NonceCommitment, error) {
	nonces, err := r.queries.GetAllNoncesForTransaction(ctx, transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get all nonces for transaction: %w", err)
	}

	result := make([]*multisig.NonceCommitment, 0, len(nonces))
	for i := range nonces {
		commitment, err := convertToNonceCommitment(&nonces[i])
		if err != nil {
			return nil, fmt.Errorf("failed to convert nonce at index %d: %w", i, err)
		}
		result = append(result, commitment)
	}

	return result, nil
}

// GetUnusedNoncesForTransaction retrieves all unused nonce commitments for a transaction.
func (r *NonceRepositorySqlc) GetUnusedNoncesForTransaction(
	ctx context.Context,
	transactionID string,
) ([]*multisig.NonceCommitment, error) {
	nonces, err := r.queries.GetUnusedNoncesForTransaction(ctx, transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get unused nonces for transaction: %w", err)
	}

	result := make([]*multisig.NonceCommitment, 0, len(nonces))
	for i := range nonces {
		commitment, err := convertToNonceCommitment(&nonces[i])
		if err != nil {
			return nil, fmt.Errorf("failed to convert nonce at index %d: %w", i, err)
		}
		result = append(result, commitment)
	}

	return result, nil
}

// MarkNonceUsed marks a nonce as used after successful signing.
func (r *NonceRepositorySqlc) MarkNonceUsed(
	ctx context.Context,
	signerID, transactionID string,
) error {
	result, err := r.queries.MarkNonceUsed(ctx, sqlcgen.MarkNonceUsedParams{
		IsUsed:        boolToInt64(true),
		UsedAt:        sql.NullString{String: time.Now().Format("2006-01-02 15:04:05"), Valid: true},
		SignerID:      signerID,
		TransactionID: transactionID,
	})
	if err != nil {
		return fmt.Errorf("failed to mark nonce as used: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return multisig.ErrNonceNotFound
	}

	return nil
}

// DeleteNoncesForTransaction removes all nonces associated with a transaction.
func (r *NonceRepositorySqlc) DeleteNoncesForTransaction(
	ctx context.Context,
	transactionID string,
) error {
	_, err := r.queries.DeleteNoncesForTransaction(ctx, transactionID)
	if err != nil {
		return fmt.Errorf("failed to delete nonces for transaction: %w", err)
	}

	return nil
}

// CleanupOldUnusedNonces removes unused nonces older than the specified time.
func (r *NonceRepositorySqlc) CleanupOldUnusedNonces(
	ctx context.Context,
	olderThan time.Time,
) error {
	createdAt := sql.NullString{String: olderThan.Format("2006-01-02 15:04:05"), Valid: true}
	_, err := r.queries.CleanupOldUnusedNonces(ctx, createdAt)
	if err != nil {
		return fmt.Errorf("failed to cleanup old unused nonces: %w", err)
	}

	return nil
}

// convertToNonceCommitment converts a SQLC model to a domain NonceCommitment.
func convertToNonceCommitment(nonce *sqlcgen.Musig2Nonce) (*multisig.NonceCommitment, error) {
	if len(nonce.PublicNonce) != 66 {
		return nil, fmt.Errorf("invalid public nonce length: expected 66 bytes, got %d", len(nonce.PublicNonce))
	}

	var publicNonce [66]byte
	copy(publicNonce[:], nonce.PublicNonce)

	return multisig.NewNonceCommitment(publicNonce, nonce.SignerID)
}
