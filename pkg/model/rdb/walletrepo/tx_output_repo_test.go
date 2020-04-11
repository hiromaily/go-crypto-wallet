package walletrepo_test

//TODO: use mock interface
//func TestInsertTxOutputSequence(t *testing.T) {
//	var receiptID int64 = 1
//
//	txOutputDetails := []TxOutput{
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
//	//1.TestInsertTxOutputForUnsigned
//	t.Run("TestInsertTxOutputForUnsigned", func(t *testing.T) {
//		for _, typ := range types {
//			//FIXME: Error 1054: Unknown column 'output_address' in 'field list'
//			err := db.InsertTxOutputForUnsigned(typ, txOutputDetails, nil, true)
//			if err != nil {
//				t.Fatal(err)
//			}
//		}
//	})
//
//	//2.GetTxOutputByReceiptID
//	t.Run("GetTxOutputByReceiptID", func(t *testing.T) {
//		for _, typ := range types {
//			//FIXME: missing destination name input_txid in *[]model.TxOutput
//			txOutputs, err := db.GetTxOutputByReceiptID(typ, receiptID)
//			if err != nil {
//				t.Fatal(err)
//			}
//			t.Logf("[ActionType:%s] txOutputs: %+v", typ, txOutputs)
//		}
//	})
//}
