//go:build integration
// +build integration

package btc_test

import (
	"testing"

	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
)

// TestGetNetworkInfo is test for GetNetworkInfo
func TestGetNetworkInfo(t *testing.T) {
	// t.SkipNow()
	bc := testutil.GetBTC()

	// GetNetworkInfo
	if res, err := bc.GetNetworkInfo(); err != nil {
		t.Errorf("fail to call GetNetworkInfo(): %v", err)
	} else {
		t.Log(res)
	}

	// bc.Close()
}

// TestGetBlockchainInfo is test for GetBlockchainInfo
func TestBlockchainInfo(t *testing.T) {
	// t.SkipNow()
	bc := testutil.GetBTC()

	// GetBlockchainInfo
	if res, err := bc.GetBlockchainInfo(); err != nil {
		t.Errorf("fail to call GetBlockchainInfo(): %v", err)
	} else {
		t.Log(res)
	}

	// bc.Close()
}
