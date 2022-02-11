//go:build integration
// +build integration

package btc_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
)

type addressTest struct {
	testutil.BTCTestSuite
}

// TestGetAddressInfo is test for GetAddressInfo
func (at *addressTest) TestGetAddressInfo() {
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
		at.T().Run(tt.name, func(t *testing.T) {
			res, err := at.BTC.GetAddressInfo(tt.args.addr)
			at.Equal(tt.want.err, err)
			if err == nil {
				t.Log(res)
			}
		})
	}
}

// TestGetAddressesByLabel is test for GetAddressesByLabel
func (at *addressTest) TestGetAddressesByLabel() {
	type args struct {
		labelName string
	}
	type want struct{}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "happy path",
			args: args{"client"},
			want: want{},
		},
	}
	for _, tt := range tests {
		at.T().Run(tt.name, func(t *testing.T) {
			got, err := at.BTC.GetAddressesByLabel(tt.args.labelName)
			at.NoError(err)
			if err == nil {
				t.Log(got)
			}
		})
	}
}

// TestValidateAddress is test for ValidateAddress
func (at *addressTest) TestValidateAddress() {
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
		at.T().Run(tt.name, func(t *testing.T) {
			got, err := at.BTC.ValidateAddress(tt.args.addr)
			at.Equal(tt.want.isErr, err != nil)
			if err == nil {
				t.Log(got)
			}
		})
	}
}

func TestAddressTestSuite(t *testing.T) {
	suite.Run(t, new(addressTest))
}
