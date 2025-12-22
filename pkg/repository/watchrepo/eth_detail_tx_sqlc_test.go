//go:build integration
// +build integration

package watchrepo_test

import (
	"log"
	"os"
	"testing"

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
	if err != nil {
		t.Fatalf("fail to create parent tx: %v", err)
	}

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
	if err := ethDetailTxRepo.Insert(ethTx); err != nil {
		t.Fatalf("fail to call Insert() %v", err)
	}

	// Get all by tx ID
	ethTxs, err := ethDetailTxRepo.GetAllByTxID(txID)
	if err != nil {
		t.Fatalf("fail to call GetAllByTxID() %v", err)
	}
	if len(ethTxs) < 1 {
		t.Errorf("GetAllByTxID() returned %d records, want at least 1", len(ethTxs))
		return
	}

	// Get one
	retrievedTx, err := ethDetailTxRepo.GetOne(ethTxs[0].ID)
	if err != nil {
		t.Fatalf("fail to call GetOne() %v", err)
	}
	if retrievedTx.UUID != uuid {
		t.Errorf("GetOne() returned UUID = %s, want %s", retrievedTx.UUID, uuid)
		return
	}

	// Update after tx sent
	signedHex := "0xsigned-hex-sqlc"
	sentHashTx := "0xsent-hash-sqlc"
	rowsAffected, err := ethDetailTxRepo.UpdateAfterTxSent(uuid, tx.TxTypeSent, signedHex, sentHashTx)
	if err != nil {
		t.Fatalf("fail to call UpdateAfterTxSent() %v", err)
	}
	if rowsAffected < 1 {
		t.Errorf("UpdateAfterTxSent() affected %d rows, want at least 1", rowsAffected)
		return
	}

	// Verify update
	updatedTx, err := ethDetailTxRepo.GetOne(retrievedTx.ID)
	if err != nil {
		t.Fatalf("fail to call GetOne() after update %v", err)
	}
	if updatedTx.SignedHexTX != signedHex {
		t.Errorf("UpdateAfterTxSent() did not update SignedHexTX, got %s, want %s", updatedTx.SignedHexTX, signedHex)
		return
	}
	if updatedTx.SentHashTX != sentHashTx {
		t.Errorf("UpdateAfterTxSent() did not update SentHashTX, got %s, want %s", updatedTx.SentHashTX, sentHashTx)
		return
	}
	if updatedTx.CurrentTXType != tx.TxTypeSent.Int8() {
		t.Errorf("UpdateAfterTxSent() did not update CurrentTXType, got %d, want %d", updatedTx.CurrentTXType, tx.TxTypeSent.Int8())
		return
	}

	// Get sent hash tx
	hashes, err := ethDetailTxRepo.GetSentHashTx(tx.TxTypeSent)
	if err != nil {
		t.Fatalf("fail to call GetSentHashTx() %v", err)
	}
	if len(hashes) < 1 {
		t.Errorf("GetSentHashTx() returned %d hashes, want at least 1", len(hashes))
		return
	}

	// Update tx type by sent hash
	rowsAffected, err = ethDetailTxRepo.UpdateTxTypeBySentHashTx(tx.TxTypeDone, sentHashTx)
	if err != nil {
		t.Fatalf("fail to call UpdateTxTypeBySentHashTx() %v", err)
	}
	if rowsAffected < 1 {
		t.Errorf("UpdateTxTypeBySentHashTx() affected %d rows, want at least 1", rowsAffected)
		return
	}

	// Verify tx type update
	verifyTx, err := ethDetailTxRepo.GetOne(retrievedTx.ID)
	if err != nil {
		t.Fatalf("fail to call GetOne() after UpdateTxTypeBySentHashTx() %v", err)
	}
	if verifyTx.CurrentTXType != tx.TxTypeDone.Int8() {
		t.Errorf("UpdateTxTypeBySentHashTx() did not update CurrentTXType, got %d, want %d", verifyTx.CurrentTXType, tx.TxTypeDone.Int8())
		return
	}

	// Update tx type by ID
	rowsAffected, err = ethDetailTxRepo.UpdateTxType(retrievedTx.ID, tx.TxTypeNotified)
	if err != nil {
		t.Fatalf("fail to call UpdateTxType() %v", err)
	}
	if rowsAffected != 1 {
		t.Errorf("UpdateTxType() affected %d rows, want 1", rowsAffected)
		return
	}

	// Verify final tx type
	finalTx, err := ethDetailTxRepo.GetOne(retrievedTx.ID)
	if err != nil {
		t.Fatalf("fail to call GetOne() after UpdateTxType() %v", err)
	}
	if finalTx.CurrentTXType != tx.TxTypeNotified.Int8() {
		t.Errorf("UpdateTxType() did not update CurrentTXType, got %d, want %d", finalTx.CurrentTXType, tx.TxTypeNotified.Int8())
		return
	}

	// Test InsertBulk
	// Create another tx record for bulk insert
	txID2, err := txRepo.InsertUnsignedTx(action.ActionTypePayment)
	if err != nil {
		t.Fatalf("fail to create second parent tx: %v", err)
	}

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

	if err := ethDetailTxRepo.InsertBulk(bulkTxs); err != nil {
		t.Fatalf("fail to call InsertBulk() %v", err)
	}

	// Verify bulk insert
	bulkRetrieved, err := ethDetailTxRepo.GetAllByTxID(txID2)
	if err != nil {
		t.Fatalf("fail to call GetAllByTxID() after InsertBulk() %v", err)
	}
	if len(bulkRetrieved) != 2 {
		t.Errorf("InsertBulk() inserted %d records, want 2", len(bulkRetrieved))
		return
	}
}
