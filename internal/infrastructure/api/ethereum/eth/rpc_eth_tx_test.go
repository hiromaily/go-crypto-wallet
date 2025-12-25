//go:build integration
// +build integration

package eth_test

import (
	"testing"

	"github.com/bookerzzz/grok"
	"github.com/stretchr/testify/suite"

	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
)

type ethTxTest struct {
	testutil.ETHTestSuite
}

// TestGetTransactionByHash is test for GetTransactionByHash
func (ett *ethTxTest) TestGetTransactionByHash() {
	type args struct {
		txHash string
	}
	type want struct{}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "happy path",
			args: args{
				txHash: "0x1b3ee7f02e622b4bbfe39a3aa9b98ca4651e75a88880c53ca6e34729b452dd9d",
			},
			want: want{},
		},
	}

	for _, tt := range tests {
		ett.T().Run(tt.name, func(t *testing.T) {
			txHash, err := ett.ETH.GetTransactionByHash(tt.args.txHash)
			ett.NoError(err)
			if err == nil {
				// t.Log(res)
				grok.Value(txHash)
			}

			txReceipt, err := ett.ETH.GetTransactionReceipt(tt.args.txHash)
			ett.NoError(err)
			if err == nil {
				// t.Log(txReceipt)
				grok.Value(txReceipt)
			}
		})
	}
}

func TestEthTxTestSuite(t *testing.T) {
	suite.Run(t, new(ethTxTest))
}
