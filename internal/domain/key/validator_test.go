package key_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiromaily/go-crypto-wallet/internal/domain/key"
)

func TestValidateKeyIndex(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		index   uint32
		wantErr bool
	}{
		{
			name:    "zero index is valid",
			index:   0,
			wantErr: false,
		},
		{
			name:    "small index is valid",
			index:   100,
			wantErr: false,
		},
		{
			name:    "max normal index is valid",
			index:   0x7FFFFFFF,
			wantErr: false,
		},
		{
			name:    "hardened index boundary exceeds normal",
			index:   0x80000000,
			wantErr: true,
		},
		{
			name:    "max uint32 exceeds normal",
			index:   0xFFFFFFFF,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := key.ValidateKeyIndex(tt.index)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateKeyCount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		count   uint32
		wantErr bool
	}{
		{
			name:    "zero count is invalid",
			count:   0,
			wantErr: true,
		},
		{
			name:    "one key is valid",
			count:   1,
			wantErr: false,
		},
		{
			name:    "typical count is valid",
			count:   100,
			wantErr: false,
		},
		{
			name:    "max count is valid",
			count:   10000,
			wantErr: false,
		},
		{
			name:    "exceeds max count",
			count:   10001,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := key.ValidateKeyCount(tt.count)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateKeyRange(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		idxFrom uint32
		count   uint32
		wantErr bool
	}{
		{
			name:    "valid range from zero",
			idxFrom: 0,
			count:   100,
			wantErr: false,
		},
		{
			name:    "valid range from middle",
			idxFrom: 1000,
			count:   500,
			wantErr: false,
		},
		{
			name:    "zero count is invalid",
			idxFrom: 0,
			count:   0,
			wantErr: true,
		},
		{
			name:    "overflow check",
			idxFrom: 0x7FFFFFFF,
			count:   2,
			wantErr: true,
		},
		{
			name:    "near max is valid",
			idxFrom: 0x7FFFFFFF - 10,
			count:   10,
			wantErr: false,
		},
		{
			name:    "exceeds max count",
			idxFrom: 0,
			count:   10001,
			wantErr: true,
		},
		{
			name:    "start index in hardened range",
			idxFrom: 0x80000000,
			count:   1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := key.ValidateKeyRange(tt.idxFrom, tt.count)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateSeed(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		seed    []byte
		wantErr bool
	}{
		{
			name:    "nil seed is invalid",
			seed:    nil,
			wantErr: true,
		},
		{
			name:    "empty seed is invalid",
			seed:    []byte{},
			wantErr: true,
		},
		{
			name:    "too short seed is invalid",
			seed:    make([]byte, 15),
			wantErr: true,
		},
		{
			name:    "minimum valid seed (16 bytes)",
			seed:    make([]byte, 16),
			wantErr: false,
		},
		{
			name:    "typical seed (32 bytes)",
			seed:    make([]byte, 32),
			wantErr: false,
		},
		{
			name:    "BIP39 seed (64 bytes)",
			seed:    make([]byte, 64),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := key.ValidateSeed(tt.seed)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateWalletKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		key     key.WalletKey
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid key with P2PKH address",
			key: key.WalletKey{
				FullPubKey: "0279BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798",
				P2PKHAddr:  "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
			},
			wantErr: false,
		},
		{
			name: "valid key with P2SH-SegWit address",
			key: key.WalletKey{
				FullPubKey:     "0279BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798",
				P2SHSegWitAddr: "3J98t1WpEZ73CNmYviecrnyiWrnqRhWNLy",
			},
			wantErr: false,
		},
		{
			name: "valid key with Bech32 address",
			key: key.WalletKey{
				FullPubKey: "0279BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798",
				Bech32Addr: "bc1qw508d6qejxtdg4y5r3zarvary0c5xw7kv8f3t4",
			},
			wantErr: false,
		},
		{
			name: "valid key with Taproot address (BIP86)",
			key: key.WalletKey{
				FullPubKey:  "0279BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798",
				TaprootAddr: "bc1p5cyxnuxmeuwuvkwfem96lqzszd02n6xdcjrs20cac6yqjjwudpxqkedrcr",
			},
			wantErr: false,
		},
		{
			name: "valid key with multiple address types including Taproot",
			key: key.WalletKey{
				FullPubKey:     "0279BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798",
				P2PKHAddr:      "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
				P2SHSegWitAddr: "3J98t1WpEZ73CNmYviecrnyiWrnqRhWNLy",
				Bech32Addr:     "bc1qw508d6qejxtdg4y5r3zarvary0c5xw7kv8f3t4",
				TaprootAddr:    "bc1p5cyxnuxmeuwuvkwfem96lqzszd02n6xdcjrs20cac6yqjjwudpxqkedrcr",
			},
			wantErr: false,
		},
		{
			name: "missing full public key",
			key: key.WalletKey{
				TaprootAddr: "bc1p5cyxnuxmeuwuvkwfem96lqzszd02n6xdcjrs20cac6yqjjwudpxqkedrcr",
			},
			wantErr: true,
			errMsg:  "wallet key must have full public key",
		},
		{
			name: "missing all address formats",
			key: key.WalletKey{
				FullPubKey: "0279BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798",
			},
			wantErr: true,
			errMsg:  "wallet key must have at least one address format",
		},
		{
			name:    "empty wallet key",
			key:     key.WalletKey{},
			wantErr: true,
			errMsg:  "wallet key must have full public key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := key.ValidateWalletKey(tt.key)
			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestValidateWalletKey_TaprootSupport specifically tests that Taproot addresses
// are properly validated in the WalletKey validation function
func TestValidateWalletKey_TaprootSupport(t *testing.T) {
	t.Parallel()

	// Test that Taproot-only keys are valid
	taprootOnlyKey := key.WalletKey{
		FullPubKey:  "0279BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798",
		TaprootAddr: "bc1p5cyxnuxmeuwuvkwfem96lqzszd02n6xdcjrs20cac6yqjjwudpxqkedrcr",
	}

	err := key.ValidateWalletKey(taprootOnlyKey)
	assert.NoError(t, err, "Taproot-only key should be valid")

	// Test various Taproot address formats
	taprootAddresses := []string{
		// Mainnet Taproot addresses (bc1p prefix)
		"bc1p5cyxnuxmeuwuvkwfem96lqzszd02n6xdcjrs20cac6yqjjwudpxqkedrcr",
		"bc1pxwww0ct9ue7e8tdnlmug5m2tamfn7q06sahstg39ys4c9f3340qqxrdu9k",
		"bc1p0xlxvlhemja6c4dqv22uapctqupfhlxm9h8z3k2e72q4k9hcz7vqzk5jj0",
		// Testnet/Signet Taproot addresses (tb1p prefix)
		"tb1pqqqqp399et2xygdj5xreqhjjvcmzhxw4aywxecjdzew6hylgvsesf3hn0c",
		"tb1p0xlxvlhemja6c4dqv22uapctqupfhlxm9h8z3k2e72q4k9hcz7vq47zagq",
	}

	for _, addr := range taprootAddresses {
		t.Run("valid_taproot_"+addr[:10], func(t *testing.T) {
			t.Parallel()
			testKey := key.WalletKey{
				FullPubKey:  "0279BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798",
				TaprootAddr: addr,
			}
			err := key.ValidateWalletKey(testKey)
			assert.NoError(t, err, "Valid Taproot address should pass validation: %s", addr)
		})
	}

	// Verify Taproot addresses have correct characteristics
	for _, addr := range taprootAddresses {
		// Check prefix
		assert.True(t,
			strings.HasPrefix(addr, "bc1p") || strings.HasPrefix(addr, "tb1p"),
			"Taproot address should have bc1p or tb1p prefix: %s", addr)

		// Check length (Taproot addresses are 62 characters)
		assert.Equal(t, 62, len(addr),
			"Taproot address should be 62 characters: %s", addr)

		// Check lowercase
		assert.Equal(t, strings.ToLower(addr), addr,
			"Taproot address should be lowercase: %s", addr)
	}
}

// TestValidateWalletKey_BackwardCompatibility ensures that Taproot support
// doesn't break validation for existing address types
func TestValidateWalletKey_BackwardCompatibility(t *testing.T) {
	t.Parallel()

	legacyFormats := []struct {
		name       string
		walletKey  key.WalletKey
		shouldPass bool
	}{
		{
			name: "Legacy P2PKH only",
			walletKey: key.WalletKey{
				FullPubKey: "0279BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798",
				P2PKHAddr:  "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
			},
			shouldPass: true,
		},
		{
			name: "P2SH-SegWit only",
			walletKey: key.WalletKey{
				FullPubKey:     "0279BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798",
				P2SHSegWitAddr: "3J98t1WpEZ73CNmYviecrnyiWrnqRhWNLy",
			},
			shouldPass: true,
		},
		{
			name: "Bech32 only (Native SegWit)",
			walletKey: key.WalletKey{
				FullPubKey: "0279BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798",
				Bech32Addr: "bc1qw508d6qejxtdg4y5r3zarvary0c5xw7kv8f3t4",
			},
			shouldPass: true,
		},
		{
			name: "All legacy formats combined",
			walletKey: key.WalletKey{
				FullPubKey:     "0279BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798",
				P2PKHAddr:      "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
				P2SHSegWitAddr: "3J98t1WpEZ73CNmYviecrnyiWrnqRhWNLy",
				Bech32Addr:     "bc1qw508d6qejxtdg4y5r3zarvary0c5xw7kv8f3t4",
			},
			shouldPass: true,
		},
	}

	for _, tc := range legacyFormats {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := key.ValidateWalletKey(tc.walletKey)
			if tc.shouldPass {
				assert.NoError(t, err, "Legacy address format should still be valid")
			} else {
				assert.Error(t, err, "Invalid legacy address format should fail")
			}
		})
	}
}
