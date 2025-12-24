//go:build integration
// +build integration

package xrp_test

import (
	"testing"

	"github.com/bookerzzz/grok"
	"github.com/stretchr/testify/suite"

	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
)

type adminKeygenTest struct {
	testutil.XRPTestSuite
}

// TestValidationCreate is test for ValidationCreate
func (akt *adminKeygenTest) TestValidationCreate() {
	type args struct {
		secret string
	}
	type want struct{}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "happy path 1",
			args: args{"ssCATR7CBvn4GLd1UuU2bqqQffHki"},
			want: want{},
		},
		{
			name: "happy path 2",
			args: args{"BAWL MAN JADE MOON DOVE GEM SON NOW HAD ADEN GLOW TIRE"},
			want: want{},
		},
	}

	for _, tt := range tests {
		akt.T().Run(tt.name, func(t *testing.T) {
			res, err := akt.XRP.ValidationCreate(tt.args.secret)
			akt.NoError(err)
			if err == nil {
				t.Log("ValidationCreate:", res)
				grok.Value(res)
			}
		})
	}
}

// TestWalletPropose is test for WalletPropose
func (akt *adminKeygenTest) TestWalletPropose() {
	type args struct {
		passphrase string
	}
	type want struct{}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "happy path 1",
			args: args{"password1"},
			want: want{},
		},
		{
			name: "happy path 2",
			args: args{"foobar"},
			want: want{},
		},
		{
			name: "happy path 3",
			args: args{"0x931D387731bBbC988B312206c74F77D004D6B84b"},
			want: want{},
		},
	}

	for _, tt := range tests {
		akt.T().Run(tt.name, func(t *testing.T) {
			res, err := akt.XRP.WalletPropose(tt.args.passphrase)
			akt.NoError(err)
			if err == nil {
				t.Log("WalletPropose:", res)
				grok.Value(res)
			}
		})
	}
}

func TestAdminKeygenTestSuite(t *testing.T) {
	suite.Run(t, new(adminKeygenTest))
}
