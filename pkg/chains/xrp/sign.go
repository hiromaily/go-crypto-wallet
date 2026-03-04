// Package xrp provides XRP (Ripple) cryptocurrency utilities including
// key generation, address encoding, and cryptographic operations.
package xrp

import (
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"fmt"
	"maps"
	"strings"

	addresscodec "github.com/Peersyst/xrpl-go/address-codec"
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
//   - txFields: Transaction fields as a map (from JSON). This map is NOT modified.
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

	// Create a copy of txFields to avoid modifying the caller's map
	txFieldsCopy := make(map[string]any, len(txFields))
	maps.Copy(txFieldsCopy, txFields)

	// Derive the key pair from the seed
	keyPair, err := s.deriveKeyPair(seed)
	if err != nil {
		return nil, fmt.Errorf("failed to derive key pair: %w", err)
	}

	// Add the signing public key to the transaction copy
	txFieldsCopy["SigningPubKey"] = keyPair.PublicKeyHex

	// Create transaction for serialization
	tx := NewTransaction(txFieldsCopy)

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

	// Add signature to transaction copy
	txFieldsCopy["TxnSignature"] = signatureHex

	// Serialize the complete signed transaction
	signedTx := NewTransaction(txFieldsCopy)
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

// decodeSeed decodes an XRP seed from Base58Check format to raw seed bytes.
// Handles both secp256k1 seeds (prefix [0x21], "s..." in Base58) and
// ed25519 seeds (prefix [0x01, 0xe1, 0x4b], "sEd..." in Base58).
func decodeSeed(seed string) ([]byte, error) {
	decoded, err := addresscodec.Base58CheckDecode(seed)
	if err != nil {
		return nil, fmt.Errorf("invalid seed encoding: %w", err)
	}

	// Check ed25519 prefix [0x01, 0xe1, 0x4b]
	if len(decoded) >= 3 && decoded[0] == 0x01 && decoded[1] == 0xe1 && decoded[2] == 0x4b {
		return decoded[3:], nil
	}
	// Check secp256k1 prefix [0x21]
	if len(decoded) >= 1 && decoded[0] == byte(RippleFamilySeed) {
		return decoded[1:], nil
	}
	return nil, errors.New("invalid seed format: unrecognized prefix")
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

	// Use shared helper to derive private key (same as keygen)
	privateKeyBytes := deriveSecp256k1PrivateKey(seedBytes)

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

// DetectAlgorithmFromSeed detects the key algorithm from the seed encoding.
// Ed25519 seeds use prefix [0x01, 0xe1, 0x4b] ("sEd..." in Base58).
// Secp256k1 seeds use prefix [0x21] ("s..." in Base58).
// Hex seeds default to Ed25519 as recommended for new accounts.
func DetectAlgorithmFromSeed(seed string) (KeyAlgorithm, error) {
	if strings.HasPrefix(seed, "s") {
		decoded, err := addresscodec.Base58CheckDecode(seed)
		if err != nil {
			return KeyAlgorithm(-1), fmt.Errorf("invalid seed format: %w", err)
		}
		if len(decoded) >= 3 && decoded[0] == 0x01 && decoded[1] == 0xe1 && decoded[2] == 0x4b {
			return AlgorithmEd25519, nil
		}
		return AlgorithmSecp256k1, nil
	}
	_, err := hex.DecodeString(seed)
	if err != nil {
		return KeyAlgorithm(-1), fmt.Errorf("invalid seed hex: %w", err)
	}
	// Default to Ed25519 for hex seeds (recommended for new accounts)
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
