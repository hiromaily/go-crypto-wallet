package walletrepo_test

//TODO: use mock interface
//TODO: procedure for test is Insert, Select, Update
//func TestInsertTxInputSequence(t *testing.T) {
//	var receiptID int64 = 1
//
//	txInputDetails := []TxInput{
//		{
//			ReceiptID:          receiptID,
//			InputTxid:          "txidxxxxxx",
//			InputVout:          0,
//			InputAddress:       "address",
//			InputAccount:       "acount",
//			InputAmount:        "0.05",
//			InputConfirmations: 6,
//		},
//		{
//			ReceiptID:          receiptID,
//			InputTxid:          "txidxxxxxx2",
//			InputVout:          1,
//			InputAddress:       "address2",
//			InputAccount:       "acount2",
//			InputAmount:        "0.051111",
//			InputConfirmations: 8,
//		},
//	}
//
//	//1.TestInsertTxInputForUnsigned
//	t.Run("TestInsertTxInputForUnsigned", func(t *testing.T) {
//		for _, typ := range types {
//			err := db.InsertTxInputForUnsigned(typ, txInputDetails, nil, true)
//			if err != nil {
//				t.Fatal(err)
//			}
//		}
//	})
//
//	//2.GetTxInputByReceiptID
//	t.Run("GetTxInputByReceiptID", func(t *testing.T) {
//		for _, typ := range types {
//			txInputs, err := db.GetTxInputByReceiptID(typ, receiptID)
//			if err != nil {
//				t.Fatal(err)
//			}
//			t.Logf("[ActionType:%s] txInputs: %+v", typ, txInputs)
//		}
//	})
//}
