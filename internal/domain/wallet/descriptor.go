package wallet

import (
	"fmt"
	"regexp"
	"strings"
)

// Descriptor represents a Bitcoin output descriptor as defined in BIP380.
//
// Descriptors are human-readable strings that describe how to derive addresses
// and scripts from public keys. They make wallet functionality explicit and portable.
//
// Example descriptors:
//   - P2PKH: pkh([fingerprint/44'/0'/0']xpub.../0/*)
//   - Bech32: wpkh([fingerprint/84'/0'/0']xpub.../0/*)
//   - Taproot: tr([fingerprint/86'/0'/0']xpub.../0/*)
type Descriptor struct {
	// Type is the descriptor type (PKH, SHWPKH, WPKH, TR, WSH)
	Type DescriptorType

	// Script is the full descriptor string
	Script string

	// Keys are the keys involved in the descriptor
	Keys []DescriptorKey

	// Checksum is the descriptor checksum (optional)
	Checksum string
}

// DescriptorType represents the type of output descriptor.
type DescriptorType int

const (
	// DescriptorTypePKH represents Pay-to-Public-Key-Hash (P2PKH) descriptors.
	// Format: pkh(KEY)
	// Example: pkh([fingerprint/44'/0'/0']xpub.../0/*)
	DescriptorTypePKH DescriptorType = iota

	// DescriptorTypeSHWPKH represents Pay-to-Script-Hash wrapping Pay-to-Witness-Public-Key-Hash.
	// Also known as P2SH-SegWit or Nested SegWit.
	// Format: sh(wpkh(KEY))
	// Example: sh(wpkh([fingerprint/49'/0'/0']xpub.../0/*))
	DescriptorTypeSHWPKH

	// DescriptorTypeWPKH represents Pay-to-Witness-Public-Key-Hash (Native SegWit, Bech32).
	// Format: wpkh(KEY)
	// Example: wpkh([fingerprint/84'/0'/0']xpub.../0/*)
	DescriptorTypeWPKH

	// DescriptorTypeTR represents Pay-to-Taproot (P2TR) descriptors.
	// Format: tr(KEY)
	// Example: tr([fingerprint/86'/0'/0']xpub.../0/*)
	DescriptorTypeTR

	// DescriptorTypeWSH represents Pay-to-Witness-Script-Hash descriptors.
	// Typically used for multisig: wsh(sortedmulti(M,KEY1,KEY2,...))
	// Example: wsh(sortedmulti(2,[fp1/48'/0'/0'/2']xpub1...,[fp2/48'/0'/0'/2']xpub2.../0/*))
	DescriptorTypeWSH

	// DescriptorTypeUnknown represents an unknown or unsupported descriptor type.
	DescriptorTypeUnknown
)

// String returns the string representation of the descriptor type.
func (d DescriptorType) String() string {
	switch d {
	case DescriptorTypePKH:
		return "pkh"
	case DescriptorTypeSHWPKH:
		return "sh(wpkh)"
	case DescriptorTypeWPKH:
		return "wpkh"
	case DescriptorTypeTR:
		return "tr"
	case DescriptorTypeWSH:
		return "wsh"
	default:
		return "unknown"
	}
}

// DescriptorKey represents a key in a descriptor with its metadata.
//
// Keys in descriptors can include:
//   - Fingerprint: 4-byte identifier of the master key (8 hex characters)
//   - DerivationPath: BIP32 derivation path (e.g., /44'/0'/0')
//   - ExtendedPubKey: BIP32 extended public key (xpub...)
type DescriptorKey struct {
	// Fingerprint is the 4-byte master key identifier (8 hex characters).
	// Example: "a1b2c3d4"
	Fingerprint string

	// DerivationPath is the BIP32 derivation path from the master key.
	// Example: "/44'/0'/0'" or "/84'/0'/0'/0/*"
	DerivationPath string

	// ExtendedPubKey is the BIP32 extended public key.
	// Example: "xpub6ERApfZwUNrhLCkDtcHTcxd75RbzS1ed54G1LkBUHQVHQKqhMkhgbmJbZRkrgZw4koxb5JaHWkY4ALHY2grBGRjaDMzQLcgJvLJuZZvRcEL"
	ExtendedPubKey string
}

// ValidateDescriptor validates a descriptor according to BIP380 rules.
//
// Validation checks:
//   - Descriptor syntax is well-formed
//   - Descriptor type is supported
//   - Keys are valid (fingerprint, derivation path, extended public key)
//   - Derivation paths follow BIP32 format
//   - Extended public keys are valid
//   - Checksum is valid (if present)
func ValidateDescriptor(desc *Descriptor) error {
	if desc == nil {
		return fmt.Errorf("descriptor is nil")
	}

	// Validate descriptor type
	if desc.Type == DescriptorTypeUnknown {
		return fmt.Errorf("unknown descriptor type")
	}

	// Validate script is not empty
	if strings.TrimSpace(desc.Script) == "" {
		return fmt.Errorf("descriptor script is empty")
	}

	// Validate keys
	if len(desc.Keys) == 0 {
		return fmt.Errorf("descriptor must have at least one key")
	}

	for i, key := range desc.Keys {
		if err := ValidateDescriptorKey(&key); err != nil {
			return fmt.Errorf("invalid key at index %d: %w", i, err)
		}
	}

	return nil
}

// ValidateDescriptorKey validates a descriptor key.
func ValidateDescriptorKey(key *DescriptorKey) error {
	if key == nil {
		return fmt.Errorf("descriptor key is nil")
	}

	// Validate fingerprint (8 hex characters)
	if key.Fingerprint != "" {
		if err := ValidateFingerprint(key.Fingerprint); err != nil {
			return fmt.Errorf("invalid fingerprint: %w", err)
		}
	}

	// Validate derivation path
	if key.DerivationPath != "" {
		if err := ValidateDerivationPath(key.DerivationPath); err != nil {
			return fmt.Errorf("invalid derivation path: %w", err)
		}
	}

	// Validate extended public key
	if key.ExtendedPubKey == "" {
		return fmt.Errorf("extended public key is empty")
	}

	if err := ValidateExtendedPubKey(key.ExtendedPubKey); err != nil {
		return fmt.Errorf("invalid extended public key: %w", err)
	}

	return nil
}

// ValidateFingerprint validates a master key fingerprint.
//
// Fingerprints must be exactly 8 hexadecimal characters (4 bytes).
// Example: "a1b2c3d4"
func ValidateFingerprint(fingerprint string) error {
	if len(fingerprint) != 8 {
		return fmt.Errorf("fingerprint must be exactly 8 characters, got %d", len(fingerprint))
	}

	fingerprintRegex := regexp.MustCompile(`^[0-9a-fA-F]{8}$`)
	if !fingerprintRegex.MatchString(fingerprint) {
		return fmt.Errorf("fingerprint must be 8 hexadecimal characters")
	}

	return nil
}

// ValidateDerivationPath validates a BIP32 derivation path.
//
// Valid paths:
//   - Must start with "/" or "m/"
//   - Segments are numbers optionally followed by ' or h (hardened)
//   - Can end with /* or /<number>/* for wildcards
//
// Examples:
//   - "/44'/0'/0'"
//   - "/84'/0'/0'/0/*"
//   - "m/86'/0'/0'"
func ValidateDerivationPath(path string) error {
	if path == "" {
		return fmt.Errorf("derivation path is empty")
	}

	// Remove leading "m/" if present
	path = strings.TrimPrefix(path, "m")

	// Path must start with /
	if !strings.HasPrefix(path, "/") {
		return fmt.Errorf("derivation path must start with /")
	}

	// Split into segments
	segments := strings.Split(path, "/")[1:] // Skip empty first element

	if len(segments) == 0 {
		return fmt.Errorf("derivation path has no segments")
	}

	// Regex for valid segment: number optionally followed by ' or h, or *
	segmentRegex := regexp.MustCompile(`^(\d+['h]?|\*)$`)

	for i, segment := range segments {
		if !segmentRegex.MatchString(segment) {
			return fmt.Errorf("invalid derivation path segment at position %d: %s", i, segment)
		}

		// Wildcard can only be in the last segment
		if segment == "*" && i != len(segments)-1 {
			return fmt.Errorf("wildcard * can only be in the last segment")
		}
	}

	return nil
}

// ValidateExtendedPubKey validates a BIP32 extended public key.
//
// Valid extended public keys:
//   - Start with xpub, tpub, ypub, zpub, Ypub, or Zpub
//   - Are base58-encoded
//   - Have correct length (111 characters for standard xpub)
func ValidateExtendedPubKey(xpub string) error {
	if xpub == "" {
		return fmt.Errorf("extended public key is empty")
	}

	// Check prefix
	validPrefixes := []string{"xpub", "tpub", "ypub", "zpub", "Ypub", "Zpub"}
	hasValidPrefix := false
	for _, prefix := range validPrefixes {
		if strings.HasPrefix(xpub, prefix) {
			hasValidPrefix = true
			break
		}
	}

	if !hasValidPrefix {
		return fmt.Errorf("extended public key must start with xpub, tpub, ypub, zpub, Ypub, or Zpub")
	}

	// Check length (extended keys are typically 111 characters)
	if len(xpub) < 100 || len(xpub) > 120 {
		return fmt.Errorf("extended public key has invalid length: %d", len(xpub))
	}

	// Check base58 characters (alphanumeric, no 0, O, I, l)
	base58Regex := regexp.MustCompile(`^[1-9A-HJ-NP-Za-km-z]+$`)
	if !base58Regex.MatchString(xpub) {
		return fmt.Errorf("extended public key contains invalid base58 characters")
	}

	return nil
}
