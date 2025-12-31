package key

import (
	"errors"
	"fmt"
	"regexp"
)

// Fingerprint represents a BIP32 master key fingerprint.
//
// According to BIP32, a fingerprint is a 4-byte identifier derived from
// the HASH160 of the master public key. It uniquely identifies an HD wallet
// master key and is used in descriptor syntax to specify which wallet a key belongs to.
//
// Format: 8 hexadecimal characters (4 bytes)
// Example: "a1b2c3d4"
type Fingerprint string

// fingerprintRegex validates fingerprint format (8 hex characters)
var fingerprintRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}$`)

// NewFingerprint creates a new Fingerprint from a hex string.
// The fingerprint must be exactly 8 hexadecimal characters (4 bytes).
func NewFingerprint(hex string) (Fingerprint, error) {
	if err := ValidateFingerprint(hex); err != nil {
		return "", fmt.Errorf("invalid fingerprint: %w", err)
	}
	return Fingerprint(hex), nil
}

// String returns the fingerprint as a lowercase hex string.
func (f Fingerprint) String() string {
	return string(f)
}

// Bytes returns the fingerprint as a byte slice (4 bytes).
func (f Fingerprint) Bytes() ([]byte, error) {
	if len(f) != 8 {
		return nil, errors.New("fingerprint must be 8 hex characters")
	}

	bytes := make([]byte, 4)
	for i := 0; i < 4; i++ {
		var b byte
		_, err := fmt.Sscanf(string(f[i*2:i*2+2]), "%02x", &b)
		if err != nil {
			return nil, fmt.Errorf("invalid hex character at position %d: %w", i*2, err)
		}
		bytes[i] = b
	}
	return bytes, nil
}

// ValidateFingerprint validates a fingerprint string.
//
// A valid fingerprint must be:
//   - Exactly 8 characters long
//   - Contain only hexadecimal characters (0-9, a-f, A-F)
//
// This function is used by descriptor validation and is also available
// in the descriptor package for backward compatibility.
func ValidateFingerprint(fingerprint string) error {
	if len(fingerprint) != 8 {
		return fmt.Errorf("fingerprint must be exactly 8 characters, got %d", len(fingerprint))
	}

	if !fingerprintRegex.MatchString(fingerprint) {
		return errors.New("fingerprint must be 8 hexadecimal characters")
	}

	return nil
}
