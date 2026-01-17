//go:build integration

package sqlite_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite" // Pure Go SQLite driver (CGO-free)

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAuth "github.com/hiromaily/go-crypto-wallet/internal/domain/auth"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
	coldsqlite "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/cold/sqlite"
)

// setupAuthFullPubkeyTable creates the auth_fullpubkey table for testing
func setupAuthFullPubkeyTable(t *testing.T, db *sql.DB) {
	t.Helper()

	createTableSQL := `
		CREATE TABLE IF NOT EXISTS auth_fullpubkey (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			coin TEXT NOT NULL,
			auth_account TEXT NOT NULL,
			purpose INTEGER NOT NULL DEFAULT 49,
			full_public_key TEXT NOT NULL,
			extended_pubkey TEXT,
			fingerprint TEXT,
			derivation_path TEXT,
			updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(coin, auth_account, purpose)
		);
	`
	_, err := db.Exec(createTableSQL)
	require.NoError(t, err, "failed to create auth_fullpubkey table")
}

// cleanupAuthFullPubkeyTable drops the auth_fullpubkey table
func cleanupAuthFullPubkeyTable(t *testing.T, db *sql.DB) {
	t.Helper()
	_, _ = db.Exec("DROP TABLE IF EXISTS auth_fullpubkey")
}

// TestAuthFullPubkeyRepositorySqlc_Insert tests inserting auth full pubkey
func TestAuthFullPubkeyRepositorySqlc_Insert(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupAuthFullPubkeyTable(t, db)
	defer cleanupAuthFullPubkeyTable(t, db)

	repo := coldsqlite.NewAuthFullPubkeyRepositorySqlc(db, domainCoin.BTC)

	authType := domainAccount.AuthType("auth1")
	fullPubKey := "test_full_public_key"

	err = repo.Insert(authType, fullPubKey)
	require.NoError(t, err, "Insert should succeed")
}

// TestAuthFullPubkeyRepositorySqlc_GetOne tests retrieving auth full pubkey
func TestAuthFullPubkeyRepositorySqlc_GetOne(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupAuthFullPubkeyTable(t, db)
	defer cleanupAuthFullPubkeyTable(t, db)

	repo := coldsqlite.NewAuthFullPubkeyRepositorySqlc(db, domainCoin.BTC)

	authType := domainAccount.AuthType("auth1")
	fullPubKey := "test_full_public_key"

	// Insert first
	err = repo.Insert(authType, fullPubKey)
	require.NoError(t, err)

	// Retrieve
	ctx := context.Background()
	authPubkey, err := repo.GetOne(ctx, authType)
	require.NoError(t, err, "GetOne should succeed")
	require.NotNil(t, authPubkey, "auth pubkey should not be nil")

	assert.Equal(t, domainCoin.BTC, authPubkey.CoinTypeCode, "coin type should match")
	assert.Equal(t, authType, authPubkey.AuthAccount, "auth account should match")
	assert.Equal(t, fullPubKey, authPubkey.FullPublicKey, "full public key should match")
	assert.Equal(t, domainAuth.Purpose(49), authPubkey.Purpose, "purpose should default to 49")
}

// TestAuthFullPubkeyRepositorySqlc_GetOneByPurpose tests retrieving by purpose
func TestAuthFullPubkeyRepositorySqlc_GetOneByPurpose(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupAuthFullPubkeyTable(t, db)
	defer cleanupAuthFullPubkeyTable(t, db)

	repo := coldsqlite.NewAuthFullPubkeyRepositorySqlc(db, domainCoin.BTC)

	authType := domainAccount.AuthType("auth1")
	fullPubKey := "test_full_public_key"
	fingerprint, err := domainKey.NewFingerprint("12345678")
	require.NoError(t, err, "NewFingerprint should succeed")

	// Insert with extended info
	authPubkey := &domainAuth.AuthFullPubkey{
		CoinTypeCode:   domainCoin.BTC,
		AuthAccount:    authType,
		Purpose:        domainAuth.Purpose(84),
		FullPublicKey:  fullPubKey,
		ExtendedPubKey: "xpub...",
		Fingerprint:    &fingerprint,
		DerivationPath: "m/84'/1'/0'",
	}

	err = repo.InsertBulk([]*domainAuth.AuthFullPubkey{authPubkey})
	require.NoError(t, err)

	// Retrieve by purpose
	retrieved, err := repo.GetOneByPurpose(authType, domainAuth.Purpose(84))
	require.NoError(t, err, "GetOneByPurpose should succeed")
	require.NotNil(t, retrieved, "auth pubkey should not be nil")

	assert.Equal(t, domainAuth.Purpose(84), retrieved.Purpose, "purpose should match")
	assert.Equal(t, "xpub...", retrieved.ExtendedPubKey, "extended pubkey should match")
	assert.Equal(t, "m/84'/1'/0'", retrieved.DerivationPath, "derivation path should match")
}

// TestAuthFullPubkeyRepositorySqlc_InsertBulk tests bulk insert
func TestAuthFullPubkeyRepositorySqlc_InsertBulk(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupAuthFullPubkeyTable(t, db)
	defer cleanupAuthFullPubkeyTable(t, db)

	repo := coldsqlite.NewAuthFullPubkeyRepositorySqlc(db, domainCoin.BTC)

	items := []*domainAuth.AuthFullPubkey{
		{
			CoinTypeCode:  domainCoin.BTC,
			AuthAccount:   domainAccount.AuthType("auth1"),
			Purpose:       domainAuth.Purpose(44),
			FullPublicKey: "pubkey1",
		},
		{
			CoinTypeCode:  domainCoin.BTC,
			AuthAccount:   domainAccount.AuthType("auth2"),
			Purpose:       domainAuth.Purpose(49),
			FullPublicKey: "pubkey2",
		},
	}

	err = repo.InsertBulk(items)
	require.NoError(t, err, "InsertBulk should succeed")

	// Verify both were inserted
	authPubkey1, err := repo.GetOneByPurpose(domainAccount.AuthType("auth1"), domainAuth.Purpose(44))
	require.NoError(t, err)
	assert.Equal(t, "pubkey1", authPubkey1.FullPublicKey)

	authPubkey2, err := repo.GetOneByPurpose(domainAccount.AuthType("auth2"), domainAuth.Purpose(49))
	require.NoError(t, err)
	assert.Equal(t, "pubkey2", authPubkey2.FullPublicKey)
}

// TestAuthFullPubkeyRepositorySqlc_Constructor tests the repository constructor
func TestAuthFullPubkeyRepositorySqlc_Constructor(t *testing.T) {
	t.Parallel()

	db := &sql.DB{}
	repo := coldsqlite.NewAuthFullPubkeyRepositorySqlc(db, domainCoin.BTC)

	require.NotNil(t, repo, "NewAuthFullPubkeyRepositorySqlc should not return nil")
}
