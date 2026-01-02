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

// TestBTCTxOutputSqlc is integration test for TxOutputRepositorySqlc
func TestBTCTxOutputSqlc(t *testing.T) {
	// Need to create a btc_tx first
	btcTxRepo := watchTestutil.NewBTCTxRepositorySqlc()
	btcTxOutputRepo := watchTestutil.NewBTCTxOutputRepositorySqlc()

	// Create a parent tx
	txItem := domainBitcoin.NewBtcTransaction(
		domainCoin.BTC,
		domainTx.ActionTypePayment,
		domainTx.TxTypeUnsigned,
	)
	txItem.SetUnsignedTx("output-test-hex", "0.100", "0.090", "0.010")

	txID, err := btcTxRepo.InsertUnsignedTx(domainTx.ActionTypePayment, txItem)
	require.NoError(t, err, "fail to create parent tx")

	// Create test outputs
	output1, err := domainBitcoin.NewBtcTxOutput(
		txID, "output-address-sqlc-1", "receipt", "0.08", false,
	)
	require.NoError(t, err, "fail to create output1")

	output2, err := domainBitcoin.NewBtcTxOutput(
		txID, "output-address-sqlc-2", "change", "0.01", true,
	)
	require.NoError(t, err, "fail to create output2")

	outputs := []*domainBitcoin.BtcTxOutput{output1, output2}

	// Insert bulk
	err = btcTxOutputRepo.InsertBulk(outputs)
	require.NoError(t, err, "fail to call InsertBulk()")

	// Get all by tx ID
	retrievedOutputs, err := btcTxOutputRepo.GetAllByTxID(txID)
	require.NoError(t, err, "fail to call GetAllByTxID()")
	require.Len(t, retrievedOutputs, 2, "GetAllByTxID() should return 2 outputs")

	// Verify one is change and one is not
	hasChange := false
	hasNonChange := false
	for _, output := range retrievedOutputs {
		if output.IsChange {
			hasChange = true
		} else {
			hasNonChange = true
		}
	}
	require.True(t, hasChange, "GetAllByTxID() should return at least one change output")
	require.True(t, hasNonChange, "GetAllByTxID() should return at least one non-change output")

	// Get one
	oneOutput, err := btcTxOutputRepo.GetOne(retrievedOutputs[0].ID)
	require.NoError(t, err, "fail to call GetOne()")
	require.Equal(t, txID, oneOutput.TxID, "GetOne() should return output with correct TxID")

	// Insert single
	singleOutput, err := domainBitcoin.NewBtcTxOutput(
		txID, "output-address-sqlc-3", "receipt", "0.02", false,
	)
	require.NoError(t, err, "fail to create singleOutput")
	err = btcTxOutputRepo.Insert(singleOutput)
	require.NoError(t, err, "fail to call Insert()")

	// Verify count increased
	allOutputs, err := btcTxOutputRepo.GetAllByTxID(txID)
	require.NoError(t, err, "fail to call GetAllByTxID() after Insert()")
	require.Len(t, allOutputs, 3, "GetAllByTxID() should return 3 outputs after Insert()")
}
