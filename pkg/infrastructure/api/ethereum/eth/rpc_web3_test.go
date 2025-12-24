//go:build integration
// +build integration

package eth_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
)

type web3Test struct {
	testutil.ETHTestSuite
}

// TestClientVersion is test for ClientVersion
func (wt *web3Test) TestClientVersion() {
	clientVersion, err := wt.ETH.ClientVersion()
	wt.NoError(err)
	wt.T().Log("clientVersion:", clientVersion)
}

// TestSHA3 is test for SHA3
func (wt *web3Test) TestSHA3() {
	data := "0x68656c6c6f20776f726c64"

	res, err := wt.ETH.SHA3(data)
	wt.NoError(err)
	wt.T().Log("response of SHA3:", res)
}

func TestWeb3TestSuite(t *testing.T) {
	suite.Run(t, new(web3Test))
}
