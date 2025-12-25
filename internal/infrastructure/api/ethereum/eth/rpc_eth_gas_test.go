//go:build integration
// +build integration

package eth_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
)

type ethGasTest struct {
	testutil.ETHTestSuite
}

// TestGasPrice is test for GasPrice
func (egt *ethGasTest) TestGasPrice() {
	price, err := egt.ETH.GasPrice()
	egt.NoError(err)
	egt.T().Log("gasPrice:", price)
}

func TestEthGasTestSuite(t *testing.T) {
	suite.Run(t, new(ethGasTest))
}
