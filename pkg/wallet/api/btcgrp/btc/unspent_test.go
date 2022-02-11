//go:build integration
// +build integration

package btc_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/btcgrp/btc"
)

type unspentTest struct {
	testutil.BTCTestSuite
}

func findUnspentListID(unspentList []btc.ListUnspentResult, txid string) bool {
	for _, unspent := range unspentList {
		if unspent.TxID == txid {
			return true
		}
	}
	return false
}

// TestListUnspent is test for ListUnspent
func (ut *unspentTest) TestListUnspent() {
	// ListUnspent
	res, err := ut.BTC.ListUnspent(6)
	ut.NoError(err)
	if err == nil {
		ut.T().Log(res)
	}
}

// TestListUnspentByAccount is test for ListUnspentByAccount
func (ut *unspentTest) TestListUnspentByAccount() {
	// ListUnspentByAccount
	res, err := ut.BTC.ListUnspentByAccount(account.AccountTypePayment, ut.BTC.ConfirmationBlock())
	ut.NoError(err)
	if err == nil {
		ut.T().Log(res)
	}
}

// TestLockUnspent is test for LockUnspent, UnlockSpent
func (ut *unspentTest) TestLockUnspent() {
	// get unspent list
	listUnspent, err := ut.BTC.ListUnspent(ut.BTC.ConfirmationBlock())
	ut.NoError(err)

	// use one of ListUnspentResult
	if len(listUnspent) == 0 {
		ut.T().Log("unspent list is required for this test")
		return
	}
	// call LockUnspent to lock
	err = ut.BTC.LockUnspent(&listUnspent[0])
	ut.NoError(err)
	targetTxID := listUnspent[0].TxID

	// get unspent list again
	listUnspent, err = ut.BTC.ListUnspent(ut.BTC.ConfirmationBlock())
	ut.NoError(err)

	// check that tx is locked or not
	ut.False(findUnspentListID(listUnspent, targetTxID), "LockUnspent() failed to lock target unspent")

	// call UnlockUnspent to unlock
	err = ut.BTC.UnlockUnspent()
	ut.NoError(err)

	// get unspent list again
	listUnspent, err = ut.BTC.ListUnspent(ut.BTC.ConfirmationBlock())
	ut.NoError(err)

	// check unlocked or not
	ut.True(findUnspentListID(listUnspent, targetTxID), "UnlockUnspent() failed to unlock")
}

func TestUnspentTestSuite(t *testing.T) {
	suite.Run(t, new(unspentTest))
}
