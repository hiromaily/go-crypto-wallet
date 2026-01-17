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
	domainEthereum "github.com/hiromaily/go-crypto-wallet/internal/domain/ethereum"
	coldsqlite "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/cold/sqlite"
)

// setupETHAccountKeyTable creates the eth_account_key table for testing
func setupETHAccountKeyTable(t *testing.T, db *sql.DB) {
	t.Helper()

	createTableSQL := `
		CREATE TABLE IF NOT EXISTS eth_account_key (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			account TEXT NOT NULL,
			address TEXT NOT NULL,
			full_public_key TEXT NOT NULL,
			private_key TEXT NOT NULL,
			idx INTEGER NOT NULL,
			addr_status INTEGER NOT NULL DEFAULT 0,
			updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(address),
			UNIQUE(private_key)
		);
		CREATE INDEX IF NOT EXISTS idx_account ON eth_account_key(account);
	`
	_, err := db.Exec(createTableSQL)
	require.NoError(t, err, "failed to create eth_account_key table")
}

// cleanupETHAccountKeyTable drops the eth_account_key table
func cleanupETHAccountKeyTable(t *testing.T, db *sql.DB) {
	t.Helper()
	_, _ = db.Exec("DROP TABLE IF EXISTS eth_account_key")
}

// TestETHAccountKeyRepositorySqlc_InsertBulk tests bulk insert
func TestETHAccountKeyRepositorySqlc_InsertBulk(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupETHAccountKeyTable(t, db)
	defer cleanupETHAccountKeyTable(t, db)

	repo := coldsqlite.NewETHAccountKeyRepositorySqlc(db)

	accountKeys := []*domainEthereum.ETHAccountKey{
		{
			Account:       domainAccount.AccountTypeClient,
			Address:       "0x1234567890123456789012345678901234567890",
			FullPublicKey: "full_public_key_1",
			PrivateKey:    "private_key_1",
			Idx:           0,
			AddrStatus:    domainAddress.AddrStatusHDKeyGenerated,
		},
		{
			Account:       domainAccount.AccountTypeClient,
			Address:       "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd",
			FullPublicKey: "full_public_key_2",
			PrivateKey:    "private_key_2",
			Idx:           1,
			AddrStatus:    domainAddress.AddrStatusHDKeyGenerated,
		},
	}

	err = repo.InsertBulk(accountKeys)
	require.NoError(t, err, "InsertBulk should succeed")
}

// TestETHAccountKeyRepositorySqlc_GetMaxIndex tests getting max index
func TestETHAccountKeyRepositorySqlc_GetMaxIndex(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupETHAccountKeyTable(t, db)
	defer cleanupETHAccountKeyTable(t, db)

	repo := coldsqlite.NewETHAccountKeyRepositorySqlc(db)

	accountKeys := []*domainEthereum.ETHAccountKey{
		{
			Account:       domainAccount.AccountTypeClient,
			Address:       "0x1234567890123456789012345678901234567890",
			FullPublicKey: "full_public_key_1",
			PrivateKey:    "private_key_1",
			Idx:           5,
			AddrStatus:    domainAddress.AddrStatusHDKeyGenerated,
		},
		{
			Account:       domainAccount.AccountTypeClient,
			Address:       "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd",
			FullPublicKey: "full_public_key_2",
			PrivateKey:    "private_key_2",
			Idx:           10,
			AddrStatus:    domainAddress.AddrStatusHDKeyGenerated,
		},
	}

	err = repo.InsertBulk(accountKeys)
	require.NoError(t, err)

	ctx := context.Background()
	maxIdx, err := repo.GetMaxIndex(ctx, domainAccount.AccountTypeClient)
	require.NoError(t, err, "GetMaxIndex should succeed")
	assert.Equal(t, int64(10), maxIdx, "max index should be 10")
}

// TestETHAccountKeyRepositorySqlc_GetOneMaxID tests getting one by max ID
func TestETHAccountKeyRepositorySqlc_GetOneMaxID(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupETHAccountKeyTable(t, db)
	defer cleanupETHAccountKeyTable(t, db)

	repo := coldsqlite.NewETHAccountKeyRepositorySqlc(db)

	accountKeys := []*domainEthereum.ETHAccountKey{
		{
			Account:       domainAccount.AccountTypeClient,
			Address:       "0x1234567890123456789012345678901234567890",
			FullPublicKey: "full_public_key_1",
			PrivateKey:    "private_key_1",
			Idx:           0,
			AddrStatus:    domainAddress.AddrStatusHDKeyGenerated,
		},
		{
			Account:       domainAccount.AccountTypeClient,
			Address:       "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd",
			FullPublicKey: "full_public_key_2",
			PrivateKey:    "private_key_2",
			Idx:           1,
			AddrStatus:    domainAddress.AddrStatusHDKeyGenerated,
		},
	}

	err = repo.InsertBulk(accountKeys)
	require.NoError(t, err)

	accountKey, err := repo.GetOneMaxID(domainAccount.AccountTypeClient)
	require.NoError(t, err, "GetOneMaxID should succeed")
	require.NotNil(t, accountKey, "account key should not be nil")

	assert.Equal(t, "full_public_key_2", accountKey.FullPublicKey, "should get the last inserted key")
}

// TestETHAccountKeyRepositorySqlc_GetAllAddrStatus tests getting all by address status
func TestETHAccountKeyRepositorySqlc_GetAllAddrStatus(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupETHAccountKeyTable(t, db)
	defer cleanupETHAccountKeyTable(t, db)

	repo := coldsqlite.NewETHAccountKeyRepositorySqlc(db)

	accountKeys := []*domainEthereum.ETHAccountKey{
		{
			Account:       domainAccount.AccountTypeClient,
			Address:       "0x1234567890123456789012345678901234567890",
			FullPublicKey: "full_public_key_1",
			PrivateKey:    "private_key_1",
			Idx:           0,
			AddrStatus:    domainAddress.AddrStatusHDKeyGenerated,
		},
		{
			Account:       domainAccount.AccountTypeClient,
			Address:       "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd",
			FullPublicKey: "full_public_key_2",
			PrivateKey:    "private_key_2",
			Idx:           1,
			AddrStatus:    domainAddress.AddrStatusAddressExported,
		},
	}

	err = repo.InsertBulk(accountKeys)
	require.NoError(t, err)

	keys, err := repo.GetAllAddrStatus(
		domainAccount.AccountTypeClient,
		domainAddress.AddrStatusHDKeyGenerated,
	)
	require.NoError(t, err, "GetAllAddrStatus should succeed")
	assert.Equal(t, 1, len(keys), "should get 1 key with status None")
	assert.Equal(t, "full_public_key_1", keys[0].FullPublicKey)
}

// TestETHAccountKeyRepositorySqlc_UpdateAddrStatus tests updating address status
func TestETHAccountKeyRepositorySqlc_UpdateAddrStatus(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupETHAccountKeyTable(t, db)
	defer cleanupETHAccountKeyTable(t, db)

	repo := coldsqlite.NewETHAccountKeyRepositorySqlc(db)

	accountKeys := []*domainEthereum.ETHAccountKey{
		{
			Account:       domainAccount.AccountTypeClient,
			Address:       "0x1234567890123456789012345678901234567890",
			FullPublicKey: "full_public_key_1",
			PrivateKey:    "private_key_1",
			Idx:           0,
			AddrStatus:    domainAddress.AddrStatusHDKeyGenerated,
		},
	}

	err = repo.InsertBulk(accountKeys)
	require.NoError(t, err)

	rowsAffected, err := repo.UpdateAddrStatus(
		domainAccount.AccountTypeClient,
		domainAddress.AddrStatusAddressExported,
		[]string{"private_key_1"},
	)
	require.NoError(t, err, "UpdateAddrStatus should succeed")
	assert.Equal(t, int64(1), rowsAffected, "should affect 1 row")
}

// TestETHAccountKeyRepositorySqlc_Constructor tests the repository constructor
func TestETHAccountKeyRepositorySqlc_Constructor(t *testing.T) {
	t.Parallel()

	db := &sql.DB{}
	repo := coldsqlite.NewETHAccountKeyRepositorySqlc(db)

	require.NotNil(t, repo, "NewETHAccountKeyRepositorySqlc should not return nil")
}
