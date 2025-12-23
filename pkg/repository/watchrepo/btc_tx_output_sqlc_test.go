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
	if err != nil {
		t.Fatalf("fail to create parent tx: %v", err)
	}

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
	if err := btcTxOutputRepo.InsertBulk(outputs); err != nil {
		t.Fatalf("fail to call InsertBulk() %v", err)
	}

	// Get all by tx ID
	retrievedOutputs, err := btcTxOutputRepo.GetAllByTxID(txID)
	if err != nil {
		t.Fatalf("fail to call GetAllByTxID() %v", err)
	}
	if len(retrievedOutputs) != 2 {
		t.Errorf("GetAllByTxID() returned %d outputs, want 2", len(retrievedOutputs))
		return
	}

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
	if !hasChange || !hasNonChange {
		t.Errorf("GetAllByTxID() should return both change and non-change outputs")
		return
	}

	// Get one
	oneOutput, err := btcTxOutputRepo.GetOne(retrievedOutputs[0].ID)
	if err != nil {
		t.Fatalf("fail to call GetOne() %v", err)
	}
	if oneOutput.TXID != txID {
		t.Errorf("GetOne() returned TXID = %d, want %d", oneOutput.TXID, txID)
		return
	}

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
	if err := btcTxOutputRepo.Insert(singleOutput); err != nil {
		t.Fatalf("fail to call Insert() %v", err)
	}

	// Verify count increased
	allOutputs, err := btcTxOutputRepo.GetAllByTxID(txID)
	if err != nil {
		t.Fatalf("fail to call GetAllByTxID() after Insert() %v", err)
	}
	if len(allOutputs) != 3 {
		t.Errorf("GetAllByTxID() returned %d outputs, want 3", len(allOutputs))
		return
	}
}
