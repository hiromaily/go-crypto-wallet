package wallet

import (
	"errors"
	"fmt"
	"strings"
)

// inputCharset defines the valid characters for descriptor strings.
// This is the character set used in Bitcoin descriptors before checksum encoding.
const inputCharset = "0123456789()[],'/*abcdefgh@:$%{}IJKLMNOPQRSTUVWXYZ&+-.;" +
	"<=>?!^_|~ijklmnopqrstuvwxyzABCDEFGH`#\"\\ "

// checksumCharset defines the character set used for the descriptor checksum.
// This is based on the bech32 character set defined in BIP173.
const checksumCharset = "qpzry9x8gf2tvdw0s3jn54khce6mua7l"

// DescriptorBuilder provides functionality to construct Bitcoin output descriptors
// with proper BIP380 formatting and checksums.
//
// BIP380 defines output descriptors as a standard way to describe Bitcoin addresses
// and scripts in a human-readable format that includes all information needed to
// reconstruct the scriptPubKey.
//
// Example descriptors:
//   - P2PKH: pkh([fingerprint/44'/0'/0']xpub.../0/*)
//   - P2SH-SegWit: sh(wpkh([fingerprint/49'/0'/0']xpub.../0/*))
//   - Native SegWit: wpkh([fingerprint/84'/0'/0']xpub.../0/*)
//   - Taproot: tr([fingerprint/86'/0'/0']xpub.../0/*)
type DescriptorBuilder struct{}

// NewDescriptorBuilder creates a new descriptor builder.
func NewDescriptorBuilder() *DescriptorBuilder {
	return &DescriptorBuilder{}
}

// BuildDescriptor constructs a descriptor string from its components.
//
// The descriptor format follows BIP380 specification:
//   - Type: The descriptor function (pkh, wpkh, sh(wpkh), tr, etc.)
//   - Key: The key information including fingerprint, derivation path, and xpub
//
// Parameters:
//   - descType: The descriptor type (PKH, SHWPKH, WPKH, TR, WSH)
//   - key: The descriptor key with fingerprint, derivation path, and xpub
//
// Returns:
//   - Descriptor string WITHOUT checksum
//   - Error if the components are invalid
//
// Example:
//
//	key := DescriptorKey{
//	    Fingerprint: "a1b2c3d4",
//	    DerivationPath: "/49'/0'/0'",
//	    ExtendedPubKey: "xpub6ERApfZw...",
//	}
//	desc, err := builder.BuildDescriptor(DescriptorTypeSHWPKH, key)
//	// Result: "sh(wpkh([a1b2c3d4/49'/0'/0']xpub6ERApfZw...))"
func (*DescriptorBuilder) BuildDescriptor(descType DescriptorType, key DescriptorKey) (string, error) {
	// Validate key
	if err := ValidateDescriptorKey(&key); err != nil {
		return "", fmt.Errorf("invalid descriptor key: %w", err)
	}

	// Build key string: [fingerprint/path]xpub
	keyStr := ""
	if key.Fingerprint != "" || key.DerivationPath != "" {
		keyStr = "["
		if key.Fingerprint != "" {
			keyStr += key.Fingerprint
		}
		if key.DerivationPath != "" {
			keyStr += key.DerivationPath
		}
		keyStr += "]"
	}
	keyStr += key.ExtendedPubKey

	// Build descriptor based on type
	var descriptor string
	switch descType {
	case DescriptorTypePKH:
		descriptor = fmt.Sprintf("pkh(%s)", keyStr)
	case DescriptorTypeSHWPKH:
		descriptor = fmt.Sprintf("sh(wpkh(%s))", keyStr)
	case DescriptorTypeWPKH:
		descriptor = fmt.Sprintf("wpkh(%s)", keyStr)
	case DescriptorTypeTR:
		descriptor = fmt.Sprintf("tr(%s)", keyStr)
	case DescriptorTypeWSH:
		// WSH typically requires multisig, this is a simple case
		descriptor = fmt.Sprintf("wsh(%s)", keyStr)
	case DescriptorTypeUnknown:
		return "", errors.New("unknown descriptor type")
	default:
		return "", fmt.Errorf("unsupported descriptor type: %s", descType.String())
	}

	return descriptor, nil
}

// BuildDescriptorWithChecksum constructs a descriptor string with BIP380 checksum.
//
// This is a convenience method that combines BuildDescriptor and CalculateChecksum.
//
// Parameters:
//   - descType: The descriptor type
//   - key: The descriptor key
//
// Returns:
//   - Complete descriptor string WITH checksum (format: descriptor#checksum)
//   - Error if the components are invalid or checksum calculation fails
func (b *DescriptorBuilder) BuildDescriptorWithChecksum(descType DescriptorType, key DescriptorKey) (string, error) {
	// Build descriptor without checksum
	descriptor, err := b.BuildDescriptor(descType, key)
	if err != nil {
		return "", err
	}

	// Calculate and append checksum
	checksum, err := b.CalculateChecksum(descriptor)
	if err != nil {
		return "", fmt.Errorf("failed to calculate checksum: %w", err)
	}

	return descriptor + "#" + checksum, nil
}

// FormatDescriptorWithRange adds derivation path range to a descriptor.
//
// This appends a wildcard derivation path to the descriptor's key, allowing
// the descriptor to represent multiple addresses in a hierarchical deterministic wallet.
//
// Parameters:
//   - descriptor: Base descriptor string (can be with or without checksum)
//   - change: Change index (typically 0 for external, 1 for internal addresses)
//   - wildcard: Set to true to use /* wildcard, false for specific index
//
// Returns:
//   - Descriptor with range appended
//   - Error if the descriptor is malformed
//
// Example:
//
//	input:  "sh(wpkh([a1b2c3d4/49'/0'/0']xpub...))"
//	output: "sh(wpkh([a1b2c3d4/49'/0'/0']xpub.../0/*))"  (for change=0, wildcard=true)
//
// Note: If the descriptor already has a checksum, it will be removed and must be recalculated.
func (*DescriptorBuilder) FormatDescriptorWithRange(descriptor string, change uint32, wildcard bool) (string, error) {
	if strings.TrimSpace(descriptor) == "" {
		return "", errors.New("descriptor is empty")
	}

	// Remove checksum if present (it will need to be recalculated after adding range)
	if idx := strings.LastIndex(descriptor, "#"); idx != -1 {
		descriptor = descriptor[:idx]
	}

	// Find the last occurrence of xpub/tpub/etc. and insert path after it
	// This is a simplified approach - a full implementation would parse the descriptor
	xpubIdx := -1
	for _, prefix := range []string{"xpub", "tpub", "ypub", "zpub", "Ypub", "Zpub"} {
		if idx := strings.LastIndex(descriptor, prefix); idx > xpubIdx {
			xpubIdx = idx
		}
	}

	if xpubIdx == -1 {
		return "", errors.New("no extended public key found in descriptor")
	}

	// Find the end of the xpub (next non-base58 character)
	xpubEnd := xpubIdx
	for i := xpubIdx; i < len(descriptor); i++ {
		c := descriptor[i]
		// Base58 characters: 1-9, A-Z, a-z (excluding 0, O, I, l)
		if (c < '1' || c > '9') && (c < 'A' || c > 'Z') && (c < 'a' || c > 'z') {
			xpubEnd = i
			break
		}
		if i == len(descriptor)-1 {
			xpubEnd = len(descriptor)
			break
		}
	}

	// Build range suffix
	var rangeSuffix string
	if wildcard {
		rangeSuffix = fmt.Sprintf("/%d/*", change)
	} else {
		rangeSuffix = fmt.Sprintf("/%d", change)
	}

	// Insert range after xpub
	result := descriptor[:xpubEnd] + rangeSuffix + descriptor[xpubEnd:]

	return result, nil
}

// CalculateChecksum calculates the BIP380 descriptor checksum.
//
// The checksum algorithm is defined in BIP380 and uses a modified version of
// the bech32 checksum algorithm. The checksum is 8 characters long and uses
// the bech32 character set.
//
// Algorithm:
//  1. Expand each character to its numeric value based on inputCharset
//  2. Compute checksum using polymod with generator polynomial
//  3. Encode checksum using checksumCharset
//
// Parameters:
//   - descriptor: Descriptor string WITHOUT checksum
//
// Returns:
//   - 8-character checksum string
//   - Error if the descriptor contains invalid characters
//
// Reference: BIP380 - Output Script Descriptors General Operation
func (*DescriptorBuilder) CalculateChecksum(descriptor string) (string, error) {
	if strings.TrimSpace(descriptor) == "" {
		return "", errors.New("descriptor is empty")
	}

	// Convert descriptor to numeric values
	values, err := descriptorToValues(descriptor)
	if err != nil {
		return "", fmt.Errorf("failed to convert descriptor to values: %w", err)
	}

	// Calculate checksum using polymod
	checksum := descriptorChecksum(values)

	// Encode checksum as 8 characters using checksumCharset
	result := make([]byte, 8)
	for i := range 8 {
		result[i] = checksumCharset[(checksum>>(5*(7-i)))&31]
	}

	return string(result), nil
}

// descriptorToValues converts a descriptor string to numeric values for checksum calculation.
//
// Each character in the descriptor is mapped to a numeric value based on its position
// in inputCharset. Characters are then expanded according to BIP380 rules.
func descriptorToValues(descriptor string) ([]uint64, error) {
	values := make([]uint64, 0, len(descriptor))

	for _, c := range descriptor {
		// Find character position in inputCharset
		idx := strings.IndexRune(inputCharset, c)
		if idx == -1 {
			return nil, fmt.Errorf("invalid character in descriptor: %q", c)
		}

		// Expand character according to BIP380 rules
		// Characters are expanded to: (high 3 bits, low 5 bits)
		values = append(values, uint64(idx>>5))
		values = append(values, uint64(idx&31))
	}

	return values, nil
}

// descriptorChecksum calculates the descriptor checksum using polymod.
//
// This implements the BIP380 checksum algorithm using a generator polynomial
// similar to bech32 but with different constants.
//
// Generator polynomial values for descriptor checksums (from BIP380):
//
//	0xf5dee51989, 0xa9fdca3312, 0x1bab10e32d, 0x3706b1677a, 0x644d626ffd
func descriptorChecksum(values []uint64) uint64 {
	// Generator polynomial constants (from BIP380)
	const (
		generator0 = 0xf5dee51989
		generator1 = 0xa9fdca3312
		generator2 = 0x1bab10e32d
		generator3 = 0x3706b1677a
		generator4 = 0x644d626ffd
	)

	// Initialize checksum to 1
	c := uint64(1)

	// Process each value
	for _, v := range values {
		// Get top 5 bits of current checksum
		c0 := c >> 35

		// Shift and XOR with new value
		c = ((c & 0x7ffffffff) << 5) ^ v

		// XOR with generator polynomial terms based on top bits
		if c0&1 != 0 {
			c ^= generator0
		}
		if c0&2 != 0 {
			c ^= generator1
		}
		if c0&4 != 0 {
			c ^= generator2
		}
		if c0&8 != 0 {
			c ^= generator3
		}
		if c0&16 != 0 {
			c ^= generator4
		}
	}

	// Final XOR with 1 as specified in BIP380
	return c ^ 1
}
