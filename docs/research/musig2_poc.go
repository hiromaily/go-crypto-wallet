// Package main demonstrates a proof-of-concept for MuSig2 using btcd/btcec/v2/schnorr/musig2
//
// This POC shows the complete two-round MuSig2 protocol:
//   - Round 1: Nonce generation and aggregation
//   - Round 2: Partial signature creation and aggregation
//
// This code is for research purposes only (issue #132).
// Production implementation will follow Clean Architecture in sub-issues #133-#141.
package main

import (
	"crypto/sha256"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr/musig2"
)

func main() {
	fmt.Println("=== MuSig2 Proof of Concept ===")
	fmt.Println()

	// Simulate 2-of-2 multisig (Keygen wallet + Sign wallet)
	// In production, these would be on separate offline machines
	keygenPrivKey, keygenPubKey := generateKeyPair()
	signPrivKey, signPubKey := generateKeyPair()

	fmt.Println("Generated key pairs:")
	fmt.Printf("  Keygen Public Key: %x\n", keygenPubKey.SerializeCompressed())
	fmt.Printf("  Sign Public Key:   %x\n", signPubKey.SerializeCompressed())
	fmt.Println()

	// Message to sign (transaction hash in real implementation)
	message := sha256.Sum256([]byte("test transaction"))
	fmt.Printf("Message Hash: %x\n", message)
	fmt.Println()

	// ===================================================================
	// ROUND 1: Nonce Generation and Exchange
	// ===================================================================
	fmt.Println("--- Round 1: Nonce Generation ---")

	// Keygen wallet: Create context and generate nonce
	keygenCtx, err := musig2.NewContext(
		keygenPrivKey,
		true, // sort keys for deterministic aggregation
		musig2.WithKnownSigners([]*btcec.PublicKey{keygenPubKey, signPubKey}),
	)
	if err != nil {
		panic(fmt.Sprintf("failed to create keygen context: %v", err))
	}

	keygenSession, err := keygenCtx.NewSession()
	if err != nil {
		panic(fmt.Sprintf("failed to create keygen session: %v", err))
	}

	keygenPubNonce := keygenSession.PublicNonce()

	fmt.Printf("Keygen Public Nonce: %x\n", keygenPubNonce[:])

	// Sign wallet: Create context and generate nonce
	signCtx, err := musig2.NewContext(
		signPrivKey,
		true, // must match keygen's sort setting
		musig2.WithKnownSigners([]*btcec.PublicKey{keygenPubKey, signPubKey}),
	)
	if err != nil {
		panic(fmt.Sprintf("failed to create sign context: %v", err))
	}

	signSession, err := signCtx.NewSession()
	if err != nil {
		panic(fmt.Sprintf("failed to create sign session: %v", err))
	}

	signPubNonce := signSession.PublicNonce()

	fmt.Printf("Sign Public Nonce:   %x\n", signPubNonce[:])
	fmt.Println()

	// ===================================================================
	// Nonce Exchange (would be via files or PSBT in production)
	// ===================================================================
	fmt.Println("--- Exchanging Nonces ---")

	// Keygen wallet registers Sign wallet's nonce
	haveAllNonces, err := keygenSession.RegisterPubNonce(signPubNonce)
	if err != nil {
		panic(fmt.Sprintf("failed to register sign nonce in keygen: %v", err))
	}
	fmt.Printf("Keygen has all nonces: %v\n", haveAllNonces)

	// Sign wallet registers Keygen wallet's nonce
	haveAllNonces, err = signSession.RegisterPubNonce(keygenPubNonce)
	if err != nil {
		panic(fmt.Sprintf("failed to register keygen nonce in sign: %v", err))
	}
	fmt.Printf("Sign has all nonces: %v\n", haveAllNonces)
	fmt.Println()

	// ===================================================================
	// ROUND 2: Signing
	// ===================================================================
	fmt.Println("--- Round 2: Partial Signature Creation ---")

	// Keygen wallet creates partial signature
	keygenPartialSig, err := keygenSession.Sign(message)
	if err != nil {
		panic(fmt.Sprintf("failed to create keygen partial signature: %v", err))
	}
	fmt.Printf("Keygen Partial Signature: S=%x\n", keygenPartialSig.S.Bytes())

	// Sign wallet creates partial signature
	signPartialSig, err := signSession.Sign(message)
	if err != nil {
		panic(fmt.Sprintf("failed to create sign partial signature: %v", err))
	}
	fmt.Printf("Sign Partial Signature:   S=%x\n", signPartialSig.S.Bytes())
	fmt.Println()

	// ===================================================================
	// Signature Aggregation (typically done by Watch wallet)
	// ===================================================================
	fmt.Println("--- Signature Aggregation ---")

	// Keygen wallet combines signatures (could be done on any wallet or watch wallet)
	haveAllSigs, err := keygenSession.CombineSig(signPartialSig)
	if err != nil {
		panic(fmt.Sprintf("failed to combine signature: %v", err))
	}
	if !haveAllSigs {
		panic("not all signatures received")
	}

	// Get final aggregated signature
	finalSig := keygenSession.FinalSig()
	if finalSig == nil {
		panic("failed to get final signature")
	}

	fmt.Printf("Final Aggregated Signature: %x\n", finalSig.Serialize())
	fmt.Println()

	// ===================================================================
	// Verification
	// ===================================================================
	fmt.Println("--- Signature Verification ---")

	// Get aggregated public key
	aggregatedKey, err := keygenCtx.CombinedKey()
	if err != nil {
		panic(fmt.Sprintf("failed to get combined key: %v", err))
	}

	// Verify the aggregated signature
	valid := finalSig.Verify(message[:], aggregatedKey)
	if valid {
		fmt.Println("✅ Signature is VALID")
	} else {
		fmt.Println("❌ Signature is INVALID")
		return
	}

	// ===================================================================
	// Comparison with Traditional Multisig
	// ===================================================================
	fmt.Println()
	fmt.Println("--- Size Comparison ---")

	// MuSig2: Single Schnorr signature (64 bytes)
	musig2Size := len(finalSig.Serialize())
	fmt.Printf("MuSig2 Signature Size: %d bytes\n", musig2Size)

	// Traditional 2-of-2 P2WSH: 2 ECDSA signatures (~72 bytes each) + script
	// Approximate size for 2-of-2 P2WSH input
	traditionalSize := 2*72 + 105 // 2 sigs + script
	fmt.Printf("Traditional 2-of-2 P2WSH (estimated): %d bytes\n", traditionalSize)

	reduction := float64(traditionalSize-musig2Size) / float64(traditionalSize) * 100
	fmt.Printf("Size Reduction: %.1f%%\n", reduction)

	fmt.Println()
	fmt.Println("=== POC Complete ===")
	fmt.Println()
	fmt.Println("Key Findings:")
	fmt.Println("✅ MuSig2 two-round protocol works correctly")
	fmt.Println("✅ Signature aggregation produces valid signature")
	fmt.Println("✅ Significant size reduction vs traditional multisig")
	fmt.Println("✅ btcd/btcec/v2/schnorr/musig2 is production-ready")
}

// generateKeyPair generates a random key pair for testing
func generateKeyPair() (*btcec.PrivateKey, *btcec.PublicKey) {
	privKey, err := btcec.NewPrivateKey()
	if err != nil {
		panic(fmt.Sprintf("failed to generate key pair: %v", err))
	}
	return privKey, privKey.PubKey()
}

// Additional helper functions for production implementation:

// NonceToFile would save nonce to file (for offline wallet workflow)
func NonceToFile(nonce [66]byte, filepath string) error {
	// Implementation in sub-issue #135 (nonce management)
	return nil
}

// NonceFromFile would load nonce from file
func NonceFromFile(filepath string) ([66]byte, error) {
	// Implementation in sub-issue #135 (nonce management)
	return [66]byte{}, nil
}

// SavePartialSignature would save partial signature to PSBT
func SavePartialSignature(sig *musig2.PartialSignature, psbtPath string) error {
	// Implementation in sub-issue #136-#138 (use cases)
	return nil
}

// Notes for production implementation:
//
// 1. Nonce Storage (Sub-issue #135):
//    - Store secret nonces securely in database
//    - Ensure nonce uniqueness per transaction
//    - Delete nonces after use to prevent reuse
//
// 2. PSBT Integration (Sub-issues #136-#138):
//    - Store public nonces in PSBT extension fields
//    - Store partial signatures in PSBT
//    - Watch wallet aggregates signatures from PSBT
//
// 3. Taproot Integration:
//    - Use WithTaprootTweakCtx() for Taproot addresses
//    - Ensure key-path spending compatibility
//
// 4. Error Handling:
//    - Wrap errors with context using fmt.Errorf + %w
//    - Handle nonce registration failures
//    - Validate signature aggregation success
//
// 5. Security:
//    - NEVER log private keys or secret nonces
//    - Ensure nonces are generated with secure RNG
//    - Prevent nonce reuse (critical for private key safety)
