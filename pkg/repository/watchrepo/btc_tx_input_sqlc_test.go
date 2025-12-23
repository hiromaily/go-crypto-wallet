//go:build integration
// +build integration

package watchrepo_test

import (
	"testing"

	"github.com/quagmt/udecimal"
	_ "github.com/go-sql-driver/mysql"

	"github.com/hiromaily/go-crypto-wallet/pkg/action"
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
		Action:            action.ActionTypePayment.String(),
		UnsignedHexTX:     "input-test-hex",
		TotalInputAmount:  inputAmt,
		TotalOutputAmount: outputAmt,
		Fee:               feeAmt,
	}
	txID, err := btcTxRepo.InsertUnsignedTx(action.ActionTypePayment, txItem)
	if err != nil {
		t.Fatalf("fail to create parent tx: %v", err)
	}

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
	if err := btcTxInputRepo.InsertBulk(inputs); err != nil {
		t.Fatalf("fail to call InsertBulk() %v", err)
	}

	// Get all by tx ID
	retrievedInputs, err := btcTxInputRepo.GetAllByTxID(txID)
	if err != nil {
		t.Fatalf("fail to call GetAllByTxID() %v", err)
	}
	if len(retrievedInputs) != 2 {
		t.Errorf("GetAllByTxID() returned %d inputs, want 2", len(retrievedInputs))
		return
	}

	// Get one
	oneInput, err := btcTxInputRepo.GetOne(retrievedInputs[0].ID)
	if err != nil {
		t.Fatalf("fail to call GetOne() %v", err)
	}
	if oneInput.InputTxid != "input-txid-sqlc-1" {
		t.Errorf("GetOne() returned InputTxid = %s, want input-txid-sqlc-1", oneInput.InputTxid)
		return
	}

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
	if err := btcTxInputRepo.Insert(singleInput); err != nil {
		t.Fatalf("fail to call Insert() %v", err)
	}

	// Verify count increased
	allInputs, err := btcTxInputRepo.GetAllByTxID(txID)
	if err != nil {
		t.Fatalf("fail to call GetAllByTxID() after Insert() %v", err)
	}
	if len(allInputs) != 3 {
		t.Errorf("GetAllByTxID() returned %d inputs, want 3", len(allInputs))
		return
	}
}
