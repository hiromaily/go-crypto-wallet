package watchrepo_test

import (
	"testing"

	"github.com/ericlagergren/decimal"
	_ "github.com/go-sql-driver/mysql"
	"github.com/volatiletech/sqlboiler/types"

	"github.com/hiromaily/go-crypto-wallet/pkg/action"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
	"github.com/hiromaily/go-crypto-wallet/pkg/tx"
)

// TestTx is test for any data operation
func TestTx(t *testing.T) {
	//boil.DebugMode = true
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
	actionType := action.ActionTypePayment
	txItem := &models.BTCTX{
		Action:            actionType.String(),
		UnsignedHexTX:     hex,
		TotalInputAmount:  inputAmt,
		TotalOutputAmount: outputAmt,
		Fee:               feeAmt,
	}
	id, err := txRepo.InsertUnsignedTx(actionType, txItem)
	if err != nil {
		t.Fatalf("fail to call InsertUnsignedTx() %v", err)
	}
	// check inserted record
	tmpTx, err := txRepo.GetOne(id)
	if err != nil {
		t.Fatalf("fail to call GetOne() %v", err)
	}
	if tmpTx.UnsignedHexTX != hex {
		t.Errorf("InsertUnsignedTx() = %s, want %s", tmpTx.UnsignedHexTX, hex)
		return
	}
	// check Count
	cnt, err := txRepo.GetCountByUnsignedHex(actionType, hex)
	if err != nil {
		t.Fatalf("fail to call GetCount() %v", err)
	}
	if cnt != 1 {
		t.Errorf("GetCount() = %d, want %d", cnt, 1)
		return
	}

	// Update only UnsignedHexTX
	hex2 := "unsigned-hex2"
	txItem.UnsignedHexTX = hex2
	_, err = txRepo.Update(txItem)
	if err != nil {
		t.Fatalf("fail to call UpdateTx() %v", err)
	}
	// check updated unsigned hex tx
	tmpTx, err = txRepo.GetOne(txItem.ID)
	if err != nil {
		t.Fatalf("fail to call GetOne() %v", err)
	}
	if tmpTx.UnsignedHexTX != hex2 {
		t.Errorf("Update() = %s, want %s", tmpTx.UnsignedHexTX, hex2)
		return
	}

	// Update like after tx sent
	// TODO: how to update partially??
	// => object should includes all, base object should be retrieved for updating first
	//    not good performance
	signedHex := "signed-hex"
	sentHashTx := "sent-hash-tx"
	_, err = txRepo.UpdateAfterTxSent(txItem.ID, tx.TxTypeSent, signedHex, sentHashTx)
	if err != nil {
		t.Fatalf("fail to call UpdateTx() %v", err)
	}
	// check updated record
	tmpTx, err = txRepo.GetOne(txItem.ID)
	if err != nil {
		t.Fatalf("fail to call GetOne() %v", err)
	}
	if tmpTx.SignedHexTX != signedHex {
		t.Errorf("Update() = %s, want %s", tmpTx.SignedHexTX, signedHex)
		return
	}
	// sent_hash_tx should be retrieved
	hashes, err := txRepo.GetSentHashTx(actionType, tx.TxTypeSent)
	if err != nil {
		t.Fatalf("fail to call GetSentHashTx() %v", err)
	}
	if len(hashes) != 1 {
		t.Errorf("GetSentHashTx() = %d, want %d", len(hashes), 1)
		return
	}

	// update txType
	_, err = txRepo.UpdateTxTypeBySentHashTx(actionType, tx.TxTypeDone, sentHashTx)
	if err != nil {
		t.Fatalf("fail to call UpdateTxTypeBySentHashTx() %v", err)
	}
	// check updated record
	tmpTx, err = txRepo.GetOne(txItem.ID)
	if err != nil {
		t.Fatalf("fail to call GetOne() %v", err)
	}
	if tmpTx.CurrentTXType != tx.TxTypeDone.Int8() {
		t.Errorf("Update() = %d, want %d", tmpTx.CurrentTXType, tx.TxTypeDone.Int8())
		return
	}

	// update txType
	_, err = txRepo.UpdateTxType(txItem.ID, tx.TxTypeNotified)
	if err != nil {
		t.Fatalf("fail to call UpdateTxType() %v", err)
	}
	// check updated record
	tmpTx, err = txRepo.GetOne(txItem.ID)
	if err != nil {
		t.Fatalf("fail to call GetOne() %v", err)
	}
	if tmpTx.CurrentTXType != tx.TxTypeNotified.Int8() {
		t.Errorf("Update() = %d, want %d", tmpTx.CurrentTXType, tx.TxTypeNotified.Int8())
		return
	}

}
