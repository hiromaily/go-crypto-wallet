package eth_test

import (
	"testing"

	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
)

// TestGetTotalBalance is test for GetTotalBalance
func TestGetTotalBalance(t *testing.T) {
	//t.SkipNow()
	et := testutil.GetETH()

	type args struct {
		addrs []string
	}
	type want struct {
		total uint64
		isErr bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "happy path",
			args: args{[]string{
				"0x967B50a5E4d1D35Fa9aAf7DB8A391b0546209fD2",
				"0x16e486ED0148A59C58400232D5b4AF6bE4dC1ef0",
				"0x57033fC5434F3271e83B4695dAd726F348b854c2",
				"0x3727eE9FA88B21a0703946f9afEE3930f5980c15",
				"0xe933a3318C3f5D94c2A3B2BEAEF772F67a45311c",
			}},
			want: want{100, false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			total, userAmounts := et.GetTotalBalance(tt.args.addrs)
			t.Log(total)
			t.Log(userAmounts)
		})
	}
	//et.Close()
}
