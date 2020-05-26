package eth_test

import (
	"testing"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/ethgrp/eth"
)

// TestImportRawKey is test for ImportRawKey
func TestImportRawKey(t *testing.T) {
	//t.SkipNow()
	et := testutil.GetETH()

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
		t.Run(tt.name, func(t *testing.T) {
			addr, err := et.ImportRawKey(tt.args.key, pw)
			if (err == nil) == tt.want.isErr {
				t.Errorf("ImportRawKey() = %v, want error = %v", err, tt.want.isErr)
				return
			}
			if err == nil {
				t.Log("address:", addr)
			}
		})
	}
}

// TestListAccounts is test for ListAccounts
func TestListAccounts(t *testing.T) {
	//t.SkipNow()
	et := testutil.GetETH()

	addrs, err := et.ListAccounts()
	if err != nil {
		t.Fatal(err)
	}
	for _, addr := range addrs {
		t.Log("address:", addr)
	}
}

// TestNewAccount is test for ListAccounts
func TestNewAccount(t *testing.T) {
	//t.SkipNow()
	et := testutil.GetETH()

	addr, err := et.NewAccount(eth.Password, account.AccountTypeClient)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("address:", addr)
}

func TestLockAccount(t *testing.T) {
	//t.SkipNow()
	et := testutil.GetETH()
	addr := "0x852d4ae6bfa5ae9d44d3ac03122674bcb32a0861"

	// unlock
	isUnlocked, err := et.UnlockAccount(addr, eth.Password, uint64(1))
	if err != nil {
		t.Fatal(err)
	}
	if !isUnlocked {
		t.Error("address is not unlocked")
		return
	}

	// lock
	err = et.LockAccount(addr)
	if err != nil {
		t.Fatal(err)
	}
}
