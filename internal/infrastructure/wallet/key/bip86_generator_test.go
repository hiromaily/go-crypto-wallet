package key

import (
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/address"
)

func TestBIP86Generator(t *testing.T) {
	t.Parallel()

	// Test seed (for testing only, never use in production)
	seed := []byte("test seed for bip86 taproot key generation testing")

	tests := []struct {
		name         string
		coinTypeCode domainCoin.CoinTypeCode
		conf         *chaincfg.Params
		accountType  domainAccount.AccountType
	}{
		{
			name:         "Bitcoin Mainnet Client",
			coinTypeCode: domainCoin.BTC,
			conf:         &chaincfg.MainNetParams,
			accountType:  domainAccount.AccountTypeClient,
		},
		{
			name:         "Bitcoin Testnet Deposit",
			coinTypeCode: domainCoin.BTC,
			conf:         &chaincfg.TestNet3Params,
			accountType:  domainAccount.AccountTypeDeposit,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create BIP86 generator
			generator := NewBIP86Generator(tt.coinTypeCode, tt.conf)

			// Verify generator properties
			assert.Equal(t, domainKey.KeyTypeBIP86, generator.KeyType(), "should return BIP86 key type")
			assert.True(t, generator.SupportsAddressType(address.AddrTypeTaproot), "should support Taproot")
			assert.False(t, generator.SupportsAddressType(address.AddrTypeLegacy), "should not support Legacy")
			assert.False(t, generator.SupportsAddressType(address.AddrTypeP2shSegwit), "should not support P2SH-SegWit")
			assert.False(t, generator.SupportsAddressType(address.AddrTypeBech32), "should not support Bech32")

			// Generate keys
			keys, err := generator.CreateKey(seed, tt.accountType, 0, 5)
			require.NoError(t, err, "should generate keys without error")
			require.Len(t, keys, 5, "should generate 5 keys")

			// Verify keys structure
			for i, key := range keys {
				assert.NotEmpty(t, key.WIF, "WIF should not be empty for key %d", i)
				assert.NotEmpty(t, key.TaprootAddr, "Taproot address should not be empty for key %d", i)
				assert.NotEmpty(t, key.FullPubKey, "Full public key should not be empty for key %d", i)

				// Verify Taproot address format (bech32m)
				if tt.conf == &chaincfg.MainNetParams {
					assert.Contains(t, key.TaprootAddr, "bc1p", "Mainnet address should start with bc1p, key %d", i)
				} else if tt.conf == &chaincfg.TestNet3Params {
					assert.Contains(t, key.TaprootAddr, "tb1p", "Testnet address should start with tb1p, key %d", i)
				}
			}

			// Verify derivation path format
			derivationPath := generator.GetDerivationPath(tt.accountType, 0)
			expectedPrefix := "m/86'/"
			assert.Contains(t, derivationPath, expectedPrefix, "derivation path should start with m/86'/")
		})
	}
}

func TestBIP86GeneratorConsistency(t *testing.T) {
	t.Parallel()

	// Test that BIP86 generates consistent keys
	seed := []byte("test seed for consistency check")
	generator := NewBIP86Generator(domainCoin.BTC, &chaincfg.MainNetParams)

	// Generate keys twice
	keys1, err := generator.CreateKey(seed, domainAccount.AccountTypeClient, 0, 3)
	require.NoError(t, err)

	keys2, err := generator.CreateKey(seed, domainAccount.AccountTypeClient, 0, 3)
	require.NoError(t, err)

	// Verify consistency
	require.Equal(t, len(keys1), len(keys2), "should generate same number of keys")

	for i := range keys1 {
		assert.Equal(t, keys1[i].WIF, keys2[i].WIF, "WIF should be consistent, key %d", i)
		assert.Equal(t, keys1[i].TaprootAddr, keys2[i].TaprootAddr, "Taproot address should match, key %d", i)
		assert.Equal(t, keys1[i].FullPubKey, keys2[i].FullPubKey, "Full public key should match, key %d", i)
	}
}

func TestBIP86vsHDKeyEquivalence(t *testing.T) {
	t.Parallel()

	// Verify that BIP86Generator produces same results as HDKey with purpose 86
	seed := []byte("test seed for equivalence check")
	conf := &chaincfg.MainNetParams
	accountType := domainAccount.AccountTypeClient

	bip86Gen := NewBIP86Generator(domainCoin.BTC, conf)
	hdKey := NewHDKey(PurposeTypeBIP86, domainCoin.BTC, conf)

	keys1, err := bip86Gen.CreateKey(seed, accountType, 0, 3)
	require.NoError(t, err)

	keys2, err := hdKey.CreateKey(seed, accountType, 0, 3)
	require.NoError(t, err)

	require.Equal(t, len(keys1), len(keys2), "should generate same number of keys")

	for i := range keys1 {
		assert.Equal(t, keys1[i].WIF, keys2[i].WIF, "WIF should match for key %d", i)
		assert.Equal(t, keys1[i].TaprootAddr, keys2[i].TaprootAddr, "Taproot address should match for key %d", i)
		assert.Equal(t, keys1[i].FullPubKey, keys2[i].FullPubKey, "Full public key should match for key %d", i)
	}
}
