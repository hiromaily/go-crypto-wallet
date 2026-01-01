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

// TestBTCTxInputSqlc is integration test for TxInputRepositorySqlc
func TestBTCTxInputSqlc(t *testing.T) {
	// Need to create a btc_tx first
	btcTxRepo := watchTestutil.NewBTCTxRepositorySqlc()
	btcTxInputRepo := watchTestutil.NewBTCTxInputRepositorySqlc()

	// Create a parent tx
	txItem := &sqlcgen.BtcTx{
		Coin:              sqlcgen.BtcTxCoinBtc,
		Action:            sqlcgen.BtcTxActionPayment,
		UnsignedHexTx:     "input-test-hex",
		TotalInputAmount:  "0.100",
		TotalOutputAmount: "0.090",
		Fee:               "0.010",
	}
	txID, err := btcTxRepo.InsertUnsignedTx(domainTx.ActionTypePayment, txItem)
	require.NoError(t, err, "fail to create parent tx")

	// Create test inputs
	inputs := []*sqlcgen.BtcTxInput{
		{
			TxID:               txID,
			InputTxid:          "input-txid-sqlc-1",
			InputVout:          0,
			InputAddress:       "input-address-sqlc-1",
			InputAccount:       "client",
			InputAmount:        "0.05",
			InputConfirmations: 6,
		},
		{
			TxID:               txID,
			InputTxid:          "input-txid-sqlc-2",
			InputVout:          1,
			InputAddress:       "input-address-sqlc-2",
			InputAccount:       "client",
			InputAmount:        "0.05",
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
	singleInput := &sqlcgen.BtcTxInput{
		TxID:               txID,
		InputTxid:          "input-txid-sqlc-3",
		InputVout:          2,
		InputAddress:       "input-address-sqlc-3",
		InputAccount:       "client",
		InputAmount:        "0.03",
		InputConfirmations: 6,
	}
	err = btcTxInputRepo.Insert(singleInput)
	require.NoError(t, err, "fail to call Insert()")

	// Verify count increased
	allInputs, err := btcTxInputRepo.GetAllByTxID(txID)
	require.NoError(t, err, "fail to call GetAllByTxID() after Insert()")
	require.Equal(t, 3, len(allInputs), "GetAllByTxID() should return 3 inputs")
}
