package xrp

import (
	"errors"
	"time"
)

// XRPRegularKey represents a regular key assignment for an XRP account.
// Regular keys allow accounts to sign transactions without using the master key,
// providing enhanced security by keeping the master key offline.
//
// Reference: https://xrpl.org/docs/concepts/accounts/cryptographic-keys#regular-key-pair
type XRPRegularKey struct {
	ID                int64
	AccountID         string     // XRP account address (r...) that owns this regular key
	RegularKeyAddress string     // Regular key address (r...) authorized to sign for account
	PublicKey         string     // Regular key public key
	PublicKeyHex      string     // Regular key public key in hex format
	IsActive          bool       // True if this regular key is currently active for signing
	SetTxHash         *string    // Transaction hash of SetRegularKey that activated this key
	CreatedAt         *time.Time // When the regular key was created
	RotatedAt         *time.Time // When this key was rotated out (set inactive)
}

// NewXRPRegularKey creates a new XRPRegularKey entity.
func NewXRPRegularKey(
	accountID string,
	regularKeyAddress string,
	publicKey string,
	publicKeyHex string,
) (*XRPRegularKey, error) {
	if accountID == "" {
		return nil, errors.New("account ID cannot be empty")
	}
	if regularKeyAddress == "" {
		return nil, errors.New("regular key address cannot be empty")
	}
	if publicKey == "" {
		return nil, errors.New("public key cannot be empty")
	}
	if publicKeyHex == "" {
		return nil, errors.New("public key hex cannot be empty")
	}

	now := time.Now()
	return &XRPRegularKey{
		AccountID:         accountID,
		RegularKeyAddress: regularKeyAddress,
		PublicKey:         publicKey,
		PublicKeyHex:      publicKeyHex,
		IsActive:          true,
		CreatedAt:         &now,
	}, nil
}

// SetTxHashValue sets the transaction hash that activated this regular key.
func (k *XRPRegularKey) SetTxHashValue(txHash string) {
	k.SetTxHash = &txHash
}

// Deactivate marks this regular key as inactive (rotated out).
func (k *XRPRegularKey) Deactivate() {
	k.IsActive = false
	now := time.Now()
	k.RotatedAt = &now
}

// Activate marks this regular key as active.
func (k *XRPRegularKey) Activate() {
	k.IsActive = true
	k.RotatedAt = nil
}
