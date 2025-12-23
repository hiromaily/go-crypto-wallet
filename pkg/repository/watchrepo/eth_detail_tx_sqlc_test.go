//go:build integration
// +build integration

package watchrepo_test

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "github.com/go-sql-driver/mysql"

	"github.com/hiromaily/go-crypto-wallet/pkg/action"
	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	mysql "github.com/hiromaily/go-crypto-wallet/pkg/db/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/repository/watchrepo"
	"github.com/hiromaily/go-crypto-wallet/pkg/tx"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

// TestEthDetailTxSqlc is integration test for EthDetailTxInputRepositorySqlc
func TestEthDetailTxSqlc(t *testing.T) {
	// Create ETH repositories
	projPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-crypto-wallet"
	confPath := projPath + "/data/config/eth_watch.toml"
	conf, err := config.NewWallet(confPath, wallet.WalletTypeWatchOnly, coin.ETH)
	if err != nil {
		log.Fatalf("fail to create config: %v", err)
	}
	zapLog := logger.NewSlogFromConfig(conf.Logger.Env, conf.Logger.Level, conf.Logger.Service)
	db, err := mysql.NewMySQL(&conf.MySQL)
	if err != nil {
		log.Fatalf("fail to create db: %v", err)
	}

	ethDetailTxRepo := watchrepo.NewEthDetailTxInputRepositorySqlc(db, coin.ETH, zapLog)
	txRepo := watchrepo.NewTxRepositorySqlc(db, coin.ETH, zapLog)

	// Clean up any existing test data
	_, _ = db.Exec("DELETE FROM eth_detail_tx WHERE uuid LIKE 'eth-uuid-%'")
	_, _ = db.Exec("DELETE FROM tx WHERE coin = 'eth'")

	// Create a tx record first (eth_detail_tx joins with tx table)
	txID, err := txRepo.InsertUnsignedTx(action.ActionTypePayment)
	require.NoError(t, err, "fail to create parent tx")

	// Create test eth detail tx
	uuid := "eth-uuid-sqlc-test"
	ethTx := &models.EthDetailTX{
		TXID:            txID,
		UUID:            uuid,
		CurrentTXType:   tx.TxTypeUnsigned.Int8(),
		SenderAccount:   "deposit",
		SenderAddress:   "0xsender-sqlc",
		ReceiverAccount: "client",
		ReceiverAddress: "0xreceiver-sqlc",
		Amount:          1000000000,
		Fee:             21000,
		GasLimit:        21000,
		Nonce:           1,
		UnsignedHexTX:   "0xunsigned-hex-sqlc",
	}

	// Insert
	err = ethDetailTxRepo.Insert(ethTx)
	require.NoError(t, err, "fail to call Insert()")

	// Get all by tx ID
	ethTxs, err := ethDetailTxRepo.GetAllByTxID(txID)
	require.NoError(t, err, "fail to call GetAllByTxID()")
	require.GreaterOrEqual(t, len(ethTxs), 1, "GetAllByTxID() should return at least 1 record")

	// Get one
	retrievedTx, err := ethDetailTxRepo.GetOne(ethTxs[0].ID)
	require.NoError(t, err, "fail to call GetOne()")
	require.Equal(t, uuid, retrievedTx.UUID, "GetOne() should return correct UUID")

	// Update after tx sent
	signedHex := "0xsigned-hex-sqlc"
	sentHashTx := "0xsent-hash-sqlc"
	rowsAffected, err := ethDetailTxRepo.UpdateAfterTxSent(uuid, tx.TxTypeSent, signedHex, sentHashTx)
	require.NoError(t, err, "fail to call UpdateAfterTxSent()")
	require.GreaterOrEqual(t, rowsAffected, int64(1), "UpdateAfterTxSent() should affect at least 1 row")

	// Verify update
	updatedTx, err := ethDetailTxRepo.GetOne(retrievedTx.ID)
	require.NoError(t, err, "fail to call GetOne() after update")
	require.Equal(t, signedHex, updatedTx.SignedHexTX, "UpdateAfterTxSent() should update SignedHexTX")
	require.Equal(t, sentHashTx, updatedTx.SentHashTX, "UpdateAfterTxSent() should update SentHashTX")
	require.Equal(t, tx.TxTypeSent.Int8(), updatedTx.CurrentTXType, "UpdateAfterTxSent() should update CurrentTXType")

	// Get sent hash tx
	hashes, err := ethDetailTxRepo.GetSentHashTx(tx.TxTypeSent)
	require.NoError(t, err, "fail to call GetSentHashTx()")
	require.GreaterOrEqual(t, len(hashes), 1, "GetSentHashTx() should return at least 1 hash")

	// Update tx type by sent hash
	rowsAffected, err = ethDetailTxRepo.UpdateTxTypeBySentHashTx(tx.TxTypeDone, sentHashTx)
	require.NoError(t, err, "fail to call UpdateTxTypeBySentHashTx()")
	require.GreaterOrEqual(t, rowsAffected, int64(1), "UpdateTxTypeBySentHashTx() should affect at least 1 row")

	// Verify tx type update
	verifyTx, err := ethDetailTxRepo.GetOne(retrievedTx.ID)
	require.NoError(t, err, "fail to call GetOne() after UpdateTxTypeBySentHashTx()")
	require.Equal(t, tx.TxTypeDone.Int8(), verifyTx.CurrentTXType, "UpdateTxTypeBySentHashTx() should update CurrentTXType to TxTypeDone")

	// Update tx type by ID
	rowsAffected, err = ethDetailTxRepo.UpdateTxType(retrievedTx.ID, tx.TxTypeNotified)
	require.NoError(t, err, "fail to call UpdateTxType()")
	require.Equal(t, int64(1), rowsAffected, "UpdateTxType() should affect 1 row")

	// Verify final tx type
	finalTx, err := ethDetailTxRepo.GetOne(retrievedTx.ID)
	require.NoError(t, err, "fail to call GetOne() after UpdateTxType()")
	require.Equal(t, tx.TxTypeNotified.Int8(), finalTx.CurrentTXType, "UpdateTxType() should update CurrentTXType to TxTypeNotified")

	// Test InsertBulk
	// Create another tx record for bulk insert
	txID2, err := txRepo.InsertUnsignedTx(action.ActionTypePayment)
	require.NoError(t, err, "fail to create second parent tx")

	bulkTxs := []*models.EthDetailTX{
		{
			TXID:            txID2,
			UUID:            "eth-uuid-bulk-1",
			CurrentTXType:   tx.TxTypeUnsigned.Int8(),
			SenderAccount:   "deposit",
			SenderAddress:   "0xsender-bulk-1",
			ReceiverAccount: "client",
			ReceiverAddress: "0xreceiver-bulk-1",
			Amount:          2000000000,
			Fee:             21000,
			GasLimit:        21000,
			Nonce:           2,
			UnsignedHexTX:   "0xunsigned-bulk-1",
		},
		{
			TXID:            txID2,
			UUID:            "eth-uuid-bulk-2",
			CurrentTXType:   tx.TxTypeUnsigned.Int8(),
			SenderAccount:   "deposit",
			SenderAddress:   "0xsender-bulk-2",
			ReceiverAccount: "client",
			ReceiverAddress: "0xreceiver-bulk-2",
			Amount:          3000000000,
			Fee:             21000,
			GasLimit:        21000,
			Nonce:           3,
			UnsignedHexTX:   "0xunsigned-bulk-2",
		},
	}

	err = ethDetailTxRepo.InsertBulk(bulkTxs)
	require.NoError(t, err, "fail to call InsertBulk()")

	// Verify bulk insert
	bulkRetrieved, err := ethDetailTxRepo.GetAllByTxID(txID2)
	require.NoError(t, err, "fail to call GetAllByTxID() after InsertBulk()")
	assert.Len(t, bulkRetrieved, 2, "InsertBulk() should insert 2 records")
}
