package api_test

import (
	"testing"
)

// TestEstimateSmartFee
func TestEstimateSmartFee(t *testing.T) {
	fee, err := wlt.BTC.EstimateSmartFee()
	if err != nil {
		t.Errorf("Unexpectedly error occorred. %v", err)
	}
	//この値は変動するので、目視確認
	t.Log(fee)
	//0.00250122

	//estimatesmartfee 6
	//0.00250122
}

// TestGetTransactionFee
func TestGetTransactionFee(t *testing.T) {
	var tests = []struct {
		hex   string //未署名のトランザクションhex
		txNum int    //utxo数
		size  int    //msgTx.SerializeSize()で算出されたサイズ
	}{
		{
			"02000000022e0183cd8e082c185030b8eed4bf19bace65936960fe79736dc21f3b0586b7640100000000ffffffff8afd01d2ecdfeb1657ae7a0ecee9e89b86feb58ed10803cdf6c95d25354161ff0100000000ffffffff01a0f69e0b0000000017a9148191d41a7415a6a1f6ee14337e039f50b949e80e8700000000",
			2,
			124,
		},
	}
	for _, val := range tests {
		// Hexからトランザクションを取得
		msgTx, err := wlt.BTC.ToMsgTx(val.hex)
		if err != nil {
			t.Fatal("Hex can not be converted as *wire.MsgTx")
		}

		//size
		t.Logf("msgTx.SerializeSize() %d", msgTx.SerializeSize())

		fee, err := wlt.BTC.GetTransactionFee(msgTx)
		if err != nil {
			t.Fatalf("GetTransactionFee() must be fixed. error: %v", err)
		}
		t.Logf("fee: %v", fee)
		//0.00264256(bit/kb)
		//0.000328 BTC
	}
}
