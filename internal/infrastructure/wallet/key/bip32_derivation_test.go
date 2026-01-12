package key_test

import (
	"bytes"
	"testing"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/require"

	infraKey "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/wallet/key"
)

func TestParseBIP32DerivationPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		path           string
		expectedIndex  uint32
		expectedError  bool
		errorSubstring string
	}{
		{
			name:          "valid path with m prefix",
			path:          "m/49'/1'/1/0/5",
			expectedIndex: 5,
		},
		{
			name:          "valid path without m prefix",
			path:          "49'/1'/1/0/10",
			expectedIndex: 10,
		},
		{
			name:          "valid path - index 0",
			path:          "m/84'/0'/0/0/0",
			expectedIndex: 0,
		},
		{
			name:          "valid path - index 999",
			path:          "m/84'/0'/2/0/999",
			expectedIndex: 999,
		},
		{
			name:           "invalid path - too few components",
			path:           "m/49'/1'/1",
			expectedError:  true,
			errorSubstring: "expected at least 5 components",
		},
		{
			name:           "invalid path - non-numeric index",
			path:           "m/49'/1'/1/0/abc",
			expectedError:  true,
			errorSubstring: "failed to parse address index",
		},
		{
			name:          "valid path with hardened address index (unusual but valid)",
			path:          "m/49'/1'/1/0/5'",
			expectedIndex: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			index, err := infraKey.ParseBIP32DerivationPath(tt.path)

			if tt.expectedError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errorSubstring)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedIndex, index)
			}
		})
	}
}

func TestDeriveChildPrivateKey(t *testing.T) {
	t.Parallel()

	// Generate test seed and master key
	testSeed := bytes.Repeat([]byte{0x01}, hdkeychain.RecommendedSeedLen)
	masterKey, err := hdkeychain.NewMaster(testSeed, &chaincfg.TestNet3Params)
	require.NoError(t, err)

	// Derive account-level key: m/49'/1'/1
	purpose, err := masterKey.Derive(hdkeychain.HardenedKeyStart + 49)
	require.NoError(t, err)
	coinType, err := purpose.Derive(hdkeychain.HardenedKeyStart + 1)
	require.NoError(t, err)
	accountKey, err := coinType.Derive(hdkeychain.HardenedKeyStart + 1)
	require.NoError(t, err)

	accountXpriv := accountKey.String()

	tests := []struct {
		name          string
		accountXpriv  string
		change        uint32
		addressIndex  uint32
		expectedError bool
	}{
		{
			name:         "derive external address index 0",
			accountXpriv: accountXpriv,
			change:       0,
			addressIndex: 0,
		},
		{
			name:         "derive external address index 5",
			accountXpriv: accountXpriv,
			change:       0,
			addressIndex: 5,
		},
		{
			name:         "derive internal (change) address index 10",
			accountXpriv: accountXpriv,
			change:       1,
			addressIndex: 10,
		},
		{
			name:         "derive large address index",
			accountXpriv: accountXpriv,
			change:       0,
			addressIndex: 999,
		},
		{
			name:          "invalid xpriv",
			accountXpriv:  "invalid_xpriv",
			change:        0,
			addressIndex:  0,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			childKey, err := infraKey.DeriveChildPrivateKey(tt.accountXpriv, tt.change, tt.addressIndex)

			if tt.expectedError {
				require.Error(t, err)
				require.Nil(t, childKey)
			} else {
				require.NoError(t, err)
				require.NotNil(t, childKey)

				// Verify the child key can be used to get private key
				privKey, err := childKey.ECPrivKey()
				require.NoError(t, err)
				require.NotNil(t, privKey)

				// Verify the derivation is consistent
				childKey2, err := infraKey.DeriveChildPrivateKey(tt.accountXpriv, tt.change, tt.addressIndex)
				require.NoError(t, err)
				require.Equal(t, childKey.String(), childKey2.String(), "derivation should be deterministic")
			}
		})
	}
}

func TestExtractAddressIndexFromPSBTInput(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		bip32Derivations []struct {
			PubKey               []byte
			MasterKeyFingerprint uint32
			Bip32Path            []uint32
		}
		expectedAddressIndex uint32
		expectedChange       uint32
		expectedError        bool
		errorSubstring       string
	}{
		{
			name: "valid derivation - external address index 5",
			bip32Derivations: []struct {
				PubKey               []byte
				MasterKeyFingerprint uint32
				Bip32Path            []uint32
			}{
				{
					PubKey:               []byte{0x02, 0x03, 0x04},
					MasterKeyFingerprint: 0x12345678,
					Bip32Path:            []uint32{49 | 0x80000000, 1 | 0x80000000, 1 | 0x80000000, 0, 5},
				},
			},
			expectedAddressIndex: 5,
			expectedChange:       0,
		},
		{
			name: "valid derivation - change address index 10",
			bip32Derivations: []struct {
				PubKey               []byte
				MasterKeyFingerprint uint32
				Bip32Path            []uint32
			}{
				{
					PubKey:               []byte{0x02, 0x03, 0x04},
					MasterKeyFingerprint: 0x12345678,
					Bip32Path:            []uint32{84 | 0x80000000, 0 | 0x80000000, 0 | 0x80000000, 1, 10},
				},
			},
			expectedAddressIndex: 10,
			expectedChange:       1,
		},
		{
			name: "no derivation information",
			bip32Derivations: []struct {
				PubKey               []byte
				MasterKeyFingerprint uint32
				Bip32Path            []uint32
			}{},
			expectedError:  true,
			errorSubstring: "no BIP32 derivation information",
		},
		{
			name: "invalid path length",
			bip32Derivations: []struct {
				PubKey               []byte
				MasterKeyFingerprint uint32
				Bip32Path            []uint32
			}{
				{
					PubKey:               []byte{0x02, 0x03, 0x04},
					MasterKeyFingerprint: 0x12345678,
					Bip32Path:            []uint32{49 | 0x80000000, 1 | 0x80000000},
				},
			},
			expectedError:  true,
			errorSubstring: "invalid BIP32 path length",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			addressIndex, change, err := infraKey.ExtractAddressIndexFromPSBTInput(tt.bip32Derivations)

			if tt.expectedError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errorSubstring)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedAddressIndex, addressIndex)
				require.Equal(t, tt.expectedChange, change)
			}
		})
	}
}

// TestDeriveChildPrivateKey_Integration tests the full derivation flow
func TestDeriveChildPrivateKey_Integration(t *testing.T) {
	t.Parallel()

	// Generate test seed
	testSeed := bytes.Repeat([]byte{0x01}, hdkeychain.RecommendedSeedLen)
	masterKey, err := hdkeychain.NewMaster(testSeed, &chaincfg.TestNet3Params)
	require.NoError(t, err)

	// Derive account-level key: m/49'/1'/1
	purpose, err := masterKey.Derive(hdkeychain.HardenedKeyStart + 49)
	require.NoError(t, err)
	coinType, err := purpose.Derive(hdkeychain.HardenedKeyStart + 1)
	require.NoError(t, err)
	accountKey, err := coinType.Derive(hdkeychain.HardenedKeyStart + 1)
	require.NoError(t, err)

	accountXpriv := accountKey.String()

	// Test deriving multiple indices and verify they are different
	child0, err := infraKey.DeriveChildPrivateKey(accountXpriv, 0, 0)
	require.NoError(t, err)

	child5, err := infraKey.DeriveChildPrivateKey(accountXpriv, 0, 5)
	require.NoError(t, err)

	// Verify that different indices produce different keys
	require.NotEqual(t, child0.String(), child5.String(), "different address indices should produce different keys")

	// Verify public keys are also different
	pubKey0, err := child0.ECPubKey()
	require.NoError(t, err)
	pubKey5, err := child5.ECPubKey()
	require.NoError(t, err)
	require.NotEqual(t, pubKey0.SerializeCompressed(), pubKey5.SerializeCompressed())
}
