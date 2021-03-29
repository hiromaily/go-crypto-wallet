package eth

import "testing"

// TestValidationAddr is test for ValidationAddr
func TestValidationAddr(t *testing.T) {
	type args struct {
		addr string
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
			args: args{"0x967B50a5E4d1D35Fa9aAf7DB8A391b0546209fD2"},
			want: want{false},
		},
		{
			name: "no 0x is OK",
			args: args{"967B50a5E4d1D35Fa9aAf7DB8A391b0546209fD2"},
			want: want{false},
		},
		{
			name: "invalid address",
			args: args{"0xafaljjl3Jd7DB8A391b0546209fD2"},
			want: want{true},
		},
	}
	et := new(Ethereum)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := et.ValidationAddr(tt.args.addr)
			if (err == nil) == tt.want.isErr {
				t.Errorf("ValidationAddr() = %v, want error = %v", err, tt.want.isErr)
			}
		})
	}
	// et.Close()
}

// TestValidationAddr is test for ValidationAddr
func TestConvertToWei(t *testing.T) {
	type args struct {
		wei        int64
		gwei       int64
		floatEther float64
	}
	type want struct {
		wei uint64
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "happy path",
			args: args{
				wei:        1580000000000000000,
				gwei:       1580000000,
				floatEther: 1.58,
			},
			want: want{wei: 1580000000000000000},
		},
	}
	et := new(Ethereum)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := et.FromWei(tt.args.wei).Uint64(); got != tt.want.wei {
				t.Errorf("called FromWei() got: %d, want: %d", got, tt.want.wei)
			}
			if got := et.FromGWei(tt.args.gwei).Uint64(); got != tt.want.wei {
				t.Errorf("called FromGWei() got: %d, want: %d", got, tt.want.wei)
			}
			if got := et.FromFloatEther(tt.args.floatEther).Uint64(); got != tt.want.wei {
				t.Errorf("called FromFloatEther() got: %d, want: %d", got, tt.want.wei)
			}
		})
	}
	// et.Close()
}
