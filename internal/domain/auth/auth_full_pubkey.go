package auth

import (
	"errors"
	"time"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
)

// AuthFullPubkey represents an authentication full public key in the domain layer.
// This stores the full public key for authentication accounts used in multisig operations.
type AuthFullPubkey struct {
	ID            int16
	CoinTypeCode  domainCoin.CoinTypeCode
	AuthAccount   domainAccount.AuthType
	FullPublicKey string
	Fingerprint   *domainKey.Fingerprint // BIP32 master key fingerprint - nullable
	UpdatedAt     *time.Time
}

// NewAuthFullPubkey creates a new AuthFullPubkey entity.
func NewAuthFullPubkey(
	coinTypeCode domainCoin.CoinTypeCode,
	authAccount domainAccount.AuthType,
	fullPublicKey string,
) (*AuthFullPubkey, error) {
	if fullPublicKey == "" {
		return nil, errors.New("full public key cannot be empty")
	}

	now := time.Now()
	return &AuthFullPubkey{
		CoinTypeCode:  coinTypeCode,
		AuthAccount:   authAccount,
		FullPublicKey: fullPublicKey,
		UpdatedAt:     &now,
	}, nil
}

// SetFingerprint sets the BIP32 master key fingerprint.
func (k *AuthFullPubkey) SetFingerprint(fingerprint domainKey.Fingerprint) {
	k.Fingerprint = &fingerprint
	k.updateTimestamp()
}

func (k *AuthFullPubkey) updateTimestamp() {
	now := time.Now()
	k.UpdatedAt = &now
}
