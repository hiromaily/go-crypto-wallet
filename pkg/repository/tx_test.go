package repository_test

import (
	"testing"

	"github.com/ericlagergren/decimal"
	"github.com/volatiletech/sqlboiler/types"

	"github.com/hiromaily/go-bitcoin/pkg/action"
	models "github.com/hiromaily/go-bitcoin/pkg/models/rdb"
	"github.com/hiromaily/go-bitcoin/pkg/testutil"
)

// TestTx is test for any data operation
func TestTx(t *testing.T) {

	//tx, err := db.Begin()
	// Rollback or commit
	//tx.Commit()
	//tx.Rollback()

	txRepo := testutil.NewTxRepository()

	// Delete records
	if _, err := txRepo.DeleteAll(); err != nil {
		t.Fatalf("fail to call DeleteAll() %v", err)
	}

	// Insert
	inputAmt := types.Decimal{Big: new(decimal.Big)}
	inputAmt.Big, _ = inputAmt.SetString("0.100")
	outputAmt := types.Decimal{Big: new(decimal.Big)}
	outputAmt.Big, _ = outputAmt.SetString("0.090")
	feeAmt := types.Decimal{Big: new(decimal.Big)}
	feeAmt.Big, _ = feeAmt.SetString("0.010")

	hex := "unsigned-hex"
	txItem := &models.TX{
		Action:            action.ActionTypePayment.String(),
		UnsignedHexTX:     hex,
		TotalInputAmount:  inputAmt,
		TotalOutputAmount: outputAmt,
		Fee:               feeAmt,
	}
	id, err := txRepo.InsertUnsignedTx(txItem)
	if err != nil {
		t.Fatalf("fail to call InsertUnsignedTx() %v", err)
	}
	// Search
	tmpTx, err := txRepo.GetOne(id)
	if err != nil {
		t.Fatalf("fail to call GetOne() %v", err)
	}
	if tmpTx.UnsignedHexTX != hex {
		t.Errorf("InsertUnsignedTx() = %s, want %s", tmpTx.UnsignedHexTX, hex)
		return
	}

	// Update
	hex2 := "unsigned-hex2"
	txItem.UnsignedHexTX = hex2
	_, err = txRepo.Update(txItem)
	if err != nil {
		t.Fatalf("fail to call UpdateTx() %v", err)
	}
}
