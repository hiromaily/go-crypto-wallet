//go:build integration

package sqlite_test

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainXRP "github.com/hiromaily/go-crypto-wallet/internal/domain/xrp"
	coldsqlite "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/cold/sqlite"
)

// setupXRPAccountKeyTable creates the xrp_account_key table for testing
func setupXRPAccountKeyTable(t *testing.T, db *sql.DB) {
	t.Helper()

	createTableSQL := `
		CREATE TABLE IF NOT EXISTS xrp_account_key (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			coin TEXT NOT NULL,
			account TEXT NOT NULL,
			account_id TEXT NOT NULL,
			key_type INTEGER NOT NULL DEFAULT 0,
			master_key TEXT NOT NULL,
			master_seed TEXT NOT NULL,
			master_seed_hex TEXT NOT NULL,
			public_key TEXT NOT NULL,
			public_key_hex TEXT NOT NULL,
			is_regular_key_pair INTEGER NOT NULL DEFAULT 0,
			allocated_id INTEGER NOT NULL DEFAULT 0,
			addr_status INTEGER NOT NULL DEFAULT 0,
			updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(account_id),
			UNIQUE(master_seed)
		);
		CREATE INDEX IF NOT EXISTS idx_account ON xrp_account_key(account);
	`
	_, err := db.Exec(createTableSQL)
	require.NoError(t, err, "failed to create xrp_account_key table")
}

// cleanupXRPAccountKeyTable drops the xrp_account_key table
func cleanupXRPAccountKeyTable(t *testing.T, db *sql.DB) {
	t.Helper()
	_, _ = db.Exec("DROP TABLE IF EXISTS xrp_account_key")
}

// TestXRPAccountKeyRepositorySqlc_InsertBulk tests bulk insert
func TestXRPAccountKeyRepositorySqlc_InsertBulk(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupXRPAccountKeyTable(t, db)
	defer cleanupXRPAccountKeyTable(t, db)

	repo := coldsqlite.NewXRPAccountKeyRepositorySqlc(db, domainCoin.XRP)

	accountKeys := []*domainXRP.XRPAccountKey{
		{
			CoinTypeCode:     domainCoin.XRP,
			Account:          domainAccount.AccountTypeClient,
			AccountID:        "rAccountID1",
			KeyType:          0,
			MasterKey:        "master_key_1",
			MasterSeed:       "master_seed_1",
			MasterSeedHex:    "master_seed_hex_1",
			PublicKey:        "public_key_1",
			PublicKeyHex:     "public_key_hex_1",
			IsRegularKeyPair: false,
			AllocatedID:      0,
			AddrStatus:       domainAddress.AddrStatusHDKeyGenerated,
		},
		{
			CoinTypeCode:     domainCoin.XRP,
			Account:          domainAccount.AccountTypeClient,
			AccountID:        "rAccountID2",
			KeyType:          0,
			MasterKey:        "master_key_2",
			MasterSeed:       "master_seed_2",
			MasterSeedHex:    "master_seed_hex_2",
			PublicKey:        "public_key_2",
			PublicKeyHex:     "public_key_hex_2",
			IsRegularKeyPair: false,
			AllocatedID:      0,
			AddrStatus:       domainAddress.AddrStatusHDKeyGenerated,
		},
	}

	err = repo.InsertBulk(context.Background(), accountKeys)
	require.NoError(t, err, "InsertBulk should succeed")
}

// TestXRPAccountKeyRepositorySqlc_GetAllAddrStatus tests getting all by address status
func TestXRPAccountKeyRepositorySqlc_GetAllAddrStatus(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupXRPAccountKeyTable(t, db)
	defer cleanupXRPAccountKeyTable(t, db)

	repo := coldsqlite.NewXRPAccountKeyRepositorySqlc(db, domainCoin.XRP)

	accountKeys := []*domainXRP.XRPAccountKey{
		{
			CoinTypeCode:     domainCoin.XRP,
			Account:          domainAccount.AccountTypeClient,
			AccountID:        "rAccountID1",
			KeyType:          0,
			MasterKey:        "master_key_1",
			MasterSeed:       "master_seed_1",
			MasterSeedHex:    "master_seed_hex_1",
			PublicKey:        "public_key_1",
			PublicKeyHex:     "public_key_hex_1",
			IsRegularKeyPair: false,
			AllocatedID:      0,
			AddrStatus:       domainAddress.AddrStatusHDKeyGenerated,
		},
		{
			CoinTypeCode:     domainCoin.XRP,
			Account:          domainAccount.AccountTypeClient,
			AccountID:        "rAccountID2",
			KeyType:          0,
			MasterKey:        "master_key_2",
			MasterSeed:       "master_seed_2",
			MasterSeedHex:    "master_seed_hex_2",
			PublicKey:        "public_key_2",
			PublicKeyHex:     "public_key_hex_2",
			IsRegularKeyPair: false,
			AllocatedID:      1,
			AddrStatus:       domainAddress.AddrStatusAddressExported,
		},
	}

	err = repo.InsertBulk(context.Background(), accountKeys)
	require.NoError(t, err)

	keys, err := repo.GetAllAddrStatus(
		context.Background(),
		domainAccount.AccountTypeClient,
		domainAddress.AddrStatusHDKeyGenerated,
	)
	require.NoError(t, err, "GetAllAddrStatus should succeed")
	assert.Equal(t, 1, len(keys), "should get 1 key with status HDKeyGenerated")
	assert.Equal(t, "public_key_1", keys[0].PublicKey)
}

// TestXRPAccountKeyRepositorySqlc_UpdateAddrStatus tests updating address status
func TestXRPAccountKeyRepositorySqlc_UpdateAddrStatus(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupXRPAccountKeyTable(t, db)
	defer cleanupXRPAccountKeyTable(t, db)

	repo := coldsqlite.NewXRPAccountKeyRepositorySqlc(db, domainCoin.XRP)

	accountKeys := []*domainXRP.XRPAccountKey{
		{
			CoinTypeCode:     domainCoin.XRP,
			Account:          domainAccount.AccountTypeClient,
			AccountID:        "rAccountID1",
			KeyType:          0,
			MasterKey:        "master_key_1",
			MasterSeed:       "master_seed_1",
			MasterSeedHex:    "master_seed_hex_1",
			PublicKey:        "public_key_1",
			PublicKeyHex:     "public_key_hex_1",
			IsRegularKeyPair: false,
			AllocatedID:      0,
			AddrStatus:       domainAddress.AddrStatusHDKeyGenerated,
		},
	}

	err = repo.InsertBulk(context.Background(), accountKeys)
	require.NoError(t, err)

	rowsAffected, err := repo.UpdateAddrStatus(
		context.Background(),
		domainAccount.AccountTypeClient,
		domainAddress.AddrStatusAddressExported,
		[]string{"rAccountID1"},
	)
	require.NoError(t, err, "UpdateAddrStatus should succeed")
	assert.Equal(t, int64(1), rowsAffected, "should affect 1 row")
}

// TestXRPAccountKeyRepositorySqlc_GetSecret tests getting secret
func TestXRPAccountKeyRepositorySqlc_GetSecret(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupXRPAccountKeyTable(t, db)
	defer cleanupXRPAccountKeyTable(t, db)

	repo := coldsqlite.NewXRPAccountKeyRepositorySqlc(db, domainCoin.XRP)

	accountKeys := []*domainXRP.XRPAccountKey{
		{
			CoinTypeCode:     domainCoin.XRP,
			Account:          domainAccount.AccountTypeClient,
			AccountID:        "rAccountID1",
			KeyType:          0,
			MasterKey:        "master_key_1",
			MasterSeed:       "master_seed_1",
			MasterSeedHex:    "master_seed_hex_1",
			PublicKey:        "public_key_1",
			PublicKeyHex:     "public_key_hex_1",
			IsRegularKeyPair: false,
			AllocatedID:      0,
			AddrStatus:       domainAddress.AddrStatusHDKeyGenerated,
		},
	}

	err = repo.InsertBulk(context.Background(), accountKeys)
	require.NoError(t, err)

	secret, err := repo.GetSecret(
		context.Background(),
		domainAccount.AccountTypeClient,
		"rAccountID1",
	)
	require.NoError(t, err, "GetSecret should succeed")
	assert.Equal(t, "master_seed_1", secret, "secret should match")
}

// TestXRPAccountKeyRepositorySqlc_Constructor tests the repository constructor
func TestXRPAccountKeyRepositorySqlc_Constructor(t *testing.T) {
	t.Parallel()

	db := &sql.DB{}
	repo := coldsqlite.NewXRPAccountKeyRepositorySqlc(db, domainCoin.XRP)

	require.NotNil(t, repo, "NewXRPAccountKeyRepositorySqlc should not return nil")
}
