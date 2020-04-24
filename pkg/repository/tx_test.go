package repository_test

import (
	"github.com/ericlagergren/decimal"
	"github.com/hiromaily/go-bitcoin/pkg/action"
	"github.com/volatiletech/sqlboiler/types"
	"testing"

	models "github.com/hiromaily/go-bitcoin/pkg/models/rdb"
	"github.com/hiromaily/go-bitcoin/pkg/testutil"
)

// TestTx is test for any data operation
func TestTx(t *testing.T) {
	txRepo := testutil.NewTxRepository()
	//bc := testutil.GetBTC()

	inputAmt := types.Decimal{Big: new(decimal.Big)}
	inputAmt.Big, _ = inputAmt.SetString("0.100")
	outputAmt := types.Decimal{Big: new(decimal.Big)}
	outputAmt.Big, _ = outputAmt.SetString("0.090")
	feeAmt := types.Decimal{Big: new(decimal.Big)}
	feeAmt.Big, _ = feeAmt.SetString("0.010")

	// Insert
	txItem := &models.TX{
		Action:            action.ActionTypePayment.String(),
		UnsignedHexTX:     "unsigned-hex",
		TotalInputAmount:  inputAmt,
		TotalOutputAmount: outputAmt,
		Fee:               feeAmt,
	}
	if err := txRepo.InsertUnsignedTx(txItem); err != nil {
		t.Fatalf("fail to call InsertUnsignedTx() %v", err)
	}

}
