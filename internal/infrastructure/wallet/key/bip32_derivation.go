package key

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
)

// ParseBIP32DerivationPath parses a BIP32 derivation path and extracts the address index.
// Path format: m/purpose'/coin_type'/account'/change/address_index
// Example: m/49'/1'/1/0/5 -> address_index = 5
//
// Returns:
//   - addressIndex: The address index (last component of the path)
//   - err: Error if path is invalid or cannot be parsed
func ParseBIP32DerivationPath(path string) (addressIndex uint32, err error) {
	// Remove "m/" prefix if present
	path = strings.TrimPrefix(path, "m/")

	// Split path into components
	components := strings.Split(path, "/")
	if len(components) < 5 {
		return 0, fmt.Errorf("invalid BIP32 derivation path: %s (expected at least 5 components)", path)
	}

	// Last component is the address index
	indexStr := components[len(components)-1]

	// Remove hardened marker (') if present
	indexStr = strings.TrimSuffix(indexStr, "'")

	// Parse as uint32
	index, err := strconv.ParseUint(indexStr, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("failed to parse address index from path %s: %w", path, err)
	}

	return uint32(index), nil
}

// DeriveChildPrivateKey derives a child private key at the specified address index
// from an account-level extended private key.
//
// Parameters:
//   - accountXpriv: Account-level extended private key (xpriv format)
//   - change: Change index (0 for external/receiving, 1 for internal/change)
//   - addressIndex: Address index to derive
//
// Returns:
//   - childPrivKey: Child private key at m/account/change/addressIndex
//   - err: Error if derivation fails
//
// Path: account_xpriv/change/addressIndex
// Example: If accountXpriv is at m/49'/1'/1, derives m/49'/1'/1/0/5
func DeriveChildPrivateKey(
	accountXpriv string,
	change uint32,
	addressIndex uint32,
) (*hdkeychain.ExtendedKey, error) {
	// Parse account-level extended private key
	accountKey, err := hdkeychain.NewKeyFromString(accountXpriv)
	if err != nil {
		return nil, fmt.Errorf("failed to parse account extended private key: %w", err)
	}

	// Derive change level (m/account/change)
	changeKey, err := accountKey.Derive(change)
	if err != nil {
		return nil, fmt.Errorf("failed to derive change key at index %d: %w", change, err)
	}

	// Derive address index (m/account/change/addressIndex)
	childKey, err := changeKey.Derive(addressIndex)
	if err != nil {
		return nil, fmt.Errorf("failed to derive child key at index %d: %w", addressIndex, err)
	}

	return childKey, nil
}

// ExtractAddressIndexFromPSBTInput extracts the address index from a PSBT input's
// BIP32 derivation information.
//
// This function looks at the Bip32Derivation field of a PSBT input and extracts
// the address index from the derivation path. It returns the first valid address
// index found, or an error if no valid derivation paths are found.
//
// Parameters:
//   - bip32Derivations: BIP32 derivation information from PSBT input
//
// Returns:
//   - addressIndex: The address index extracted from the derivation path
//   - change: The change index (0 for external, 1 for internal)
//   - err: Error if no valid derivation path found
func ExtractAddressIndexFromPSBTInput(
	bip32Derivations []struct {
		PubKey               []byte
		MasterKeyFingerprint uint32
		Bip32Path            []uint32
	},
) (addressIndex uint32, change uint32, err error) {
	if len(bip32Derivations) == 0 {
		return 0, 0, fmt.Errorf("no BIP32 derivation information in PSBT input")
	}

	// Use the first derivation (all keys should have same address index)
	deriv := bip32Derivations[0]
	if len(deriv.Bip32Path) < 5 {
		return 0, 0, fmt.Errorf("invalid BIP32 path length: %d (expected at least 5)", len(deriv.Bip32Path))
	}

	// BIP32 path format: [purpose, coin_type, account, change, address_index]
	// All values are in non-hardened form in Bip32Path
	change = deriv.Bip32Path[len(deriv.Bip32Path)-2]       // Second to last component
	addressIndex = deriv.Bip32Path[len(deriv.Bip32Path)-1] // Last component

	return addressIndex, change, nil
}
