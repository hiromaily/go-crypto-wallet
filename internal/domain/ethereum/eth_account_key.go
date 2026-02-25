package ethereum

import (
	"errors"
	"fmt"
	"time"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
)

// ETHAccountKey represents an Ethereum account key with associated address in the domain layer.
//
// SECURITY NOTE: The PrivateKey and AccountExtendedPrivkey fields contain sensitive key material.
// These must NEVER be logged or exposed in error messages. Handle with extreme care.
type ETHAccountKey struct {
	ID                     int64
	Account                domainAccount.AccountType
	Address                string // Ethereum address (0x...)
	FullPublicKey          string // Full public key (uncompressed, 65 bytes hex)
	PrivateKey             string // Private key (hex encoded) - NEVER log this field
	AccountExtendedPrivkey string // BIP32 account-level extended private key - NEVER log this field
	Idx                    int64  // Index for HD wallet
	AddrStatus             domainAddress.AddrStatus
	UpdatedAt              *time.Time
}

// NewETHAccountKey creates a new ETHAccountKey entity for key generation.
//
// SECURITY: The privateKey parameter contains the private key in hexadecimal format.
// Ensure this value is handled securely and never logged.
func NewETHAccountKey(
	account domainAccount.AccountType,
	address string,
	fullPublicKey string,
	privateKey string,
	idx int64,
) (*ETHAccountKey, error) {
	if address == "" {
		return nil, errors.New("address cannot be empty")
	}
	if fullPublicKey == "" {
		return nil, errors.New("full public key cannot be empty")
	}
	if privateKey == "" {
		return nil, errors.New("private key cannot be empty")
	}

	now := time.Now()
	return &ETHAccountKey{
		Account:       account,
		Address:       address,
		FullPublicKey: fullPublicKey,
		PrivateKey:    privateKey,
		Idx:           idx,
		AddrStatus:    domainAddress.AddrStatusHDKeyGenerated,
		UpdatedAt:     &now,
	}, nil
}

// UpdateAddrStatus updates the address status.
func (k *ETHAccountKey) UpdateAddrStatus(status domainAddress.AddrStatus) {
	k.AddrStatus = status
	k.updateTimestamp()
}

// GetPrivateKey returns the private key (hexadecimal format).
// SECURITY: This returns sensitive private key data. Never log the return value.
func (k *ETHAccountKey) GetPrivateKey() string {
	return k.PrivateKey
}

// SetAccountExtendedPrivkey sets the BIP32 account-level extended private key.
// SECURITY: This stores sensitive private key material. Never log the value parameter.
func (k *ETHAccountKey) SetAccountExtendedPrivkey(xpriv string) {
	k.AccountExtendedPrivkey = xpriv
}

// DeriveAccountXpub derives the account-level extended public key (xpub) from the stored xpriv.
// The returned xpub can be shared with the Watch wallet for child address derivation
// without exposing any private key material.
func (k *ETHAccountKey) DeriveAccountXpub() (string, error) {
	if k.AccountExtendedPrivkey == "" {
		return "", errors.New("accountXpriv not found for this account; run key generation first")
	}
	extKey, err := hdkeychain.NewKeyFromString(k.AccountExtendedPrivkey)
	if err != nil {
		return "", fmt.Errorf("failed to parse accountXpriv: %w", err)
	}
	pubKey, err := extKey.Neuter()
	if err != nil {
		return "", fmt.Errorf("failed to derive accountXpub from accountXpriv: %w", err)
	}
	return pubKey.String(), nil
}

func (k *ETHAccountKey) updateTimestamp() {
	now := time.Now()
	k.UpdatedAt = &now
}
