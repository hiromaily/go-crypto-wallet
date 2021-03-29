package xrp_test

import (
	"testing"

	"github.com/bookerzzz/grok"

	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
)

// TestGenerateAddress is test for GenerateAddress
func TestGenerateAddress(t *testing.T) {
	// t.SkipNow()
	xr := testutil.GetXRP()

	addressInfo, err := xr.GenerateAddress()
	if err != nil {
		t.Fatal(err)
	}
	grok.Value(addressInfo)
}

// TestGenerateXAddress is test for GenerateXAddress
func TestGenerateXAddress(t *testing.T) {
	// t.SkipNow()
	xr := testutil.GetXRP()

	addressInfo, err := xr.GenerateXAddress()
	if err != nil {
		t.Fatal(err)
	}
	grok.Value(addressInfo)
}

// TestIsValidAddress is test for IsValidAddress
func TestIsValidAddress(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			// PrepareTransaction
			accountInfo, err := xr.IsValidAddress(tt.args.address)
			if err != nil {
				t.Fatal(err)
			}
			grok.Value(accountInfo)
		})
	}
}
