package eth_test

import (
	"testing"
	"time"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/testutil"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api/ethgrp/eth"
)

// TestCreateRawTransaction is test for CreateRawTransaction
func TestCreateRawTransaction(t *testing.T) {
	//t.SkipNow()
	et := testutil.GetETH()

	type args struct {
		senderAddr   string
		receiverAddr string
		amount       uint64
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
			name: "happy path, send all",
			args: args{
				senderAddr:   "0xe52307Deb1a7dC3985D2873b45AE23b91D57a36d",
				receiverAddr: "0x72cCC7a7C3fa28C79aaC4f834168767A5762a7D0",
				amount:       0,
			},
			want: want{false},
		},
		{
			name: "happy path, send specific amount",
			args: args{
				senderAddr:   "0xe52307Deb1a7dC3985D2873b45AE23b91D57a36d",
				receiverAddr: "0x72cCC7a7C3fa28C79aaC4f834168767A5762a7D0",
				amount:       40000000000000000, // 0.04 Ether
			},
			want: want{false},
		},
		{
			name: "sender balance is insufficient",
			args: args{
				senderAddr:   "0xe52307Deb1a7dC3985D2873b45AE23b91D57a36d",
				receiverAddr: "0x72cCC7a7C3fa28C79aaC4f834168767A5762a7D0",
				amount:       250000000000000000, // 0.25 Ether
			},
			want: want{true},
		},
		{
			name: "sender doesn't have amount",
			args: args{
				senderAddr:   "0x0Dd4d77D8b3bf210974332d1E16275bbEDdbF1CE",
				receiverAddr: "0x72cCC7a7C3fa28C79aaC4f834168767A5762a7D0",
				amount:       0,
			},
			want: want{true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rawTx, txDetail, err := et.CreateRawTransaction(tt.args.senderAddr, tt.args.receiverAddr, tt.args.amount)
			if (err == nil) == tt.want.isErr {
				t.Errorf("CreateRawTransaction() = %v, want error = %v", err, tt.want.isErr)
				return
			}
			if rawTx != nil {
				t.Log(rawTx)
				//grok.Value(rawTx)
			}
			if txDetail != nil {
				t.Log(txDetail)
				//grok.Value(txDetail)
			}
		})
	}
	//et.Close()
}

// TestSignAndSendRawTransaction is test for SignOnRawTransaction and SendSignedRawTransaction
func TestSignAndSendRawTransaction(t *testing.T) {
	//t.SkipNow()
	et := testutil.GetETH()

	type args struct {
		senderAddr    string
		receiverAddr  string
		amount        uint64
		senderAccount account.AccountType
		password      string
	}
	type want struct {
		isSignErr bool
		isSendErr bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "wrong password",
			args: args{
				senderAddr:    "0xe52307Deb1a7dC3985D2873b45AE23b91D57a36d",
				receiverAddr:  "0x72cCC7a7C3fa28C79aaC4f834168767A5762a7D0",
				amount:        0,
				senderAccount: account.AccountTypeClient,
				password:      "foobar",
			},
			want: want{true, false},
		},
		{
			name: "happy path",
			args: args{
				senderAddr:    "0xe52307Deb1a7dC3985D2873b45AE23b91D57a36d",
				receiverAddr:  "0x72cCC7a7C3fa28C79aaC4f834168767A5762a7D0",
				amount:        0,
				senderAccount: account.AccountTypeClient,
				password:      eth.Password,
			},
			want: want{false, false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create raw transaction
			rawTx, _, err := et.CreateRawTransaction(tt.args.senderAddr, tt.args.receiverAddr, tt.args.amount)
			if err != nil {
				t.Fatal(err)
			}
			// sign on raw transaction
			signedTx, err := et.SignOnRawTransaction(rawTx, tt.args.password, tt.args.senderAccount)
			if (err == nil) == tt.want.isSignErr {
				t.Errorf("SignOnRawTransaction() = %v, want error = %v", err, tt.want.isSignErr)
				return
			}
			if err != nil {
				return
			}
			if signedTx != nil {
				t.Log(signedTx)
			}
			// send signed transaction
			txHash, err := et.SendSignedRawTransaction(signedTx.TxHex)
			if (err == nil) == tt.want.isSendErr {
				t.Errorf("SendSignedRawTransaction() = %v, want error = %v", err, tt.want.isSignErr)
				return
			}
			if txHash != "" {
				t.Log(txHash)
				//eth_getTransactionByHash
				//0x1b3ee7f02e622b4bbfe39a3aa9b98ca4651e75a88880c53ca6e34729b452dd9d
				time.Sleep(3 * time.Second)
				res, err := et.GetTransactionByHash(txHash)
				if err != nil {
					t.Fatal(err)
				}
				t.Log(res)

				//TODO:check balance

			}
		})
	}
	//et.Close()
}
