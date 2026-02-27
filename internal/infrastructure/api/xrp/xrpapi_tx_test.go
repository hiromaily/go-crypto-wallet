//go:build integration

package xrp_test

import (
	"context"
	"strings"
	"testing"

	"github.com/bookerzzz/grok"
	"github.com/stretchr/testify/suite"

	dtoxrp "github.com/hiromaily/go-crypto-wallet/internal/application/dto/xrp"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/xrp/testutil"
)

type apiTxTest struct {
	testutil.XRPTestSuite
}

// TestTransaction is test for sequential transaction
func (att *apiTxTest) TestTransaction() {
	type args struct {
		sernderAccount  string
		senderSecret    string
		receiverAccount string
		amount          float64
		instructions    *dtoxrp.Instructions
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
				sernderAccount:  "rxcip8BgcUMbMr9QEm15WMDUdG7oGEwT8",
				senderSecret:    "sswVaSDUNnLd5pB4F9oPT9u7jLf2X",
				receiverAccount: "rnkZMhbXQZ8GTfSihmdTqNUtvUAAqwkLWN",
				amount:          100,
				instructions: &dtoxrp.Instructions{
					MaxLedgerVersionOffset: 2,
				},
			},
			want: want{},
		},
	}

	for _, tt := range tests {
		att.T().Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			// PrepareTransaction
			txJSON, _, err := att.XRP.PrepareTransaction(
				ctx, tt.args.sernderAccount, tt.args.receiverAccount, tt.args.amount, tt.args.instructions,
			)
			att.NoError(err)
			grok.Value(txJSON)
			//- creating raw transaction
			// LastLedgerSequence:    8153687
			// Sequence:              8153441
			//- after sending
			// earlistLedgerVersion:             8153673
			// sentTx.TxJSON.LastLedgerSequence: 8153687

			// SingTransaction
			txID, txBlob, err := att.XRP.SignTransaction(ctx, txJSON, tt.args.senderSecret)
			att.NoError(err)

			// SendTransaction
			sentTx, earlistLedgerVersion, err := att.XRP.SubmitTransaction(ctx, txBlob)
			att.NoError(err)
			if strings.Contains(sentTx.ResultCode, "UNFUNDED_PAYMENT") {
				t.Errorf(
					"fail to call SubmitTransaction. resultCode: %s, resultMessage: %s",
					sentTx.ResultCode, sentTx.ResultMessage,
				)
				return
			}

			// validate transaction
			ledgerVer, err := att.XRP.WaitValidation(ctx, sentTx.TxJSON.LastLedgerSequence)
			att.NoError(err)
			t.Log("currentLedgerVersion: ", ledgerVer)

			// get transaction info
			txInfo, err := att.XRP.GetTransaction(ctx, txID, earlistLedgerVersion)
			att.NoError(err)
			t.Log("GetTransaction: ", txInfo)

			// TODO: sender account info
			// TODO: receiver account info
		})
	}
}

// TestGetTransaction is test for GetTransaction
func (att *apiTxTest) TestGetTransaction() {
	type args struct {
		txID                 string
		earlistLedgerVersion uint64
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
				txID:                 "AA672DB687B52DA733FDD211D435BB8632E2E4B471C282ABA912E8D799D16B6B",
				earlistLedgerVersion: 8007165,
			},
			want: want{},
		},
		{
			name: "wrong txID",
			args: args{
				txID:                 "AA672DB687B52DA733FDD211D435BB8632E2E4B471C282ABA912E8D799D16B6Babcde",
				earlistLedgerVersion: 8007165,
			},
			want: want{},
		},
		{
			name: "wrong ledger version",
			args: args{
				txID:                 "AA672DB687B52DA733FDD211D435BB8632E2E4B471C282ABA912E8D799D16B6B",
				earlistLedgerVersion: 99999999999999,
			},
			want: want{},
		},
	}
	for _, tt := range tests {
		att.T().Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			txInfo, err := att.XRP.GetTransaction(ctx, tt.args.txID, tt.args.earlistLedgerVersion)
			att.NoError(err)
			if err == nil {
				t.Log(txInfo)
			}
		})
	}
}

func TestAPITxTestSuite(t *testing.T) {
	suite.Run(t, new(apiTxTest))
}
