package multisig

import (
	"bytes"
	"errors"
	"fmt"
)

// ValidateSignerCount validates the number of signers for MuSig2.
// MuSig2 requires at least 2 signers.
func ValidateSignerCount(count int) error {
	if count < 2 {
		return fmt.Errorf("MuSig2 requires at least 2 signers: got %d", count)
	}

	// Practical limit check (same as traditional multisig)
	if count > 15 {
		return fmt.Errorf("signer count cannot exceed 15: got %d", count)
	}

	return nil
}

// ValidateNonceUniqueness validates that all nonce commitments are unique.
// Nonce reuse is a critical security vulnerability that can leak private keys.
func ValidateNonceUniqueness(nonces []*NonceCommitment) error {
	if len(nonces) == 0 {
		return errors.New("no nonces provided")
	}

	seen := make(map[[66]byte]bool, len(nonces))
	seenSigners := make(map[string]bool, len(nonces))

	for i, nonce := range nonces {
		if nonce == nil {
			return fmt.Errorf("nonce at index %d is nil", i)
		}

		// Check for duplicate nonces (critical security issue)
		publicNonce := nonce.PublicNonce()
		if seen[publicNonce] {
			return fmt.Errorf("duplicate nonce detected at index %d: nonce reuse is forbidden", i)
		}
		seen[publicNonce] = true

		// Check for duplicate signer IDs
		signerID := nonce.SignerID()
		if signerID == "" {
			return fmt.Errorf("nonce at index %d has empty signer ID", i)
		}
		if seenSigners[signerID] {
			return fmt.Errorf("duplicate signer ID %s detected at index %d", signerID, i)
		}
		seenSigners[signerID] = true

		// Validate nonce is not zero
		var zero [66]byte
		if bytes.Equal(publicNonce[:], zero[:]) {
			return fmt.Errorf("nonce at index %d is zero", i)
		}
	}

	return nil
}

// ValidatePartialSignatures validates partial signatures.
func ValidatePartialSignatures(signatures []*PartialSignature) error {
	if len(signatures) == 0 {
		return errors.New("no partial signatures provided")
	}

	seenSigners := make(map[string]bool, len(signatures))

	for i, sig := range signatures {
		if sig == nil {
			return fmt.Errorf("partial signature at index %d is nil", i)
		}

		// Check for duplicate signer IDs
		signerID := sig.SignerID()
		if signerID == "" {
			return fmt.Errorf("partial signature at index %d has empty signer ID", i)
		}
		if seenSigners[signerID] {
			return fmt.Errorf("duplicate signer ID %s detected at index %d", signerID, i)
		}
		seenSigners[signerID] = true

		// Validate signature scalar is not zero
		var zero [32]byte
		scalar := sig.SignatureScalar()
		if bytes.Equal(scalar[:], zero[:]) {
			return fmt.Errorf("signature scalar at index %d is zero", i)
		}
	}

	return nil
}

// ValidatePublicKeysForMuSig2 validates public keys for MuSig2 key aggregation.
func ValidatePublicKeysForMuSig2(publicKeys [][]byte) error {
	if len(publicKeys) < 2 {
		return fmt.Errorf("at least 2 public keys required for MuSig2: got %d", len(publicKeys))
	}

	if len(publicKeys) > 15 {
		return fmt.Errorf("public key count cannot exceed 15: got %d", len(publicKeys))
	}

	seen := make(map[string]bool, len(publicKeys))

	for i, pk := range publicKeys {
		if len(pk) == 0 {
			return fmt.Errorf("public key at index %d is empty", i)
		}

		// Public keys should be 33 bytes (compressed) or 65 bytes (uncompressed)
		if len(pk) != 33 && len(pk) != 65 {
			return fmt.Errorf("public key at index %d has invalid length: got %d bytes, expected 33 or 65", i, len(pk))
		}

		// Check for duplicate public keys
		pkStr := string(pk)
		if seen[pkStr] {
			return fmt.Errorf("duplicate public key detected at index %d", i)
		}
		seen[pkStr] = true

		// Validate that public key is not all zeros
		allZero := true
		for _, b := range pk {
			if b != 0 {
				allZero = false
				break
			}
		}
		if allZero {
			return fmt.Errorf("public key at index %d is all zeros", i)
		}
	}

	return nil
}

// ValidateAggregatedPublicKey validates an aggregated public key.
func ValidateAggregatedPublicKey(publicKey []byte) error {
	if len(publicKey) == 0 {
		return errors.New("aggregated public key is empty")
	}

	// Aggregated keys should be compressed (33 bytes)
	if len(publicKey) != 33 {
		return fmt.Errorf("aggregated public key has invalid length: got %d bytes, expected 33", len(publicKey))
	}

	// Validate that public key is not all zeros
	allZero := true
	for _, b := range publicKey {
		if b != 0 {
			allZero = false
			break
		}
	}
	if allZero {
		return errors.New("aggregated public key is all zeros")
	}

	return nil
}

// ValidateSigningSessionComplete validates that a signing session has all required data.
func ValidateSigningSessionComplete(session *SigningSession) error {
	if session == nil {
		return errors.New("signing session is nil")
	}

	if !session.HasAllNonces() {
		return fmt.Errorf("signing session %s is missing nonce commitments: have %d, need %d",
			session.SessionID(), len(session.nonceCommitments), session.ParticipantCount())
	}

	if !session.HasAllSignatures() {
		return fmt.Errorf("signing session %s is missing partial signatures: have %d, need %d",
			session.SessionID(), len(session.partialSignatures), session.ParticipantCount())
	}

	return nil
}
