//go:build integration
// +build integration

package xrp_test

import (
	"context"
	"testing"

	"github.com/bookerzzz/grok"
	"github.com/stretchr/testify/suite"

	dtoRipple "github.com/hiromaily/go-crypto-wallet/internal/application/dto/ripple"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/xrp/testutil"
	apixrpimpl "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/xrp/xrp"
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
		instructions    *dtoRipple.Instructions
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
				instructions: &dtoRipple.Instructions{
					MaxLedgerVersionOffset: apixrpimpl.MaxLedgerVersionOffset,
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
		// },
	}

	for _, tt := range tests {
		txt.T().Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			// PrepareTransaction
			txJSON, _, err := txt.XRP.CreateRawTransaction(
				ctx, tt.args.sernderAccount, tt.args.receiverAccount, tt.args.amount, tt.args.instructions,
			)
			txt.NoError(err)
			grok.Value(txJSON)
		})
	}
}

func TestTxTestSuite(t *testing.T) {
	suite.Run(t, new(txTest))
}
