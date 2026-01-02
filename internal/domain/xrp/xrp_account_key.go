package xrp

import (
	"errors"
	"time"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
)

// XRPKeyType represents XRP key types
type XRPKeyType int8

const (
	XRPKeyTypeSecp256k1 XRPKeyType = 0
	XRPKeyTypeEd25519   XRPKeyType = 1
)

// XRPAccountKey represents an XRP account key with associated address in the domain layer.
//
// SECURITY NOTE: The MasterSeed and MasterSeedHex fields contain sensitive key material.
// These must NEVER be logged or exposed in error messages. Handle with extreme care.
type XRPAccountKey struct {
	ID               int64
	CoinTypeCode     domainCoin.CoinTypeCode
	Account          domainAccount.AccountType
	AccountID        string     // XRP account ID (address)
	KeyType          XRPKeyType // Key type: secp256k1 or ed25519
	MasterKey        string     // DEPRECATED field
	MasterSeed       string     // Master seed - NEVER log this field
	MasterSeedHex    string     // Master seed (hex) - NEVER log this field
	PublicKey        string     // Public key
	PublicKeyHex     string     // Public key (hex)
	IsRegularKeyPair bool       // True if this key is for regular key pair
	AllocatedID      int64      // Index for HD wallet
	AddrStatus       domainAddress.AddrStatus
	UpdatedAt        *time.Time
}

// NewXRPAccountKey creates a new XRPAccountKey entity for key generation.
//
// SECURITY: The masterSeed and masterSeedHex parameters contain sensitive key material.
// Ensure these values are handled securely and never logged.
func NewXRPAccountKey(
	coinTypeCode domainCoin.CoinTypeCode,
	account domainAccount.AccountType,
	accountID string,
	keyType XRPKeyType,
	masterSeed string,
	masterSeedHex string,
	publicKey string,
	publicKeyHex string,
	isRegularKeyPair bool,
	allocatedID int64,
) (*XRPAccountKey, error) {
	if accountID == "" {
		return nil, errors.New("account ID cannot be empty")
	}
	if masterSeed == "" {
		return nil, errors.New("master seed cannot be empty")
	}
	if masterSeedHex == "" {
		return nil, errors.New("master seed hex cannot be empty")
	}
	if publicKey == "" {
		return nil, errors.New("public key cannot be empty")
	}

	now := time.Now()
	return &XRPAccountKey{
		CoinTypeCode:     coinTypeCode,
		Account:          account,
		AccountID:        accountID,
		KeyType:          keyType,
		MasterSeed:       masterSeed,
		MasterSeedHex:    masterSeedHex,
		PublicKey:        publicKey,
		PublicKeyHex:     publicKeyHex,
		IsRegularKeyPair: isRegularKeyPair,
		AllocatedID:      allocatedID,
		AddrStatus:       domainAddress.AddrStatusHDKeyGenerated,
		UpdatedAt:        &now,
	}, nil
}

// UpdateAddrStatus updates the address status.
func (k *XRPAccountKey) UpdateAddrStatus(status domainAddress.AddrStatus) {
	k.AddrStatus = status
	k.updateTimestamp()
}

// GetMasterSeed returns the master seed.
// SECURITY: This returns sensitive key material. Never log the return value.
func (k *XRPAccountKey) GetMasterSeed() string {
	return k.MasterSeed
}

// GetMasterSeedHex returns the master seed in hexadecimal format.
// SECURITY: This returns sensitive key material. Never log the return value.
func (k *XRPAccountKey) GetMasterSeedHex() string {
	return k.MasterSeedHex
}

func (k *XRPAccountKey) updateTimestamp() {
	now := time.Now()
	k.UpdatedAt = &now
}
