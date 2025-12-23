package key

import (
	"errors"
	"fmt"
)

// ValidateKeyIndex validates that the key index is within acceptable bounds.
// Key indices in HD wallets typically use uint32, but we validate against practical limits.
func ValidateKeyIndex(index uint32) error {
	// Hardened key derivation starts at 2^31 (0x80000000)
	// Normal derivation is 0 to 2^31-1
	const maxNormalIndex = 0x7FFFFFFF

	if index > maxNormalIndex {
		return fmt.Errorf("key index %d exceeds maximum normal derivation index %d", index, maxNormalIndex)
	}

	return nil
}

// ValidateKeyCount validates that the number of keys to generate is reasonable.
func ValidateKeyCount(count uint32) error {
	if count == 0 {
		return errors.New("key count must be at least 1")
	}

	// Practical limit to prevent resource exhaustion
	const maxKeyCount = 10000
	if count > maxKeyCount {
		return fmt.Errorf("key count %d exceeds maximum allowed %d", count, maxKeyCount)
	}

	return nil
}

// ValidateKeyRange validates that a key index range is valid.
func ValidateKeyRange(idxFrom, count uint32) error {
	if err := ValidateKeyIndex(idxFrom); err != nil {
		return fmt.Errorf("invalid start index: %w", err)
	}

	if err := ValidateKeyCount(count); err != nil {
		return fmt.Errorf("invalid key count: %w", err)
	}

	// Check for overflow
	const maxNormalIndex = 0x7FFFFFFF
	if idxFrom+count-1 > maxNormalIndex {
		return fmt.Errorf("key range [%d, %d] exceeds maximum index %d", idxFrom, idxFrom+count-1, maxNormalIndex)
	}

	return nil
}

// ValidateSeed validates that a seed is not empty and meets minimum length requirements.
func ValidateSeed(seed []byte) error {
	if len(seed) == 0 {
		return errors.New("seed cannot be empty")
	}

	// BIP39 seeds are typically 512 bits (64 bytes)
	// Minimum reasonable seed length is 128 bits (16 bytes)
	const minSeedLength = 16
	if len(seed) < minSeedLength {
		return fmt.Errorf("seed length %d is too short, minimum is %d bytes", len(seed), minSeedLength)
	}

	return nil
}

// ValidateWalletKey validates that a WalletKey has required fields populated.
func ValidateWalletKey(wk WalletKey) error {
	if wk.FullPubKey == "" {
		return errors.New("wallet key must have full public key")
	}

	// At least one address format should be present
	if wk.P2PKHAddr == "" && wk.P2SHSegWitAddr == "" && wk.Bech32Addr == "" {
		return errors.New("wallet key must have at least one address format")
	}

	return nil
}
