package eth_test

import (
	"testing"

	"github.com/bookerzzz/grok"

	"github.com/hiromaily/go-bitcoin/pkg/testutil"
)

// TestGetTransactionByHash is test for GetTransactionByHash
func TestGetTransactionByHash(t *testing.T) {
	//t.SkipNow()
	et := testutil.GetETH()

	type args struct {
		txHash string
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
			args: args{
				txHash: "0x1b3ee7f02e622b4bbfe39a3aa9b98ca4651e75a88880c53ca6e34729b452dd9d",
			},
			want: want{false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := et.GetTransactionByHash(tt.args.txHash)
			if (err == nil) == tt.want.isErr {
				t.Errorf("GetTransactionByHash() = %v, want error = %v", err, tt.want.isErr)
				return
			}
			if res != nil {
				//t.Log(res)
				grok.Value(res)
			}
		})
	}
	//et.Close()
}
