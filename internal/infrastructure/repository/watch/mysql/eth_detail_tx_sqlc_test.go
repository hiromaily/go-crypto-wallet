//go:build integration
// +build integration

package mysql_test

import (
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainEth "github.com/hiromaily/go-crypto-wallet/internal/domain/ethereum"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
	watchmysql "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/watch/mysql"
	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	pkgmysql "github.com/hiromaily/go-crypto-wallet/pkg/db/mysql"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// TestETHDetailTXSqlc is integration test for ETHDetailTXInputRepositorySqlc
func TestETHDetailTXSqlc(t *testing.T) {
	// Create ETH repositories
	projPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-crypto-wallet"
	confPath := projPath + "/config/wallet/eth/watch.yaml"
	conf, err := config.NewWallet(confPath, domainWallet.WalletTypeWatchOnly, domainCoin.ETH)
	if err != nil {
		log.Fatalf("fail to create config: %v", err)
	}
	_ = logger.NewSlogLogger(conf.Logger.Format, conf.Logger.Level, conf.Logger.Service, "")
	db, err := pkgmysql.NewMySQL(&conf.Database.MySQL)
	if err != nil {
		log.Fatalf("fail to create db: %v", err)
	}

	ethDetailTXRepo := watchmysql.NewETHDetailTXInputRepositorySqlc(db, domainCoin.ETH)
	txRepo := watchmysql.NewTxRepositorySqlc(db, domainCoin.ETH)

	// Clean up any existing test data
	_, _ = db.Exec("DELETE FROM eth_detail_tx WHERE uuid LIKE 'eth-uuid-%'")
	_, _ = db.Exec("DELETE FROM tx WHERE coin = 'eth'")

	// Create a tx record first (eth_detail_tx joins with tx table)
	txID, err := txRepo.InsertUnsignedTx(domainTx.ActionTypePayment)
	require.NoError(t, err, "fail to create parent tx")

	// Create test eth detail tx
	uuid := "eth-uuid-sqlc-test"
	ethTx, err := domainEth.NewETHDetailTx(
		txID,
		uuid,
		domainTx.TxTypeUnsigned,
		"deposit",
		"0xsender-sqlc",
		"client",
		"0xreceiver-sqlc",
		1000000000,
		21000,
		21000,
		1,
	)
	require.NoError(t, err, "fail to create ETHDetailTx")
	ethTx.UnsignedHexTx = "0xunsigned-hex-sqlc"

	// Insert
	err = ethDetailTXRepo.Insert(ethTx)
	require.NoError(t, err, "fail to call Insert()")

	// Get all by tx ID
	ethTxs, err := ethDetailTXRepo.GetAllByTxID(txID)
	require.NoError(t, err, "fail to call GetAllByTxID()")
	require.GreaterOrEqual(t, len(ethTxs), 1, "GetAllByTxID() should return at least 1 record")

	// Get one
	retrievedTx, err := ethDetailTXRepo.GetOne(ethTxs[0].ID)
	require.NoError(t, err, "fail to call GetOne()")
	require.Equal(t, uuid, retrievedTx.UUID, "GetOne() should return correct UUID")

	// Update after tx sent
	signedHex := "0xsigned-hex-sqlc"
	sentHashTx := "0xsent-hash-sqlc"
	rowsAffected, err := ethDetailTXRepo.UpdateAfterTxSent(uuid, domainTx.TxTypeSent, signedHex, sentHashTx)
	require.NoError(t, err, "fail to call UpdateAfterTxSent()")
	require.GreaterOrEqual(t, rowsAffected, int64(1), "UpdateAfterTxSent() should affect at least 1 row")

	// Verify update
	updatedTx, err := ethDetailTXRepo.GetOne(retrievedTx.ID)
	require.NoError(t, err, "fail to call GetOne() after update")
	require.Equal(t, signedHex, updatedTx.SignedHexTx, "UpdateAfterTxSent() should update SignedHexTx")
	require.Equal(t, sentHashTx, updatedTx.SentHashTx, "UpdateAfterTxSent() should update SentHashTx")
	require.Equal(
		t, domainTx.TxTypeSent, updatedTx.CurrentTxType, "UpdateAfterTxSent() should update CurrentTxType",
	)

	// Get sent hash tx
	hashes, err := ethDetailTXRepo.GetSentHashTx(domainTx.TxTypeSent)
	require.NoError(t, err, "fail to call GetSentHashTx()")
	require.GreaterOrEqual(t, len(hashes), 1, "GetSentHashTx() should return at least 1 hash")

	// Update tx type by sent hash
	rowsAffected, err = ethDetailTXRepo.UpdateTxTypeBySentHashTx(domainTx.TxTypeDone, sentHashTx)
	require.NoError(t, err, "fail to call UpdateTxTypeBySentHashTx()")
	require.GreaterOrEqual(t, rowsAffected, int64(1), "UpdateTxTypeBySentHashTx() should affect at least 1 row")

	// Verify tx type update
	verifyTx, err := ethDetailTXRepo.GetOne(retrievedTx.ID)
	require.NoError(t, err, "fail to call GetOne() after UpdateTxTypeBySentHashTx()")
	require.Equal(
		t,
		domainTx.TxTypeDone,
		verifyTx.CurrentTxType,
		"UpdateTxTypeBySentHashTx() should update CurrentTxType to TxTypeDone",
	)

	// Update tx type by ID
	rowsAffected, err = ethDetailTXRepo.UpdateTxType(retrievedTx.ID, domainTx.TxTypeNotified)
	require.NoError(t, err, "fail to call UpdateTxType()")
	require.Equal(t, int64(1), rowsAffected, "UpdateTxType() should affect 1 row")

	// Verify final tx type
	finalTx, err := ethDetailTXRepo.GetOne(retrievedTx.ID)
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

	bulkTx1, err := domainEth.NewETHDetailTx(
		txID2,
		"eth-uuid-bulk-1",
		domainTx.TxTypeUnsigned,
		"deposit",
		"0xsender-bulk-1",
		"client",
		"0xreceiver-bulk-1",
		2000000000,
		21000,
		21000,
		2,
	)
	require.NoError(t, err, "fail to create bulk ETHDetailTx 1")
	bulkTx1.UnsignedHexTx = "0xunsigned-bulk-1"

	bulkTx2, err := domainEth.NewETHDetailTx(
		txID2,
		"eth-uuid-bulk-2",
		domainTx.TxTypeUnsigned,
		"deposit",
		"0xsender-bulk-2",
		"client",
		"0xreceiver-bulk-2",
		3000000000,
		21000,
		21000,
		3,
	)
	require.NoError(t, err, "fail to create bulk ETHDetailTx 2")
	bulkTx2.UnsignedHexTx = "0xunsigned-bulk-2"

	bulkTxs := []*domainEth.ETHDetailTx{bulkTx1, bulkTx2}

	err = ethDetailTXRepo.InsertBulk(bulkTxs)
	require.NoError(t, err, "fail to call InsertBulk()")

	// Verify bulk insert
	bulkRetrieved, err := ethDetailTXRepo.GetAllByTxID(txID2)
	require.NoError(t, err, "fail to call GetAllByTxID() after InsertBulk()")
	assert.Len(t, bulkRetrieved, 2, "InsertBulk() should insert 2 records")
}
