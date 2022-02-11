//go:build integration
// +build integration

package eth_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/ethgrp/eth"
)

type personalTest struct {
	testutil.ETHTestSuite
}

// TestImportRawKey is test for ImportRawKey
func (pt *personalTest) TestImportRawKey() {
	pw := eth.Password

	type args struct {
		key string
	}
	type want struct {
		address string
		isErr   bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "happy path 1",
			args: args{key: "0x11e9c4671932ae288e50c40a9aecc09b28788f519812e0dfc7fe9b5119074b7b"},
			want: want{"", false},
		},
		{
			name: "happy path 2",
			args: args{key: "0xa8dc1250a7801a660893e9cfa7440aa96d47e6f768afd5f244f9876880f9fee4"},
			want: want{"", false},
		},
		{
			name: "wrong key",
			args: args{key: "0x0913da3091587e34320e5001a0fd582b09e6eed9eb6679d84f754318488bfd68"},
			want: want{"", true},
		},
	}
	for _, tt := range tests {
		pt.T().Run(tt.name, func(t *testing.T) {
			addr, err := pt.ETH.ImportRawKey(tt.args.key, pw)
			pt.Equal(tt.want.isErr, err != nil)
			if err == nil {
				t.Log("address:", addr)
			}
		})
	}
}

// TestListAccounts is test for ListAccounts
func (pt *personalTest) TestListAccounts() {
	addrs, err := pt.ETH.ListAccounts()
	pt.NoError(err)
	for _, addr := range addrs {
		pt.T().Log("address:", addr)
	}
}

// TestNewAccount is test for ListAccounts
func (pt *personalTest) TestNewAccount() {
	addr, err := pt.ETH.NewAccount(eth.Password, account.AccountTypeClient)
	pt.NoError(err)
	pt.T().Log("address:", addr)
}

func (pt *personalTest) TestLockAccount() {
	addr := "0x852d4ae6bfa5ae9d44d3ac03122674bcb32a0861"

	// unlock
	isUnlocked, err := pt.ETH.UnlockAccount(addr, eth.Password, uint64(1))
	pt.NoError(err)
	if !isUnlocked {
		pt.T().Error("address is not unlocked")
		return
	}

	// lock
	err = pt.ETH.LockAccount(addr)
	pt.NoError(err)
}

func TestPersonalTestSuite(t *testing.T) {
	suite.Run(t, new(personalTest))
}
