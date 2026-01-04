package btc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const testDerivationP2SH = "/49'/0'/0'"

func TestGenerateP2SHSegWitDescriptor(t *testing.T) {
	service := NewDescriptorService()
	xpub := mustNewExtendedKey(t, testMainnetXpub)

	t.Run("receive descriptor", func(t *testing.T) {
		descriptor, err := service.GenerateP2SHSegWitDescriptor(
			testFingerprint,
			testDerivationP2SH,
			xpub,
			false,
		)
		require.NoError(t, err)
		require.Equal(
			t,
			"sh(wpkh([a1b2c3d4/49'/0'/0']"+testMainnetXpub+"/0/*))",
			descriptor,
		)
	})

	t.Run("change descriptor", func(t *testing.T) {
		descriptor, err := service.GenerateP2SHSegWitDescriptor(
			testFingerprint,
			testDerivationP2SH,
			xpub,
			true,
		)
		require.NoError(t, err)
		require.Equal(
			t,
			"sh(wpkh([a1b2c3d4/49'/0'/0']"+testMainnetXpub+"/1/*))",
			descriptor,
		)
	})
}

func TestGenerateP2SHSegWitDescriptor_NormalizesPath(t *testing.T) {
	service := NewDescriptorService()
	xpub := mustNewExtendedKey(t, testMainnetXpub)

	descriptor, err := service.GenerateP2SHSegWitDescriptor(
		testFingerprint,
		" m/49'/0'/0' ",
		xpub,
		false,
	)
	require.NoError(t, err)
	require.Contains(t, descriptor, "[a1b2c3d4/49'/0'/0']")
	require.Contains(t, descriptor, "/0/*))")
}
