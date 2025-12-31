package key

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"

	domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
)

// CalculateMasterFingerprint calculates the BIP32 fingerprint of a master public key.
//
// According to BIP32, the fingerprint is calculated as:
//  1. Serialize the public key in compressed format (33 bytes)
//  2. Calculate HASH160 (SHA256 then RIPEMD160)
//  3. Take the first 4 bytes
//  4. Encode as lowercase hex string (8 characters)
//
// The fingerprint uniquely identifies an HD wallet master key and is used
// in descriptor syntax to specify which wallet a key belongs to.
//
// Example: "a1b2c3d4"
func CalculateMasterFingerprint(masterPubKey *btcec.PublicKey) (domainKey.Fingerprint, error) {
	if masterPubKey == nil {
		return "", errors.New("master public key is nil")
	}

	// Serialize the public key (compressed format: 33 bytes)
	pubKeyBytes := masterPubKey.SerializeCompressed()

	// Calculate HASH160 (SHA256 then RIPEMD160)
	hash := btcutil.Hash160(pubKeyBytes)

	// Take first 4 bytes
	fingerprint := hash[:4]

	// Encode as lowercase hex string
	hexString := hex.EncodeToString(fingerprint)

	// Validate and create domain fingerprint
	return domainKey.NewFingerprint(hexString)
}

// FingerprintFromExtendedKey extracts the fingerprint from an extended public key (xpub).
//
// This function:
//  1. Parses the extended public key string
//  2. Extracts the public key
//  3. Calculates the fingerprint using CalculateMasterFingerprint
//
// Note: For master keys, this calculates the fingerprint directly.
// For derived keys, this calculates the fingerprint of that specific key,
// not the master key. To get the master fingerprint, you need the master key.
func FingerprintFromExtendedKey(extendedKey string) (domainKey.Fingerprint, error) {
	if strings.TrimSpace(extendedKey) == "" {
		return "", errors.New("extended key is empty")
	}

	// Parse the extended key
	key, err := hdkeychain.NewKeyFromString(extendedKey)
	if err != nil {
		return "", fmt.Errorf("failed to parse extended key: %w", err)
	}

	// Get the public key
	pubKey, err := key.ECPubKey()
	if err != nil {
		return "", fmt.Errorf("failed to get public key from extended key: %w", err)
	}

	// Calculate fingerprint
	return CalculateMasterFingerprint(pubKey)
}

// FingerprintFromSeed calculates the master fingerprint from a seed.
//
// This function:
//  1. Derives the master key from the seed
//  2. Extracts the master public key
//  3. Calculates the fingerprint
//
// This is the most common way to get the master fingerprint for an HD wallet.
func FingerprintFromSeed(seed []byte, params *chaincfg.Params) (domainKey.Fingerprint, error) {
	if len(seed) == 0 {
		return "", errors.New("seed is empty")
	}

	if params == nil {
		return "", errors.New("network params is nil")
	}

	// Create master key from seed
	masterKey, err := hdkeychain.NewMaster(seed, params)
	if err != nil {
		return "", fmt.Errorf("failed to create master key from seed: %w", err)
	}

	// Get the master public key
	pubKey, err := masterKey.ECPubKey()
	if err != nil {
		return "", fmt.Errorf("failed to get master public key: %w", err)
	}

	// Calculate fingerprint
	return CalculateMasterFingerprint(pubKey)
}
