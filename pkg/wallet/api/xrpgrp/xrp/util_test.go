package xrp_test

import (
	"testing"

	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/xrpgrp/xrp"
)

// TestValidateAddress is test for ValidateAddress
func TestValidateAddress(t *testing.T) {
	type args struct {
		addr string
	}
	type want struct {
		isValid bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "happy path 1",
			args: args{"rEoRcMBfg7VUryw5xSyw883bXU74T8eoYj"},
			want: want{true},
		},
		{
			name: "happy path 2",
			args: args{"rHudMs7gbhag2gzSrKDxpCNGnPMbDkSDqh"},
			want: want{true},
		},
		{
			name: "eth address",
			args: args{"0x931D387731bBbC988B312206c74F77D004D6B84b"},
			want: want{false},
		},
		{
			name: "btc address",
			args: args{"bc1qar0srrr7xfkvy5l643lydnw9re59gtzzwf5mdq"},
			want: want{false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := xrp.ValidateAddress(tt.args.addr)
			if res != tt.want.isValid {
				t.Errorf("ValidateAddress() = %v, want result = %v", res, tt.want.isValid)
				return
			}
		})
	}
}
