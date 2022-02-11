//go:build integration
// +build integration

package xrp_test

import (
	"testing"

	"github.com/bookerzzz/grok"

	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
)

// TestGetAccountInfo is test for GetAccountInfo
func TestGetAccountInfo(t *testing.T) {
	// t.SkipNow()
	xr := testutil.GetXRP()

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
				address: "rss1EZUwTCPZSTyJiDKvhBfCXjTxffcArZ",
			},
			want: want{false},
		},
		{
			name: "happy path 2",
			args: args{
				address: "rNajCSDNXZLCioutY6xk4r3mYWMGYAorcN",
			},
			want: want{false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// PrepareTransaction
			accountInfo, err := xr.GetAccountInfo(tt.args.address)
			if err != nil {
				t.Fatal(err)
			}
			grok.Value(accountInfo)
		})
	}
}
