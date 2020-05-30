package rippleapi_test

import (
	"testing"

	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
)

// TestPrepareTransaction is test for PrepareTransaction
func TestPrepareTransaction(t *testing.T) {
	//t.SkipNow()
	api := testutil.GetRippleAPI()

	var (
		sernderAccount          = "rN7n7otQDd6FczFgLdSqtcsAUxDkw6fzRH"
		receiverAccount         = "rf1BiGeXwwQoi8Z2ueFYTEXSwuJYfV2Jpn"
		amount          float64 = 10.5
	)
	err := api.PrepareTransaction(sernderAccount, receiverAccount, amount)
	if err != nil {
		t.Fatal(err)
	}
}
