package key_test

import (
	"github.com/tyler-smith/go-bip39"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/key"
)

func TestGenerateSeed(t *testing.T) {
	type args struct {
		passphrase string
	}
	type want struct {
		isError bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "happy path",
			args: args{
				passphrase: "password",
			},
			want: want{
				isError: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seed, err := key.GenerateSeed()
			if !tt.want.isError {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
			seed2, mnemonic, err := key.GenerateMnemonic(tt.args.passphrase)
			if !tt.want.isError {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
			t.Log("seed:", seed)
			t.Log("seed2:", seed2)
			t.Log("mnemonic:", mnemonic)

			if !bip39.IsMnemonicValid(mnemonic) {
				t.Error(err)
			}
		})
	}
}
