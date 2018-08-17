package model_test

import (
	"testing"
	"time"

	. "github.com/hiromaily/go-bitcoin/pkg/model"
)

func TestGetTxReceiptByID(t *testing.T) {
	hexTx := "02000000ss2b5085ddcbe61200c54b29c2d664df31341cd72834ec03a6c0b71bba7054429cb0100000000ffffffffb9401d39321d17fe1ec07668256820b0ccd2184b9ad4a8083c9a7295641d52220100000000ffffffff0114ba9e0b0000000017a9148191d41a7415a6a1f6ee14337e039f50b949e80e8700000000"
	count, err := db.GetTxReceiptByUnsignedHex(db.TableNameReceipt(), hexTx)

	if err != nil {
		t.Fatal(err)
	}
	t.Log(count)
}

func TestInsertTxReceiptForUnsigned(t *testing.T) {
	txReceipt := TxReceipt{}
	txReceipt.UnsignedHexTx = "12345"
	txReceipt.TotalAmount = "1.5"
	txReceipt.Fee = "0.2"
	txReceipt.ReceiverAddress = "address"
	txReceipt.TxType = 1

	id, err := db.InsertTxReceiptForUnsigned(db.TableNameReceipt(), &txReceipt, nil, true)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("ID is %d", id)
}

func TestUpdateTxReceiptForSent(t *testing.T) {
	tm := time.Now()
	txReceipt := TxReceipt{}
	txReceipt.ID = 1
	txReceipt.SignedHexTx = "signedHex"
	txReceipt.SentHexTx = "sentTxID"
	txReceipt.SentUpdatedAt = &tm
	txReceipt.TxType = 3 //未署名:TODO:Constとして定義しておく

	affected, err := db.UpdateTxReceiptForSent(db.TableNameReceipt(), &txReceipt, nil, true)
	if err != nil {
		t.Fatal(err)
	}
	if affected == 0 {
		t.Fatal("table was not updated")
	}
}
