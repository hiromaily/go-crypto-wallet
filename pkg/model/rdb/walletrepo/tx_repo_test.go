package walletrepo_test

//TODO: use mock interface
//TODO: procedure of test is Insert, Select, Update, Select
//var (
//	types = []ctype.ActionType{ctype.ActionTypeReceipt, ctype.ActionTypePayment}
//)
//
//
//func TestTxTableSequence(t *testing.T) {
//
//	var (
//		//テストで利用するシーケンシャルなデータ
//		testTableID  = make(map[ctype.ActionType]int64)
//		testTableHEX = make(map[ctype.ActionType]string)
//		err          error
//	)
//
//	//1.TestInsertTxReceiptForUnsigned Insert
//	t.Run("TestInsertTxReceiptForUnsigned", func(t *testing.T) {
//		txTable := TxTable{}
//		txTable.UnsignedHexTx = "test12345"
//		txTable.TotalInputAmount = "1.5"
//		txTable.TotalOutputAmount = "1.3"
//		txTable.Fee = "0.2"
//		txTable.TxType = ctype.TxTypeValue[ctype.TxTypeUnsigned]
//
//		for _, typ := range types {
//			testTableID[typ], err = db.InsertTxForUnsigned(typ, &txTable, nil, true)
//			if err != nil {
//				t.Fatal(err)
//			}
//
//			t.Logf("[ActionType:%s] ID is %d", typ, testTableID[typ])
//		}
//	})
//
//	//2.TestGetTxTableByID Get
//	t.Run("TestGetTxTableByID", func(t *testing.T) {
//		for _, typ := range types {
//			if testTableID[typ] == 0 {
//				t.Fatalf("testReceiptID[%s] should be set", typ)
//			}
//			txTable, err := db.GetTxByID(typ, testTableID[typ])
//			if err != nil {
//				t.Fatal(err)
//			}
//
//			t.Logf("[ActionType:%s] txTable is %+v", typ, txTable)
//			testTableHEX[typ] = txTable.UnsignedHexTx
//		}
//	})
//
//	//3.TestGetTxTableByUnsignedHex Get
//	t.Run("TestGetTxTableByUnsignedHex", func(t *testing.T) {
//		for _, typ := range types {
//			if testTableHEX[typ] == "" {
//				t.Fatalf("testReceiptHEX[%s] should be set", typ)
//			}
//
//			count, err := db.GetTxCountByUnsignedHex(typ, testTableHEX[typ])
//			if err != nil {
//				t.Fatal(err)
//			}
//			t.Logf("[ActionType:%s] count is %d", typ, count)
//		}
//	})
//
//	//4.TestUpdateTxTableForSent Update
//	t.Run("TestUpdateTxTableForSent", func(t *testing.T) {
//		for _, typ := range types {
//			tm := time.Now()
//			txTable := TxTable{}
//			txTable.ID = testTableID[typ] //更新のキーとなる
//			txTable.SignedHexTx = "signedHex"
//			txTable.SentHashTx = "sentHashID"
//			txTable.SentUpdatedAt = &tm
//			txTable.TxType = ctype.TxTypeValue[ctype.TxTypeSent]
//
//			affected, err := db.UpdateTxAfterSent(typ, &txTable, nil, true)
//			if err != nil {
//				t.Fatal(err)
//			}
//			if affected == 0 {
//				t.Fatalf("[ActionType:%s] table was not updated", typ)
//			}
//		}
//	})
//}
