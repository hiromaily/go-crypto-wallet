//go:build integration

package sqlite_test

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainAuth "github.com/hiromaily/go-crypto-wallet/internal/domain/auth"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	coldsqlite "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/cold/sqlite"
)

// setupAuthAccountKeyTable creates the auth_account_key table for testing
func setupAuthAccountKeyTable(t *testing.T, db *sql.DB) {
	t.Helper()

	createTableSQL := `
		CREATE TABLE IF NOT EXISTS auth_account_key (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			coin TEXT NOT NULL,
			key_type TEXT NOT NULL DEFAULT "bip44",
			auth_account TEXT NOT NULL,
			account TEXT NOT NULL,
			p2pkh_address TEXT NOT NULL,
			p2sh_segwit_address TEXT NOT NULL,
			bech32_address TEXT NOT NULL,
			taproot_address TEXT,
			full_public_key TEXT NOT NULL,
			multisig_address TEXT NOT NULL DEFAULT "",
			redeem_script TEXT NOT NULL DEFAULT "",
			wallet_import_format TEXT NOT NULL,
			account_extended_privkey TEXT,
			idx INTEGER NOT NULL,
			addr_status INTEGER NOT NULL DEFAULT 0,
			updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(wallet_import_format)
		);
	`
	_, err := db.Exec(createTableSQL)
	require.NoError(t, err, "failed to create auth_account_key table")
}

// cleanupAuthAccountKeyTable drops the auth_account_key table
func cleanupAuthAccountKeyTable(t *testing.T, db *sql.DB) {
	t.Helper()
	_, _ = db.Exec("DROP TABLE IF EXISTS auth_account_key")
}

// TestAuthAccountKeyRepositorySqlc_Insert tests inserting auth account key
func TestAuthAccountKeyRepositorySqlc_Insert(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupAuthAccountKeyTable(t, db)
	defer cleanupAuthAccountKeyTable(t, db)

	repo := coldsqlite.NewAuthAccountKeyRepositorySqlc(db, domainCoin.BTC)

	authKey := &domainAuth.AuthAccountKey{
		CoinTypeCode:       domainCoin.BTC,
		KeyType:            "bip44",
		AuthAccount:        domainAccount.AuthType("auth1"),
		Account:            domainAccount.AccountTypeClient,
		P2pkhAddress:       "p2pkh_address",
		P2shSegwitAddress:  "p2sh_segwit_address",
		Bech32Address:      "bech32_address",
		FullPublicKey:      "full_public_key",
		WalletImportFormat: "test_wif",
		Idx:                0,
		AddrStatus:         domainAddress.AddrStatusHDKeyGenerated,
	}

	err = repo.Insert(authKey)
	require.NoError(t, err, "Insert should succeed")
}

// TestAuthAccountKeyRepositorySqlc_GetOne tests retrieving auth account key
func TestAuthAccountKeyRepositorySqlc_GetOne(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupAuthAccountKeyTable(t, db)
	defer cleanupAuthAccountKeyTable(t, db)

	repo := coldsqlite.NewAuthAccountKeyRepositorySqlc(db, domainCoin.BTC)

	authKey := &domainAuth.AuthAccountKey{
		CoinTypeCode:       domainCoin.BTC,
		KeyType:            "bip44",
		AuthAccount:        domainAccount.AuthType("auth1"),
		Account:            domainAccount.AccountTypeClient,
		P2pkhAddress:       "p2pkh_address",
		P2shSegwitAddress:  "p2sh_segwit_address",
		Bech32Address:      "bech32_address",
		FullPublicKey:      "full_public_key",
		WalletImportFormat: "test_wif",
		Idx:                0,
		AddrStatus:         domainAddress.AddrStatusHDKeyGenerated,
	}

	// Insert first
	err = repo.Insert(authKey)
	require.NoError(t, err)

	// Retrieve
	retrieved, err := repo.GetOne(domainAccount.AuthType("auth1"))
	require.NoError(t, err, "GetOne should succeed")
	require.NotNil(t, retrieved, "auth key should not be nil")

	assert.Equal(t, domainCoin.BTC, retrieved.CoinTypeCode, "coin type should match")
	assert.Equal(t, domainAccount.AuthType("auth1"), retrieved.AuthAccount, "auth account should match")
	assert.Equal(t, domainAccount.AccountTypeClient, retrieved.Account, "account type should match")
	assert.Equal(t, "full_public_key", retrieved.FullPublicKey, "full public key should match")
}

// TestAuthAccountKeyRepositorySqlc_GetByAccount tests retrieving by account
func TestAuthAccountKeyRepositorySqlc_GetByAccount(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupAuthAccountKeyTable(t, db)
	defer cleanupAuthAccountKeyTable(t, db)

	repo := coldsqlite.NewAuthAccountKeyRepositorySqlc(db, domainCoin.BTC)

	authKey := &domainAuth.AuthAccountKey{
		CoinTypeCode:       domainCoin.BTC,
		KeyType:            "bip44",
		AuthAccount:        domainAccount.AuthType("auth1"),
		Account:            domainAccount.AccountTypeDeposit,
		P2pkhAddress:       "p2pkh_address",
		P2shSegwitAddress:  "p2sh_segwit_address",
		Bech32Address:      "bech32_address",
		FullPublicKey:      "full_public_key",
		WalletImportFormat: "test_wif_deposit",
		Idx:                0,
		AddrStatus:         domainAddress.AddrStatusHDKeyGenerated,
	}

	// Insert first
	err = repo.Insert(authKey)
	require.NoError(t, err)

	// Retrieve by auth type and account type
	retrieved, err := repo.GetByAccount(
		domainAccount.AuthType("auth1"),
		domainAccount.AccountTypeDeposit,
	)
	require.NoError(t, err, "GetByAccount should succeed")
	require.NotNil(t, retrieved, "auth key should not be nil")

	assert.Equal(t, domainAccount.AccountTypeDeposit, retrieved.Account, "account type should match")
}

// TestAuthAccountKeyRepositorySqlc_UpdateAddrStatus tests updating address status
func TestAuthAccountKeyRepositorySqlc_UpdateAddrStatus(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupAuthAccountKeyTable(t, db)
	defer cleanupAuthAccountKeyTable(t, db)

	repo := coldsqlite.NewAuthAccountKeyRepositorySqlc(db, domainCoin.BTC)

	authKey := &domainAuth.AuthAccountKey{
		CoinTypeCode:       domainCoin.BTC,
		KeyType:            "bip44",
		AuthAccount:        domainAccount.AuthType("auth1"),
		Account:            domainAccount.AccountTypeClient,
		P2pkhAddress:       "p2pkh_address",
		P2shSegwitAddress:  "p2sh_segwit_address",
		Bech32Address:      "bech32_address",
		FullPublicKey:      "full_public_key",
		WalletImportFormat: "test_wif",
		Idx:                0,
		AddrStatus:         domainAddress.AddrStatusHDKeyGenerated,
	}

	// Insert first
	err = repo.Insert(authKey)
	require.NoError(t, err)

	// Update status
	rowsAffected, err := repo.UpdateAddrStatus(
		domainAddress.AddrStatusAddressExported,
		"test_wif",
	)
	require.NoError(t, err, "UpdateAddrStatus should succeed")
	assert.Equal(t, int64(1), rowsAffected, "should affect 1 row")

	// Verify update
	retrieved, err := repo.GetOne(domainAccount.AuthType("auth1"))
	require.NoError(t, err)
	assert.Equal(t, domainAddress.AddrStatusAddressExported, retrieved.AddrStatus, "status should be updated")
}

// TestAuthAccountKeyRepositorySqlc_Constructor tests the repository constructor
func TestAuthAccountKeyRepositorySqlc_Constructor(t *testing.T) {
	t.Parallel()

	db := &sql.DB{}
	repo := coldsqlite.NewAuthAccountKeyRepositorySqlc(db, domainCoin.BTC)

	require.NotNil(t, repo, "NewAuthAccountKeyRepositorySqlc should not return nil")
}
