package xrp_test

import (
	"strings"
	"testing"

	"github.com/bookerzzz/grok"

	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
)

// TestTransaction is test for sequential transaction
func TestTransaction(t *testing.T) {
	//t.SkipNow()
	xr := testutil.GetXRP()

	type args struct {
		sernderAccount  string
		senderSecret    string
		receiverAccount string
		amount          float64
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
				sernderAccount:  "rHMbddzSv3nvKx3dhQ2xThH7BoXsjsdbA9",
				senderSecret:    "sn4MKv1CqaYAJ4LJRkin6q9xJPzQu",
				receiverAccount: "raWG2eo1tEXwN4HtGFJCagvukC2nBuiHxC",
				amount:          100,
			},
			want: want{false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// PrepareTransaction
			txJSON, _, err := xr.PrepareTransaction(tt.args.sernderAccount, tt.args.receiverAccount, tt.args.amount)
			if err != nil {
				t.Fatal(err)
			}
			grok.Value(txJSON)

			// SingTransaction
			txID, txBlob, err := xr.SignTransaction(txJSON, tt.args.senderSecret)
			if err != nil {
				t.Fatal(err)
			}
			t.Log("txID: ", txID)
			t.Log("txBlob: ", txBlob)

			// SendTransaction
			sentTx, earlistLedgerVersion, err := xr.SubmitTransaction(txBlob)
			if err != nil {
				t.Fatal(err)
			}
			t.Log("earlistLedgerVersion: ", earlistLedgerVersion)
			grok.Value(sentTx)
			if strings.Contains(sentTx.ResultCode, "UNFUNDED_PAYMENT") {
				t.Errorf("fail to call SubmitTransaction. resultCode: %s, resultMessage: %s", sentTx.ResultCode, sentTx.ResultMessage)
				return
			}

			// validate transaction
			ledgerVer, err := xr.WaitValidation(sentTx.TxJSON.LastLedgerSequence)
			if err != nil {
				t.Fatal(err)
			}
			t.Log("currentLedgerVersion: ", ledgerVer)

			// get transaction info
			grok.Value(txID)
			grok.Value(earlistLedgerVersion)
			txInfo, err := xr.GetTransaction(txID, earlistLedgerVersion)
			if err != nil {
				t.Fatal(err)
			}
			grok.Value(txInfo)

			//TODO: sender account info
			//TODO: receiver account info

		})
	}

}

// TestGetTransaction is test for GetTransaction
func TestGetTransaction(t *testing.T) {
	//t.SkipNow()
	xr := testutil.GetXRP()

	type args struct {
		txID                 string
		earlistLedgerVersion uint64
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
				txID:                 "AA672DB687B52DA733FDD211D435BB8632E2E4B471C282ABA912E8D799D16B6B",
				earlistLedgerVersion: 8007165,
			},
			want: want{false},
		},
		{
			name: "wrong txID",
			args: args{
				txID:                 "AA672DB687B52DA733FDD211D435BB8632E2E4B471C282ABA912E8D799D16B6Babcde",
				earlistLedgerVersion: 8007165,
			},
			want: want{false},
		},
		{
			name: "wrong ledger version",
			args: args{
				txID:                 "AA672DB687B52DA733FDD211D435BB8632E2E4B471C282ABA912E8D799D16B6B",
				earlistLedgerVersion: 99999999999999,
			},
			want: want{false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txInfo, err := xr.GetTransaction(tt.args.txID, tt.args.earlistLedgerVersion)
			if err != nil {
				t.Error(err)
			}
			if txInfo != nil {
				t.Log(txInfo)
			}
		})
	}

}
