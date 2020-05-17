package eth_test

import (
	"github.com/hiromaily/go-bitcoin/pkg/testutil"
	"testing"
)

// TestClientVersion is test for ClientVersion
func TestClientVersion(t *testing.T) {
	//t.SkipNow()
	et := testutil.GetETH()

	clientVersion, err := et.ClientVersion()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("clientVersion:", clientVersion)
}

// TestSHA3 is test for SHA3
func TestSHA3(t *testing.T) {
	//t.SkipNow()
	et := testutil.GetETH()
	data := "0x68656c6c6f20776f726c64"

	res, err := et.SHA3(data)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("response of SHA3:", res)
}
