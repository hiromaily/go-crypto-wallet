package xrp_test

import (
	"testing"

	"github.com/bookerzzz/grok"

	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
)

// TestServerInfo is test for ServerInfo
func TestServerInfo(t *testing.T) {
	// t.SkipNow()
	xr := testutil.GetXRP()

	res, err := xr.ServerInfo()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ServerInfo:", res)
	grok.Value(res)
}
