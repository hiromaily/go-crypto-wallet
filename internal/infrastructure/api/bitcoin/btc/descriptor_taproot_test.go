package btc

import (
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/require"
)

const testDerivationTaproot = "/86'/0'/0'"

func TestGenerateTaprootDescriptor(t *testing.T) {
	service := NewDescriptorService(&chaincfg.MainNetParams)
	xpub := mustNewExtendedKey(t, testMainnetXpub)

	t.Run("receive descriptor", func(t *testing.T) {
		descriptor, err := service.GenerateTaprootDescriptor(
			testFingerprint,
			testDerivationTaproot,
			xpub,
			false,
		)
		require.NoError(t, err)
		require.Equal(
			t,
			"tr([a1b2c3d4/86'/0'/0']"+testMainnetXpub+"/0/*)",
			descriptor,
		)
	})

	t.Run("change descriptor", func(t *testing.T) {
		descriptor, err := service.GenerateTaprootDescriptor(
			testFingerprint,
			testDerivationTaproot,
			xpub,
			true,
		)
		require.NoError(t, err)
		require.Equal(
			t,
			"tr([a1b2c3d4/86'/0'/0']"+testMainnetXpub+"/1/*)",
			descriptor,
		)
	})
}

func TestGenerateTaprootDescriptor_NormalizesPath(t *testing.T) {
	service := NewDescriptorService(&chaincfg.MainNetParams)
	xpub := mustNewExtendedKey(t, testMainnetXpub)

	descriptor, err := service.GenerateTaprootDescriptor(
		testFingerprint,
		" m/86'/0'/0' ",
		xpub,
		false,
	)
	require.NoError(t, err)
	require.Equal(t, "tr([a1b2c3d4/86'/0'/0']"+testMainnetXpub+"/0/*)", descriptor)
}

func TestGenerateTaprootDescriptor_NetworkMismatch(t *testing.T) {
	service := NewDescriptorService(&chaincfg.MainNetParams)
	tpub := mustNewExtendedKey(t, testTestnetTpub)

	_, err := service.GenerateTaprootDescriptor(
		testFingerprint,
		testDerivationTaproot,
		tpub,
		false,
	)
	require.ErrorContains(t, err, "network mismatch")
}
