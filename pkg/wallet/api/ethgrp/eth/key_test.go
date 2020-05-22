package eth_test

import (
	"testing"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/ethgrp/eth"
)

// TestGetPrivKey is test for GetPrivKey
// Note: keydir in config must be fullpath when testing
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
				addr: "0xd4EC46122b3f0afc0287144Adcca5d65B22B799c",
				acnt: account.AccountTypeClient,
			},
			want: want{false},
		},
		{
			name: "wrong address",
			args: args{
				addr: "0x5357135e0D3CbBD37cFCeb9F06257Bb133548Exx",
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
