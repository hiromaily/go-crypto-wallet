package btc_test

import (
	"testing"

	"github.com/hiromaily/go-bitcoin/pkg/testutil"
)

// TestGetBlockCount is test for GetBlockCount
func TestGetBlockCount(t *testing.T) {
	//t.SkipNow()
	bc := testutil.GetBTC()

	// GetBalance
	if res, err := bc.GetBlockCount(); err != nil {
		t.Errorf("fail to call GetBlockCount(): %v", err)
	} else {
		t.Log(res)
	}

	//bc.Close()
}
