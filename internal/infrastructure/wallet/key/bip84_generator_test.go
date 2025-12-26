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

func TestBIP84Generator(t *testing.T) {
	t.Parallel()

	// Test seed (for testing only, never use in production)
	seed := []byte("test seed for bip84 key generation testing")

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

			// Create BIP84 generator
			generator := NewBIP84Generator(tt.coinTypeCode, tt.conf)

			// Verify generator properties
			assert.Equal(t, domainKey.KeyTypeBIP84, generator.KeyType(), "should return BIP84 key type")
			assert.True(t, generator.SupportsAddressType(address.AddrTypeBech32), "should support Bech32")
			assert.False(t, generator.SupportsAddressType(address.AddrTypeLegacy), "should not support Legacy")
			assert.False(t, generator.SupportsAddressType(address.AddrTypeP2shSegwit), "should not support P2SH-SegWit")

			// Generate keys
			keys, err := generator.CreateKey(seed, tt.accountType, 0, 5)
			require.NoError(t, err, "should generate keys without error")
			require.Len(t, keys, 5, "should generate 5 keys")

			// Verify keys structure
			for i, key := range keys {
				assert.NotEmpty(t, key.WIF, "WIF should not be empty for key %d", i)
				assert.NotEmpty(t, key.Bech32Addr, "Bech32 address should not be empty for key %d", i)
				assert.NotEmpty(t, key.FullPubKey, "Full public key should not be empty for key %d", i)

				// Verify Bech32 address format
				if tt.conf == &chaincfg.MainNetParams {
					assert.Contains(t, key.Bech32Addr, "bc1", "Mainnet address should start with bc1, key %d", i)
				} else if tt.conf == &chaincfg.TestNet3Params {
					assert.Contains(t, key.Bech32Addr, "tb1", "Testnet address should start with tb1, key %d", i)
				}
			}

			// Verify derivation path format
			derivationPath := generator.GetDerivationPath(tt.accountType, 0)
			expectedPrefix := "m/84'/"
			assert.Contains(t, derivationPath, expectedPrefix, "derivation path should start with m/84'/")
		})
	}
}

func TestBIP84GeneratorConsistency(t *testing.T) {
	t.Parallel()

	// Test that BIP84 generates consistent keys
	seed := []byte("test seed for consistency check")
	generator := NewBIP84Generator(domainCoin.BTC, &chaincfg.MainNetParams)

	// Generate keys twice
	keys1, err := generator.CreateKey(seed, domainAccount.AccountTypeClient, 0, 3)
	require.NoError(t, err)

	keys2, err := generator.CreateKey(seed, domainAccount.AccountTypeClient, 0, 3)
	require.NoError(t, err)

	// Verify consistency
	require.Equal(t, len(keys1), len(keys2), "should generate same number of keys")

	for i := range keys1 {
		assert.Equal(t, keys1[i].WIF, keys2[i].WIF, "WIF should be consistent for key %d", i)
		assert.Equal(t, keys1[i].Bech32Addr, keys2[i].Bech32Addr, "Bech32 address should be consistent for key %d", i)
		assert.Equal(t, keys1[i].FullPubKey, keys2[i].FullPubKey, "Full public key should be consistent for key %d", i)
	}
}

func TestBIP84vsHDKeyEquivalence(t *testing.T) {
	t.Parallel()

	// Verify that BIP84Generator produces same results as HDKey with purpose 84
	seed := []byte("test seed for equivalence check")
	conf := &chaincfg.MainNetParams
	accountType := domainAccount.AccountTypeClient

	bip84Gen := NewBIP84Generator(domainCoin.BTC, conf)
	hdKey := NewHDKey(PurposeTypeBIP84, domainCoin.BTC, conf)

	keys1, err := bip84Gen.CreateKey(seed, accountType, 0, 3)
	require.NoError(t, err)

	keys2, err := hdKey.CreateKey(seed, accountType, 0, 3)
	require.NoError(t, err)

	require.Equal(t, len(keys1), len(keys2), "should generate same number of keys")

	for i := range keys1 {
		assert.Equal(t, keys1[i].WIF, keys2[i].WIF, "WIF should match for key %d", i)
		assert.Equal(t, keys1[i].Bech32Addr, keys2[i].Bech32Addr, "Bech32 address should match for key %d", i)
		assert.Equal(t, keys1[i].FullPubKey, keys2[i].FullPubKey, "Full public key should match for key %d", i)
	}
}
