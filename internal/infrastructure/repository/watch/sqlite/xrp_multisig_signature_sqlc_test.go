//go:build integration

package sqlite_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"

	domainXRP "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/xrp"
	watchsqlite "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/watch/sqlite"
)

func setupXRPMultisigSignatureTable(t *testing.T, db *sql.DB) {
	t.Helper()
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS xrp_multisig_signature (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			pending_multisig_id INTEGER NOT NULL,
			signer_account TEXT NOT NULL,
			signed_tx_blob TEXT NOT NULL,
			signer_weight INTEGER NOT NULL,
			signed_at TEXT
		)
	`)
	require.NoError(t, err)
}

func cleanupXRPMultisigSignatureTable(t *testing.T, db *sql.DB) {
	t.Helper()
	_, _ = db.Exec("DROP TABLE IF EXISTS xrp_multisig_signature")
}

func insertTestSignature(
	t *testing.T,
	repo *watchsqlite.XRPMultisigSignatureRepositorySqlc,
	signerAccount, blob string, weight uint32,
) int64 {
	t.Helper()
	id, err := repo.Insert(context.Background(), &domainXRP.XRPMultisigSignature{
		PendingMultisigID: 1,
		SignerAccount:     signerAccount,
		SignedTxBlob:      blob,
		SignerWeight:      weight,
	})
	require.NoError(t, err)
	return id
}

func TestXRPMultisigSignatureRepositorySqlc_Constructor(t *testing.T) {
	t.Parallel()

	db := &sql.DB{}
	repo := watchsqlite.NewXRPMultisigSignatureRepositorySqlc(db)
	require.NotNil(t, repo)
}

func TestXRPMultisigSignatureRepositorySqlc_InsertAndGetByID(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupXRPMultisigSignatureTable(t, db)
	defer cleanupXRPMultisigSignatureTable(t, db)

	repo := watchsqlite.NewXRPMultisigSignatureRepositorySqlc(db)
	id := insertTestSignature(
		t, repo, "rPEPPER7kfTD9w2To4CQk6UCfuHM9c6GDY", "BLOB001", 1,
	)
	assert.Positive(t, id)

	got, err := repo.GetByID(context.Background(), id)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, id, got.ID)
	assert.Equal(t, "rPEPPER7kfTD9w2To4CQk6UCfuHM9c6GDY", got.SignerAccount)
	assert.Equal(t, uint32(1), got.SignerWeight)
}

func TestXRPMultisigSignatureRepositorySqlc_GetByID_NotFound(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupXRPMultisigSignatureTable(t, db)
	defer cleanupXRPMultisigSignatureTable(t, db)

	repo := watchsqlite.NewXRPMultisigSignatureRepositorySqlc(db)

	got, err := repo.GetByID(context.Background(), 9999)
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestXRPMultisigSignatureRepositorySqlc_GetByPendingID(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupXRPMultisigSignatureTable(t, db)
	defer cleanupXRPMultisigSignatureTable(t, db)

	repo := watchsqlite.NewXRPMultisigSignatureRepositorySqlc(db)

	insertTestSignature(t, repo, "rPEPPER7kfTD9w2To4CQk6UCfuHM9c6GDY", "BLOB001", 1)
	insertTestSignature(t, repo, "rG31cLyErnqeVj2eomEjBZtq7PYaupGYzL", "BLOB002", 1)

	sigs, err := repo.GetByPendingID(context.Background(), 1)
	require.NoError(t, err)
	assert.Len(t, sigs, 2)
}

func TestXRPMultisigSignatureRepositorySqlc_GetByPendingAndSigner(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupXRPMultisigSignatureTable(t, db)
	defer cleanupXRPMultisigSignatureTable(t, db)

	repo := watchsqlite.NewXRPMultisigSignatureRepositorySqlc(db)

	insertTestSignature(t, repo, "rPEPPER7kfTD9w2To4CQk6UCfuHM9c6GDY", "BLOB001", 1)

	got, err := repo.GetByPendingAndSigner(
		context.Background(), 1, "rPEPPER7kfTD9w2To4CQk6UCfuHM9c6GDY",
	)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "rPEPPER7kfTD9w2To4CQk6UCfuHM9c6GDY", got.SignerAccount)

	notFound, err := repo.GetByPendingAndSigner(context.Background(), 1, "rNonExistent")
	require.NoError(t, err)
	assert.Nil(t, notFound)
}

func TestXRPMultisigSignatureRepositorySqlc_GetSignatureCount(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupXRPMultisigSignatureTable(t, db)
	defer cleanupXRPMultisigSignatureTable(t, db)

	repo := watchsqlite.NewXRPMultisigSignatureRepositorySqlc(db)

	count, err := repo.GetSignatureCount(context.Background(), 1)
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)

	insertTestSignature(t, repo, "rPEPPER7kfTD9w2To4CQk6UCfuHM9c6GDY", "BLOB001", 1)
	insertTestSignature(t, repo, "rG31cLyErnqeVj2eomEjBZtq7PYaupGYzL", "BLOB002", 1)

	count, err = repo.GetSignatureCount(context.Background(), 1)
	require.NoError(t, err)
	assert.Equal(t, int64(2), count)
}

func TestXRPMultisigSignatureRepositorySqlc_GetTotalWeight(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupXRPMultisigSignatureTable(t, db)
	defer cleanupXRPMultisigSignatureTable(t, db)

	repo := watchsqlite.NewXRPMultisigSignatureRepositorySqlc(db)

	insertTestSignature(t, repo, "rPEPPER7kfTD9w2To4CQk6UCfuHM9c6GDY", "BLOB001", 2)
	insertTestSignature(t, repo, "rG31cLyErnqeVj2eomEjBZtq7PYaupGYzL", "BLOB002", 3)

	total, err := repo.GetTotalWeight(context.Background(), 1)
	require.NoError(t, err)
	assert.Equal(t, uint32(5), total)
}

func TestXRPMultisigSignatureRepositorySqlc_DeleteByPendingID(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupXRPMultisigSignatureTable(t, db)
	defer cleanupXRPMultisigSignatureTable(t, db)

	repo := watchsqlite.NewXRPMultisigSignatureRepositorySqlc(db)

	insertTestSignature(t, repo, "rPEPPER7kfTD9w2To4CQk6UCfuHM9c6GDY", "BLOB001", 1)
	insertTestSignature(t, repo, "rG31cLyErnqeVj2eomEjBZtq7PYaupGYzL", "BLOB002", 1)

	err = repo.DeleteByPendingID(context.Background(), 1)
	require.NoError(t, err)

	sigs, err := repo.GetByPendingID(context.Background(), 1)
	require.NoError(t, err)
	assert.Empty(t, sigs)
}
