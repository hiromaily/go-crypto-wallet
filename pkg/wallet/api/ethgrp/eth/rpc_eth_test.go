package eth_test

import (
	"testing"

	"github.com/hiromaily/go-bitcoin/pkg/testutil"
)

// TestSyncing is test for Syncing
func TestSyncing(t *testing.T) {
	//t.SkipNow()
	et := testutil.GetETH()

	res, isSyncing, err := et.Syncing()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("resMap:", res)
	t.Log("isSyncing:", isSyncing)
}

// TestProtocolVersion is test for ProtocolVersion
func TestProtocolVersion(t *testing.T) {
	//t.SkipNow()
	et := testutil.GetETH()

	protocolVer, err := et.ProtocolVersion()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ProtocolVersion:", protocolVer)
}

