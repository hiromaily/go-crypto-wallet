package btc_test

import (
	"github.com/hiromaily/go-bitcoin/pkg/testutil"
	"testing"
)

// TestListAccounts is test for ListAccounts
func TestGetBalance(t *testing.T) {
	//t.SkipNow()
	bc := testutil.GetBTC()

	if res, err := bc.GetBalance(); err != nil {
		t.Errorf("fail to call GetBalance(): %v", err)
	} else {
		t.Log(res)
	}

	bc.Close()
}
