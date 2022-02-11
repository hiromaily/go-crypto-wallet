//go:build integration
// +build integration

package xrp_test

import (
	"testing"

	"github.com/bookerzzz/grok"
	"github.com/stretchr/testify/suite"

	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
)

type addressTest struct {
	testutil.XRPTestSuite
}

// TestGenerateAddress is test for GenerateAddress
func (at *addressTest) TestGenerateAddress() {
	addressInfo, err := at.XRP.GenerateAddress()
	at.NoError(err)
	grok.Value(addressInfo)
}

// TestGenerateXAddress is test for GenerateXAddress
func (at *addressTest) TestGenerateXAddress() {
	addressInfo, err := at.XRP.GenerateXAddress()
	at.NoError(err)
	grok.Value(addressInfo)
}

// TestIsValidAddress is test for IsValidAddress
func (at *addressTest) TestIsValidAddress() {
	type args struct {
		address string
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
			name: "happy path 1",
			args: args{
				address: "XV9StxHCQ5meDLmRkw2ifV97iy7KiSW22Aku1D4UKqRKXwR",
			},
			want: want{false},
		},
		{
			name: "happy path 2",
			args: args{
				address: "rss1EZUwTCPZSTyJiDKvhBfCXjTxffcArZ",
			},
			want: want{false},
		},
		{
			name: "happy path 3",
			args: args{
				address: "X7vq1EiQAv1K4miEEqJLWcwsbbTENyVZLc96rGUd8XJSX7C",
			},
			want: want{false},
		},
		{
			name: "happy path 4",
			args: args{
				address: "r94FEwKytUn6qf4hxTL1wTeMP3sQfLGac9",
			},
			want: want{false},
		},
		{
			name: "wrong address",
			args: args{
				address: "0xabc12345",
			},
			want: want{true},
		},
	}
	for _, tt := range tests {
		at.T().Run(tt.name, func(t *testing.T) {
			accountInfo, err := at.XRP.IsValidAddress(tt.args.address)
			at.Equal(tt.want.isErr, err != nil)
			if err != nil {
				grok.Value(accountInfo)
			}
		})
	}
}

func TestAddressTestSuite(t *testing.T) {
	suite.Run(t, new(addressTest))
}
