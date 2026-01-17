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
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainBitcoin "github.com/hiromaily/go-crypto-wallet/internal/domain/bitcoin"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	coldsqlite "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/cold/sqlite"
)

// setupBTCAccountKeyTable creates the btc_account_key table for testing
func setupBTCAccountKeyTable(t *testing.T, db *sql.DB) {
	t.Helper()

	createTableSQL := `
		CREATE TABLE IF NOT EXISTS btc_account_key (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			coin TEXT NOT NULL,
			key_type TEXT NOT NULL DEFAULT "bip44",
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
		CREATE INDEX IF NOT EXISTS idx_account ON btc_account_key(account);
		CREATE INDEX IF NOT EXISTS idx_multisig_address ON btc_account_key(multisig_address);
	`
	_, err := db.Exec(createTableSQL)
	require.NoError(t, err, "failed to create btc_account_key table")
}

// cleanupBTCAccountKeyTable drops the btc_account_key table
func cleanupBTCAccountKeyTable(t *testing.T, db *sql.DB) {
	t.Helper()
	_, _ = db.Exec("DROP TABLE IF EXISTS btc_account_key")
}

// TestBTCAccountKeyRepositorySqlc_InsertBulk tests bulk insert
func TestBTCAccountKeyRepositorySqlc_InsertBulk(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupBTCAccountKeyTable(t, db)
	defer cleanupBTCAccountKeyTable(t, db)

	repo := coldsqlite.NewBTCAccountKeyRepositorySqlc(db, domainCoin.BTC)

	accountKeys := []*domainBitcoin.BtcAccountKey{
		{
			CoinTypeCode:       domainCoin.BTC,
			KeyType:            "bip44",
			Account:            domainAccount.AccountTypeClient,
			P2pkhAddress:       "p2pkh_address_1",
			P2shSegwitAddress:  "p2sh_segwit_address_1",
			Bech32Address:      "bech32_address_1",
			FullPublicKey:      "full_public_key_1",
			WalletImportFormat: "wif_1",
			Idx:                0,
			AddrStatus:         domainAddress.AddrStatusHDKeyGenerated,
		},
		{
			CoinTypeCode:       domainCoin.BTC,
			KeyType:            "bip44",
			Account:            domainAccount.AccountTypeClient,
			P2pkhAddress:       "p2pkh_address_2",
			P2shSegwitAddress:  "p2sh_segwit_address_2",
			Bech32Address:      "bech32_address_2",
			FullPublicKey:      "full_public_key_2",
			WalletImportFormat: "wif_2",
			Idx:                1,
			AddrStatus:         domainAddress.AddrStatusHDKeyGenerated,
		},
	}

	err = repo.InsertBulk(accountKeys)
	require.NoError(t, err, "InsertBulk should succeed")
}

// TestBTCAccountKeyRepositorySqlc_GetMaxIndex tests getting max index
func TestBTCAccountKeyRepositorySqlc_GetMaxIndex(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupBTCAccountKeyTable(t, db)
	defer cleanupBTCAccountKeyTable(t, db)

	repo := coldsqlite.NewBTCAccountKeyRepositorySqlc(db, domainCoin.BTC)

	accountKeys := []*domainBitcoin.BtcAccountKey{
		{
			CoinTypeCode:       domainCoin.BTC,
			KeyType:            "bip44",
			Account:            domainAccount.AccountTypeClient,
			P2pkhAddress:       "p2pkh_address_1",
			P2shSegwitAddress:  "p2sh_segwit_address_1",
			Bech32Address:      "bech32_address_1",
			FullPublicKey:      "full_public_key_1",
			WalletImportFormat: "wif_1",
			Idx:                5,
			AddrStatus:         domainAddress.AddrStatusHDKeyGenerated,
		},
		{
			CoinTypeCode:       domainCoin.BTC,
			KeyType:            "bip44",
			Account:            domainAccount.AccountTypeClient,
			P2pkhAddress:       "p2pkh_address_2",
			P2shSegwitAddress:  "p2sh_segwit_address_2",
			Bech32Address:      "bech32_address_2",
			FullPublicKey:      "full_public_key_2",
			WalletImportFormat: "wif_2",
			Idx:                10,
			AddrStatus:         domainAddress.AddrStatusHDKeyGenerated,
		},
	}

	err = repo.InsertBulk(accountKeys)
	require.NoError(t, err)

	ctx := context.Background()
	maxIdx, err := repo.GetMaxIndex(ctx, domainAccount.AccountTypeClient)
	require.NoError(t, err, "GetMaxIndex should succeed")
	assert.Equal(t, int64(10), maxIdx, "max index should be 10")
}

// TestBTCAccountKeyRepositorySqlc_GetOneMaxID tests getting one by max ID
func TestBTCAccountKeyRepositorySqlc_GetOneMaxID(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupBTCAccountKeyTable(t, db)
	defer cleanupBTCAccountKeyTable(t, db)

	repo := coldsqlite.NewBTCAccountKeyRepositorySqlc(db, domainCoin.BTC)

	accountKeys := []*domainBitcoin.BtcAccountKey{
		{
			CoinTypeCode:       domainCoin.BTC,
			KeyType:            "bip44",
			Account:            domainAccount.AccountTypeClient,
			P2pkhAddress:       "p2pkh_address_1",
			P2shSegwitAddress:  "p2sh_segwit_address_1",
			Bech32Address:      "bech32_address_1",
			FullPublicKey:      "full_public_key_1",
			WalletImportFormat: "wif_1",
			Idx:                0,
			AddrStatus:         domainAddress.AddrStatusHDKeyGenerated,
		},
		{
			CoinTypeCode:       domainCoin.BTC,
			KeyType:            "bip44",
			Account:            domainAccount.AccountTypeClient,
			P2pkhAddress:       "p2pkh_address_2",
			P2shSegwitAddress:  "p2sh_segwit_address_2",
			Bech32Address:      "bech32_address_2",
			FullPublicKey:      "full_public_key_2",
			WalletImportFormat: "wif_2",
			Idx:                1,
			AddrStatus:         domainAddress.AddrStatusHDKeyGenerated,
		},
	}

	err = repo.InsertBulk(accountKeys)
	require.NoError(t, err)

	accountKey, err := repo.GetOneMaxID(domainAccount.AccountTypeClient)
	require.NoError(t, err, "GetOneMaxID should succeed")
	require.NotNil(t, accountKey, "account key should not be nil")

	assert.Equal(t, "full_public_key_2", accountKey.FullPublicKey, "should get the last inserted key")
}

// TestBTCAccountKeyRepositorySqlc_GetAllAddrStatus tests getting all by address status
func TestBTCAccountKeyRepositorySqlc_GetAllAddrStatus(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupBTCAccountKeyTable(t, db)
	defer cleanupBTCAccountKeyTable(t, db)

	repo := coldsqlite.NewBTCAccountKeyRepositorySqlc(db, domainCoin.BTC)

	accountKeys := []*domainBitcoin.BtcAccountKey{
		{
			CoinTypeCode:       domainCoin.BTC,
			KeyType:            "bip44",
			Account:            domainAccount.AccountTypeClient,
			P2pkhAddress:       "p2pkh_address_1",
			P2shSegwitAddress:  "p2sh_segwit_address_1",
			Bech32Address:      "bech32_address_1",
			FullPublicKey:      "full_public_key_1",
			WalletImportFormat: "wif_1",
			Idx:                0,
			AddrStatus:         domainAddress.AddrStatusHDKeyGenerated,
		},
		{
			CoinTypeCode:       domainCoin.BTC,
			KeyType:            "bip44",
			Account:            domainAccount.AccountTypeClient,
			P2pkhAddress:       "p2pkh_address_2",
			P2shSegwitAddress:  "p2sh_segwit_address_2",
			Bech32Address:      "bech32_address_2",
			FullPublicKey:      "full_public_key_2",
			WalletImportFormat: "wif_2",
			Idx:                1,
			AddrStatus:         domainAddress.AddrStatusAddressExported,
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

// TestBTCAccountKeyRepositorySqlc_UpdateAddrStatus tests updating address status
func TestBTCAccountKeyRepositorySqlc_UpdateAddrStatus(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupBTCAccountKeyTable(t, db)
	defer cleanupBTCAccountKeyTable(t, db)

	repo := coldsqlite.NewBTCAccountKeyRepositorySqlc(db, domainCoin.BTC)

	accountKeys := []*domainBitcoin.BtcAccountKey{
		{
			CoinTypeCode:       domainCoin.BTC,
			KeyType:            "bip44",
			Account:            domainAccount.AccountTypeClient,
			P2pkhAddress:       "p2pkh_address_1",
			P2shSegwitAddress:  "p2sh_segwit_address_1",
			Bech32Address:      "bech32_address_1",
			FullPublicKey:      "full_public_key_1",
			WalletImportFormat: "wif_1",
			Idx:                0,
			AddrStatus:         domainAddress.AddrStatusHDKeyGenerated,
		},
	}

	err = repo.InsertBulk(accountKeys)
	require.NoError(t, err)

	rowsAffected, err := repo.UpdateAddrStatus(
		domainAccount.AccountTypeClient,
		domainAddress.AddrStatusAddressExported,
		[]string{"wif_1"},
	)
	require.NoError(t, err, "UpdateAddrStatus should succeed")
	assert.Equal(t, int64(1), rowsAffected, "should affect 1 row")
}

// TestBTCAccountKeyRepositorySqlc_UpdateMultisigAddr tests updating multisig address
func TestBTCAccountKeyRepositorySqlc_UpdateMultisigAddr(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupBTCAccountKeyTable(t, db)
	defer cleanupBTCAccountKeyTable(t, db)

	repo := coldsqlite.NewBTCAccountKeyRepositorySqlc(db, domainCoin.BTC)

	accountKeys := []*domainBitcoin.BtcAccountKey{
		{
			CoinTypeCode:       domainCoin.BTC,
			KeyType:            "bip44",
			Account:            domainAccount.AccountTypeClient,
			P2pkhAddress:       "p2pkh_address_1",
			P2shSegwitAddress:  "p2sh_segwit_address_1",
			Bech32Address:      "bech32_address_1",
			FullPublicKey:      "full_public_key_1",
			WalletImportFormat: "wif_1",
			Idx:                0,
			AddrStatus:         domainAddress.AddrStatusHDKeyGenerated,
		},
	}

	err = repo.InsertBulk(accountKeys)
	require.NoError(t, err)

	accountKeys[0].MultisigAddress = "multisig_address_1"
	accountKeys[0].RedeemScript = "redeem_script_1"
	accountKeys[0].AddrStatus = domainAddress.AddrStatusMultisigAddressGenerated

	rowsAffected, err := repo.UpdateMultisigAddr(
		domainAccount.AccountTypeClient,
		accountKeys[0],
	)
	require.NoError(t, err, "UpdateMultisigAddr should succeed")
	assert.Equal(t, int64(1), rowsAffected, "should affect 1 row")
}

// TestBTCAccountKeyRepositorySqlc_GetAllMultiAddr tests getting all by multisig address
func TestBTCAccountKeyRepositorySqlc_GetAllMultiAddr(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupBTCAccountKeyTable(t, db)
	defer cleanupBTCAccountKeyTable(t, db)

	repo := coldsqlite.NewBTCAccountKeyRepositorySqlc(db, domainCoin.BTC)

	accountKeys := []*domainBitcoin.BtcAccountKey{
		{
			CoinTypeCode:       domainCoin.BTC,
			KeyType:            "bip44",
			Account:            domainAccount.AccountTypeClient,
			P2pkhAddress:       "p2pkh_address_1",
			P2shSegwitAddress:  "p2sh_segwit_address_1",
			Bech32Address:      "bech32_address_1",
			FullPublicKey:      "full_public_key_1",
			MultisigAddress:    "multisig_address_1",
			WalletImportFormat: "wif_1",
			Idx:                0,
			AddrStatus:         domainAddress.AddrStatusHDKeyGenerated,
		},
		{
			CoinTypeCode:       domainCoin.BTC,
			KeyType:            "bip44",
			Account:            domainAccount.AccountTypeClient,
			P2pkhAddress:       "p2pkh_address_2",
			P2shSegwitAddress:  "p2sh_segwit_address_2",
			Bech32Address:      "bech32_address_2",
			FullPublicKey:      "full_public_key_2",
			MultisigAddress:    "multisig_address_2",
			WalletImportFormat: "wif_2",
			Idx:                1,
			AddrStatus:         domainAddress.AddrStatusHDKeyGenerated,
		},
	}

	err = repo.InsertBulk(accountKeys)
	require.NoError(t, err)

	keys, err := repo.GetAllMultiAddr(
		domainAccount.AccountTypeClient,
		[]string{"multisig_address_1", "multisig_address_2"},
	)
	require.NoError(t, err, "GetAllMultiAddr should succeed")
	assert.Equal(t, 2, len(keys), "should get 2 keys")
}

// TestBTCAccountKeyRepositorySqlc_Constructor tests the repository constructor
func TestBTCAccountKeyRepositorySqlc_Constructor(t *testing.T) {
	t.Parallel()

	db := &sql.DB{}
	repo := coldsqlite.NewBTCAccountKeyRepositorySqlc(db, domainCoin.BTC)

	require.NotNil(t, repo, "NewBTCAccountKeyRepositorySqlc should not return nil")
}
