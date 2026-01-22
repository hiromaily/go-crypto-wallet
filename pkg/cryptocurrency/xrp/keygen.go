// Package xrp provides XRP (Ripple) cryptocurrency utilities including
// key generation, address encoding, and cryptographic operations.
//
// This package implements fully offline key generation for XRP,
// supporting both Ed25519 and secp256k1 algorithms.
package xrp

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
)

// KeyAlgorithm represents the cryptographic algorithm used for XRP keys.
type KeyAlgorithm int

const (
	// AlgorithmSecp256k1 uses the secp256k1 elliptic curve (same as Bitcoin).
	AlgorithmSecp256k1 KeyAlgorithm = iota
	// AlgorithmEd25519 uses the Ed25519 signature scheme (recommended for new accounts).
	AlgorithmEd25519
)

// String returns the string representation of the key algorithm.
func (a KeyAlgorithm) String() string {
	switch a {
	case AlgorithmSecp256k1:
		return "secp256k1"
	case AlgorithmEd25519:
		return "ed25519"
	default:
		return "unknown"
	}
}

// ParseKeyAlgorithm parses a string into a KeyAlgorithm.
// Returns an error if the algorithm string is not recognized.
// Valid values are "secp256k1" and "ed25519".
func ParseKeyAlgorithm(s string) (KeyAlgorithm, error) {
	switch s {
	case "secp256k1":
		return AlgorithmSecp256k1, nil
	case "ed25519":
		return AlgorithmEd25519, nil
	default:
		// Return -1 as invalid value to make it clear the parse failed
		return KeyAlgorithm(-1), fmt.Errorf("unknown key algorithm: %s (valid values: secp256k1, ed25519)", s)
	}
}

// XRPKeyPair represents a complete XRP key pair with all derived values.
//
// SECURITY NOTE: The Seed, SeedHex, and PrivateKey fields contain sensitive
// key material. These must NEVER be logged or exposed in error messages.
type XRPKeyPair struct {
	// Seed is the master seed in XRP's standard format (s... prefix)
	Seed string
	// SeedHex is the master seed in hexadecimal format
	SeedHex string
	// PublicKey is the public key in XRP's standard format
	PublicKey string
	// PublicKeyHex is the public key in hexadecimal format
	PublicKeyHex string
	// ClassicAddress is the XRP address in classic format (r... prefix)
	ClassicAddress string
	// Algorithm is the cryptographic algorithm used
	Algorithm KeyAlgorithm
}

// KeyGenerator provides methods for generating XRP key pairs.
type KeyGenerator struct {
	algorithm KeyAlgorithm
}

// NewKeyGenerator creates a new KeyGenerator with the specified algorithm.
func NewKeyGenerator(algorithm KeyAlgorithm) *KeyGenerator {
	return &KeyGenerator{algorithm: algorithm}
}

// GenerateRandom generates a new random XRP key pair.
//
// This function generates cryptographically secure random entropy
// and derives an XRP key pair from it.
func (g *KeyGenerator) GenerateRandom() (*XRPKeyPair, error) {
	// Generate 16 bytes of random entropy for the seed
	entropy := make([]byte, 16)
	if _, err := rand.Read(entropy); err != nil {
		return nil, fmt.Errorf("failed to generate random entropy: %w", err)
	}

	return g.GenerateFromEntropy(entropy)
}

// GenerateFromEntropy generates an XRP key pair from the given entropy.
//
// The entropy should be 16 bytes (128 bits) of random data.
// This allows for deterministic key generation from a known seed.
func (g *KeyGenerator) GenerateFromEntropy(entropy []byte) (*XRPKeyPair, error) {
	if len(entropy) < 16 {
		return nil, errors.New("entropy must be at least 16 bytes")
	}

	// Use first 16 bytes as the seed
	seedBytes := entropy[:16]

	switch g.algorithm {
	case AlgorithmSecp256k1:
		return g.generateSecp256k1(seedBytes)
	case AlgorithmEd25519:
		return g.generateEd25519(seedBytes)
	default:
		return nil, fmt.Errorf("unsupported algorithm: %d", g.algorithm)
	}
}

// deriveSecp256k1PrivateKey derives a secp256k1 private key from seed bytes.
// This is used by both key generation and signing.
//
// The derivation follows XRP's method:
// 1. SHA512Half of the seed bytes
// 2. SHA512Half of (seedHash || sequence) where sequence is 0 for root key
func deriveSecp256k1PrivateKey(seedBytes []byte) []byte {
	// Generate the family seed hash
	seedHash := Sha512Half(seedBytes)

	// Derive the private key from the seed
	// XRP uses SHA512Half of (seed || sequence) for key derivation
	// For the root key, sequence is 0
	sequence := make([]byte, 4)
	keyGenData := make([]byte, 0, len(seedHash)+len(sequence))
	keyGenData = append(keyGenData, seedHash...)
	keyGenData = append(keyGenData, sequence...)
	return Sha512Half(keyGenData)
}

// generateSecp256k1 generates an XRP key pair using secp256k1.
func (*KeyGenerator) generateSecp256k1(seedBytes []byte) (*XRPKeyPair, error) {
	privateKeyBytes := deriveSecp256k1PrivateKey(seedBytes)

	// Create secp256k1 private key
	// btcec.PrivKeyFromBytes returns (*PrivateKey, *PublicKey)
	privKey, _ := btcec.PrivKeyFromBytes(privateKeyBytes)
	if privKey == nil {
		return nil, errors.New("failed to create secp256k1 private key from seed")
	}

	// Get compressed public key (33 bytes)
	pubKeyBytes := privKey.PubKey().SerializeCompressed()

	// Generate account ID (address hash)
	accountID := Sha256RipeMD160(pubKeyBytes)

	// Create the seed in XRP format
	seed, err := NewFamilySeed(seedBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create family seed: %w", err)
	}

	// Create the account address
	address, err := NewAccountID(accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to create account ID: %w", err)
	}

	// Create public key hash (for XRP public key format)
	pubKeyHash, err := NewAccountPublicKey(pubKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create account public key: %w", err)
	}

	return &XRPKeyPair{
		Seed:           seed.String(),
		SeedHex:        hex.EncodeToString(seedBytes),
		PublicKey:      pubKeyHash.String(),
		PublicKeyHex:   hex.EncodeToString(pubKeyBytes),
		ClassicAddress: address.String(),
		Algorithm:      AlgorithmSecp256k1,
	}, nil
}

// generateEd25519 generates an XRP key pair using Ed25519.
func (*KeyGenerator) generateEd25519(seedBytes []byte) (*XRPKeyPair, error) {
	// For Ed25519, XRP uses SHA512Half of the seed to derive the private key
	privateKeyBytes := Sha512Half(seedBytes)

	// Ed25519 uses 32 bytes for the seed/private key
	edPrivateKey := ed25519.NewKeyFromSeed(privateKeyBytes)
	edPublicKey := edPrivateKey.Public().(ed25519.PublicKey)

	// XRP Ed25519 public keys are prefixed with 0xED
	pubKeyWithPrefix := append([]byte{0xED}, edPublicKey...)

	// Generate account ID (address hash) from prefixed public key
	accountID := Sha256RipeMD160(pubKeyWithPrefix)

	// Create the seed in XRP format
	seed, err := NewFamilySeed(seedBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create family seed: %w", err)
	}

	// Create the account address
	address, err := NewAccountID(accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to create account ID: %w", err)
	}

	// Create public key hash (for XRP public key format)
	pubKeyHash, err := NewAccountPublicKey(pubKeyWithPrefix)
	if err != nil {
		return nil, fmt.Errorf("failed to create account public key: %w", err)
	}

	return &XRPKeyPair{
		Seed:           seed.String(),
		SeedHex:        hex.EncodeToString(seedBytes),
		PublicKey:      pubKeyHash.String(),
		PublicKeyHex:   hex.EncodeToString(pubKeyWithPrefix),
		ClassicAddress: address.String(),
		Algorithm:      AlgorithmEd25519,
	}, nil
}

// GenerateFromPassphrase generates an XRP key pair from a passphrase.
//
// This provides deterministic key generation compatible with rippled's
// wallet_propose passphrase functionality.
//
// SECURITY NOTE: Using a passphrase for key generation is less secure
// than using random entropy. Only use this for testing or when
// deterministic key generation is specifically required.
func (g *KeyGenerator) GenerateFromPassphrase(passphrase string) (*XRPKeyPair, error) {
	if passphrase == "" {
		return nil, errors.New("passphrase cannot be empty")
	}

	// XRP derives the seed from SHA512Quarter of the passphrase
	// This matches rippled's wallet_propose behavior
	seedBytes := Sha512Quarter([]byte(passphrase))

	return g.GenerateFromEntropy(seedBytes)
}

// DeriveKeyFromHDKey generates an XRP key pair from an HD wallet private key.
//
// This is useful for integrating XRP key generation with BIP32/BIP44
// hierarchical deterministic wallets.
func (g *KeyGenerator) DeriveKeyFromHDKey(hdPrivateKey []byte) (*XRPKeyPair, error) {
	if len(hdPrivateKey) < 32 {
		return nil, errors.New("HD private key must be at least 32 bytes")
	}

	// Use first 16 bytes of HD private key hash as entropy
	entropy := Sha512Quarter(hdPrivateKey)

	return g.GenerateFromEntropy(entropy)
}
