//go:build integration
// +build integration

package watchrepo_test

import (
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
	mysql "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/mysql"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/watch"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
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
	zapLog := logger.NewSlogFromConfig(conf.Logger.Env, conf.Logger.Level, conf.Logger.Service)
	db, err := mysql.NewMySQL(&conf.MySQL)
	if err != nil {
		log.Fatalf("fail to create db: %v", err)
	}

	xrpDetailTxRepo := watch.NewXrpDetailTxInputRepositorySqlc(db, domainCoin.XRP, zapLog)
	txRepo := watch.NewTxRepositorySqlc(db, domainCoin.XRP, zapLog)

	// Clean up any existing test data
	_, _ = db.Exec("DELETE FROM xrp_detail_tx WHERE uuid LIKE 'xrp-uuid-%'")
	_, _ = db.Exec("DELETE FROM tx WHERE coin = 'xrp'")

	// Create a tx record first (xrp_detail_tx joins with tx table)
	txID, err := txRepo.InsertUnsignedTx(domainTx.ActionTypePayment)
	require.NoError(t, err, "fail to create parent tx")

	// Create test xrp detail tx
	uuid := "xrp-uuid-sqlc-test"
	xrpTx := &models.XRPDetailTX{
		TXID:                  txID,
		UUID:                  uuid,
		CurrentTXType:         domainTx.TxTypeUnsigned.Int8(),
		SenderAccount:         "deposit",
		SenderAddress:         "rSender-sqlc",
		ReceiverAccount:       "client",
		ReceiverAddress:       "rReceiver-sqlc",
		Amount:                "1000000",
		XRPTXType:             "Payment",
		Fee:                   "12",
		Flags:                 0,
		LastLedgerSequence:    12345,
		Sequence:              1,
		SigningPubkey:         "pubkey-sqlc",
		TXNSignature:          "",
		Hash:                  "",
		EarliestLedgerVersion: 0,
		SignedTXID:            "",
		TXBlob:                "",
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
	require.Equal(t, uuid, retrievedTx.UUID, "GetOne() should return correct UUID")

	// Update after tx sent
	signedTxID := "signed-txid-sqlc"
	txBlob := "tx-blob-sqlc"
	earliestLedgerVersion := uint64(12340)
	rowsAffected, err := xrpDetailTxRepo.UpdateAfterTxSent(uuid, domainTx.TxTypeSent, signedTxID, txBlob, earliestLedgerVersion)
	require.NoError(t, err, "fail to call UpdateAfterTxSent()")
	require.GreaterOrEqual(t, rowsAffected, int64(1), "UpdateAfterTxSent() should affect at least 1 row")

	// Verify update
	updatedTx, err := xrpDetailTxRepo.GetOne(retrievedTx.ID)
	require.NoError(t, err, "fail to call GetOne() after update")
	require.Equal(t, signedTxID, updatedTx.SignedTXID, "UpdateAfterTxSent() should update SignedTXID")
	require.Equal(t, txBlob, updatedTx.TXBlob, "UpdateAfterTxSent() should update TXBlob")
	require.Equal(t, domainTx.TxTypeSent.Int8(), updatedTx.CurrentTXType, "UpdateAfterTxSent() should update CurrentTXType")
	require.Equal(t, earliestLedgerVersion, updatedTx.EarliestLedgerVersion, "UpdateAfterTxSent() should update EarliestLedgerVersion")

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
	require.Equal(t, domainTx.TxTypeDone.Int8(), verifyTx.CurrentTXType, "UpdateTxTypeBySentHashTx() should update CurrentTXType to TxTypeDone")

	// Update tx type by ID
	rowsAffected, err = xrpDetailTxRepo.UpdateTxType(retrievedTx.ID, domainTx.TxTypeNotified)
	require.NoError(t, err, "fail to call UpdateTxType()")
	require.Equal(t, int64(1), rowsAffected, "UpdateTxType() should affect 1 row")

	// Verify final tx type
	finalTx, err := xrpDetailTxRepo.GetOne(retrievedTx.ID)
	require.NoError(t, err, "fail to call GetOne() after UpdateTxType()")
	require.Equal(t, domainTx.TxTypeNotified.Int8(), finalTx.CurrentTXType, "UpdateTxType() should update CurrentTXType to TxTypeNotified")

	// Test InsertBulk
	// Create another tx record for bulk insert
	txID2, err := txRepo.InsertUnsignedTx(domainTx.ActionTypePayment)
	require.NoError(t, err, "fail to create second parent tx")

	bulkTxs := []*models.XRPDetailTX{
		{
			TXID:                  txID2,
			UUID:                  "xrp-uuid-bulk-1",
			CurrentTXType:         domainTx.TxTypeUnsigned.Int8(),
			SenderAccount:         "deposit",
			SenderAddress:         "rSender-bulk-1",
			ReceiverAccount:       "client",
			ReceiverAddress:       "rReceiver-bulk-1",
			Amount:                "2000000",
			XRPTXType:             "Payment",
			Fee:                   "12",
			Flags:                 0,
			LastLedgerSequence:    12346,
			Sequence:              2,
			SigningPubkey:         "pubkey-bulk-1",
			TXNSignature:          "",
			Hash:                  "",
			EarliestLedgerVersion: 0,
			SignedTXID:            "",
			TXBlob:                "",
		},
		{
			TXID:                  txID2,
			UUID:                  "xrp-uuid-bulk-2",
			CurrentTXType:         domainTx.TxTypeUnsigned.Int8(),
			SenderAccount:         "deposit",
			SenderAddress:         "rSender-bulk-2",
			ReceiverAccount:       "client",
			ReceiverAddress:       "rReceiver-bulk-2",
			Amount:                "3000000",
			XRPTXType:             "Payment",
			Fee:                   "12",
			Flags:                 0,
			LastLedgerSequence:    12347,
			Sequence:              3,
			SigningPubkey:         "pubkey-bulk-2",
			TXNSignature:          "",
			Hash:                  "",
			EarliestLedgerVersion: 0,
			SignedTXID:            "",
			TXBlob:                "",
		},
	}

	err = xrpDetailTxRepo.InsertBulk(bulkTxs)
	require.NoError(t, err, "fail to call InsertBulk()")

	// Verify bulk insert
	bulkRetrieved, err := xrpDetailTxRepo.GetAllByTxID(txID2)
	require.NoError(t, err, "fail to call GetAllByTxID() after InsertBulk()")
	assert.Len(t, bulkRetrieved, 2, "InsertBulk() should insert 2 records")
}
