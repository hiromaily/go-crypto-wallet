//go:build integration
// +build integration

package xrp_test

import (
	"testing"

	"github.com/bookerzzz/grok"
	"github.com/stretchr/testify/suite"

	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
)

type accountTest struct {
	testutil.XRPTestSuite
}

// TestGetAccountInfo is test for GetAccountInfo
func (at *accountTest) TestGetAccountInfo() {
	type args struct {
		address string
	}
	type want struct{}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "happy path 1",
			args: args{
				address: "rss1EZUwTCPZSTyJiDKvhBfCXjTxffcArZ",
			},
			want: want{},
		},
		{
			name: "happy path 2",
			args: args{
				address: "rNajCSDNXZLCioutY6xk4r3mYWMGYAorcN",
			},
			want: want{},
		},
	}
	for _, tt := range tests {
		at.T().Run(tt.name, func(t *testing.T) {
			// PrepareTransaction
			accountInfo, err := at.XRP.GetAccountInfo(tt.args.address)
			at.NoError(err)
			grok.Value(accountInfo)
		})
	}
}

func TestAccountTestSuite(t *testing.T) {
	suite.Run(t, new(accountTest))
}
