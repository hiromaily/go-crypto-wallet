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
		return 44 // MuSig2 doesn't have a BIP purpose number, default to BIP44
	default:
		return 44 // Default to BIP44
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
