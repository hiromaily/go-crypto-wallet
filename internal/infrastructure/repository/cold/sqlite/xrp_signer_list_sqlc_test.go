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
	coldsqlite "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/cold/sqlite"
)

// setupXRPSignerListTable creates the xrp_signer_list table for testing.
func setupXRPSignerListTable(t *testing.T, db *sql.DB) {
	t.Helper()
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS xrp_signer_list (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			account_id TEXT NOT NULL,
			signer_quorum INTEGER NOT NULL,
			is_active INTEGER NOT NULL DEFAULT 1,
			set_tx_hash TEXT,
			created_at TEXT,
			updated_at TEXT
		)
	`)
	require.NoError(t, err)
}

// cleanupXRPSignerListTable drops the table after the test.
func cleanupXRPSignerListTable(t *testing.T, db *sql.DB) {
	t.Helper()
	_, _ = db.Exec("DROP TABLE IF EXISTS xrp_signer_list")
}

func TestXRPSignerListRepositorySqlc_Constructor(t *testing.T) {
	t.Parallel()

	db := &sql.DB{}
	repo := coldsqlite.NewXRPSignerListRepositorySqlc(db)
	require.NotNil(t, repo)
}

func TestXRPSignerListRepositorySqlc_InsertAndGetByAccountID(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupXRPSignerListTable(t, db)
	defer cleanupXRPSignerListTable(t, db)

	repo := coldsqlite.NewXRPSignerListRepositorySqlc(db)

	signerList := &domainXRP.XRPSignerList{
		AccountID:    "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh",
		SignerQuorum: 2,
		IsActive:     true,
	}

	id, err := repo.Insert(context.Background(), signerList)
	require.NoError(t, err)
	assert.Positive(t, id)

	got, err := repo.GetByAccountID(context.Background(), "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh")
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh", got.AccountID)
	assert.Equal(t, uint32(2), got.SignerQuorum)
	assert.True(t, got.IsActive)
}

func TestXRPSignerListRepositorySqlc_GetByAccountID_NotFound(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupXRPSignerListTable(t, db)
	defer cleanupXRPSignerListTable(t, db)

	repo := coldsqlite.NewXRPSignerListRepositorySqlc(db)

	got, err := repo.GetByAccountID(context.Background(), "rNonExistent")
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestXRPSignerListRepositorySqlc_GetByID(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupXRPSignerListTable(t, db)
	defer cleanupXRPSignerListTable(t, db)

	repo := coldsqlite.NewXRPSignerListRepositorySqlc(db)

	id, err := repo.Insert(context.Background(), &domainXRP.XRPSignerList{
		AccountID:    "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh",
		SignerQuorum: 1,
		IsActive:     true,
	})
	require.NoError(t, err)

	got, err := repo.GetByID(context.Background(), id)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, id, got.ID)
	assert.Equal(t, uint32(1), got.SignerQuorum)
}

func TestXRPSignerListRepositorySqlc_GetByID_NotFound(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupXRPSignerListTable(t, db)
	defer cleanupXRPSignerListTable(t, db)

	repo := coldsqlite.NewXRPSignerListRepositorySqlc(db)

	got, err := repo.GetByID(context.Background(), 9999)
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestXRPSignerListRepositorySqlc_GetHistoryByAccountID(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupXRPSignerListTable(t, db)
	defer cleanupXRPSignerListTable(t, db)

	repo := coldsqlite.NewXRPSignerListRepositorySqlc(db)
	account := "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh"
	ctx := context.Background()

	_, err = repo.Insert(ctx, &domainXRP.XRPSignerList{
		AccountID: account, SignerQuorum: 1, IsActive: true,
	})
	require.NoError(t, err)
	_, err = repo.Insert(ctx, &domainXRP.XRPSignerList{
		AccountID: account, SignerQuorum: 2, IsActive: false,
	})
	require.NoError(t, err)

	history, err := repo.GetHistoryByAccountID(context.Background(), account)
	require.NoError(t, err)
	assert.Len(t, history, 2)
}

func TestXRPSignerListRepositorySqlc_UpdateStatus(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupXRPSignerListTable(t, db)
	defer cleanupXRPSignerListTable(t, db)

	repo := coldsqlite.NewXRPSignerListRepositorySqlc(db)

	id, err := repo.Insert(context.Background(), &domainXRP.XRPSignerList{
		AccountID:    "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh",
		SignerQuorum: 2,
		IsActive:     true,
	})
	require.NoError(t, err)

	err = repo.UpdateStatus(context.Background(), id, false)
	require.NoError(t, err)

	got, err := repo.GetByID(context.Background(), id)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.False(t, got.IsActive)
}

func TestXRPSignerListRepositorySqlc_UpdateTxHash(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupXRPSignerListTable(t, db)
	defer cleanupXRPSignerListTable(t, db)

	repo := coldsqlite.NewXRPSignerListRepositorySqlc(db)

	id, err := repo.Insert(context.Background(), &domainXRP.XRPSignerList{
		AccountID:    "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh",
		SignerQuorum: 2,
		IsActive:     true,
	})
	require.NoError(t, err)

	txHash := "ABCDEF1234567890"
	err = repo.UpdateTxHash(context.Background(), id, txHash)
	require.NoError(t, err)

	got, err := repo.GetByID(context.Background(), id)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.NotNil(t, got.SetTxHash)
	assert.Equal(t, txHash, *got.SetTxHash)
}

func TestXRPSignerListRepositorySqlc_DeactivateByAccountID(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupXRPSignerListTable(t, db)
	defer cleanupXRPSignerListTable(t, db)

	repo := coldsqlite.NewXRPSignerListRepositorySqlc(db)
	account := "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh"

	_, err = repo.Insert(context.Background(), &domainXRP.XRPSignerList{
		AccountID: account, SignerQuorum: 2, IsActive: true,
	})
	require.NoError(t, err)

	// Confirm active before deactivation
	got, err := repo.GetByAccountID(context.Background(), account)
	require.NoError(t, err)
	require.NotNil(t, got)

	err = repo.DeactivateByAccountID(context.Background(), account)
	require.NoError(t, err)

	// Should return nil after deactivation (no active list)
	got, err = repo.GetByAccountID(context.Background(), account)
	require.NoError(t, err)
	assert.Nil(t, got)
}
