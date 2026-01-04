package btc

import (
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/require"
)

const testDerivationBech32 = "/84'/0'/0'"

func TestGenerateBech32Descriptor(t *testing.T) {
	service := NewDescriptorService(&chaincfg.MainNetParams)
	xpub := mustNewExtendedKey(t, testMainnetXpub)

	t.Run("receive descriptor", func(t *testing.T) {
		descriptor, err := service.GenerateBech32Descriptor(
			testFingerprint,
			testDerivationBech32,
			xpub,
			false,
		)
		require.NoError(t, err)
		require.Equal(
			t,
			"wpkh([a1b2c3d4/84'/0'/0']"+testMainnetXpub+"/0/*)",
			descriptor,
		)
	})

	t.Run("change descriptor", func(t *testing.T) {
		descriptor, err := service.GenerateBech32Descriptor(
			testFingerprint,
			testDerivationBech32,
			xpub,
			true,
		)
		require.NoError(t, err)
		require.Equal(
			t,
			"wpkh([a1b2c3d4/84'/0'/0']"+testMainnetXpub+"/1/*)",
			descriptor,
		)
	})
}

func TestGenerateBech32Descriptor_NormalizesPath(t *testing.T) {
	service := NewDescriptorService(&chaincfg.MainNetParams)
	xpub := mustNewExtendedKey(t, testMainnetXpub)

	descriptor, err := service.GenerateBech32Descriptor(
		testFingerprint,
		" m/84'/0'/0' ",
		xpub,
		false,
	)
	require.NoError(t, err)
	require.Contains(t, descriptor, "[a1b2c3d4/84'/0'/0']")
	require.Contains(t, descriptor, "/0/*)")
}

func TestGenerateBech32Descriptor_NetworkMismatch(t *testing.T) {
	service := NewDescriptorService(&chaincfg.MainNetParams)
	tpub := mustNewExtendedKey(t, testTestnetTpub)

	_, err := service.GenerateBech32Descriptor(
		testFingerprint,
		testDerivationBech32,
		tpub,
		false,
	)
	require.ErrorContains(t, err, "network mismatch")
}
