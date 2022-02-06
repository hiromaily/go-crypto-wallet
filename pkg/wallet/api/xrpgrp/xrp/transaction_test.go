package xrp_test

import (
	"testing"

	"github.com/bookerzzz/grok"

	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/xrpgrp/xrp"
)

// TestCreateRawTransaction is test for CreateRawTransaction
func TestCreateRawTransaction(t *testing.T) {
	// t.SkipNow()
	xr := testutil.GetXRP()

	type args struct {
		sernderAccount  string
		receiverAccount string
		amount          float64
		instructions    *xrp.Instructions
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
				sernderAccount:  "rKXvsrd5H6MQNVpYgdeffFYjfGq4VdDogd",
				receiverAccount: "rpBzBQ6aWJhuatJCkQgfE3VJT67ukBQopf",
				amount:          50,
				instructions: &xrp.Instructions{
					MaxLedgerVersionOffset: xrp.MaxLedgerVersionOffset,
				},
			},
			want: want{false},
		},
		//{
		//	name: "happy path 2",
		//	args: args{
		//		sernderAccount:  "rEoRcMBfg7VUryw5xSyw883bXU74T8eoYj",
		//		receiverAccount: "raWG2eo1tEXwN4HtGFJCagvukC2nBuiHxC",
		//		amount:          0,
		//	},
		//	want: want{false},
		//},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// PrepareTransaction
			txJSON, _, err := xr.CreateRawTransaction(tt.args.sernderAccount, tt.args.receiverAccount, tt.args.amount, tt.args.instructions)
			if err != nil {
				t.Fatal(err)
			}
			grok.Value(txJSON)
		})
	}
}
