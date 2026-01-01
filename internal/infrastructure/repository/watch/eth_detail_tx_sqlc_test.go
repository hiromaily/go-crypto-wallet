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

// TestETHDetailTXSqlc is integration test for ETHDetailTXInputRepositorySqlc
func TestETHDetailTXSqlc(t *testing.T) {
	// Create ETH repositories
	projPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-crypto-wallet"
	confPath := projPath + "/data/config/eth_watch.toml"
	conf, err := config.NewWallet(confPath, domainWallet.WalletTypeWatchOnly, domainCoin.ETH)
	if err != nil {
		log.Fatalf("fail to create config: %v", err)
	}
	_ = logger.NewSlogFromConfig(conf.Logger.Env, conf.Logger.Level, conf.Logger.Service)
	db, err := mysql.NewMySQL(&conf.MySQL)
	if err != nil {
		log.Fatalf("fail to create db: %v", err)
	}

	ETHDetailTXRepo := watch.NewETHDetailTXInputRepositorySqlc(db, domainCoin.ETH)
	txRepo := watch.NewTxRepositorySqlc(db, domainCoin.ETH)

	// Clean up any existing test data
	_, _ = db.Exec("DELETE FROM eth_detail_tx WHERE uuid LIKE 'eth-uuid-%'")
	_, _ = db.Exec("DELETE FROM tx WHERE coin = 'eth'")

	// Create a tx record first (eth_detail_tx joins with tx table)
	txID, err := txRepo.InsertUnsignedTx(domainTx.ActionTypePayment)
	require.NoError(t, err, "fail to create parent tx")

	// Create test eth detail tx
	uuid := "eth-uuid-sqlc-test"
	ethTx := &sqlc.ETHDetailTX{
		TxID:            txID,
		Uuid:            uuid,
		CurrentTxType:   domainTx.TxTypeUnsigned.Int8(),
		SenderAccount:   "deposit",
		SenderAddress:   "0xsender-sqlc",
		ReceiverAccount: "client",
		ReceiverAddress: "0xreceiver-sqlc",
		Amount:          1000000000,
		Fee:             21000,
		GasLimit:        21000,
		Nonce:           1,
		UnsignedHexTx:   "0xunsigned-hex-sqlc",
	}

	// Insert
	err = ETHDetailTXRepo.Insert(ethTx)
	require.NoError(t, err, "fail to call Insert()")

	// Get all by tx ID
	ethTxs, err := ETHDetailTXRepo.GetAllByTxID(txID)
	require.NoError(t, err, "fail to call GetAllByTxID()")
	require.GreaterOrEqual(t, len(ethTxs), 1, "GetAllByTxID() should return at least 1 record")

	// Get one
	retrievedTx, err := ETHDetailTXRepo.GetOne(ethTxs[0].ID)
	require.NoError(t, err, "fail to call GetOne()")
	require.Equal(t, uuid, retrievedTx.Uuid, "GetOne() should return correct Uuid")

	// Update after tx sent
	signedHex := "0xsigned-hex-sqlc"
	sentHashTx := "0xsent-hash-sqlc"
	rowsAffected, err := ETHDetailTXRepo.UpdateAfterTxSent(uuid, domainTx.TxTypeSent, signedHex, sentHashTx)
	require.NoError(t, err, "fail to call UpdateAfterTxSent()")
	require.GreaterOrEqual(t, rowsAffected, int64(1), "UpdateAfterTxSent() should affect at least 1 row")

	// Verify update
	updatedTx, err := ETHDetailTXRepo.GetOne(retrievedTx.ID)
	require.NoError(t, err, "fail to call GetOne() after update")
	require.Equal(t, signedHex, updatedTx.SignedHexTx, "UpdateAfterTxSent() should update SignedHexTx")
	require.Equal(t, sentHashTx, updatedTx.SentHashTx, "UpdateAfterTxSent() should update SentHashTx")
	require.Equal(
		t, domainTx.TxTypeSent.Int8(), updatedTx.CurrentTxType, "UpdateAfterTxSent() should update CurrentTxType",
	)

	// Get sent hash tx
	hashes, err := ETHDetailTXRepo.GetSentHashTx(domainTx.TxTypeSent)
	require.NoError(t, err, "fail to call GetSentHashTx()")
	require.GreaterOrEqual(t, len(hashes), 1, "GetSentHashTx() should return at least 1 hash")

	// Update tx type by sent hash
	rowsAffected, err = ETHDetailTXRepo.UpdateTxTypeBySentHashTx(domainTx.TxTypeDone, sentHashTx)
	require.NoError(t, err, "fail to call UpdateTxTypeBySentHashTx()")
	require.GreaterOrEqual(t, rowsAffected, int64(1), "UpdateTxTypeBySentHashTx() should affect at least 1 row")

	// Verify tx type update
	verifyTx, err := ETHDetailTXRepo.GetOne(retrievedTx.ID)
	require.NoError(t, err, "fail to call GetOne() after UpdateTxTypeBySentHashTx()")
	require.Equal(
		t,
		domainTx.TxTypeDone.Int8(),
		verifyTx.CurrentTxType,
		"UpdateTxTypeBySentHashTx() should update CurrentTxType to TxTypeDone",
	)

	// Update tx type by ID
	rowsAffected, err = ETHDetailTXRepo.UpdateTxType(retrievedTx.ID, domainTx.TxTypeNotified)
	require.NoError(t, err, "fail to call UpdateTxType()")
	require.Equal(t, int64(1), rowsAffected, "UpdateTxType() should affect 1 row")

	// Verify final tx type
	finalTx, err := ETHDetailTXRepo.GetOne(retrievedTx.ID)
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

	bulkTxs := []*sqlc.ETHDetailTX{
		{
			TxID:            txID2,
			Uuid:            "eth-uuid-bulk-1",
			CurrentTxType:   domainTx.TxTypeUnsigned.Int8(),
			SenderAccount:   "deposit",
			SenderAddress:   "0xsender-bulk-1",
			ReceiverAccount: "client",
			ReceiverAddress: "0xreceiver-bulk-1",
			Amount:          2000000000,
			Fee:             21000,
			GasLimit:        21000,
			Nonce:           2,
			UnsignedHexTx:   "0xunsigned-bulk-1",
		},
		{
			TxID:            txID2,
			Uuid:            "eth-uuid-bulk-2",
			CurrentTxType:   domainTx.TxTypeUnsigned.Int8(),
			SenderAccount:   "deposit",
			SenderAddress:   "0xsender-bulk-2",
			ReceiverAccount: "client",
			ReceiverAddress: "0xreceiver-bulk-2",
			Amount:          3000000000,
			Fee:             21000,
			GasLimit:        21000,
			Nonce:           3,
			UnsignedHexTx:   "0xunsigned-bulk-2",
		},
	}

	err = ETHDetailTXRepo.InsertBulk(bulkTxs)
	require.NoError(t, err, "fail to call InsertBulk()")

	// Verify bulk insert
	bulkRetrieved, err := ETHDetailTXRepo.GetAllByTxID(txID2)
	require.NoError(t, err, "fail to call GetAllByTxID() after InsertBulk()")
	assert.Len(t, bulkRetrieved, 2, "InsertBulk() should insert 2 records")
}
