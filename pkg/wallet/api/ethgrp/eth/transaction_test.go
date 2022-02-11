//go:build integration
// +build integration

package eth_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/ethgrp/eth"
)

type transactionTest struct {
	testutil.ETHTestSuite
}

// TestCreateRawTransaction is test for CreateRawTransaction
func (txt *transactionTest) TestCreateRawTransaction() {
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
		txt.T().Run(tt.name, func(t *testing.T) {
			rawTx, txDetail, err := txt.ETH.CreateRawTransaction(tt.args.senderAddr, tt.args.receiverAddr, tt.args.amount, 0)
			txt.Equal(tt.want.isErr, err != nil)
			if err == nil {
				t.Log(rawTx)
				t.Log(txDetail)
				// grok.Value(rawTx)
				// grok.Value(txDetail)
			}
		})
	}
}

// TestSignAndSendRawTransaction is test for SignOnRawTransaction and SendSignedRawTransaction
func (txt *transactionTest) TestSignAndSendRawTransaction() {
	type args struct {
		senderAddr   string
		receiverAddr string
		amount       uint64
		password     string
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
				senderAddr:   "0xe52307Deb1a7dC3985D2873b45AE23b91D57a36d",
				receiverAddr: "0x72cCC7a7C3fa28C79aaC4f834168767A5762a7D0",
				amount:       0,
				password:     "foobar",
			},
			want: want{true, false},
		},
		{
			name: "happy path",
			args: args{
				senderAddr:   "0xe52307Deb1a7dC3985D2873b45AE23b91D57a36d",
				receiverAddr: "0x72cCC7a7C3fa28C79aaC4f834168767A5762a7D0",
				amount:       0,
				password:     eth.Password,
			},
			want: want{false, false},
		},
	}

	for _, tt := range tests {
		txt.T().Run(tt.name, func(t *testing.T) {
			// create raw transaction
			rawTx, _, err := txt.ETH.CreateRawTransaction(tt.args.senderAddr, tt.args.receiverAddr, tt.args.amount, 0)
			txt.NoError(err)

			// sign on raw transaction
			signedTx, err := txt.ETH.SignOnRawTransaction(rawTx, tt.args.password)
			txt.Equal(tt.want.isSignErr, err != nil)
			if err == nil {
				t.Log(signedTx)
			}

			// send signed transaction
			txHash, err := txt.ETH.SendSignedRawTransaction(signedTx.TxHex)
			txt.Equal(tt.want.isSendErr, err != nil)
			if txHash != "" {
				t.Logf("txHash: %s", txHash)

				// check transaction
				time.Sleep(3 * time.Second)
				tx, err := txt.ETH.GetTransactionByHash(txHash)
				txt.NoError(err)
				t.Logf("tx: %v", tx)

				// check balance
				balance, err := txt.ETH.GetBalance(tt.args.receiverAddr, eth.QuantityTagPending)
				txt.NoError(err)
				txt.NotEqual(0, balance.Uint64())

				// check confirmation
				confirmNum, err := txt.ETH.GetConfirmation(txHash)
				txt.NoError(err)
				t.Logf("confirmation is %d", confirmNum)
			}
		})
	}
}

func TestTransactionTestSuite(t *testing.T) {
	suite.Run(t, new(transactionTest))
}
