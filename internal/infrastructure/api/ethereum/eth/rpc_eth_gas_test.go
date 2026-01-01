//go:build integration
// +build integration

package eth_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/ethereum/testutil"
)

type ethGasTest struct {
	testutil.ETHTestSuite
}

// TestGasPrice is test for GasPrice
func (egt *ethGasTest) TestGasPrice() {
	ctx := context.Background()
	price, err := egt.ETH.GasPrice(ctx)
	egt.NoError(err)
	egt.T().Log("gasPrice:", price)
}

func TestEthGasTestSuite(t *testing.T) {
	suite.Run(t, new(ethGasTest))
}
