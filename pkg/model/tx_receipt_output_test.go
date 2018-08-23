package model_test

//
//import (
//	"testing"
//
//	. "github.com/hiromaily/go-bitcoin/pkg/model"
//)
//
////TODO:テストの順序はInsert, Select, Update
//
//func TestInsertTxReceiptDetailForUnsigned(t *testing.T) {
//	txReceiptOutputs := []TxOutput{
//		{
//			ReceiptID:     1,
//			OutputAddress: "output-address",
//			OutputAccount: "output-acount",
//			OutputAmount:  "0.05",
//			IsChange:      false,
//		},
//		{
//			ReceiptID:     1,
//			OutputAddress: "output-address2",
//			OutputAccount: "output-acount2",
//			OutputAmount:  "0.25",
//			IsChange:      true,
//		},
//	}
//
//	err := db.InsertTxReceiptOutputForUnsigned(txReceiptOutputs, nil, true)
//	if err != nil {
//		t.Fatal(err)
//	}
//}
//
//func TestGetTxReceiptOutputByReceiptID(t *testing.T) {
//	txReceiptOutputs, err := db.GetTxReceiptOutputByReceiptID(1)
//	if err != nil {
//		t.Fatal(err)
//	}
//	t.Log(txReceiptOutputs)
//}
