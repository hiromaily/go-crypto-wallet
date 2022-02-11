//go:build integration
// +build integration

package btc_test

import (
	"testing"

	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
)

// TODO: mock of bitcoin interface is required to test

// TestGetAddressInfo is test for GetAddressInfo
func TestGetAddressInfo(t *testing.T) {
	// t.SkipNow()
	bc := testutil.GetBTC()

	type args struct {
		addr string
	}
	type want struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "happy path",
			args: args{"mvTRCKpKVUUv3QgMEn838xXDDZS5SSEhnj"},
			want: want{nil},
		},
		{
			name: "happy path",
			args: args{"n3f97rFX5p1vbwKqkdhjT6QjaiqBw6TfxQ"},
			want: want{nil},
		},
		{
			name: "happy path",
			args: args{"n3f97rFX5p1vbwKqkdhjT6QjaiqBw6TfxQ"},
			want: want{nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := bc.GetAddressInfo(tt.args.addr)
			if err != tt.want.err {
				t.Errorf("GetAddressInfo() = %v, want %v", err, tt.want.err)
			}
			t.Log(res)
		})
	}
	// bc.Close()()()
}

// TestGetAddressesByLabel is test for GetAddressesByLabel
func TestGetAddressesByLabel(t *testing.T) {
	// t.SkipNow()
	bc := testutil.GetBTC()
	type args struct {
		labelName string
	}
	type want struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "happy path",
			args: args{"client"},
			want: want{nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := bc.GetAddressesByLabel(tt.args.labelName); err != tt.want.err {
				t.Errorf("GetAddressesByLabel() = %v, want %v", got, tt.want)
			} else {
				t.Log(got)
			}
		})
	}
	// bc.Close()()()
}

// TestValidateAddress is test for ValidateAddress
func TestValidateAddress(t *testing.T) {
	bc := testutil.GetBTC()

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
			args: args{"2NFXSXxw8Fa6P6CSovkdjXE6UF4hupcTHtr"},
			want: want{false},
		},
		{
			name: "happy path",
			args: args{"2NDGkbQTwg2v1zP6yHZw3UJhmsBh9igsSos"},
			want: want{false},
		},
		{
			name: "wrong address",
			args: args{"4VHGkbQTGg2vN5P6yHZw3UJhmsBh9igsSos"},
			want: want{true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := bc.ValidateAddress(tt.args.addr); (err != nil) != tt.want.isErr {
				t.Errorf("ValidateAddress() = %v, isErr %v", err, tt.want.isErr)
			} else {
				t.Log(got)
			}
		})
	}

	// bc.Close()
}
