//go:build integration
// +build integration

package btc_test

import (
	"testing"

	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
)

// TestEstimateSmartFee is test for EstimateSmartFee
func TestEstimateSmartFee(t *testing.T) {
	// t.SkipNow()
	bc := testutil.GetBTC()

	// EstimateSmartFee
	if res, err := bc.EstimateSmartFee(); err != nil {
		t.Errorf("fail to call EstimateSmartFee(): %v", err)
	} else {
		t.Logf("%f", res)
	}

	// bc.Close()
}
