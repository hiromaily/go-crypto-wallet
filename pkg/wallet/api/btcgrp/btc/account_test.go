//go:build integration
// +build integration

package btc_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
)

type accountTest struct {
	testutil.BTCTestSuite
}

// TestGetAccount is test for GetAccount
func (at *accountTest) TestGetAccount() {
	type args struct {
		addr string
	}
	type want struct {
		account string
		err     error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "happy path",
			args: args{"2N4TcHSCteXwiF2dj8SQijj3w2HieR4x6r5"},
			want: want{"deposit", nil},
		},
		{
			name: "happy path",
			args: args{"2N6DcSuPo8NoLrCPTSqrwjnuLYoN7xDMSzX"},
			want: want{"payment", nil},
		},
	}

	for _, tt := range tests {
		at.T().Run(tt.name, func(t *testing.T) {
			res, err := at.BTC.GetAccount(tt.args.addr)
			at.Equal(tt.want.account, res)
			at.Equal(tt.want.err, err)
			if err == nil {
				t.Log(res)
			}
		})
	}
}

func TestAccountTestSuite(t *testing.T) {
	suite.Run(t, new(accountTest))
}
