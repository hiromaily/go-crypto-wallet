package model_test

import (
	"testing"
	"time"

	"github.com/hiromaily/go-bitcoin/pkg/enum"
	. "github.com/hiromaily/go-bitcoin/pkg/model"
)

var (
	testPaymentID  int64
	testPaymentHEX string
)

//TODO:テストの順序はInsert, Select, Update

func TestInsertTxPaymentForUnsigned(t *testing.T) {
	txPayment := TxTable{}
	txPayment.UnsignedHexTx = "12345"
	txPayment.TotalInputAmount = "1.5"
	txPayment.TotalOutputAmount = "1.3"
	txPayment.Fee = "0.2"
	txPayment.TxType = enum.TxTypeValue[enum.TxTypeUnsigned]

	testPaymentID, err := db.InsertTxPaymentForUnsigned(&txPayment, nil, true)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("ID is %d", testPaymentID)
}

func TestGetTxPaymentByID(t *testing.T) {
	//TODO:testPaymentIDは引き継げない？？
	testPaymentID = 1 //FIXME:とりあえず
	if testPaymentID == 0 {
		t.Fatal("testPaymentID should be set")
	}
	txPayment, err := db.GetTxPaymentByID(testPaymentID)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("txPayment is %+v", txPayment)
	testPaymentHEX = txPayment.UnsignedHexTx
}

func TestGetTxPaymentByUnsignedHex(t *testing.T) {
	//TODO:testPaymentIDは引き継げない？？
	testPaymentHEX = "12345" //FIXME:とりあえず
	if testPaymentHEX == "" {
		t.Fatal("testPaymentHEX should be set")
	}

	//hexTx := "02000000ss2b5085ddcbe61200c54b29c2d664df31341cd72834ec03a6c0b71bba7054429cb0100000000ffffffffb9401d39321d17fe1ec07668256820b0ccd2184b9ad4a8083c9a7295641d52220100000000ffffffff0114ba9e0b0000000017a9148191d41a7415a6a1f6ee14337e039f50b949e80e8700000000"
	count, err := db.GetTxPaymentCountByUnsignedHex(testPaymentHEX)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(count)
}

func TestUpdateTxPaymentForSent(t *testing.T) {
	tm := time.Now()
	txPayment := TxTable{}
	txPayment.ID = 1
	txPayment.SignedHexTx = "signedHex"
	txPayment.SentHashTx = "sentHashID"
	txPayment.SentUpdatedAt = &tm
	txPayment.TxType = enum.TxTypeValue[enum.TxTypeSent]

	affected, err := db.UpdateTxPaymentForSent(&txPayment, nil, true)
	if err != nil {
		t.Fatal(err)
	}
	if affected == 0 {
		t.Fatal("table was not updated")
	}
}
