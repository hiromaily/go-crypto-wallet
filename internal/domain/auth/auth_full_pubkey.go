package auth

import (
	"fmt"
	"time"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
)

// Purpose represents BIP purpose for HD wallet derivation paths.
// Each purpose corresponds to a different address type and script format.
type Purpose uint8

const (
	// PurposeBIP44 is for legacy P2PKH addresses (m/44')
	PurposeBIP44 Purpose = 44
	// PurposeBIP49 is for P2SH-wrapped SegWit addresses (m/49')
	PurposeBIP49 Purpose = 49
	// PurposeBIP84 is for native SegWit (Bech32) addresses (m/84')
	PurposeBIP84 Purpose = 84
	// PurposeBIP86 is for Taproot addresses (m/86')
	PurposeBIP86 Purpose = 86
)

// String returns string representation of Purpose.
func (p Purpose) String() string {
	switch p {
	case PurposeBIP44:
		return "BIP44"
	case PurposeBIP49:
		return "BIP49"
	case PurposeBIP84:
		return "BIP84"
	case PurposeBIP86:
		return "BIP86"
	default:
		return fmt.Sprintf("Unknown(%d)", p)
	}
}

// Validate validates the purpose value.
func (p Purpose) Validate() error {
	switch p {
	case PurposeBIP44, PurposeBIP49, PurposeBIP84, PurposeBIP86:
		return nil
	default:
		return fmt.Errorf("invalid BIP purpose: %d (must be 44, 49, 84, or 86)", p)
	}
}

// PurposeFromUint8 creates Purpose from uint8 with validation.
func PurposeFromUint8(value uint8) (Purpose, error) {
	p := Purpose(value)
	if err := p.Validate(); err != nil {
		return 0, err
	}
	return p, nil
}

//----------------------------------------------------
// Address Type to BIP Purpose Mapping
//----------------------------------------------------

// PurposeForAddressType returns the BIP purpose for a given address type string.
// This mapping is used to determine which extended public key (xpub) to use
// when generating descriptors for different address types.
//
// Mapping:
// - "legacy" (P2PKH) → BIP44 (m/44')
// - "p2sh-segwit" → BIP49 (m/49')
// - "bech32" (Native SegWit) → BIP84 (m/84')
// - "taproot" → BIP86 (m/86')
// - "bch-cashaddr" → BIP44 (m/44') - BCH uses BIP44 path
func PurposeForAddressType(addrType string) (Purpose, error) {
	switch addrType {
	case "legacy":
		return PurposeBIP44, nil
	case "p2sh-segwit":
		return PurposeBIP49, nil
	case "bech32":
		return PurposeBIP84, nil
	case "taproot":
		return PurposeBIP86, nil
	case "bch-cashaddr":
		// BCH uses BIP44 path (no SegWit/Taproot support)
		return PurposeBIP44, nil
	case "eth-address":
		// Ethereum doesn't use BIP44/49/84/86 purposes in the same way,
		// but we return BIP44 as a default for consistency
		return PurposeBIP44, nil
	default:
		return 0, fmt.Errorf("unknown address type: %s", addrType)
	}
}

// AuthFullPubkey represents an authentication full public key in the domain layer.
// This stores the full public key for authentication accounts used in multisig operations.
// Supports both legacy format (compressed pubkey only) and extended format (xpub + metadata).
type AuthFullPubkey struct {
	ID             int16
	CoinTypeCode   domainCoin.CoinTypeCode
	AuthAccount    domainAccount.AuthType
	Purpose        Purpose                // BIP purpose (44, 49, 84, 86) for address type
	FullPublicKey  string                 // Legacy: compressed public key (33 bytes hex)
	ExtendedPubKey string                 // Extended public key (xpub/tpub format) - optional, new format
	Fingerprint    *domainKey.Fingerprint // BIP32 master key fingerprint - nullable
	DerivationPath string                 // BIP32 derivation path (e.g., m/49'/1'/0') - empty for legacy
	UpdatedAt      *time.Time
}

// NewAuthFullPubkey creates a new AuthFullPubkey entity.
// For legacy format, fullPublicKey contains compressed public key.
// For extended format, fullPublicKey can be empty if extended key is set via SetExtendedPubKey.
func NewAuthFullPubkey(
	coinTypeCode domainCoin.CoinTypeCode,
	authAccount domainAccount.AuthType,
	purpose Purpose,
	fullPublicKey string,
) (*AuthFullPubkey, error) {
	// Validate purpose
	if err := purpose.Validate(); err != nil {
		return nil, fmt.Errorf("invalid purpose: %w", err)
	}

	// Allow empty fullPublicKey for extended format (extended key will be set separately)
	// Legacy format requires compressed pubkey

	now := time.Now()
	return &AuthFullPubkey{
		CoinTypeCode:  coinTypeCode,
		AuthAccount:   authAccount,
		Purpose:       purpose,
		FullPublicKey: fullPublicKey, // Can be empty for extended format
		UpdatedAt:     &now,
	}, nil
}

// SetFingerprint sets the BIP32 master key fingerprint.
func (k *AuthFullPubkey) SetFingerprint(fingerprint domainKey.Fingerprint) {
	k.Fingerprint = &fingerprint
	k.updateTimestamp()
}

// SetExtendedPubKey sets the BIP32 extended public key.
func (k *AuthFullPubkey) SetExtendedPubKey(extendedPubKey string) {
	k.ExtendedPubKey = extendedPubKey
	k.updateTimestamp()
}

// SetDerivationPath sets the BIP32 derivation path.
func (k *AuthFullPubkey) SetDerivationPath(derivationPath string) {
	k.DerivationPath = derivationPath
	k.updateTimestamp()
}

func (k *AuthFullPubkey) updateTimestamp() {
	now := time.Now()
	k.UpdatedAt = &now
}
