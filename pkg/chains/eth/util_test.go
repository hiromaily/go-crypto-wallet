package eth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidationAddr(t *testing.T) {
	tests := []struct {
		name    string
		addr    string
		wantErr bool
	}{
		{"happy path", "0x967B50a5E4d1D35Fa9aAf7DB8A391b0546209fD2", false},
		{"no 0x is OK", "967B50a5E4d1D35Fa9aAf7DB8A391b0546209fD2", false},
		{"invalid address", "0xafaljjl3Jd7DB8A391b0546209fD2", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAddr(tt.addr)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConvertToWei(t *testing.T) {
	tests := []struct {
		name       string
		wei        int64
		gwei       int64
		floatEther float64
		wantWei    uint64
	}{
		{"happy path", 1580000000000000000, 1580000000, 1.58, 1580000000000000000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantWei, FromWei(tt.wei).Uint64(), "FromWei() result mismatch")
			assert.Equal(t, tt.wantWei, FromGWei(tt.gwei).Uint64(), "FromGWei() result mismatch")
			assert.Equal(t, tt.wantWei, FromFloatEther(tt.floatEther).Uint64(), "FromFloatEther() result mismatch")
		})
	}
}
