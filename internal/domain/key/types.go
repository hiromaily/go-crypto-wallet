package key

import "fmt"

// KeyType represents the type of key generation standard
type KeyType string

const (
	// KeyTypeBIP44 represents BIP44 (Legacy P2PKH)
	KeyTypeBIP44 KeyType = "bip44"

	// KeyTypeBIP49 represents BIP49 (P2SH-SegWit)
	KeyTypeBIP49 KeyType = "bip49"

	// KeyTypeBIP84 represents BIP84 (Native SegWit P2WPKH)
	KeyTypeBIP84 KeyType = "bip84"

	// KeyTypeBIP86 represents BIP86 (Taproot)
	KeyTypeBIP86 KeyType = "bip86"

	// KeyTypeMuSig2 represents MuSig2 aggregated keys
	KeyTypeMuSig2 KeyType = "musig2"
)

// String returns the string representation of the key type
func (k KeyType) String() string {
	return string(k)
}

// Purpose returns the BIP purpose number for the key type
// Panics if called on key types that don't have a BIP purpose number (e.g., MuSig2)
// or on unhandled key types (programming error)
func (k KeyType) Purpose() uint32 {
	switch k {
	case KeyTypeBIP44:
		return 44
	case KeyTypeBIP49:
		return 49
	case KeyTypeBIP84:
		return 84
	case KeyTypeBIP86:
		return 86
	case KeyTypeMuSig2:
		// MuSig2 is a signing scheme, not a derivation path standard, and has no purpose number.
		panic(fmt.Sprintf("key type %q does not have a BIP purpose number", k))
	default:
		// This case should be unreachable if Validate() is called, but as a safeguard:
		panic(fmt.Sprintf("unhandled key type: %s", k))
	}
}

// Validate validates the key type
func (k KeyType) Validate() error {
	switch k {
	case KeyTypeBIP44, KeyTypeBIP49, KeyTypeBIP84, KeyTypeBIP86, KeyTypeMuSig2:
		return nil
	default:
		return fmt.Errorf("invalid key type: %s", k)
	}
}
