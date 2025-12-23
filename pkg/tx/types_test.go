package tx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCastBlank(t *testing.T) {
	// OK
	tt := TxType("")
	ttStr := tt.String()
	t.Log(ttStr)

	// Even if no definition, it's ok
	tt2 := TxType("abc")
	ttStr2 := tt2.String()
	t.Log(ttStr2)
}

func TestValidateTxType(t *testing.T) {
	type args struct {
		txType string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "TxTypeUnsigned",
			args: args{TxTypeUnsigned.String()},
			want: true,
		},
		{
			name: "blank",
			args: args{""},
			want: false,
		},
		{
			name: "random",
			args: args{"abc"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateTxType(tt.args.txType)
			assert.Equal(t, tt.want, got, "ValidateTxType() result mismatch")
		})
	}
}
