//go:build integration
// +build integration

package btc_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
)

type feeTest struct {
	testutil.BTCTestSuite
}

// TestEstimateSmartFee is test for EstimateSmartFee
func (ft *feeTest) TestEstimateSmartFee() {
	// EstimateSmartFee
	res, err := ft.BTC.EstimateSmartFee()
	ft.NoError(err)
	if err == nil {
		ft.T().Logf("%f", res)
	}
}

func TestFeeTestSuite(t *testing.T) {
	suite.Run(t, new(feeTest))
}
