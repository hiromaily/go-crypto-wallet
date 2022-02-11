//go:build integration
// +build integration

package btc_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/key"
)

type labelTest struct {
	testutil.BTCTestSuite
}

// TestSetLabel is test for SetLabel
// Note: this test will contaminate wallet.dat
func (lt *labelTest) TestSetLabel() {
	// create address
	_, pubKey, err := key.GenerateWIF(lt.BTC.GetChainConf())
	lt.NoError(err)

	// import address
	err = lt.BTC.ImportAddressWithoutReScan(pubKey)
	lt.NoError(err)

	// set label
	err = lt.BTC.SetLabel(pubKey, account.AccountTypeTest.String())
	lt.NoError(err)

	// check addr
	res, err := lt.BTC.GetAddressInfo(pubKey)
	lt.NoError(err)
	if err == nil {
		lt.T().Log(res.Labels)
	}
}

func TestLabelTestSuite(t *testing.T) {
	suite.Run(t, new(labelTest))
}
