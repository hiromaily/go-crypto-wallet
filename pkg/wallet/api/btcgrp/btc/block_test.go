//go:build integration
// +build integration

package btc_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
)

type blockTest struct {
	testutil.BTCTestSuite
}

// TestGetBlockCount is test for GetBlockCount
func (bt *blockTest) TestGetBlockCount(t *testing.T) {
	// GetBalance
	res, err := bt.BTC.GetBlockCount()
	bt.NoError(err)
	if err == nil {
		t.Log(res)
	}
}

func TestBlockTestSuite(t *testing.T) {
	suite.Run(t, new(blockTest))
}
