package btc_test

import (
	"testing"

	"github.com/hiromaily/go-bitcoin/pkg/testutil"
)

// TestGetNetworkInfo is test for GetNetworkInfo
func TestGetNetworkInfo(t *testing.T) {
	//t.SkipNow()
	bc := testutil.GetBTC()

	// GetNetworkInfo
	if res, err := bc.GetNetworkInfo(); err != nil {
		t.Errorf("fail to call GetNetworkInfo(): %v", err)
	} else {
		t.Log(res)
	}

	//bc.Close()
}
