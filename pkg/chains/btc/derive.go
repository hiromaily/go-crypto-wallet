package btc

import (
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
)

// DeriveChildPrivateKey derives a child private key at the specified address index
// from an account-level extended private key.
//
// Parameters:
//   - accountXpriv: Account-level extended private key (xpriv format)
//   - change: Change index (0 for external/receiving, 1 for internal/change)
//   - addressIndex: Address index to derive
//
// Path: account_xpriv/change/addressIndex
func DeriveChildPrivateKey(
	accountXpriv string,
	change uint32,
	addressIndex uint32,
) (*hdkeychain.ExtendedKey, error) {
	accountKey, err := hdkeychain.NewKeyFromString(accountXpriv)
	if err != nil {
		return nil, fmt.Errorf("failed to parse account extended private key: %w", err)
	}

	changeKey, err := accountKey.Derive(change)
	if err != nil {
		return nil, fmt.Errorf("failed to derive change key at index %d: %w", change, err)
	}

	childKey, err := changeKey.Derive(addressIndex)
	if err != nil {
		return nil, fmt.Errorf("failed to derive child key at index %d: %w", addressIndex, err)
	}

	return childKey, nil
}

// DeriveAccountKey derives an account-specific key from a coin-level extended public/private key.
//
// Parameters:
//   - coinLevelExtendedKey: Extended key at m/purpose'/coin' level
//   - accountIndex: Account index to derive (non-hardened)
func DeriveAccountKey(coinLevelExtendedKey string, accountIndex uint32) (*hdkeychain.ExtendedKey, error) {
	coinLevelKey, err := hdkeychain.NewKeyFromString(coinLevelExtendedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse coin-level extended key: %w", err)
	}

	accountKey, err := coinLevelKey.Derive(accountIndex)
	if err != nil {
		return nil, fmt.Errorf("failed to derive account key at index %d: %w", accountIndex, err)
	}

	return accountKey, nil
}

// FingerprintFromExtendedKey extracts the fingerprint from an extended public key.
//
// This function:
//  1. Parses the extended key string
//  2. Extracts the public key
//  3. Calculates the HASH160 of the compressed public key
//  4. Returns the first 4 bytes as a lowercase hex string (8 characters)
func FingerprintFromExtendedKey(extendedKey string) (string, error) {
	key, err := hdkeychain.NewKeyFromString(extendedKey)
	if err != nil {
		return "", fmt.Errorf("failed to parse extended key: %w", err)
	}

	pubKey, err := key.ECPubKey()
	if err != nil {
		return "", fmt.Errorf("failed to get public key from extended key: %w", err)
	}

	pubKeyBytes := pubKey.SerializeCompressed()
	hash := btcutil.Hash160(pubKeyBytes)
	fingerprint := hash[:4]

	return hex.EncodeToString(fingerprint), nil
}
