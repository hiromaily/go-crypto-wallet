package xrp_test

import (
	"testing"

	"github.com/bookerzzz/grok"

	"github.com/hiromaily/go-crypto-wallet/pkg/testutil"
)

// TestValidationCreate is test for ValidationCreate
func TestValidationCreate(t *testing.T) {
	//t.SkipNow()
	xr := testutil.GetXRP()

	type args struct {
		secret string
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
			args: args{"ssCATR7CBvn4GLd1UuU2bqqQffHki"},
			want: want{false},
		},
		{
			name: "happy path 2",
			args: args{"BAWL MAN JADE MOON DOVE GEM SON NOW HAD ADEN GLOW TIRE"},
			want: want{false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := xr.ValidationCreate(tt.args.secret)
			if (err == nil) == tt.want.isErr {
				t.Errorf("ValidationCreate() = %v, want error = %v", err, tt.want.isErr)
				return
			}
			if res != nil {
				t.Log("response:", res)
				grok.Value(res)
			}
		})
	}
}

// TestWalletPropose is test for WalletPropose
func TestWalletPropose(t *testing.T) {
	//t.SkipNow()
	xr := testutil.GetXRP()

	type args struct {
		passphrase string
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
			args: args{"password1"},
			want: want{false},
		},
		{
			name: "happy path 2",
			args: args{"foobar"},
			want: want{false},
		},
		{
			name: "happy path 3",
			args: args{"0x931D387731bBbC988B312206c74F77D004D6B84b"},
			want: want{false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := xr.WalletPropose(tt.args.passphrase)
			if (err == nil) == tt.want.isErr {
				t.Errorf("WalletPropose() = %v, want error = %v", err, tt.want.isErr)
				return
			}
			if res != nil {
				t.Log("response:", res)
				grok.Value(res)
			}
		})
	}
}
