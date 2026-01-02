//go:build integration
// +build integration

package watch_test

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"

	domainBitcoin "github.com/hiromaily/go-crypto-wallet/internal/domain/bitcoin"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	watchTestutil "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/watch/testutil"
)

// TestBTCTxInputSqlc is integration test for TxInputRepositorySqlc
func TestBTCTxInputSqlc(t *testing.T) {
	// Need to create a btc_tx first
	btcTxRepo := watchTestutil.NewBTCTxRepositorySqlc()
	btcTxInputRepo := watchTestutil.NewBTCTxInputRepositorySqlc()

	// Create a parent tx
	txItem := domainBitcoin.NewBtcTransaction(
		domainCoin.BTC,
		domainTx.ActionTypePayment,
		domainTx.TxTypeUnsigned,
	)
	txItem.SetUnsignedTx("input-test-hex", "0.100", "0.090", "0.010")

	txID, err := btcTxRepo.InsertUnsignedTx(domainTx.ActionTypePayment, txItem)
	require.NoError(t, err, "fail to create parent tx")

	// Create test inputs
	input1, err := domainBitcoin.NewBtcTxInput(
		txID, "input-txid-sqlc-1", 0, "input-address-sqlc-1",
		"client", "0.05", 6,
	)
	require.NoError(t, err, "fail to create input1")

	input2, err := domainBitcoin.NewBtcTxInput(
		txID, "input-txid-sqlc-2", 1, "input-address-sqlc-2",
		"client", "0.05", 6,
	)
	require.NoError(t, err, "fail to create input2")

	inputs := []*domainBitcoin.BtcTxInput{input1, input2}

	// Insert bulk
	err = btcTxInputRepo.InsertBulk(inputs)
	require.NoError(t, err, "fail to call InsertBulk()")

	// Get all by tx ID
	retrievedInputs, err := btcTxInputRepo.GetAllByTxID(txID)
	require.NoError(t, err, "fail to call GetAllByTxID()")
	require.Equal(t, 2, len(retrievedInputs), "GetAllByTxID() should return 2 inputs")

	// Get one
	oneInput, err := btcTxInputRepo.GetOne(retrievedInputs[0].ID)
	require.NoError(t, err, "fail to call GetOne()")
	require.Equal(t, "input-txid-sqlc-1", oneInput.InputTxid, "GetOne() InputTxid mismatch")

	// Insert single
	singleInput, err := domainBitcoin.NewBtcTxInput(
		txID, "input-txid-sqlc-3", 2, "input-address-sqlc-3",
		"client", "0.03", 6,
	)
	require.NoError(t, err, "fail to create singleInput")
	err = btcTxInputRepo.Insert(singleInput)
	require.NoError(t, err, "fail to call Insert()")

	// Verify count increased
	allInputs, err := btcTxInputRepo.GetAllByTxID(txID)
	require.NoError(t, err, "fail to call GetAllByTxID() after Insert()")
	require.Equal(t, 3, len(allInputs), "GetAllByTxID() should return 3 inputs")
}
