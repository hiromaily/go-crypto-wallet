//go:build integration

package eth_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/eth/testutil"
)

type netTest struct {
	testutil.ETHTestSuite
}

// TestNetVersion is test for NetVersion
func (nt *netTest) TestNetVersion() {
	ctx := context.Background()
	netVersion, err := nt.ETH.NetVersion(ctx)
	nt.NoError(err)
	nt.T().Log("netVersion:", netVersion)
}

// TestNetListening is test for NetListening
func (nt *netTest) TestNetListening() {
	ctx := context.Background()
	isListening, err := nt.ETH.NetListening(ctx)
	nt.NoError(err)
	nt.T().Log("isListening:", isListening)
}

// TestNetPeerCount is test for NetPeerCount
func (nt *netTest) TestNetPeerCount() {
	ctx := context.Background()
	peerCount, err := nt.ETH.NetPeerCount(ctx)
	nt.NoError(err)
	if err == nil {
		nt.T().Log("peerCount:", peerCount.Uint64())
	}
}

func TestNetTestSuite(t *testing.T) {
	suite.Run(t, new(netTest))
}
