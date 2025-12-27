package key

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
)

// TestBIP86IntegrationRealWalletScenario tests BIP86 with real wallet operations
func TestBIP86IntegrationRealWalletScenario(t *testing.T) {
	t.Parallel()

	// Use a test seed (this would come from secure storage in production)
	// Must be 16-64 bytes (128-512 bits) for BIP32
	testSeed := make([]byte, 32)
	copy(testSeed, []byte("test seed for bip86 integration testing with real wallet scenario"))

	scenarios := []struct {
		name           string
		network        *chaincfg.Params
		accountType    domainAccount.AccountType
		expectedPrefix string
	}{
		{
			name:           "Mainnet Client Account",
			network:        &chaincfg.MainNetParams,
			accountType:    domainAccount.AccountTypeClient,
			expectedPrefix: "bc1p",
		},
		{
			name:           "Testnet3 Deposit Account",
			network:        &chaincfg.TestNet3Params,
			accountType:    domainAccount.AccountTypeDeposit,
			expectedPrefix: "tb1p",
		},
		{
			name:           "Signet Payment Account",
			network:        &chaincfg.SigNetParams,
			accountType:    domainAccount.AccountTypePayment,
			expectedPrefix: "tb1p",
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			t.Parallel()

			// Create BIP86 generator
			generator := NewBIP86Generator(domainCoin.BTC, scenario.network)

			// Generate multiple keys as would happen in a real wallet operation
			const numKeys = 10
			keys, err := generator.CreateKey(testSeed, scenario.accountType, 0, numKeys)
			require.NoError(t, err, "key generation should succeed")
			require.Len(t, keys, numKeys, "should generate requested number of keys")

			t.Logf("Testing %s with %d keys", scenario.name, numKeys)

			for i, key := range keys {
				// Verify all required fields are populated
				assert.NotEmpty(t, key.WIF, "WIF should not be empty for key %d", i)
				assert.NotEmpty(t, key.TaprootAddr, "Taproot address should not be empty for key %d", i)
				assert.NotEmpty(t, key.FullPubKey, "Full public key should not be empty for key %d", i)

				// Verify Taproot address format
				assert.True(t,
					strings.HasPrefix(key.TaprootAddr, scenario.expectedPrefix),
					"Taproot address %s should start with %s",
					key.TaprootAddr,
					scenario.expectedPrefix,
				)

				// Verify address length (Taproot addresses are typically 62 characters)
				assert.Equal(t, 62, len(key.TaprootAddr),
					"Taproot address should be 62 characters long")

				// Verify public key is hex encoded and correct length
				pubKeyBytes, err := hex.DecodeString(key.FullPubKey)
				require.NoError(t, err, "Full public key should be valid hex")
				assert.Equal(t, 33, len(pubKeyBytes),
					"Compressed public key should be 33 bytes")

				// Log first few keys for manual verification
				if i < 3 {
					t.Logf("  Key %d:", i)
					t.Logf("    Taproot Address: %s", key.TaprootAddr)
					t.Logf("    Public Key: %s", key.FullPubKey[:20]+"...")

					// Verify derivation path
					derivationPath := generator.GetDerivationPath(scenario.accountType, uint32(i))
					t.Logf("    Derivation Path: %s", derivationPath)
					assert.Contains(t, derivationPath, "m/86'/",
						"derivation path should use BIP86")
				}
			}

			// Verify key uniqueness
			addressSet := make(map[string]bool)
			for _, key := range keys {
				assert.False(t, addressSet[key.TaprootAddr],
					"addresses should be unique: %s", key.TaprootAddr)
				addressSet[key.TaprootAddr] = true
			}
		})
	}
}

// TestBIP86IntegrationKeyConsistency tests that BIP86 generates consistent keys
// across multiple calls, which is critical for wallet recovery
func TestBIP86IntegrationKeyConsistency(t *testing.T) {
	t.Parallel()

	seed := make([]byte, 32)
	copy(seed, []byte("consistency test seed for bip86 wallet recovery"))
	network := &chaincfg.MainNetParams
	accountType := domainAccount.AccountTypeClient

	generator := NewBIP86Generator(domainCoin.BTC, network)

	// Generate keys multiple times
	const iterations = 5
	var previousKeys []string

	for i := 0; i < iterations; i++ {
		keys, err := generator.CreateKey(seed, accountType, 0, 3)
		require.NoError(t, err, "iteration %d: key generation should succeed", i)

		// Extract addresses
		var addresses []string
		for _, key := range keys {
			addresses = append(addresses, key.TaprootAddr)
		}

		if i == 0 {
			previousKeys = addresses
			t.Logf("Initial key generation:")
			for idx, addr := range addresses {
				t.Logf("  Key %d: %s", idx, addr)
			}
		} else {
			// Verify consistency with first generation
			require.Equal(t, previousKeys, addresses,
				"iteration %d: keys should match initial generation", i)
		}
	}

	t.Log("✓ All iterations produced identical keys - wallet recovery would work correctly")
}

// TestBIP86IntegrationMultipleAccounts tests BIP86 with different account types
// as would be used in a multi-account wallet
// Note: Cannot use t.Parallel() here because subtests share allAddresses map
func TestBIP86IntegrationMultipleAccounts(t *testing.T) {
	seed := make([]byte, 32)
	copy(seed, []byte("multi-account test seed for bip86"))
	network := &chaincfg.MainNetParams
	generator := NewBIP86Generator(domainCoin.BTC, network)

	accounts := []struct {
		accountType domainAccount.AccountType
		name        string
	}{
		{domainAccount.AccountTypeClient, "Client"},
		{domainAccount.AccountTypeDeposit, "Deposit"},
		{domainAccount.AccountTypePayment, "Payment"},
		{domainAccount.AccountTypeStored, "Stored"},
	}

	// Track addresses across accounts to ensure no collisions
	allAddresses := make(map[string]string)

	for _, account := range accounts {
		t.Run(account.name, func(t *testing.T) {
			keys, err := generator.CreateKey(seed, account.accountType, 0, 5)
			require.NoError(t, err, "should generate keys for %s account", account.name)

			t.Logf("%s account keys:", account.name)
			for i, key := range keys {
				// Verify address is unique across all accounts
				if existingAccount, exists := allAddresses[key.TaprootAddr]; exists {
					t.Errorf("Address collision: %s appears in both %s and %s accounts",
						key.TaprootAddr, existingAccount, account.name)
				}
				allAddresses[key.TaprootAddr] = account.name

				t.Logf("  Key %d: %s", i, key.TaprootAddr)

				// Verify derivation path is well-formed
				path := generator.GetDerivationPath(account.accountType, uint32(i))
				assert.Contains(t, path, "m/86'/",
					"derivation path should use BIP86 format")
			}
		})
	}

	t.Logf("✓ Generated %d unique addresses across %d accounts",
		len(allAddresses), len(accounts))
}

// TestBIP86IntegrationAddressValidation tests that generated addresses
// are valid for the configured network
func TestBIP86IntegrationAddressValidation(t *testing.T) {
	t.Parallel()

	seed := make([]byte, 32)
	copy(seed, []byte("address validation test seed"))

	tests := []struct {
		name           string
		network        *chaincfg.Params
		expectedPrefix string
		expectedHRP    string // Human Readable Part for bech32m
	}{
		{
			name:           "Bitcoin Mainnet",
			network:        &chaincfg.MainNetParams,
			expectedPrefix: "bc1p",
			expectedHRP:    "bc",
		},
		{
			name:           "Bitcoin Testnet",
			network:        &chaincfg.TestNet3Params,
			expectedPrefix: "tb1p",
			expectedHRP:    "tb",
		},
		{
			name:           "Bitcoin Signet",
			network:        &chaincfg.SigNetParams,
			expectedPrefix: "tb1p",
			expectedHRP:    "tb",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			generator := NewBIP86Generator(domainCoin.BTC, tt.network)
			keys, err := generator.CreateKey(seed, domainAccount.AccountTypeClient, 0, 1)
			require.NoError(t, err)
			require.Len(t, keys, 1)

			address := keys[0].TaprootAddr
			t.Logf("Generated address: %s", address)

			// Verify correct prefix
			assert.True(t, strings.HasPrefix(address, tt.expectedPrefix),
				"address should start with %s for %s", tt.expectedPrefix, tt.name)

			// Verify bech32m format characteristics
			assert.True(t, len(address) == 62,
				"Taproot address should be 62 characters")

			// Verify lowercase (bech32m requirement)
			assert.Equal(t, strings.ToLower(address), address,
				"address should be lowercase")

			// Verify character set (bech32 uses specific characters)
			validChars := "qpzry9x8gf2tvdw0s3jn54khce6mua7l"
			for _, char := range address[4:] { // Skip prefix
				assert.Contains(t, validChars, string(char),
					"address should only contain valid bech32 characters")
			}
		})
	}
}
