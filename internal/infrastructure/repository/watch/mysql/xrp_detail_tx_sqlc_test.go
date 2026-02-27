//go:build integration

package mysql_test

import (
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domainXRP "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/xrp"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
	watchmysql "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/watch/mysql"
	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	pkgmysql "github.com/hiromaily/go-crypto-wallet/pkg/db/mysql"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// TestXRPDetailTxSqlc is integration test for XRPDetailTxInputRepositorySqlc
func TestXRPDetailTxSqlc(t *testing.T) {
	// Create XRP repositories
	projPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-crypto-wallet"
	confPath := projPath + "/config/wallet/xrp/watch.yaml"
	conf, err := config.NewWallet(confPath, domainWallet.WalletTypeWatchOnly, domainCoin.XRP)
	if err != nil {
		log.Fatalf("fail to create config: %v", err)
	}
	_ = logger.NewSlogLogger(conf.Logger.Format, conf.Logger.Level, conf.Logger.Service, "")
	db, err := pkgmysql.NewMySQL(&conf.Database.MySQL)
	if err != nil {
		log.Fatalf("fail to create db: %v", err)
	}

	xrpDetailTxRepo := watchmysql.NewXRPDetailTxInputRepositorySqlc(db, domainCoin.XRP)
	txRepo := watchmysql.NewTxRepositorySqlc(db, domainCoin.XRP)

	// Clean up any existing test data
	_, _ = db.Exec("DELETE FROM xrp_detail_tx WHERE uuid LIKE 'xrp-uuid-%'")
	_, _ = db.Exec("DELETE FROM tx WHERE coin = 'xrp'")

	// Create a tx record first (xrp_detail_tx joins with tx table)
	txID, err := txRepo.InsertUnsignedTx(domainTx.ActionTypePayment)
	require.NoError(t, err, "fail to create parent tx")

	// Create test xrp detail tx
	uuid := "xrp-uuid-sqlc-test"
	xrpTx, err := domainXRP.NewXRPDetailTx(
		txID,
		uuid,
		domainTx.TxTypeUnsigned,
		"deposit",
		"rSender-sqlc",
		"client",
		"rReceiver-sqlc",
		"1000000",
		"Payment",
		"12",
		0,
		12345,
		1,
	)
	require.NoError(t, err, "fail to create XRPDetailTx")
	xrpTx.SigningPubkey = "pubkey-sqlc"

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
		domainTx.TxTypeSent,
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
		domainTx.TxTypeDone,
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
		domainTx.TxTypeNotified,
		finalTx.CurrentTxType,
		"UpdateTxType() should update CurrentTxType to TxTypeNotified",
	)

	// Test InsertBulk
	// Create another tx record for bulk insert
	txID2, err := txRepo.InsertUnsignedTx(domainTx.ActionTypePayment)
	require.NoError(t, err, "fail to create second parent tx")

	bulkTx1, err := domainXRP.NewXRPDetailTx(
		txID2,
		"xrp-uuid-bulk-1",
		domainTx.TxTypeUnsigned,
		"deposit",
		"rSender-bulk-1",
		"client",
		"rReceiver-bulk-1",
		"2000000",
		"Payment",
		"12",
		0,
		12346,
		2,
	)
	require.NoError(t, err, "fail to create bulk XRPDetailTx 1")
	bulkTx1.SigningPubkey = "pubkey-bulk-1"

	bulkTx2, err := domainXRP.NewXRPDetailTx(
		txID2,
		"xrp-uuid-bulk-2",
		domainTx.TxTypeUnsigned,
		"deposit",
		"rSender-bulk-2",
		"client",
		"rReceiver-bulk-2",
		"3000000",
		"Payment",
		"12",
		0,
		12347,
		3,
	)
	require.NoError(t, err, "fail to create bulk XRPDetailTx 2")
	bulkTx2.SigningPubkey = "pubkey-bulk-2"

	bulkTxs := []*domainXRP.XRPDetailTx{bulkTx1, bulkTx2}

	err = xrpDetailTxRepo.InsertBulk(bulkTxs)
	require.NoError(t, err, "fail to call InsertBulk()")

	// Verify bulk insert
	bulkRetrieved, err := xrpDetailTxRepo.GetAllByTxID(txID2)
	require.NoError(t, err, "fail to call GetAllByTxID() after InsertBulk()")
	assert.Len(t, bulkRetrieved, 2, "InsertBulk() should insert 2 records")
}
