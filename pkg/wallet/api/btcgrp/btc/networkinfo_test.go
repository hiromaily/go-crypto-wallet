//go:build integration
// +build integration

package btc_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
)

type networkTest struct {
	testutil.BTCTestSuite
}

// TestGetNetworkInfo is test for GetNetworkInfo
func (nt *networkTest) TestGetNetworkInfo() {
	// GetNetworkInfo
	res, err := nt.BTC.GetNetworkInfo()
	nt.NoError(err)
	if err == nil {
		nt.T().Log(res)
	}
}

// TestGetBlockchainInfo is test for GetBlockchainInfo
func (nt *networkTest) TestBlockchainInfo() {
	// GetBlockchainInfo
	res, err := nt.BTC.GetBlockchainInfo()
	nt.NoError(err)
	if err == nil {
		nt.T().Log(res)
	}
}

func TestNetworkTestSuite(t *testing.T) {
	suite.Run(t, new(networkTest))
}
