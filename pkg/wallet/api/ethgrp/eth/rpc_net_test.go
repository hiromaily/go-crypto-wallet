package eth_test

import (
	"testing"

	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
)

// TestNetVersion is test for NetVersion
func TestNetVersion(t *testing.T) {
	// t.SkipNow()
	et := testutil.GetETH()

	netVersion, err := et.NetVersion()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("netVersion:", netVersion)
}

// TestNetListening is test for NetListening
func TestNetListening(t *testing.T) {
	// t.SkipNow()
	et := testutil.GetETH()

	isLitening, err := et.NetListening()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("isLitening:", isLitening)
}

// TestNetPeerCount is test for NetPeerCount
func TestNetPeerCount(t *testing.T) {
	// t.SkipNow()
	et := testutil.GetETH()

	peerCount, err := et.NetPeerCount()
	if err != nil {
		t.Fatal(err)
	}
	if peerCount != nil {
		t.Log("peerCount:", peerCount.Uint64())
	}
}
