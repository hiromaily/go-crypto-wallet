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
