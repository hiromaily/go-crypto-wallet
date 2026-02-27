//go:build integration

package postgres_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainBTC "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/btc"
	domainETH "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/eth"
	domainXRP "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/xrp"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	coldpostgres "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/cold/postgres"
	dbtest "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/testutil"
	"github.com/hiromaily/go-crypto-wallet/pkg/chains/btc/multisig"
)

// openKeygenDB opens a connection to the keygen PostgreSQL database for testing.
func openKeygenDB(t *testing.T) *sql.DB {
	t.Helper()
	return dbtest.OpenTestDB(t, dbtest.NewTestPostgresConf("POSTGRES_KEYGEN_DB", "keygen"))
}

// uniqueStr generates a test string with a unique nano-timestamp suffix.
func uniqueStr(prefix string) string {
	return dbtest.UniqueStr(prefix)
}

// cleanupBTCAccountKeys deletes test BTC account keys matching the given p2pkh addresses.
func cleanupBTCAccountKeys(t *testing.T, db *sql.DB, p2pkhAddrs ...string) {
	t.Helper()
	for _, addr := range p2pkhAddrs {
		_, err := db.ExecContext(context.Background(),
			"DELETE FROM btc_account_key WHERE p2pkh_address = $1", addr)
		require.NoError(t, err, "failed to clean up BTC account key %q", addr)
	}
}

// cleanupETHAccountKeys deletes test ETH account keys matching the given addresses.
func cleanupETHAccountKeys(t *testing.T, db *sql.DB, addresses ...string) {
	t.Helper()
	for _, addr := range addresses {
		_, err := db.ExecContext(context.Background(),
			"DELETE FROM eth_account_key WHERE address = $1", addr)
		require.NoError(t, err, "failed to clean up ETH account key %q", addr)
	}
}

// cleanupXRPAccountKeys deletes test XRP account keys matching the given account IDs.
func cleanupXRPAccountKeys(t *testing.T, db *sql.DB, accountIDs ...string) {
	t.Helper()
	for _, id := range accountIDs {
		_, err := db.ExecContext(context.Background(),
			"DELETE FROM xrp_account_key WHERE account_id = $1", id)
		require.NoError(t, err, "failed to clean up XRP account key %q", id)
	}
}

// cleanupSeed deletes the seed record for the given coin.
func cleanupSeed(t *testing.T, db *sql.DB, coin string) {
	t.Helper()
	_, err := db.ExecContext(context.Background(),
		"DELETE FROM seed WHERE coin = $1", coin)
	require.NoError(t, err, "failed to clean up seed for coin %q", coin)
}

// cleanupNonces deletes all nonce records for the given transaction ID.
func cleanupNonces(t *testing.T, db *sql.DB, transactionID string) {
	t.Helper()
	_, err := db.ExecContext(context.Background(),
		"DELETE FROM musig2_nonces WHERE transaction_id = $1", transactionID)
	require.NoError(t, err, "failed to clean up nonces for tx %q", transactionID)
}

// cleanupSignerLists deletes test XRP signer lists for the given account ID.
func cleanupSignerLists(t *testing.T, db *sql.DB, accountID string) {
	t.Helper()
	_, err := db.ExecContext(context.Background(),
		"DELETE FROM xrp_signer_list WHERE account_id = $1", accountID)
	require.NoError(t, err, "failed to clean up signer list for account %q", accountID)
}

// ── BTC Account Key ──────────────────────────────────────────────────────────

// TestBTCAccountKeyRepository_InsertBulk_and_GetAllAddrStatus verifies that
// InsertBulk stores keys and GetAllAddrStatus retrieves them filtered by status.
func TestBTCAccountKeyRepository_InsertBulk_and_GetAllAddrStatus(t *testing.T) {
	db := openKeygenDB(t)

	p2pkh1 := uniqueStr("test_btc_p2pkh_1")
	p2pkh2 := uniqueStr("test_btc_p2pkh_2")
	wif1 := uniqueStr("test_btc_wif_1")
	wif2 := uniqueStr("test_btc_wif_2")
	t.Cleanup(func() { cleanupBTCAccountKeys(t, db, p2pkh1, p2pkh2) })

	repo := coldpostgres.NewBTCAccountKeyRepositorySqlc(db, domainCoin.BTC)

	items := []*domainBTC.BTCAccountKey{
		{
			CoinTypeCode:       domainCoin.BTC,
			KeyType:            "bip84",
			Account:            domainAccount.AccountTypeClient,
			P2pkhAddress:       p2pkh1,
			FullPublicKey:      "03" + uniqueStr("pubkey1"),
			WalletImportFormat: wif1,
			Idx:                1,
			AddrStatus:         domainAddress.AddrStatusHDKeyGenerated,
		},
		{
			CoinTypeCode:       domainCoin.BTC,
			KeyType:            "bip84",
			Account:            domainAccount.AccountTypeClient,
			P2pkhAddress:       p2pkh2,
			FullPublicKey:      "03" + uniqueStr("pubkey2"),
			WalletImportFormat: wif2,
			Idx:                2,
			AddrStatus:         domainAddress.AddrStatusHDKeyGenerated,
		},
	}

	err := repo.InsertBulk(items)
	require.NoError(t, err, "InsertBulk should succeed for valid BTC account keys")

	got, err := repo.GetAllAddrStatus(
		domainAccount.AccountTypeClient,
		domainAddress.AddrStatusHDKeyGenerated,
	)
	require.NoError(t, err)

	addrSet := make(map[string]bool, len(got))
	for _, k := range got {
		addrSet[k.P2pkhAddress] = true
	}
	assert.True(t, addrSet[p2pkh1], "GetAllAddrStatus should return the first inserted key")
	assert.True(t, addrSet[p2pkh2], "GetAllAddrStatus should return the second inserted key")
}

// TestBTCAccountKeyRepository_GetOneMaxID verifies that GetOneMaxID returns the
// record with the highest ID for the given coin and account type.
func TestBTCAccountKeyRepository_GetOneMaxID(t *testing.T) {
	db := openKeygenDB(t)

	p2pkh := uniqueStr("test_btc_maxid")
	wif := uniqueStr("test_btc_wif_maxid")
	t.Cleanup(func() { cleanupBTCAccountKeys(t, db, p2pkh) })

	repo := coldpostgres.NewBTCAccountKeyRepositorySqlc(db, domainCoin.BTC)

	err := repo.InsertBulk([]*domainBTC.BTCAccountKey{
		{
			CoinTypeCode:       domainCoin.BTC,
			KeyType:            "bip84",
			Account:            domainAccount.AccountTypeDeposit,
			P2pkhAddress:       p2pkh,
			FullPublicKey:      "03" + uniqueStr("pubkey_maxid"),
			WalletImportFormat: wif,
			Idx:                10,
			AddrStatus:         domainAddress.AddrStatusHDKeyGenerated,
		},
	})
	require.NoError(t, err)

	got, err := repo.GetOneMaxID(domainAccount.AccountTypeDeposit)
	require.NoError(t, err)
	require.NotNil(t, got)

	// The retrieved key should have valid coin and account fields.
	assert.Equal(t, domainCoin.BTC, got.CoinTypeCode)
	assert.Equal(t, domainAccount.AccountTypeDeposit, got.Account)
}

// TestBTCAccountKeyRepository_UpdateAddrStatus verifies that UpdateAddrStatus
// correctly updates the addr_status field for the given WIF keys.
func TestBTCAccountKeyRepository_UpdateAddrStatus(t *testing.T) {
	db := openKeygenDB(t)

	p2pkh := uniqueStr("test_btc_update_status")
	wif := uniqueStr("test_btc_wif_update")
	t.Cleanup(func() { cleanupBTCAccountKeys(t, db, p2pkh) })

	repo := coldpostgres.NewBTCAccountKeyRepositorySqlc(db, domainCoin.BTC)

	err := repo.InsertBulk([]*domainBTC.BTCAccountKey{
		{
			CoinTypeCode:       domainCoin.BTC,
			KeyType:            "bip84",
			Account:            domainAccount.AccountTypeClient,
			P2pkhAddress:       p2pkh,
			FullPublicKey:      "03" + uniqueStr("pubkey_update"),
			WalletImportFormat: wif,
			Idx:                1,
			AddrStatus:         domainAddress.AddrStatusHDKeyGenerated,
		},
	})
	require.NoError(t, err)

	rowsAffected, err := repo.UpdateAddrStatus(
		domainAccount.AccountTypeClient,
		domainAddress.AddrStatusPrivKeyImported,
		[]string{wif},
	)
	require.NoError(t, err)
	assert.EqualValues(t, 1, rowsAffected, "exactly one row should be updated")

	// Confirm the status was updated.
	got, err := repo.GetAllAddrStatus(
		domainAccount.AccountTypeClient,
		domainAddress.AddrStatusPrivKeyImported,
	)
	require.NoError(t, err)
	found := false
	for _, k := range got {
		if k.P2pkhAddress == p2pkh {
			found = true
			assert.Equal(t, domainAddress.AddrStatusPrivKeyImported, k.AddrStatus)
		}
	}
	assert.True(t, found, "updated key should appear in GetAllAddrStatus with new status")
}

// ── ETH Account Key ──────────────────────────────────────────────────────────

// TestETHAccountKeyRepository_InsertBulk_and_GetAllAddrStatus verifies that
// InsertBulk stores ETH keys and GetAllAddrStatus retrieves them.
func TestETHAccountKeyRepository_InsertBulk_and_GetAllAddrStatus(t *testing.T) {
	db := openKeygenDB(t)

	addr1 := uniqueStr("0xtest_eth_addr_1")
	addr2 := uniqueStr("0xtest_eth_addr_2")
	priv1 := uniqueStr("test_eth_priv_1")
	priv2 := uniqueStr("test_eth_priv_2")
	t.Cleanup(func() { cleanupETHAccountKeys(t, db, addr1, addr2) })

	repo := coldpostgres.NewETHAccountKeyRepositorySqlc(db)

	items := []*domainETH.ETHAccountKey{
		{
			Account:       domainAccount.AccountTypeClient,
			Address:       addr1,
			FullPublicKey: "04" + uniqueStr("eth_pubkey_1"),
			PrivateKey:    priv1,
			Idx:           1,
			AddrStatus:    domainAddress.AddrStatusHDKeyGenerated,
		},
		{
			Account:       domainAccount.AccountTypeClient,
			Address:       addr2,
			FullPublicKey: "04" + uniqueStr("eth_pubkey_2"),
			PrivateKey:    priv2,
			Idx:           2,
			AddrStatus:    domainAddress.AddrStatusHDKeyGenerated,
		},
	}

	err := repo.InsertBulk(items)
	require.NoError(t, err, "InsertBulk should succeed for valid ETH account keys")

	got, err := repo.GetAllAddrStatus(
		domainAccount.AccountTypeClient,
		domainAddress.AddrStatusHDKeyGenerated,
	)
	require.NoError(t, err)

	addrSet := make(map[string]bool, len(got))
	for _, k := range got {
		addrSet[k.Address] = true
	}
	assert.True(t, addrSet[addr1], "GetAllAddrStatus should return the first inserted ETH key")
	assert.True(t, addrSet[addr2], "GetAllAddrStatus should return the second inserted ETH key")
}

// TestETHAccountKeyRepository_GetByAddress verifies that GetByAddress retrieves
// the ETH key for a specific address and that the private key round-trips correctly.
func TestETHAccountKeyRepository_GetByAddress(t *testing.T) {
	db := openKeygenDB(t)

	addr := uniqueStr("0xtest_eth_getbyaddr")
	priv := uniqueStr("test_eth_priv_getbyaddr")
	t.Cleanup(func() { cleanupETHAccountKeys(t, db, addr) })

	repo := coldpostgres.NewETHAccountKeyRepositorySqlc(db)

	err := repo.InsertBulk([]*domainETH.ETHAccountKey{
		{
			Account:       domainAccount.AccountTypePayment,
			Address:       addr,
			FullPublicKey: "04" + uniqueStr("eth_pubkey_getbyaddr"),
			PrivateKey:    priv,
			Idx:           5,
			AddrStatus:    domainAddress.AddrStatusHDKeyGenerated,
		},
	})
	require.NoError(t, err)

	got, err := repo.GetByAddress(addr)
	require.NoError(t, err)
	require.NotNil(t, got)

	assert.Equal(t, addr, got.Address)
	assert.Equal(t, domainAccount.AccountTypePayment, got.Account)
	assert.Equal(t, priv, got.PrivateKey, "private key should round-trip correctly")
}

// TestETHAccountKeyRepository_UpdateAddrStatus verifies that UpdateAddrStatus
// correctly updates the status field for the given private keys.
func TestETHAccountKeyRepository_UpdateAddrStatus(t *testing.T) {
	db := openKeygenDB(t)

	addr := uniqueStr("0xtest_eth_update_status")
	priv := uniqueStr("test_eth_priv_update")
	t.Cleanup(func() { cleanupETHAccountKeys(t, db, addr) })

	repo := coldpostgres.NewETHAccountKeyRepositorySqlc(db)

	err := repo.InsertBulk([]*domainETH.ETHAccountKey{
		{
			Account:       domainAccount.AccountTypeClient,
			Address:       addr,
			FullPublicKey: "04" + uniqueStr("eth_pubkey_status"),
			PrivateKey:    priv,
			Idx:           1,
			AddrStatus:    domainAddress.AddrStatusHDKeyGenerated,
		},
	})
	require.NoError(t, err)

	rowsAffected, err := repo.UpdateAddrStatus(
		domainAccount.AccountTypeClient,
		domainAddress.AddrStatusPrivKeyImported,
		[]string{priv},
	)
	require.NoError(t, err)
	assert.EqualValues(t, 1, rowsAffected, "exactly one row should be updated")

	got, err := repo.GetByAddress(addr)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, domainAddress.AddrStatusPrivKeyImported, got.AddrStatus,
		"addr_status should be updated to privkey_imported")
}

// ── XRP Account Key ──────────────────────────────────────────────────────────

// TestXRPAccountKeyRepository_InsertBulk_and_GetAllAddrStatus verifies that
// InsertBulk stores XRP keys and GetAllAddrStatus retrieves them filtered by status.
func TestXRPAccountKeyRepository_InsertBulk_and_GetAllAddrStatus(t *testing.T) {
	db := openKeygenDB(t)

	accountID1 := uniqueStr("rTestXRPAccountID1")
	accountID2 := uniqueStr("rTestXRPAccountID2")
	t.Cleanup(func() { cleanupXRPAccountKeys(t, db, accountID1, accountID2) })

	repo := coldpostgres.NewXRPAccountKeyRepositorySqlc(db, domainCoin.XRP)

	items := []*domainXRP.XRPAccountKey{
		{
			CoinTypeCode:  domainCoin.XRP,
			Account:       domainAccount.AccountTypeClient,
			AccountID:     accountID1,
			KeyType:       domainXRP.XRPKeyTypeSecp256k1,
			MasterSeed:    uniqueStr("seed1"),
			MasterSeedHex: uniqueStr("seedhex1"),
			PublicKey:     uniqueStr("pubkey1"),
			PublicKeyHex:  uniqueStr("pubkeyhex1"),
			AllocatedID:   1,
			AddrStatus:    domainAddress.AddrStatusHDKeyGenerated,
		},
		{
			CoinTypeCode:  domainCoin.XRP,
			Account:       domainAccount.AccountTypeClient,
			AccountID:     accountID2,
			KeyType:       domainXRP.XRPKeyTypeSecp256k1,
			MasterSeed:    uniqueStr("seed2"),
			MasterSeedHex: uniqueStr("seedhex2"),
			PublicKey:     uniqueStr("pubkey2"),
			PublicKeyHex:  uniqueStr("pubkeyhex2"),
			AllocatedID:   2,
			AddrStatus:    domainAddress.AddrStatusHDKeyGenerated,
		},
	}

	ctx := context.Background()
	err := repo.InsertBulk(ctx, items)
	require.NoError(t, err, "InsertBulk should succeed for valid XRP account keys")

	got, err := repo.GetAllAddrStatus(
		ctx,
		domainAccount.AccountTypeClient,
		domainAddress.AddrStatusHDKeyGenerated,
	)
	require.NoError(t, err)

	idSet := make(map[string]bool, len(got))
	for _, k := range got {
		idSet[k.AccountID] = true
	}
	assert.True(t, idSet[accountID1], "GetAllAddrStatus should return the first inserted XRP key")
	assert.True(t, idSet[accountID2], "GetAllAddrStatus should return the second inserted XRP key")
}

// TestXRPAccountKeyRepository_GetSecret verifies that GetSecret retrieves the
// master seed stored for a given account ID.
func TestXRPAccountKeyRepository_GetSecret(t *testing.T) {
	db := openKeygenDB(t)

	accountID := uniqueStr("rTestXRPGetSecret")
	masterSeed := uniqueStr("test_master_seed_secret")
	t.Cleanup(func() { cleanupXRPAccountKeys(t, db, accountID) })

	repo := coldpostgres.NewXRPAccountKeyRepositorySqlc(db, domainCoin.XRP)

	ctx := context.Background()
	err := repo.InsertBulk(ctx, []*domainXRP.XRPAccountKey{
		{
			CoinTypeCode:  domainCoin.XRP,
			Account:       domainAccount.AccountTypeDeposit,
			AccountID:     accountID,
			KeyType:       domainXRP.XRPKeyTypeSecp256k1,
			MasterSeed:    masterSeed,
			MasterSeedHex: uniqueStr("seedhex_secret"),
			PublicKey:     uniqueStr("pubkey_secret"),
			PublicKeyHex:  uniqueStr("pubkeyhex_secret"),
			AllocatedID:   1,
			AddrStatus:    domainAddress.AddrStatusHDKeyGenerated,
		},
	})
	require.NoError(t, err)

	// GetSecret returns the master_seed for the given account_id.
	secret, err := repo.GetSecret(ctx, domainAccount.AccountTypeDeposit, accountID)
	require.NoError(t, err)
	assert.Equal(t, masterSeed, secret, "GetSecret should return the stored master seed")
}

// ── Seed Repository ──────────────────────────────────────────────────────────

// TestSeedRepository_Insert_and_GetOne verifies that Insert stores an encrypted seed
// and GetOne retrieves it by coin type.
func TestSeedRepository_Insert_and_GetOne(t *testing.T) {
	db := openKeygenDB(t)

	// Remove any pre-existing seed for BTC to avoid unique constraint violations.
	_, _ = db.ExecContext(context.Background(), "DELETE FROM seed WHERE coin = $1", "btc")
	t.Cleanup(func() { cleanupSeed(t, db, "btc") })

	repo := coldpostgres.NewSeedRepositorySqlc(db, domainCoin.BTC)

	testSeed := uniqueStr("encrypted_seed")
	ctx := context.Background()

	err := repo.Insert(ctx, testSeed)
	require.NoError(t, err, "Insert should succeed for a valid BTC seed")

	got, err := repo.GetOne(ctx)
	require.NoError(t, err)
	require.NotNil(t, got)

	assert.Equal(t, domainCoin.BTC, got.CoinTypeCode, "CoinTypeCode should round-trip correctly")
	assert.Equal(t, testSeed, got.Seed, "seed value should round-trip correctly")
}

// ── Nonce Repository ─────────────────────────────────────────────────────────

// testPublicNonce returns a 66-byte public nonce with non-zero bytes.
func testPublicNonce() [66]byte {
	var nonce [66]byte
	for i := range nonce {
		nonce[i] = byte(i + 1)
	}
	return nonce
}

// TestNonceRepository_SaveAndGet verifies that SaveNonce stores a nonce and
// GetNonce retrieves it for the correct signer and transaction.
func TestNonceRepository_SaveAndGet(t *testing.T) {
	db := openKeygenDB(t)

	txID := uniqueStr("test_tx_save_get")
	signerID := uniqueStr("signer_1")
	t.Cleanup(func() { cleanupNonces(t, db, txID) })

	repo := coldpostgres.NewNonceRepositorySqlc(db)
	ctx := context.Background()

	nonce := testPublicNonce()
	err := repo.SaveNonce(ctx, signerID, txID, nonce)
	require.NoError(t, err, "SaveNonce should succeed")

	got, err := repo.GetNonce(ctx, signerID, txID)
	require.NoError(t, err)
	require.NotNil(t, got)

	assert.Equal(t, signerID, got.SignerID())
	assert.Equal(t, nonce, got.PublicNonce(), "public nonce should round-trip correctly")
}

// TestNonceRepository_MarkNonceUsed verifies that MarkNonceUsed transitions a
// nonce from unused to used and it no longer appears in GetUnusedNoncesForTransaction.
func TestNonceRepository_MarkNonceUsed(t *testing.T) {
	db := openKeygenDB(t)

	txID := uniqueStr("test_tx_mark_used")
	signerID := uniqueStr("signer_mark")
	t.Cleanup(func() { cleanupNonces(t, db, txID) })

	repo := coldpostgres.NewNonceRepositorySqlc(db)
	ctx := context.Background()

	err := repo.SaveNonce(ctx, signerID, txID, testPublicNonce())
	require.NoError(t, err)

	// Nonce should appear in unused list before marking.
	unused, err := repo.GetUnusedNoncesForTransaction(ctx, txID)
	require.NoError(t, err)
	assert.Len(t, unused, 1, "one unused nonce should exist before marking")

	err = repo.MarkNonceUsed(ctx, signerID, txID)
	require.NoError(t, err, "MarkNonceUsed should succeed")

	// Nonce should no longer appear in unused list after marking.
	unused, err = repo.GetUnusedNoncesForTransaction(ctx, txID)
	require.NoError(t, err)
	assert.Empty(t, unused, "no unused nonces should remain after marking")

	// All nonces (used + unused) still includes the nonce.
	all, err := repo.GetAllNoncesForTransaction(ctx, txID)
	require.NoError(t, err)
	assert.Len(t, all, 1, "nonce should still exist in the all-nonces list")
}

// TestNonceRepository_DuplicateNonce verifies that saving a duplicate nonce
// (same signer and transaction) returns ErrNonceDuplicate.
func TestNonceRepository_DuplicateNonce(t *testing.T) {
	db := openKeygenDB(t)

	txID := uniqueStr("test_tx_duplicate")
	signerID := uniqueStr("signer_dup")
	t.Cleanup(func() { cleanupNonces(t, db, txID) })

	repo := coldpostgres.NewNonceRepositorySqlc(db)
	ctx := context.Background()

	err := repo.SaveNonce(ctx, signerID, txID, testPublicNonce())
	require.NoError(t, err, "first SaveNonce should succeed")

	err = repo.SaveNonce(ctx, signerID, txID, testPublicNonce())
	require.ErrorIs(t, err, multisig.ErrNonceDuplicate,
		"duplicate SaveNonce should return ErrNonceDuplicate")
}

// TestNonceRepository_DeleteNoncesForTransaction verifies that
// DeleteNoncesForTransaction removes all nonces associated with a transaction.
func TestNonceRepository_DeleteNoncesForTransaction(t *testing.T) {
	db := openKeygenDB(t)

	txID := uniqueStr("test_tx_delete")
	signer1 := uniqueStr("signer_del_1")
	signer2 := uniqueStr("signer_del_2")
	t.Cleanup(func() { cleanupNonces(t, db, txID) })

	repo := coldpostgres.NewNonceRepositorySqlc(db)
	ctx := context.Background()

	// Insert two nonces for the same transaction.
	nonce1 := testPublicNonce()
	nonce2 := testPublicNonce()
	nonce2[0] = 0xFF // make nonce2 distinct from nonce1

	require.NoError(t, repo.SaveNonce(ctx, signer1, txID, nonce1))
	require.NoError(t, repo.SaveNonce(ctx, signer2, txID, nonce2))

	all, err := repo.GetAllNoncesForTransaction(ctx, txID)
	require.NoError(t, err)
	assert.Len(t, all, 2, "two nonces should exist before deletion")

	err = repo.DeleteNoncesForTransaction(ctx, txID)
	require.NoError(t, err, "DeleteNoncesForTransaction should succeed")

	all, err = repo.GetAllNoncesForTransaction(ctx, txID)
	require.NoError(t, err)
	assert.Empty(t, all, "no nonces should remain after deletion")
}

// ── XRP Signer List ──────────────────────────────────────────────────────────

// TestXRPSignerListRepository_Insert_and_GetByAccountID verifies that Insert stores
// a signer list and GetByAccountID retrieves the active one.
func TestXRPSignerListRepository_Insert_and_GetByAccountID(t *testing.T) {
	db := openKeygenDB(t)

	accountID := uniqueStr("rTestSignerList")
	t.Cleanup(func() { cleanupSignerLists(t, db, accountID) })

	repo := coldpostgres.NewXRPSignerListRepositorySqlc(db)
	ctx := context.Background()

	signerList, err := domainXRP.NewXRPSignerList(accountID, 2)
	require.NoError(t, err)

	id, err := repo.Insert(ctx, signerList)
	require.NoError(t, err, "Insert should succeed for a valid signer list")
	assert.Positive(t, id, "Insert should return a positive auto-increment ID")

	got, err := repo.GetByAccountID(ctx, accountID)
	require.NoError(t, err)
	require.NotNil(t, got, "GetByAccountID should return the inserted signer list")

	assert.Equal(t, accountID, got.AccountID)
	assert.Equal(t, uint32(2), got.SignerQuorum)
	assert.True(t, got.IsActive, "newly inserted signer list should be active")
}

// TestXRPSignerListRepository_UpdateStatus verifies that UpdateStatus correctly
// toggles the is_active flag on a signer list.
func TestXRPSignerListRepository_UpdateStatus(t *testing.T) {
	db := openKeygenDB(t)

	accountID := uniqueStr("rTestSignerStatus")
	t.Cleanup(func() { cleanupSignerLists(t, db, accountID) })

	repo := coldpostgres.NewXRPSignerListRepositorySqlc(db)
	ctx := context.Background()

	signerList, err := domainXRP.NewXRPSignerList(accountID, 1)
	require.NoError(t, err)

	id, err := repo.Insert(ctx, signerList)
	require.NoError(t, err)

	err = repo.UpdateStatus(ctx, id, false)
	require.NoError(t, err, "UpdateStatus to inactive should succeed")

	// GetByAccountID returns nil when no active list exists.
	got, err := repo.GetByAccountID(ctx, accountID)
	require.NoError(t, err)
	assert.Nil(t, got, "GetByAccountID should return nil after deactivation")

	// Re-activate and confirm.
	err = repo.UpdateStatus(ctx, id, true)
	require.NoError(t, err, "re-activating the signer list should succeed")

	got, err = repo.GetByAccountID(ctx, accountID)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.True(t, got.IsActive, "re-activated signer list should be active")
}

// TestXRPSignerListRepository_GetHistoryByAccountID verifies that
// GetHistoryByAccountID returns all signer lists (active and inactive).
func TestXRPSignerListRepository_GetHistoryByAccountID(t *testing.T) {
	db := openKeygenDB(t)

	accountID := uniqueStr("rTestSignerHistory")
	t.Cleanup(func() { cleanupSignerLists(t, db, accountID) })

	repo := coldpostgres.NewXRPSignerListRepositorySqlc(db)
	ctx := context.Background()

	// Insert two signer lists for the same account.
	list1, err := domainXRP.NewXRPSignerList(accountID, 1)
	require.NoError(t, err)
	id1, err := repo.Insert(ctx, list1)
	require.NoError(t, err)

	list2, err := domainXRP.NewXRPSignerList(accountID, 2)
	require.NoError(t, err)
	_, err = repo.Insert(ctx, list2)
	require.NoError(t, err)

	// Deactivate the first one.
	err = repo.UpdateStatus(ctx, id1, false)
	require.NoError(t, err)

	// History should return both records.
	history, err := repo.GetHistoryByAccountID(ctx, accountID)
	require.NoError(t, err)
	assert.Len(t, history, 2, "GetHistoryByAccountID should return both signer lists")
}
