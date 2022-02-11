//go:build integration
// +build integration

package eth_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/ethgrp/eth"
)

type keyTest struct {
	testutil.ETHTestSuite
}

// TestGetPrivKey is test for GetPrivKey
// Note: keydir in config must be fullpath when testing
func (kt *keyTest) TestGetPrivKey() {
	type args struct {
		addr string
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
			},
			want: want{false},
		},
		{
			name: "wrong address",
			args: args{
				addr: "0x5357135e0D3CbBD37cFCeb9F06257Bb133548Exx",
			},
			want: want{true},
		},
	}

	for _, tt := range tests {
		kt.T().Run(tt.name, func(t *testing.T) {
			prikey, err := kt.ETH.GetPrivKey(tt.args.addr, eth.Password)
			kt.Equal(tt.want.isErr, err != nil)
			if err == nil && prikey == nil {
				t.Error("prikey is nil")
			}
		})
	}
}

func TestKeyTestSuite(t *testing.T) {
	suite.Run(t, new(keyTest))
}
