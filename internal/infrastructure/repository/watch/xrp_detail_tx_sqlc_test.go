//go:build integration
// +build integration

package watch_test

import (
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/sqlc"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/watch"
	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	mysql "github.com/hiromaily/go-crypto-wallet/pkg/db/mysql"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// TestXrpDetailTxSqlc is integration test for XrpDetailTxInputRepositorySqlc
func TestXrpDetailTxSqlc(t *testing.T) {
	// Create XRP repositories
	projPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-crypto-wallet"
	confPath := projPath + "/data/config/xrp_watch.toml"
	conf, err := config.NewWallet(confPath, domainWallet.WalletTypeWatchOnly, domainCoin.XRP)
	if err != nil {
		log.Fatalf("fail to create config: %v", err)
	}
	_ = logger.NewSlogFromConfig(conf.Logger.Env, conf.Logger.Level, conf.Logger.Service)
	db, err := mysql.NewMySQL(&conf.MySQL)
	if err != nil {
		log.Fatalf("fail to create db: %v", err)
	}

	xrpDetailTxRepo := watch.NewXrpDetailTxInputRepositorySqlc(db, domainCoin.XRP)
	txRepo := watch.NewTxRepositorySqlc(db, domainCoin.XRP)

	// Clean up any existing test data
	_, _ = db.Exec("DELETE FROM xrp_detail_tx WHERE uuid LIKE 'xrp-uuid-%'")
	_, _ = db.Exec("DELETE FROM tx WHERE coin = 'xrp'")

	// Create a tx record first (xrp_detail_tx joins with tx table)
	txID, err := txRepo.InsertUnsignedTx(domainTx.ActionTypePayment)
	require.NoError(t, err, "fail to create parent tx")

	// Create test xrp detail tx
	uuid := "xrp-uuid-sqlc-test"
	xrpTx := &sqlc.XrpDetailTx{
		TxID:                  txID,
		Uuid:                  uuid,
		CurrentTxType:         domainTx.TxTypeUnsigned.Int8(),
		SenderAccount:         "deposit",
		SenderAddress:         "rSender-sqlc",
		ReceiverAccount:       "client",
		ReceiverAddress:       "rReceiver-sqlc",
		Amount:                "1000000",
		XrpTxType:             "Payment",
		Fee:                   "12",
		Flags:                 0,
		LastLedgerSequence:    12345,
		Sequence:              1,
		SigningPubkey:         "pubkey-sqlc",
		TxnSignature:          "",
		Hash:                  "",
		EarliestLedgerVersion: 0,
		SignedTxID:            "",
		TxBlob:                "",
	}

	// Insert
	err = xrpDetailTxRepo.Insert(xrpTx)
	require.NoError(t, err, "fail to call Insert()")

	// Get all by tx ID
	xrpTxs, err := xrpDetailTxRepo.GetAllByTxID(txID)
	require.NoError(t, err, "fail to call GetAllByTxID()")
	require.GreaterOrEqual(t, len(xrpTxs), 1, "GetAllByTxID() should return at least 1 record")

	// Get one
	retrievedTx, err := xrpDetailTxRepo.GetOne(xrpTxs[0].ID)
	require.NoError(t, err, "fail to call GetOne()")
	require.Equal(t, uuid, retrievedTx.Uuid, "GetOne() should return correct Uuid")

	// Update after tx sent
	signedTxID := "signed-txid-sqlc"
	txBlob := "tx-blob-sqlc"
	earliestLedgerVersion := uint64(12340)
	rowsAffected, err := xrpDetailTxRepo.UpdateAfterTxSent(
		uuid, domainTx.TxTypeSent, signedTxID, txBlob, earliestLedgerVersion,
	)
	require.NoError(t, err, "fail to call UpdateAfterTxSent()")
	require.GreaterOrEqual(t, rowsAffected, int64(1), "UpdateAfterTxSent() should affect at least 1 row")

	// Verify update
	updatedTx, err := xrpDetailTxRepo.GetOne(retrievedTx.ID)
	require.NoError(t, err, "fail to call GetOne() after update")
	require.Equal(t, signedTxID, updatedTx.SignedTxID, "UpdateAfterTxSent() should update SignedTxID")
	require.Equal(t, txBlob, updatedTx.TxBlob, "UpdateAfterTxSent() should update TxBlob")
	require.Equal(
		t,
		domainTx.TxTypeSent.Int8(),
		updatedTx.CurrentTxType,
		"UpdateAfterTxSent() should update CurrentTxType",
	)
	require.Equal(
		t,
		earliestLedgerVersion,
		updatedTx.EarliestLedgerVersion,
		"UpdateAfterTxSent() should update EarliestLedgerVersion",
	)

	// Get sent hash tx (for XRP, this is tx_blob)
	blobs, err := xrpDetailTxRepo.GetSentHashTx(domainTx.TxTypeSent)
	require.NoError(t, err, "fail to call GetSentHashTx()")
	require.GreaterOrEqual(t, len(blobs), 1, "GetSentHashTx() should return at least 1 blob")

	// Update tx type by sent hash tx (tx_blob)
	rowsAffected, err = xrpDetailTxRepo.UpdateTxTypeBySentHashTx(domainTx.TxTypeDone, txBlob)
	require.NoError(t, err, "fail to call UpdateTxTypeBySentHashTx()")
	require.GreaterOrEqual(t, rowsAffected, int64(1), "UpdateTxTypeBySentHashTx() should affect at least 1 row")

	// Verify tx type update
	verifyTx, err := xrpDetailTxRepo.GetOne(retrievedTx.ID)
	require.NoError(t, err, "fail to call GetOne() after UpdateTxTypeBySentHashTx()")
	require.Equal(
		t,
		domainTx.TxTypeDone.Int8(),
		verifyTx.CurrentTxType,
		"UpdateTxTypeBySentHashTx() should update CurrentTxType to TxTypeDone",
	)

	// Update tx type by ID
	rowsAffected, err = xrpDetailTxRepo.UpdateTxType(retrievedTx.ID, domainTx.TxTypeNotified)
	require.NoError(t, err, "fail to call UpdateTxType()")
	require.Equal(t, int64(1), rowsAffected, "UpdateTxType() should affect 1 row")

	// Verify final tx type
	finalTx, err := xrpDetailTxRepo.GetOne(retrievedTx.ID)
	require.NoError(t, err, "fail to call GetOne() after UpdateTxType()")
	require.Equal(
		t,
		domainTx.TxTypeNotified.Int8(),
		finalTx.CurrentTxType,
		"UpdateTxType() should update CurrentTxType to TxTypeNotified",
	)

	// Test InsertBulk
	// Create another tx record for bulk insert
	txID2, err := txRepo.InsertUnsignedTx(domainTx.ActionTypePayment)
	require.NoError(t, err, "fail to create second parent tx")

	bulkTxs := []*sqlc.XrpDetailTx{
		{
			TxID:                  txID2,
			Uuid:                  "xrp-uuid-bulk-1",
			CurrentTxType:         domainTx.TxTypeUnsigned.Int8(),
			SenderAccount:         "deposit",
			SenderAddress:         "rSender-bulk-1",
			ReceiverAccount:       "client",
			ReceiverAddress:       "rReceiver-bulk-1",
			Amount:                "2000000",
			XrpTxType:             "Payment",
			Fee:                   "12",
			Flags:                 0,
			LastLedgerSequence:    12346,
			Sequence:              2,
			SigningPubkey:         "pubkey-bulk-1",
			TxnSignature:          "",
			Hash:                  "",
			EarliestLedgerVersion: 0,
			SignedTxID:            "",
			TxBlob:                "",
		},
		{
			TxID:                  txID2,
			Uuid:                  "xrp-uuid-bulk-2",
			CurrentTxType:         domainTx.TxTypeUnsigned.Int8(),
			SenderAccount:         "deposit",
			SenderAddress:         "rSender-bulk-2",
			ReceiverAccount:       "client",
			ReceiverAddress:       "rReceiver-bulk-2",
			Amount:                "3000000",
			XrpTxType:             "Payment",
			Fee:                   "12",
			Flags:                 0,
			LastLedgerSequence:    12347,
			Sequence:              3,
			SigningPubkey:         "pubkey-bulk-2",
			TxnSignature:          "",
			Hash:                  "",
			EarliestLedgerVersion: 0,
			SignedTxID:            "",
			TxBlob:                "",
		},
	}

	err = xrpDetailTxRepo.InsertBulk(bulkTxs)
	require.NoError(t, err, "fail to call InsertBulk()")

	// Verify bulk insert
	bulkRetrieved, err := xrpDetailTxRepo.GetAllByTxID(txID2)
	require.NoError(t, err, "fail to call GetAllByTxID() after InsertBulk()")
	assert.Len(t, bulkRetrieved, 2, "InsertBulk() should insert 2 records")
}
