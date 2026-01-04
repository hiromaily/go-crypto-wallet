package btc

import (
	"bytes"
	"testing"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/require"
)

const (
	testFingerprint     = "a1b2c3d4"
	testDerivationP2PKH = "/44'/0'/0'"
	//nolint:revive // Descriptor test vector needs the full xpub string
	testMainnetXpub = "xpub6ERApfZwUNrhLCkDtcHTcxd75RbzS1ed54G1LkBUHQVHQKqhMkhgbmJbZRkrgZw4koxb5JaHWkY4ALHY2grBGRjaDMzQLcgJvLJuZZvRcEL"
	//nolint:revive // Descriptor test vector needs the full tpub string
	testTestnetTpub = "tpubD6NzVbkrYhZ4XgiXtGrdW5XDAPFCL9h7we1vwNCpn8tGbBcgfVYjXyhWo4E1xkh56hjod1RhGjxbaTLV3X4FyWuejifB9jusQ46QzG87VKp"
)

func TestGenerateP2PKHDescriptor(t *testing.T) {
	service := NewDescriptorService(&chaincfg.MainNetParams)
	xpub := mustNewExtendedKey(t, testMainnetXpub)

	t.Run("receive descriptor", func(t *testing.T) {
		descriptor, err := service.GenerateP2PKHDescriptor(
			testFingerprint,
			testDerivationP2PKH,
			xpub,
			false,
		)
		require.NoError(t, err)
		require.Equal(
			t,
			"pkh([a1b2c3d4/44'/0'/0']"+testMainnetXpub+"/0/*)",
			descriptor,
		)
	})

	t.Run("change descriptor", func(t *testing.T) {
		descriptor, err := service.GenerateP2PKHDescriptor(
			testFingerprint,
			testDerivationP2PKH,
			xpub,
			true,
		)
		require.NoError(t, err)
		require.Equal(
			t,
			"pkh([a1b2c3d4/44'/0'/0']"+testMainnetXpub+"/1/*)",
			descriptor,
		)
	})

	t.Run("normalizes fingerprint and derivation path", func(t *testing.T) {
		descriptor, err := service.GenerateP2PKHDescriptor(
			"A1B2C3D4",
			" m/44'/0'/0' ",
			xpub,
			false,
		)
		require.NoError(t, err)
		require.Contains(t, descriptor, "[a1b2c3d4/44'/0'/0']")
		require.Contains(t, descriptor, "/0/*)")
	})
}

func TestGenerateP2PKHDescriptorValidation(t *testing.T) {
	service := NewDescriptorService(&chaincfg.MainNetParams)
	xpub := mustNewExtendedKey(t, testMainnetXpub)

	tests := []struct {
		name           string
		fingerprint    string
		derivationPath string
		key            *hdkeychain.ExtendedKey
		isChange       bool
		wantErr        string
	}{
		{
			name:           "nil extended key",
			fingerprint:    testFingerprint,
			derivationPath: testDerivationP2PKH,
			key:            nil,
			wantErr:        "extended public key is nil",
		},
		{
			name:           "private key rejected",
			fingerprint:    testFingerprint,
			derivationPath: testDerivationP2PKH,
			key:            mustNewPrivateExtendedKey(t),
			wantErr:        "extended key must be public",
		},
		{
			name:           "invalid fingerprint",
			fingerprint:    "123",
			derivationPath: testDerivationP2PKH,
			key:            xpub,
			wantErr:        "invalid fingerprint",
		},
		{
			name:           "invalid derivation path",
			fingerprint:    testFingerprint,
			derivationPath: "44'/0'/0'",
			key:            xpub,
			wantErr:        "invalid derivation path",
		},
		{
			name:           "double master prefix rejected",
			fingerprint:    testFingerprint,
			derivationPath: "mM/44'/0'/0'",
			key:            xpub,
			wantErr:        "invalid derivation path",
		},
		{
			name:           "network mismatch",
			fingerprint:    testFingerprint,
			derivationPath: testDerivationP2PKH,
			key:            mustNewExtendedKey(t, testTestnetTpub),
			wantErr:        "network mismatch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.GenerateP2PKHDescriptor(
				tt.fingerprint,
				tt.derivationPath,
				tt.key,
				tt.isChange,
			)
			require.Error(t, err)
			require.ErrorContains(t, err, tt.wantErr)
		})
	}
}

func mustNewExtendedKey(t *testing.T, key string) *hdkeychain.ExtendedKey {
	t.Helper()

	extendedKey, err := hdkeychain.NewKeyFromString(key)
	require.NoError(t, err)

	return extendedKey
}

func mustNewPrivateExtendedKey(t *testing.T) *hdkeychain.ExtendedKey {
	t.Helper()

	seed := bytes.Repeat([]byte{0x01}, hdkeychain.RecommendedSeedLen)
	key, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	require.NoError(t, err)

	return key
}
