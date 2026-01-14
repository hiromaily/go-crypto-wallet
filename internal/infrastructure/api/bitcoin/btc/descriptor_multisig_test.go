package btc

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/require"

	wallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
)

const testDerivationMultisig = "/48'/0'/0'/2'"

func TestGenerateMultisigDescriptor(t *testing.T) {
	service := NewDescriptorService(&chaincfg.MainNetParams)

	signers := []MultisigSigner{
		{
			Fingerprint:    "a1b2c3d4",
			DerivationPath: testDerivationMultisig,
			ExtendedKey:    mustNewExtendedKey(t, testMainnetXpub),
		},
		{
			Fingerprint:    "b2c3d4e5",
			DerivationPath: testDerivationMultisig,
			ExtendedKey:    newTestXpubFromSeed(t, 0x02),
		},
		{
			Fingerprint:    "c3d4e5f6",
			DerivationPath: testDerivationMultisig,
			ExtendedKey:    newTestXpubFromSeed(t, 0x03),
		},
	}

	t.Run("receive 2-of-3 descriptor with deterministic ordering", func(t *testing.T) {
		descriptor, err := service.GenerateMultisigDescriptor(2, signers, false, wallet.DescriptorTypeWSH)
		require.NoError(t, err)

		expectedKeys := []string{
			fmt.Sprintf("[a1b2c3d4%s]%s/0/*", testDerivationMultisig, signers[0].ExtendedKey.String()),
			fmt.Sprintf("[b2c3d4e5%s]%s/0/*", testDerivationMultisig, signers[1].ExtendedKey.String()),
			fmt.Sprintf("[c3d4e5f6%s]%s/0/*", testDerivationMultisig, signers[2].ExtendedKey.String()),
		}
		sort.Strings(expectedKeys)

		expected := "wsh(sortedmulti(2," + strings.Join(expectedKeys, ",") + "))"
		require.Equal(t, expected, descriptor)
	})

	t.Run("change descriptor uses /1/*", func(t *testing.T) {
		descriptor, err := service.GenerateMultisigDescriptor(2, signers, true, wallet.DescriptorTypeWSH)
		require.NoError(t, err)

		expectedKeys := []string{
			fmt.Sprintf("[a1b2c3d4%s]%s/1/*", testDerivationMultisig, signers[0].ExtendedKey.String()),
			fmt.Sprintf("[b2c3d4e5%s]%s/1/*", testDerivationMultisig, signers[1].ExtendedKey.String()),
			fmt.Sprintf("[c3d4e5f6%s]%s/1/*", testDerivationMultisig, signers[2].ExtendedKey.String()),
		}
		sort.Strings(expectedKeys)

		expected := "wsh(sortedmulti(2," + strings.Join(expectedKeys, ",") + "))"
		require.Equal(t, expected, descriptor)
	})
}

func TestGenerateMultisigDescriptor_Validation(t *testing.T) {
	service := NewDescriptorService(&chaincfg.MainNetParams)

	validSigner := MultisigSigner{
		Fingerprint:    "a1b2c3d4",
		DerivationPath: testDerivationMultisig,
		ExtendedKey:    mustNewExtendedKey(t, testMainnetXpub),
	}

	tests := []struct {
		name         string
		requiredSigs int
		signers      []MultisigSigner
		wantErr      string
	}{
		{
			name:         "no signers",
			requiredSigs: 1,
			signers:      nil,
			wantErr:      "no multisig signers",
		},
		{
			name:         "required greater than signers",
			requiredSigs: 2,
			signers:      []MultisigSigner{validSigner},
			wantErr:      "invalid required signatures",
		},
		{
			name:         "nil extended key",
			requiredSigs: 1,
			signers: []MultisigSigner{
				{
					Fingerprint:    "a1b2c3d4",
					DerivationPath: testDerivationMultisig,
					ExtendedKey:    nil,
				},
			},
			wantErr: "extended public key is nil",
		},
		{
			name:         "network mismatch",
			requiredSigs: 1,
			signers: []MultisigSigner{
				{
					Fingerprint:    "a1b2c3d4",
					DerivationPath: testDerivationMultisig,
					ExtendedKey:    mustNewExtendedKey(t, testTestnetTpub),
				},
			},
			wantErr: "network mismatch",
		},
		{
			name:         "invalid derivation path",
			requiredSigs: 1,
			signers: []MultisigSigner{
				{
					Fingerprint:    "a1b2c3d4",
					DerivationPath: "48'/0'/0'/2'", // missing leading slash
					ExtendedKey:    validSigner.ExtendedKey,
				},
			},
			wantErr: "invalid derivation path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.GenerateMultisigDescriptor(tt.requiredSigs, tt.signers, false, wallet.DescriptorTypeWSH)
			require.ErrorContains(t, err, tt.wantErr)
		})
	}
}

func TestGenerateMultisigDescriptor_NormalizesInputs(t *testing.T) {
	service := NewDescriptorService(&chaincfg.MainNetParams)

	signer := MultisigSigner{
		Fingerprint:    "A1B2C3D4",
		DerivationPath: " m/48'/0'/0'/2' ",
		ExtendedKey:    mustNewExtendedKey(t, testMainnetXpub),
	}

	desc, err := service.GenerateMultisigDescriptor(1, []MultisigSigner{signer}, false, wallet.DescriptorTypeWSH)
	require.NoError(t, err)

	require.Equal(
		t,
		"wsh(sortedmulti(1,[a1b2c3d4/48'/0'/0'/2']"+signer.ExtendedKey.String()+"/0/*))",
		desc,
	)
}

func newTestXpubFromSeed(t *testing.T, seedByte byte) *hdkeychain.ExtendedKey {
	t.Helper()

	seed := bytes.Repeat([]byte{seedByte}, hdkeychain.RecommendedSeedLen)
	master, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	require.NoError(t, err)

	xpub, err := master.Neuter()
	require.NoError(t, err)

	return xpub
}
