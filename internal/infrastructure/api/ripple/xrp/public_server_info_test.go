//go:build integration
// +build integration

package xrp_test

import (
	"context"
	"testing"

	"github.com/bookerzzz/grok"
	"github.com/stretchr/testify/suite"

	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
)

type publicServerInfoTest struct {
	testutil.XRPTestSuite
}

// TestServerInfo is test for ServerInfo
func (pst *publicServerInfoTest) TestServerInfo() {
	ctx := context.Background()
	res, err := pst.XRP.ServerInfo(ctx)
	pst.NoError(err)
	pst.T().Log("ServerInfo:", res)
	grok.Value(res)
}

func TestPublicServerInfoTestSuite(t *testing.T) {
	suite.Run(t, new(publicServerInfoTest))
}
