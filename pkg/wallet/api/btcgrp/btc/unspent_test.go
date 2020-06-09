package btc_test

import (
	"testing"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/btcgrp/btc"
)

func findUnspentListID(unspentList []btc.ListUnspentResult, txid string) bool {
	for _, unspent := range unspentList {
		if unspent.TxID == txid {
			return true
		}
	}
	return false
}

// TestListUnspent is test for ListUnspent
func TestListUnspent(t *testing.T) {
	//t.SkipNow()
	bc := testutil.GetBTC()

	// ListUnspent
	if res, err := bc.ListUnspent(6); err != nil {
		t.Errorf("fail to call ListUnspent(): %v", err)
	} else {
		t.Log(res)
	}

	//bc.Close()
}

// TestListUnspentByAccount is test for ListUnspentByAccount
func TestListUnspentByAccount(t *testing.T) {
	//t.SkipNow()
	bc := testutil.GetBTC()

	// ListUnspentByAccount
	if res, err := bc.ListUnspentByAccount(account.AccountTypePayment, bc.ConfirmationBlock()); err != nil {
		t.Errorf("fail to call ListUnspent(): %v", err)
	} else {
		t.Log(res)
	}

	//bc.Close()
}

// TestLockUnspent is test for LockUnspent, UnlockSpent
func TestLockUnspent(t *testing.T) {
	//t.SkipNow()
	bc := testutil.GetBTC()

	// get unspent list
	listUnspent, err := bc.ListUnspent(bc.ConfirmationBlock())
	if err != nil {
		t.Fatalf("fail to call ListUnspent(): %v", err)
	}
	//use one of ListUnspentResult
	if len(listUnspent) == 0 {
		t.Log("unspent list is required for this test")
		return
	}
	// call LockUnspent to lock
	if err := bc.LockUnspent(&listUnspent[0]); err != nil {
		t.Fatalf("fail to call LockUnspent(): %v", err)
	}
	targetTxID := listUnspent[0].TxID

	// get unspent list again
	listUnspent, err = bc.ListUnspent(bc.ConfirmationBlock())
	if err != nil {
		t.Fatalf("fail to call ListUnspent(): %v", err)
	}

	// check that tx is locked or not
	if findUnspentListID(listUnspent, targetTxID) {
		t.Error("LockUnspent() fail to lock target unspent")
	}

	// call UnlockUnspent to unlock
	if err = bc.UnlockUnspent(); err != nil {
		t.Fatalf("fail to call UnlockUnspent(): %v", err)
	}

	// get unspent list again
	listUnspent, err = bc.ListUnspent(bc.ConfirmationBlock())
	if err != nil {
		t.Fatalf("fail to call ListUnspent(): %v", err)
	}

	// check unlocked or not
	if !findUnspentListID(listUnspent, targetTxID) {
		t.Error("UnlockUnspent() fail to unlock")
	}

	//bc.Close()
}
