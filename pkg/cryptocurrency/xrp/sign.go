// Package xrp provides XRP (Ripple) cryptocurrency utilities including
// key generation, address encoding, and cryptographic operations.
package xrp

import (
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
)

// Signer handles native XRP transaction signing.
// This allows offline signing without requiring a gRPC connection.
type Signer struct {
	algorithm  KeyAlgorithm
	serializer *Serializer
}

// NewSigner creates a new XRP transaction signer with the specified algorithm.
func NewSigner(algorithm KeyAlgorithm) *Signer {
	return &Signer{
		algorithm:  algorithm,
		serializer: NewSerializer(),
	}
}

// SignResult contains the result of signing a transaction.
type SignResult struct {
	// TxBlob is the signed transaction in hex format
	TxBlob string
	// TxID is the transaction hash/ID
	TxID string
	// SigningPubKey is the public key used for signing (hex)
	SigningPubKey string
	// TxnSignature is the signature (hex)
	TxnSignature string
}

// SignTransaction signs an XRP transaction using the provided seed.
//
// Parameters:
//   - txFields: Transaction fields as a map (from JSON)
//   - seed: The master seed in XRP format (s... prefix) or hex format
//
// Returns the signed transaction blob, transaction ID, and an error if any.
//
// SECURITY NOTE: The seed contains sensitive key material. It must be
// handled securely and never logged.
func (s *Signer) SignTransaction(txFields map[string]any, seed string) (*SignResult, error) {
	if txFields == nil {
		return nil, errors.New("transaction fields cannot be nil")
	}
	if seed == "" {
		return nil, errors.New("seed cannot be empty")
	}

	// Derive the key pair from the seed
	keyPair, err := s.deriveKeyPair(seed)
	if err != nil {
		return nil, fmt.Errorf("failed to derive key pair: %w", err)
	}

	// Add the signing public key to the transaction
	txFields["SigningPubKey"] = keyPair.PublicKeyHex

	// Create transaction for serialization
	tx := NewTransaction(txFields)

	// Compute the signing hash
	signingHash, err := s.serializer.ComputeSigningHash(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to compute signing hash: %w", err)
	}

	// Sign the hash
	signature, err := s.sign(signingHash, keyPair)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	signatureHex := strings.ToUpper(hex.EncodeToString(signature))

	// Add signature to transaction
	txFields["TxnSignature"] = signatureHex

	// Serialize the complete signed transaction
	signedTx := NewTransaction(txFields)
	txBlob, err := s.serializer.SerializeForSubmission(signedTx)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize signed transaction: %w", err)
	}

	txBlobHex := strings.ToUpper(hex.EncodeToString(txBlob))

	// Compute transaction ID
	txID := ComputeTransactionID(txBlob)

	return &SignResult{
		TxBlob:        txBlobHex,
		TxID:          txID,
		SigningPubKey: strings.ToUpper(keyPair.PublicKeyHex),
		TxnSignature:  signatureHex,
	}, nil
}

// deriveKeyPair derives the XRP key pair from a seed.
// The seed can be in XRP format (s... prefix) or hex format.
func (s *Signer) deriveKeyPair(seed string) (*XRPKeyPair, error) {
	var seedBytes []byte
	var err error

	if strings.HasPrefix(seed, "s") {
		// XRP seed format (s... prefix)
		seedBytes, err = decodeSeed(seed)
		if err != nil {
			return nil, fmt.Errorf("failed to decode seed: %w", err)
		}
	} else {
		// Hex format
		seedBytes, err = hex.DecodeString(seed)
		if err != nil {
			return nil, fmt.Errorf("invalid hex seed: %w", err)
		}
	}

	// Generate key pair from seed bytes
	keyGen := NewKeyGenerator(s.algorithm)
	return keyGen.GenerateFromEntropy(seedBytes)
}

// decodeSeed decodes an XRP seed from s... format to raw bytes.
func decodeSeed(seed string) ([]byte, error) {
	hash, err := NewRippleHashCheck(seed, RippleFamilySeed)
	if err != nil {
		return nil, fmt.Errorf("invalid seed format: %w", err)
	}
	return hash.Payload(), nil
}

// sign signs a hash using the appropriate algorithm.
func (s *Signer) sign(hash []byte, keyPair *XRPKeyPair) ([]byte, error) {
	switch s.algorithm {
	case AlgorithmSecp256k1:
		return s.signSecp256k1(hash, keyPair)
	case AlgorithmEd25519:
		return s.signEd25519(hash, keyPair)
	default:
		return nil, fmt.Errorf("unsupported algorithm: %d", s.algorithm)
	}
}

// signSecp256k1 signs a hash using secp256k1 (ECDSA).
func (*Signer) signSecp256k1(hash []byte, keyPair *XRPKeyPair) ([]byte, error) {
	// Decode the seed to derive the private key
	seedBytes, err := hex.DecodeString(keyPair.SeedHex)
	if err != nil {
		return nil, fmt.Errorf("invalid seed hex: %w", err)
	}

	// Derive private key same as in keygen
	seedHash := Sha512Half(seedBytes)
	sequence := make([]byte, 4)
	keyGenData := make([]byte, 0, len(seedHash)+len(sequence))
	keyGenData = append(keyGenData, seedHash...)
	keyGenData = append(keyGenData, sequence...)
	privateKeyBytes := Sha512Half(keyGenData)

	// Create secp256k1 private key
	privKey, _ := btcec.PrivKeyFromBytes(privateKeyBytes)
	if privKey == nil {
		return nil, errors.New("failed to create secp256k1 private key")
	}

	// Sign the hash using ECDSA
	// XRP uses DER-encoded signatures
	signature := ecdsa.Sign(privKey, hash)
	return signature.Serialize(), nil
}

// signEd25519 signs a hash using Ed25519.
func (*Signer) signEd25519(hash []byte, keyPair *XRPKeyPair) ([]byte, error) {
	// Decode the seed to derive the private key
	seedBytes, err := hex.DecodeString(keyPair.SeedHex)
	if err != nil {
		return nil, fmt.Errorf("invalid seed hex: %w", err)
	}

	// Derive private key same as in keygen
	privateKeyBytes := Sha512Half(seedBytes)

	// Create Ed25519 private key
	edPrivateKey := ed25519.NewKeyFromSeed(privateKeyBytes)

	// Sign the hash
	// Ed25519 signatures are 64 bytes
	signature := ed25519.Sign(edPrivateKey, hash)
	return signature, nil
}

// VerifySignature verifies an XRP transaction signature.
//
// Parameters:
//   - txFields: Transaction fields including SigningPubKey and TxnSignature
//
// Returns true if the signature is valid.
func (s *Signer) VerifySignature(txFields map[string]any) (bool, error) {
	if txFields == nil {
		return false, errors.New("transaction fields cannot be nil")
	}

	// Get signature and public key
	sigHex, ok := txFields["TxnSignature"].(string)
	if !ok || sigHex == "" {
		return false, errors.New("missing TxnSignature")
	}

	pubKeyHex, ok := txFields["SigningPubKey"].(string)
	if !ok || pubKeyHex == "" {
		return false, errors.New("missing SigningPubKey")
	}

	signature, err := hex.DecodeString(sigHex)
	if err != nil {
		return false, fmt.Errorf("invalid signature hex: %w", err)
	}

	pubKey, err := hex.DecodeString(pubKeyHex)
	if err != nil {
		return false, fmt.Errorf("invalid public key hex: %w", err)
	}

	// Create transaction without signature for hash computation
	txCopy := make(map[string]any)
	for k, v := range txFields {
		if k != "TxnSignature" && k != "Hash" {
			txCopy[k] = v
		}
	}

	tx := NewTransaction(txCopy)
	signingHash, err := s.serializer.ComputeSigningHash(tx)
	if err != nil {
		return false, fmt.Errorf("failed to compute signing hash: %w", err)
	}

	// Verify based on algorithm (determined by public key prefix)
	if len(pubKey) == 33 && pubKey[0] == 0xED {
		// Ed25519 (prefixed with 0xED)
		return ed25519.Verify(pubKey[1:], signingHash, signature), nil
	}

	// secp256k1
	parsedPubKey, err := btcec.ParsePubKey(pubKey)
	if err != nil {
		return false, fmt.Errorf("failed to parse public key: %w", err)
	}

	parsedSig, err := ecdsa.ParseDERSignature(signature)
	if err != nil {
		return false, fmt.Errorf("failed to parse signature: %w", err)
	}

	return parsedSig.Verify(signingHash, parsedPubKey), nil
}

// DetectAlgorithm detects the key algorithm from a public key.
func DetectAlgorithm(publicKeyHex string) (KeyAlgorithm, error) {
	pubKey, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		return KeyAlgorithm(-1), fmt.Errorf("invalid public key hex: %w", err)
	}

	if len(pubKey) == 33 && pubKey[0] == 0xED {
		return AlgorithmEd25519, nil
	}

	if len(pubKey) == 33 && (pubKey[0] == 0x02 || pubKey[0] == 0x03) {
		return AlgorithmSecp256k1, nil
	}

	return KeyAlgorithm(-1), errors.New("unable to determine algorithm from public key")
}

// DetectAlgorithmFromSeed attempts to detect the algorithm from a seed.
// Note: This is not always deterministic as the same seed can be used with either algorithm.
// When in doubt, use the known algorithm or default to Ed25519 (recommended for new accounts).
func DetectAlgorithmFromSeed(seed string) (KeyAlgorithm, error) {
	// Try to decode the seed first
	if strings.HasPrefix(seed, "s") {
		_, err := decodeSeed(seed)
		if err != nil {
			return KeyAlgorithm(-1), fmt.Errorf("invalid seed format: %w", err)
		}
	} else {
		_, err := hex.DecodeString(seed)
		if err != nil {
			return KeyAlgorithm(-1), fmt.Errorf("invalid seed hex: %w", err)
		}
	}

	// Seeds don't encode algorithm information
	// Default to Ed25519 as recommended for new accounts
	return AlgorithmEd25519, nil
}

// SignTransactionWithAutoDetect signs a transaction, auto-detecting the algorithm from public key if present.
// If no public key is present, defaults to Ed25519.
func SignTransactionWithAutoDetect(txFields map[string]any, seed string) (*SignResult, error) {
	algorithm := AlgorithmEd25519 // Default

	// Try to detect from existing public key if present
	if pubKeyHex, ok := txFields["SigningPubKey"].(string); ok && pubKeyHex != "" {
		detected, err := DetectAlgorithm(pubKeyHex)
		if err == nil {
			algorithm = detected
		}
	}

	signer := NewSigner(algorithm)
	return signer.SignTransaction(txFields, seed)
}
