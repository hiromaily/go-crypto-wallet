package model_test

import (
	"testing"

	. "github.com/hiromaily/go-bitcoin/pkg/model"
)

//TODO:テストの順序はInsert, Select, Update

func TestInsertTxPaymentOutputForUnsigned(t *testing.T) {
	txPaymentDetails := []TxInput{
		{
			ReceiptID:          1,
			InputTxid:          "txidxxxxxx",
			InputVout:          0,
			InputAddress:       "address",
			InputAccount:       "acount",
			InputAmount:        "0.05",
			InputConfirmations: 6,
		},
		{
			ReceiptID:          1,
			InputTxid:          "txidxxxxxx2",
			InputVout:          1,
			InputAddress:       "address2",
			InputAccount:       "acount2",
			InputAmount:        "0.051111",
			InputConfirmations: 8,
		},
	}

	err := db.InsertTxPaymentInputForUnsigned(txPaymentDetails, nil, true)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetTxPaymentInputByReceiptID(t *testing.T) {
	txPaymentInputs, err := db.GetTxPaymentInputByReceiptID(1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(txPaymentInputs)
}
