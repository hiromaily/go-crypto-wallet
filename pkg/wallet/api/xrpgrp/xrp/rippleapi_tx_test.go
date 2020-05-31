package xrp_test

import (
	"testing"

	"github.com/bookerzzz/grok"

	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
)

// TestPrepareTransaction is test for PrepareTransaction
func TestPrepareTransaction(t *testing.T) {
	//t.SkipNow()
	xr := testutil.GetXRP()

	var (
		sernderAccount          = "rftGDSJBvdwHeqFnqtdhHAKC5TkgoWswmi"
		receiverAccount         = "rf1BiGeXwwQoi8Z2ueFYTEXSwuJYfV2Jpn"
		amount          float64 = 10.5
	)
	txJSON, err := xr.PrepareTransaction(sernderAccount, receiverAccount, amount)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(txJSON)
	grok.Value(txJSON)
}
