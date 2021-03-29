package btc_test

import (
	"testing"

	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
)

// TestGetAccount is test for GetAccount
func TestGetAccount(t *testing.T) {
	// t.SkipNow()
	bc := testutil.GetBTC()

	type args struct {
		addr string
	}
	type want struct {
		account string
		err     error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "happy path",
			args: args{"2N4TcHSCteXwiF2dj8SQijj3w2HieR4x6r5"},
			want: want{"deposit", nil},
		},
		{
			name: "happy path",
			args: args{"2N6DcSuPo8NoLrCPTSqrwjnuLYoN7xDMSzX"},
			want: want{"payment", nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := bc.GetAccount(tt.args.addr)
			if res != tt.want.account {
				t.Errorf("GetAddressInfo() = %v, want %v", res, tt.want.account)
			}
			if err != tt.want.err {
				t.Errorf("GetAddressInfo() = %v, want %v", err, tt.want.err)
			}

			t.Log(res)
		})
	}
	// bc.Close()
}
