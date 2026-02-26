package btc

import (
	"fmt"
	"sort"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/require"

	dtobtc "github.com/hiromaily/go-crypto-wallet/internal/application/dto/btc"
)

const testDerivationTaprootMultiKey = "/86'/0'/0'"

func TestGenerateTaprootScriptPathDescriptor(t *testing.T) {
	service := NewDescriptorService(&chaincfg.MainNetParams)

	signers := []dtobtc.MultisigSigner{
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

	t.Run("receive descriptor uses correct Taproot format", func(t *testing.T) {
		desc, err := service.GenerateTaprootScriptPathDescriptor(signers, false)
		require.NoError(t, err)

		expectedKeys := []string{
			fmt.Sprintf("[a1b2c3d4%s]%s/0/*", testDerivationTaprootMultiKey, signers[0].ExtendedKey.String()),
			fmt.Sprintf("[b2c3d4e5%s]%s/0/*", testDerivationTaprootMultiKey, signers[1].ExtendedKey.String()),
		}
		sort.Strings(expectedKeys)

		// Correct Taproot format: tr(internal_key,pk(script_key))
		// First key is internal key (key-path), second key is script path
		expected := fmt.Sprintf("tr(%s,pk(%s))", expectedKeys[0], expectedKeys[1])
		require.Equal(t, expected, desc)
	})

	t.Run("change descriptor uses /1/*", func(t *testing.T) {
		desc, err := service.GenerateTaprootScriptPathDescriptor(signers, true)
		require.NoError(t, err)

		expectedKeys := []string{
			fmt.Sprintf("[a1b2c3d4%s]%s/1/*", testDerivationTaprootMultiKey, signers[0].ExtendedKey.String()),
			fmt.Sprintf("[b2c3d4e5%s]%s/1/*", testDerivationTaprootMultiKey, signers[1].ExtendedKey.String()),
		}
		sort.Strings(expectedKeys)

		// Correct Taproot format: tr(internal_key,pk(script_key))
		// First key is internal key (key-path), second key is script path
		expected := fmt.Sprintf("tr(%s,pk(%s))", expectedKeys[0], expectedKeys[1])
		require.Equal(t, expected, desc)
	})

	t.Run("three signers use script tree with braces", func(t *testing.T) {
		threeSigners := []dtobtc.MultisigSigner{
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
			{
				Fingerprint:    "c3d4e5f6",
				DerivationPath: testDerivationTaprootMultiKey,
				ExtendedKey:    newTestXpubFromSeed(t, 0x05),
			},
		}

		desc, err := service.GenerateTaprootScriptPathDescriptor(threeSigners, false)
		require.NoError(t, err)

		expectedKeys := []string{
			fmt.Sprintf("[a1b2c3d4%s]%s/0/*", testDerivationTaprootMultiKey, threeSigners[0].ExtendedKey.String()),
			fmt.Sprintf("[b2c3d4e5%s]%s/0/*", testDerivationTaprootMultiKey, threeSigners[1].ExtendedKey.String()),
			fmt.Sprintf("[c3d4e5f6%s]%s/0/*", testDerivationTaprootMultiKey, threeSigners[2].ExtendedKey.String()),
		}
		sort.Strings(expectedKeys)

		// With 3+ keys: tr(internal_key,{pk(key1),pk(key2)})
		// First key is internal key, remaining keys are script paths in braces
		expected := fmt.Sprintf("tr(%s,{pk(%s),pk(%s)})", expectedKeys[0], expectedKeys[1], expectedKeys[2])
		require.Equal(t, expected, desc)
	})
}

func TestGenerateTaprootScriptPathDescriptor_Validation(t *testing.T) {
	service := NewDescriptorService(&chaincfg.MainNetParams)

	validSigner := dtobtc.MultisigSigner{
		Fingerprint:    "a1b2c3d4",
		DerivationPath: testDerivationTaprootMultiKey,
		ExtendedKey:    mustNewExtendedKey(t, testMainnetXpub),
	}

	tests := []struct {
		name    string
		signers []dtobtc.MultisigSigner
		wantErr string
	}{
		{
			name:    "not enough signers",
			signers: []dtobtc.MultisigSigner{validSigner},
			wantErr: "at least 2 signers",
		},
		{
			name: "network mismatch",
			signers: []dtobtc.MultisigSigner{
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
			signers: []dtobtc.MultisigSigner{
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
