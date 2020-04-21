package btc_test

import (
	"testing"

	"github.com/hiromaily/go-bitcoin/pkg/testutil"
)

// TestEstimateSmartFee is test for EstimateSmartFee
func TestEstimateSmartFee(t *testing.T) {
	//t.SkipNow()
	bc := testutil.GetBTC()

	// GetBalance
	if res, err := bc.EstimateSmartFee(); err != nil {
		t.Errorf("fail to call EstimateSmartFee(): %v", err)
	} else {
		t.Logf("%f", res)
	}

	bc.Close()
}
