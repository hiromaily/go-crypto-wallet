//go:build integration
// +build integration

package xrp_test

import (
	"testing"

	"github.com/bookerzzz/grok"
	"github.com/stretchr/testify/suite"

	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/xrpgrp/xrp"
)

type txTest struct {
	testutil.XRPTestSuite
}

// TestCreateRawTransaction is test for CreateRawTransaction
func (txt *txTest) TestCreateRawTransaction() {
	type args struct {
		sernderAccount  string
		receiverAccount string
		amount          float64
		instructions    *xrp.Instructions
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
				sernderAccount:  "rKXvsrd5H6MQNVpYgdeffFYjfGq4VdDogd",
				receiverAccount: "rpBzBQ6aWJhuatJCkQgfE3VJT67ukBQopf",
				amount:          50,
				instructions: &xrp.Instructions{
					MaxLedgerVersionOffset: xrp.MaxLedgerVersionOffset,
				},
			},
			want: want{},
		},
		//{
		//	name: "happy path 2",
		//	args: args{
		//		sernderAccount:  "rEoRcMBfg7VUryw5xSyw883bXU74T8eoYj",
		//		receiverAccount: "raWG2eo1tEXwN4HtGFJCagvukC2nBuiHxC",
		//		amount:          0,
		//	},
		//	want: want{},
		//},
	}

	for _, tt := range tests {
		txt.T().Run(tt.name, func(t *testing.T) {
			// PrepareTransaction
			txJSON, _, err := txt.XRP.CreateRawTransaction(tt.args.sernderAccount, tt.args.receiverAccount, tt.args.amount, tt.args.instructions)
			txt.NoError(err)
			grok.Value(txJSON)
		})
	}
}

func TestTxTestSuite(t *testing.T) {
	suite.Run(t, new(txTest))
}
