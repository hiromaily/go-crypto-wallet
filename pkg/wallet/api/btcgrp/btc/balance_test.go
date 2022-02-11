//go:build integration
// +build integration

package btc_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
)

type balanceTest struct {
	testutil.BTCTestSuite
}

// TestListAccounts is test for ListAccounts
func (bt *balanceTest) TestGetBalance() {
	// GetBalance
	res, err := bt.BTC.GetBalance()
	bt.NoError(err)
	if err == nil {
		bt.T().Log(res)
	}
}

// TestGetBalanceByAccount is test for GetBalanceByAccount
func (bt *balanceTest) TestGetBalanceByAccount() {
	type args struct {
		account account.AccountType
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
			args: args{account.AccountTypeClient},
			want: want{false},
		},
		{
			name: "happy path",
			args: args{account.AccountTypeDeposit},
			want: want{false},
		},
		{
			name: "happy path",
			args: args{account.AccountTypePayment},
			want: want{false},
		},
	}

	for _, tt := range tests {
		bt.T().Run(tt.name, func(t *testing.T) {
			got, err := bt.BTC.GetBalanceByAccount(tt.args.account, bt.BTC.ConfirmationBlock())
			if (err != nil) != tt.want.isErr {
				t.Errorf("GetBalanceByAccount() = %v, isErr %v", err, tt.want.isErr)
			}
			if err == nil {
				t.Log(got)
			}
		})
	}
}

func TestBalanceTestSuite(t *testing.T) {
	suite.Run(t, new(balanceTest))
}
