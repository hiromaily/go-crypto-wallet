//go:build integration

package eth_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/eth/testutil"
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

// TestSuggestGasTipCap is test for SuggestGasTipCap
func (egt *ethGasTest) TestSuggestGasTipCap() {
	ctx := context.Background()
	tip, err := egt.ETH.SuggestGasTipCap(ctx)
	egt.NoError(err)
	egt.NotNil(tip)
	egt.T().Log("suggestGasTipCap:", tip)
}

func TestEthGasTestSuite(t *testing.T) {
	suite.Run(t, new(ethGasTest))
}
