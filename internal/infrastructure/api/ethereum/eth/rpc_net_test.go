//go:build integration
// +build integration

package eth_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
)

type netTest struct {
	testutil.ETHTestSuite
}

// TestNetVersion is test for NetVersion
func (nt *netTest) TestNetVersion() {
	netVersion, err := nt.ETH.NetVersion()
	nt.NoError(err)
	nt.T().Log("netVersion:", netVersion)
}

// TestNetListening is test for NetListening
func (nt *netTest) TestNetListening() {
	isListening, err := nt.ETH.NetListening()
	nt.NoError(err)
	nt.T().Log("isListening:", isListening)
}

// TestNetPeerCount is test for NetPeerCount
func (nt *netTest) TestNetPeerCount() {
	peerCount, err := nt.ETH.NetPeerCount()
	nt.NoError(err)
	if err == nil {
		nt.T().Log("peerCount:", peerCount.Uint64())
	}
}

func TestNetTestSuite(t *testing.T) {
	suite.Run(t, new(netTest))
}
