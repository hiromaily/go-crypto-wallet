package btc

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/require"
)

const testDerivationTaprootMultiKey = "/86'/0'/0'"

func TestGenerateTaprootScriptPathDescriptor(t *testing.T) {
	service := NewDescriptorService(&chaincfg.MainNetParams)

	signers := []MultisigSigner{
		{
			Fingerprint:    "a1b2c3d4",
			DerivationPath: testDerivationTaprootMultiKey,
			ExtendedKey:    mustNewExtendedKey(t, testMainnetXpub),
		},
		{
			Fingerprint:    "b2c3d4e5",
			DerivationPath: testDerivationTaprootMultiKey,
			ExtendedKey:    newTestXpubFromSeed(t, 0x04),
		},
	}

	t.Run("receive descriptor sorts keys", func(t *testing.T) {
		desc, err := service.GenerateTaprootScriptPathDescriptor(signers, false)
		require.NoError(t, err)

		expectedKeys := []string{
			fmt.Sprintf("[a1b2c3d4%s]%s/0/*", testDerivationTaprootMultiKey, signers[0].ExtendedKey.String()),
			fmt.Sprintf("[b2c3d4e5%s]%s/0/*", testDerivationTaprootMultiKey, signers[1].ExtendedKey.String()),
		}
		sort.Strings(expectedKeys)

		require.Equal(t, "tr("+strings.Join(expectedKeys, ",")+")", desc)
	})

	t.Run("change descriptor uses /1/*", func(t *testing.T) {
		desc, err := service.GenerateTaprootScriptPathDescriptor(signers, true)
		require.NoError(t, err)

		expectedKeys := []string{
			fmt.Sprintf("[a1b2c3d4%s]%s/1/*", testDerivationTaprootMultiKey, signers[0].ExtendedKey.String()),
			fmt.Sprintf("[b2c3d4e5%s]%s/1/*", testDerivationTaprootMultiKey, signers[1].ExtendedKey.String()),
		}
		sort.Strings(expectedKeys)

		require.Equal(t, "tr("+strings.Join(expectedKeys, ",")+")", desc)
	})
}

func TestGenerateTaprootScriptPathDescriptor_Validation(t *testing.T) {
	service := NewDescriptorService(&chaincfg.MainNetParams)

	validSigner := MultisigSigner{
		Fingerprint:    "a1b2c3d4",
		DerivationPath: testDerivationTaprootMultiKey,
		ExtendedKey:    mustNewExtendedKey(t, testMainnetXpub),
	}

	tests := []struct {
		name    string
		signers []MultisigSigner
		wantErr string
	}{
		{
			name:    "not enough signers",
			signers: []MultisigSigner{validSigner},
			wantErr: "at least 2 signers",
		},
		{
			name: "network mismatch",
			signers: []MultisigSigner{
				{
					Fingerprint:    "a1b2c3d4",
					DerivationPath: testDerivationTaprootMultiKey,
					ExtendedKey:    mustNewExtendedKey(t, testTestnetTpub),
				},
				validSigner,
			},
			wantErr: "network mismatch",
		},
		{
			name: "invalid derivation path",
			signers: []MultisigSigner{
				{
					Fingerprint:    "a1b2c3d4",
					DerivationPath: "86'/0'/0'", // missing leading slash
					ExtendedKey:    validSigner.ExtendedKey,
				},
				validSigner,
			},
			wantErr: "invalid derivation path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.GenerateTaprootScriptPathDescriptor(tt.signers, false)
			require.ErrorContains(t, err, tt.wantErr)
		})
	}
}
