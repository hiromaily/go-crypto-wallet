package auth

import (
	"time"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
)

// AuthFullPubkey represents an authentication full public key in the domain layer.
// This stores the full public key for authentication accounts used in multisig operations.
// Supports both legacy format (compressed pubkey only) and extended format (xpub + metadata).
type AuthFullPubkey struct {
	ID             int16
	CoinTypeCode   domainCoin.CoinTypeCode
	AuthAccount    domainAccount.AuthType
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
	fullPublicKey string,
) (*AuthFullPubkey, error) {
	// Allow empty fullPublicKey for extended format (extended key will be set separately)
	// Legacy format requires compressed pubkey

	now := time.Now()
	return &AuthFullPubkey{
		CoinTypeCode:  coinTypeCode,
		AuthAccount:   authAccount,
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
