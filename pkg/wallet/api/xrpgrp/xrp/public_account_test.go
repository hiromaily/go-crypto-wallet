package xrp_test

import (
	"testing"

	"github.com/bookerzzz/grok"

	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
)

// TestAccountChannels is test for AccountChannels
func TestAccountChannels(t *testing.T) {
	//t.SkipNow()
	xr := testutil.GetXRP()

	res, err := xr.AccountChannels("rN7n7otQDd6FczFgLdSqtcsAUxDkw6fzRH", "rf1BiGeXwwQoi8Z2ueFYTEXSwuJYfV2Jpn")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("accountChannels:", res)
	grok.Value(res)
}
