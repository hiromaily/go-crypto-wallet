package key

import (
	"testing"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
)

func TestDeriveAccountKey(t *testing.T) {
	// Test xpub from testnet (m/49'/1')
	// This is a valid extended public key at coin level
	testXpub := "tpubDC8k7XivvgUy33hw2vS8JK1Fkje6mRDXGog42GTFCk9k65nCQ92P3pLYjgvRwiGyC66JtchSuvZbd" +
		"jMp8E7n6PJKKuuXqRfbQuecWpG8xou"

	tests := []struct {
		name        string
		xpub        string
		accountType domainAccount.AccountType
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "Valid derivation for deposit account (index 0)",
			xpub:        testXpub,
			accountType: domainAccount.AccountTypeDeposit, // BIP44AccountIndex = 0
			wantErr:     false,
		},
		{
			name:        "Valid derivation for payment account (index 1)",
			xpub:        testXpub,
			accountType: domainAccount.AccountTypePayment, // BIP44AccountIndex = 1
			wantErr:     false,
		},
		{
			name:        "Valid derivation for stored account (index 2)",
			xpub:        testXpub,
			accountType: domainAccount.AccountTypeStored, // BIP44AccountIndex = 2
			wantErr:     false,
		},
		{
			name:        "Invalid extended public key",
			xpub:        "invalid-xpub-string",
			accountType: domainAccount.AccountTypeDeposit,
			wantErr:     true,
			errMsg:      "failed to parse coin-level extended public key",
		},
		{
			name:        "Empty extended public key",
			xpub:        "",
			accountType: domainAccount.AccountTypeDeposit,
			wantErr:     true,
			errMsg:      "failed to parse coin-level extended public key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accountKey, err := DeriveAccountKey(tt.xpub, tt.accountType)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				assert.Nil(t, accountKey)
			} else {
				require.NoError(t, err)
				require.NotNil(t, accountKey)

				// Verify the key is an extended key
				assert.NotEmpty(t, accountKey.String())

				// Verify the key is non-hardened (since derived from xpub)
				// The key should be at account level
				assert.False(t, accountKey.IsPrivate())
			}
		})
	}
}

func TestDeriveAccountKey_ConsistentDerivation(t *testing.T) {
	// Test that deriving the same account type twice produces the same result
	testXpub := "tpubDC8k7XivvgUy33hw2vS8JK1Fkje6mRDXGog42GTFCk9k65nCQ92P3pLYjgvRwiGyC66JtchSuvZbd" +
		"jMp8E7n6PJKKuuXqRfbQuecWpG8xou"
	accountType := domainAccount.AccountTypeDeposit

	key1, err1 := DeriveAccountKey(testXpub, accountType)
	require.NoError(t, err1)

	key2, err2 := DeriveAccountKey(testXpub, accountType)
	require.NoError(t, err2)

	// Both keys should be identical
	assert.Equal(t, key1.String(), key2.String())
}

func TestDeriveAccountKey_DifferentAccounts(t *testing.T) {
	// Test that different account types produce different keys
	testXpub := "tpubDC8k7XivvgUy33hw2vS8JK1Fkje6mRDXGog42GTFCk9k65nCQ92P3pLYjgvRwiGyC66JtchSuvZbd" +
		"jMp8E7n6PJKKuuXqRfbQuecWpG8xou"

	depositKey, err := DeriveAccountKey(testXpub, domainAccount.AccountTypeDeposit)
	require.NoError(t, err)

	paymentKey, err := DeriveAccountKey(testXpub, domainAccount.AccountTypePayment)
	require.NoError(t, err)

	storedKey, err := DeriveAccountKey(testXpub, domainAccount.AccountTypeStored)
	require.NoError(t, err)

	// All keys should be different
	assert.NotEqual(t, depositKey.String(), paymentKey.String())
	assert.NotEqual(t, depositKey.String(), storedKey.String())
	assert.NotEqual(t, paymentKey.String(), storedKey.String())
}

func TestDeriveAccountKey_DerivationFromMainnetXpub(t *testing.T) {
	// Generate a mainnet xpub for testing
	// Note: This is a test xpub generated for demonstration purposes only
	// In real scenarios, xpubs should be generated from actual HD wallet seeds

	// Create a test seed (must be between 128 and 512 bits = 16 to 64 bytes)
	// Using 32 bytes (256 bits) which is the recommended seed length
	seed := make([]byte, 32)
	for i := range seed {
		seed[i] = byte(i)
	}

	// Create master key
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	require.NoError(t, err)

	// Derive to purpose level (m/49')
	purposeKey, err := masterKey.Derive(hdkeychain.HardenedKeyStart + 49)
	require.NoError(t, err)

	// Derive to coin level (m/49'/0' for mainnet BTC)
	coinKey, err := purposeKey.Derive(hdkeychain.HardenedKeyStart + 0)
	require.NoError(t, err)

	// Get the extended public key at coin level
	coinXpub, err := coinKey.Neuter()
	require.NoError(t, err)

	// Test derivation for deposit account
	accountKey, err := DeriveAccountKey(coinXpub.String(), domainAccount.AccountTypeDeposit)
	require.NoError(t, err)
	require.NotNil(t, accountKey)

	// Verify it's a mainnet key (starts with "xpub")
	assert.True(t, accountKey.String()[:4] == "xpub")
}

func TestDeriveAccountKey_DerivationFromTestnetXpub(t *testing.T) {
	// Test with actual testnet xpub
	testXpub := "tpubDC8k7XivvgUy33hw2vS8JK1Fkje6mRDXGog42GTFCk9k65nCQ92P3pLYjgvRwiGyC66JtchSuvZbd" +
		"jMp8E7n6PJKKuuXqRfbQuecWpG8xou"

	accountKey, err := DeriveAccountKey(testXpub, domainAccount.AccountTypeDeposit)
	require.NoError(t, err)
	require.NotNil(t, accountKey)

	// Verify it's a testnet key (starts with "tpub")
	assert.True(t, accountKey.String()[:4] == "tpub")
}
