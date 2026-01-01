//go:build integration
// +build integration

package watch_test

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"

	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/mysql/sqlcgen"
	watchTestutil "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/watch/testutil"
)

// TestBTCTxOutputSqlc is integration test for TxOutputRepositorySqlc
func TestBTCTxOutputSqlc(t *testing.T) {
	// Need to create a btc_tx first
	btcTxRepo := watchTestutil.NewBTCTxRepositorySqlc()
	btcTxOutputRepo := watchTestutil.NewBTCTxOutputRepositorySqlc()

	// Create a parent tx
	txItem := &sqlcgen.BtcTx{
		Coin:              sqlcgen.BtcTxCoinBtc,
		Action:            sqlcgen.BtcTxActionPayment,
		UnsignedHexTx:     "output-test-hex",
		TotalInputAmount:  "0.100",
		TotalOutputAmount: "0.090",
		Fee:               "0.010",
	}
	txID, err := btcTxRepo.InsertUnsignedTx(domainTx.ActionTypePayment, txItem)
	require.NoError(t, err, "fail to create parent tx")

	// Create test outputs
	outputs := []*sqlcgen.BtcTxOutput{
		{
			TxID:          txID,
			OutputAddress: "output-address-sqlc-1",
			OutputAccount: "receipt",
			OutputAmount:  "0.08",
			IsChange:      false,
		},
		{
			TxID:          txID,
			OutputAddress: "output-address-sqlc-2",
			OutputAccount: "change",
			OutputAmount:  "0.01",
			IsChange:      true,
		},
	}

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
	singleOutput := &sqlcgen.BtcTxOutput{
		TxID:          txID,
		OutputAddress: "output-address-sqlc-3",
		OutputAccount: "receipt",
		OutputAmount:  "0.02",
		IsChange:      false,
	}
	err = btcTxOutputRepo.Insert(singleOutput)
	require.NoError(t, err, "fail to call Insert()")

	// Verify count increased
	allOutputs, err := btcTxOutputRepo.GetAllByTxID(txID)
	require.NoError(t, err, "fail to call GetAllByTxID() after Insert()")
	require.Len(t, allOutputs, 3, "GetAllByTxID() should return 3 outputs after Insert()")
}
