package xrp_test

import (
	"testing"

	"github.com/bookerzzz/grok"

	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
)

// TestTransaction is test for PrepareTransaction
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
				sernderAccount:  "rNajCSDNXZLCioutY6xk4r3mYWMGYAorcN",
				senderSecret:    "sn2S8Gdfr9uYsX81nsjetf2Rx66Sv",
				receiverAccount: "rss1EZUwTCPZSTyJiDKvhBfCXjTxffcArZ",
				amount:          50,
			},
			want: want{false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// PrepareTransaction
			txJSON, err := xr.PrepareTransaction(tt.args.sernderAccount, tt.args.receiverAccount, tt.args.amount)
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

			// validate transaction
			leggerVer, err := xr.WaitValidation(sentTx.TxJSON.LastLedgerSequence)
			if err != nil {
				t.Fatal(err)
			}
			t.Log("currentLedgerVersion: ", leggerVer)

			// get transaction info
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
