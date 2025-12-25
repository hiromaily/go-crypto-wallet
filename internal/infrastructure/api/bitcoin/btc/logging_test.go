//go:build integration
// +build integration

package btc_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
)

type loggingTest struct {
	testutil.BTCTestSuite
}

// TestLogging is test for GetLogging
func (lt *loggingTest) TestLogging() {
	// Logging
	res, err := lt.BTC.Logging()
	lt.NoError(err)
	if err == nil {
		lt.T().Log(res)
	}
}

func TestLoggingTestSuite(t *testing.T) {
	suite.Run(t, new(loggingTest))
}
