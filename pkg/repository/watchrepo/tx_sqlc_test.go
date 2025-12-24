//go:build integration
// +build integration

package watchrepo_test

import (
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"

	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	mysql "github.com/hiromaily/go-crypto-wallet/pkg/db/rdb"
	domainCoin "github.com/hiromaily/go-crypto-wallet/pkg/domain/coin"
	domainTx "github.com/hiromaily/go-crypto-wallet/pkg/domain/transaction"
	domainWallet "github.com/hiromaily/go-crypto-wallet/pkg/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	"github.com/hiromaily/go-crypto-wallet/pkg/repository/watchrepo"
)

// TestTxSqlc is integration test for TxRepositorySqlc
func TestTxSqlc(t *testing.T) {
	// Create ETH repository (tx table is for ETH/XRP only)
	projPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-crypto-wallet"
	confPath := projPath + "/data/config/eth_watch.toml"
	conf, err := config.NewWallet(confPath, domainWallet.WalletTypeWatchOnly, domainCoin.ETH)
	if err != nil {
		log.Fatalf("fail to create config: %v", err)
	}
	zapLog := logger.NewSlogFromConfig(conf.Logger.Env, conf.Logger.Level, conf.Logger.Service)
	db, err := mysql.NewMySQL(&conf.MySQL)
	if err != nil {
		log.Fatalf("fail to create db: %v", err)
	}
	txRepo := watchrepo.NewTxRepositorySqlc(db, domainCoin.ETH, zapLog)

	// Delete all records
	_, err := txRepo.DeleteAll()
	require.NoError(t, err, "fail to call DeleteAll()")

	// Insert unsigned tx
	actionType := domainTx.ActionTypePayment
	id, err := txRepo.InsertUnsignedTx(actionType)
	require.NoError(t, err, "fail to call InsertUnsignedTx()")
	require.NotZero(t, id, "InsertUnsignedTx() should return non-zero id")

	// Get one
	tx, err := txRepo.GetOne(id)
	require.NoError(t, err, "fail to call GetOne()")
	require.Equal(t, id, tx.ID, "GetOne() should return correct id")
	require.Equal(t, actionType.String(), tx.Action, "GetOne() should return correct action")

	// Get max ID
	maxID, err := txRepo.GetMaxID(actionType)
	require.NoError(t, err, "fail to call GetMaxID()")
	require.Equal(t, id, maxID, "GetMaxID() should return the inserted id")

	// Insert another tx to test max ID
	id2, err := txRepo.InsertUnsignedTx(actionType)
	require.NoError(t, err, "fail to call InsertUnsignedTx() second time")

	// Get max ID again
	maxID2, err := txRepo.GetMaxID(actionType)
	require.NoError(t, err, "fail to call GetMaxID() second time")
	require.Equal(t, id2, maxID2, "GetMaxID() should return the second inserted id")
	require.Greater(t, id2, id, "second InsertUnsignedTx() should return id greater than first")

	// Update tx
	tx.Action = domainTx.ActionTypeDeposit.String()
	rowsAffected, err := txRepo.Update(tx)
	require.NoError(t, err, "fail to call Update()")
	require.Equal(t, int64(1), rowsAffected, "Update() should affect 1 row")

	// Verify update
	updatedTx, err := txRepo.GetOne(id)
	require.NoError(t, err, "fail to call GetOne() after update")
	require.Equal(t, domainTx.ActionTypeDeposit.String(), updatedTx.Action, "Update() should change action to Deposit")
}
