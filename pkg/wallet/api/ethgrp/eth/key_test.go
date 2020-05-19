package eth_test

import (
	"testing"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/testutil"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api/ethgrp/eth"
)

// TestGetPrivKey is test for GetPrivKey
func TestGetPrivKey(t *testing.T) {
	et := testutil.GetETH()

	type args struct {
		addr string
		acnt account.AccountType
	}
	type want struct {
		isErr bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "happy path",
			args: args{
				addr: "0xe52307Deb1a7dC3985D2873b45AE23b91D57a36d",
				acnt: account.AccountTypeClient,
			},
			want: want{false},
		},
		{
			name: "wrong address",
			args: args{
				addr: "0xe52307Deb1a7dC3985D2873b45AE23b91Daaaaaa",
				acnt: account.AccountTypeClient,
			},
			want: want{true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prikey, err := et.GetPrivKey(tt.args.addr, eth.Password, tt.args.acnt)
			if (err == nil) == tt.want.isErr {
				t.Errorf("readPrivKey() = %v, want error = %v", err, tt.want.isErr)
				return
			}
			if err == nil && prikey == nil {
				t.Error("prikey is nil")
			}
		})
	}
}
