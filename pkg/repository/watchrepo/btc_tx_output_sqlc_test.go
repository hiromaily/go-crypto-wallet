//go:build integration
// +build integration

package watchrepo_test

import (
	"testing"

	"github.com/quagmt/udecimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "github.com/go-sql-driver/mysql"

	"github.com/hiromaily/go-crypto-wallet/pkg/action"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
)

// TestBTCTxOutputSqlc is integration test for TxOutputRepositorySqlc
func TestBTCTxOutputSqlc(t *testing.T) {
	// Need to create a btc_tx first
	btcTxRepo := testutil.NewBTCTxRepositorySqlc()
	btcTxOutputRepo := testutil.NewBTCTxOutputRepositorySqlc()

	// Create a parent tx
	inputAmt, _ := udecimal.Parse("0.100")
	outputAmt, _ := udecimal.Parse("0.090")
	feeAmt, _ := udecimal.Parse("0.010")

	txItem := &models.BTCTX{
		Action:            action.ActionTypePayment.String(),
		UnsignedHexTX:     "output-test-hex",
		TotalInputAmount:  inputAmt,
		TotalOutputAmount: outputAmt,
		Fee:               feeAmt,
	}
	txID, err := btcTxRepo.InsertUnsignedTx(action.ActionTypePayment, txItem)
	require.NoError(t, err, "fail to create parent tx")

	// Create test outputs
	amount, _ := udecimal.Parse("1.5")
	amount1, _ := udecimal.Parse("0.08")
	amount, _ := udecimal.Parse("1.5")
	amount2, _ := udecimal.Parse("0.01")

	outputs := []*models.BTCTXOutput{
		{
			TXID:          txID,
			OutputAddress: "output-address-sqlc-1",
			OutputAccount: "receipt",
			OutputAmount:  amount1,
			IsChange:      false,
		},
		{
			TXID:          txID,
			OutputAddress: "output-address-sqlc-2",
			OutputAccount: "change",
			OutputAmount:  amount2,
			IsChange:      true,
		},
	}

	// Insert bulk
	err = btcTxOutputRepo.InsertBulk(outputs)
	require.NoError(t, err, "fail to call InsertBulk()")

	// Get all by tx ID
	retrievedOutputs, err := btcTxOutputRepo.GetAllByTxID(txID)
	require.NoError(t, err, "fail to call GetAllByTxID()")
	assert.Len(t, retrievedOutputs, 2, "GetAllByTxID() should return 2 outputs")

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
	assert.True(t, hasChange, "GetAllByTxID() should return at least one change output")
	assert.True(t, hasNonChange, "GetAllByTxID() should return at least one non-change output")

	// Get one
	oneOutput, err := btcTxOutputRepo.GetOne(retrievedOutputs[0].ID)
	require.NoError(t, err, "fail to call GetOne()")
	assert.Equal(t, txID, oneOutput.TXID, "GetOne() should return output with correct TXID")

	// Insert single
	amount, _ := udecimal.Parse("1.5")
	amount3, _ := udecimal.Parse("0.02")
	singleOutput := &models.BTCTXOutput{
		TXID:          txID,
		OutputAddress: "output-address-sqlc-3",
		OutputAccount: "receipt",
		OutputAmount:  amount3,
		IsChange:      false,
	}
	err = btcTxOutputRepo.Insert(singleOutput)
	require.NoError(t, err, "fail to call Insert()")

	// Verify count increased
	allOutputs, err := btcTxOutputRepo.GetAllByTxID(txID)
	require.NoError(t, err, "fail to call GetAllByTxID() after Insert()")
	assert.Len(t, allOutputs, 3, "GetAllByTxID() should return 3 outputs after Insert()")
}
