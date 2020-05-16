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

// TestCoinbase is test for Coinbase
func TestCoinbase(t *testing.T) {
	//t.SkipNow()
	et := testutil.GetETH()

	addr, err := et.Coinbase()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("coinbase address:", addr)
}

// TestAccounts is test for Accounts
func TestAccounts(t *testing.T) {
	//t.SkipNow()
	et := testutil.GetETH()

	accounts, err := et.Accounts()
	if err != nil {
		t.Fatal(err)
	}
	for _, account := range accounts {
		t.Log("address:", account)
	}
}

// TestBlockNumber is test for BlockNumber
func TestBlockNumber(t *testing.T) {
	//t.SkipNow()
	et := testutil.GetETH()

	blockNum, err := et.BlockNumber()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("blockNumber:", blockNum)
}
