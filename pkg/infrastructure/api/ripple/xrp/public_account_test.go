//go:build integration
// +build integration

package xrp_test

import (
	"testing"

	"github.com/bookerzzz/grok"
	"github.com/stretchr/testify/suite"

	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
)

type publicAccountTest struct {
	testutil.XRPTestSuite
}

// TestAccountChannels is test for AccountChannels
func (pat *publicAccountTest) TestAccountChannels() {
	res, err := pat.XRP.AccountChannels("rN7n7otQDd6FczFgLdSqtcsAUxDkw6fzRH", "rf1BiGeXwwQoi8Z2ueFYTEXSwuJYfV2Jpn")
	pat.NoError(err)
	pat.T().Log("accountChannels:", res)
	grok.Value(res)
}

// TestAccountInfo is test for AccountInfo
func (pat *publicAccountTest) TestAccountInfo() {
	res, err := pat.XRP.AccountInfo("rss1EZUwTCPZSTyJiDKvhBfCXjTxffcArZ")
	pat.NoError(err)
	pat.T().Log("accountInfo:", res)
	grok.Value(res)
}

func TestPublicAccountTestSuite(t *testing.T) {
	suite.Run(t, new(publicAccountTest))
}
