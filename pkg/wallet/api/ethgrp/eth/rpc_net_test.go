package eth_test

import (
	"testing"

	"github.com/hiromaily/go-bitcoin/pkg/testutil"
)

// TestNetVersion is test for NetVersion
func TestNetVersion(t *testing.T) {
	//t.SkipNow()
	et := testutil.GetETH()

	netVersion, err := et.NetVersion()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("netVersion:", netVersion)
}
