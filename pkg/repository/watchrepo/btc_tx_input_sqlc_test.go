//go:build integration
// +build integration

package watchrepo_test

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/quagmt/udecimal"
	"github.com/stretchr/testify/require"

	domainTx "github.com/hiromaily/go-crypto-wallet/pkg/domain/transaction"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
)

// TestBTCTxInputSqlc is integration test for TxInputRepositorySqlc
func TestBTCTxInputSqlc(t *testing.T) {
	// Need to create a btc_tx first
	btcTxRepo := testutil.NewBTCTxRepositorySqlc()
	btcTxInputRepo := testutil.NewBTCTxInputRepositorySqlc()

	// Create a parent tx
	inputAmt, _ := udecimal.Parse("0.100")
	outputAmt, _ := udecimal.Parse("0.090")
	feeAmt, _ := udecimal.Parse("0.010")

	txItem := &models.BTCTX{
		Action:            domainTx.ActionTypePayment.String(),
		UnsignedHexTX:     "input-test-hex",
		TotalInputAmount:  inputAmt,
		TotalOutputAmount: outputAmt,
		Fee:               feeAmt,
	}
	txID, err := btcTxRepo.InsertUnsignedTx(domainTx.ActionTypePayment, txItem)
	require.NoError(t, err, "fail to create parent tx")

	// Create test inputs
	amount, _ := udecimal.Parse("1.5")
	amount1, _ := udecimal.Parse("0.05")
	amount, _ := udecimal.Parse("1.5")
	amount2, _ := udecimal.Parse("0.05")

	inputs := []*models.BTCTXInput{
		{
			TXID:               txID,
			InputTxid:          "input-txid-sqlc-1",
			InputVout:          0,
			InputAddress:       "input-address-sqlc-1",
			InputAccount:       "client",
			InputAmount:        amount1,
			InputConfirmations: 6,
		},
		{
			TXID:               txID,
			InputTxid:          "input-txid-sqlc-2",
			InputVout:          1,
			InputAddress:       "input-address-sqlc-2",
			InputAccount:       "client",
			InputAmount:        amount2,
			InputConfirmations: 6,
		},
	}

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
	amount, _ := udecimal.Parse("1.5")
	amount3, _ := udecimal.Parse("0.03")
	singleInput := &models.BTCTXInput{
		TXID:               txID,
		InputTxid:          "input-txid-sqlc-3",
		InputVout:          2,
		InputAddress:       "input-address-sqlc-3",
		InputAccount:       "client",
		InputAmount:        amount3,
		InputConfirmations: 6,
	}
	err = btcTxInputRepo.Insert(singleInput)
	require.NoError(t, err, "fail to call Insert()")

	// Verify count increased
	allInputs, err := btcTxInputRepo.GetAllByTxID(txID)
	require.NoError(t, err, "fail to call GetAllByTxID() after Insert()")
	require.Equal(t, 3, len(allInputs), "GetAllByTxID() should return 3 inputs")
}
