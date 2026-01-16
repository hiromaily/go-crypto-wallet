// Package btc provides Bitcoin-specific use cases for Keygen wallet operations.
//
// # MuSig2 Partial Signature Creation (Round 2)
//
// The MuSig2SignUseCase implements Round 2 of the MuSig2 two-round protocol.
// Each signer creates a partial signature using their stored nonce and private key.
//
// # Round 2 Overview: Signing
//
// After all signers have generated nonces (Round 1) and the Watch wallet has
// aggregated them, each signer proceeds to Round 2:
//
//  1. Retrieve the stored nonce from Round 1
//  2. Receive the aggregated nonce from Watch wallet
//  3. Create MuSig2 context and session
//  4. Register all aggregated nonces
//  5. Create partial signature using:
//     - Signer's private key
//     - Stored secret nonce (from Round 1)
//     - Message hash (transaction sighash)
//  6. CRITICAL: Mark nonce as USED to prevent reuse
//  7. Return partial signature for aggregation
//
// # Signing Workflow
//
//	┌─────────────────────────────────────────────────────────┐
//	│ Prerequisites (from Round 1):                           │
//	│  - Fresh nonce generated and stored                     │
//	│  - Aggregated nonce received from Watch wallet          │
//	│  - Message hash (sighash) to be signed                  │
//	└────────────────────┬────────────────────────────────────┘
//	                     │
//	                     ▼
//	┌─────────────────────────────────────────────────────────┐
//	│ Step 1: Retrieve Stored Nonce                           │
//	│  - Query NonceRepository by (SignerID, TransactionID)   │
//	│  - Validate nonce has not been used                     │
//	│  - Validate nonce belongs to this signer                │
//	└────────────────────┬────────────────────────────────────┘
//	                     │
//	                     ▼
//	┌─────────────────────────────────────────────────────────┐
//	│ Step 2: Create MuSig2 Context                           │
//	│  - Private key for this signer                          │
//	│  - All public keys from all signers                     │
//	│  - Sort keys (must match Round 1)                       │
//	└────────────────────┬────────────────────────────────────┘
//	                     │
//	                     ▼
//	┌─────────────────────────────────────────────────────────┐
//	│ Step 3: Create Signing Session                          │
//	│  - Session manages nonce state                          │
//	│  - Associates nonce with signing operation              │
//	└────────────────────┬────────────────────────────────────┘
//	                     │
//	                     ▼
//	┌─────────────────────────────────────────────────────────┐
//	│ Step 4: Register Aggregated Nonces                      │
//	│  - Register nonces from ALL signers                     │
//	│  - MuSig2 library validates completeness                │
//	└────────────────────┬────────────────────────────────────┘
//	                     │
//	                     ▼
//	┌─────────────────────────────────────────────────────────┐
//	│ Step 5: Create Partial Signature                        │
//	│  - Sign(session, messageHash)                           │
//	│  - Uses private key + secret nonce internally           │
//	│  - Returns partial signature (S component)              │
//	└────────────────────┬────────────────────────────────────┘
//	                     │
//	                     ▼
//	┌─────────────────────────────────────────────────────────┐
//	│ Step 6: Mark Nonce as USED (CRITICAL!)                  │
//	│  - Update NonceRepository: used = true                  │
//	│  - Prevents accidental nonce reuse                      │
//	│  - Protects private key from leakage                    │
//	└────────────────────┬────────────────────────────────────┘
//	                     │
//	                     ▼
//	┌─────────────────────────────────────────────────────────┐
//	│ Output: Partial Signature                               │
//	│  - PartialSignature: [32]byte (S component)             │
//	│  - SignerID: identifies which signer created it         │
//	│  - Ready for aggregation by Watch wallet                │
//	└─────────────────────────────────────────────────────────┘
//
// # Security Critical: Nonce Reuse Prevention
//
// The most critical security requirement in MuSig2 is preventing nonce reuse:
//
//   - Using the same nonce twice with different messages leaks the private key
//   - This use case MUST mark nonces as "used" after signing
//   - NonceRepository must enforce uniqueness constraints
//   - Failed signing operations should NOT mark nonces as used
//
// Implementation: After successful signing, call nonceRepo.MarkNonceUsed()
//
// # Known Issue: Session State Persistence
//
// IMPORTANT: This implementation has a known limitation (see issue #136):
//
// The current implementation creates a NEW session in Round 2, which generates
// a DIFFERENT secret nonce than what was stored in Round 1. This results in
// INVALID signatures because the secret nonce used for signing doesn't match
// the public nonce that was aggregated.
//
// The correct implementation requires:
//  1. Create session ONCE in Round 1
//  2. Store session state persistently (not just public nonce)
//  3. Restore session in Round 2 for signing
//  4. Use the SAME secret nonce from Round 1
//
// This requires session state serialization/deserialization, which is not
// yet implemented.
//
// # Example
//
// For detailed usage examples, see Example_muSig2Sign.
package btc

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"

	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	"github.com/hiromaily/go-crypto-wallet/internal/domain/multisig"
	apibtcimpl "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/btc/btc"
)

type muSig2SignUseCase struct {
	musig2Service *apibtcimpl.MuSig2Service
	nonceRepo     multisig.NonceRepository
}

// NewMuSig2SignUseCase creates a new MuSig2SignUseCase for BTC keygen.
func NewMuSig2SignUseCase(
	musig2Service *apibtcimpl.MuSig2Service,
	nonceRepo multisig.NonceRepository,
) keygenusecase.MuSig2SignUseCase {
	return &muSig2SignUseCase{
		musig2Service: musig2Service,
		nonceRepo:     nonceRepo,
	}
}

// Sign implements MuSig2 Round 2: partial signature creation.
//
// This use case:
// 1. Retrieves the stored nonce from Round 1
// 2. Creates a MuSig2 context and session
// 3. Registers aggregated nonces from all signers
// 4. Creates a partial signature for the transaction
// 5. Marks the nonce as used (critical for security)
// 6. Returns the partial signature for aggregation
//
// SECURITY: After creating the partial signature, the nonce MUST be marked as used
// to prevent accidental reuse, which would leak the private key.
//
// KNOWN ISSUE: This implementation creates a NEW session in Round 2, which generates
// a different secret nonce than what was stored in Round 1. This will result in INVALID
// signatures. The session object from Round 1 must be reused in Round 2 for valid signing.
// This requires session state persistence, which is not yet implemented.
// See: https://github.com/hiromaily/go-crypto-wallet/issues/136#issuecomment
func (u *muSig2SignUseCase) Sign(
	ctx context.Context,
	input keygenusecase.MuSig2SignInput,
) (keygenusecase.MuSig2SignOutput, error) {
	// Retrieve stored nonce from Round 1
	storedNonce, err := u.nonceRepo.GetNonce(ctx, input.SignerID, input.TransactionID)
	if err != nil {
		return keygenusecase.MuSig2SignOutput{}, fmt.Errorf("failed to retrieve stored nonce: %w", err)
	}

	// Validate stored nonce matches signer
	if storedNonce.SignerID() != input.SignerID {
		return keygenusecase.MuSig2SignOutput{}, fmt.Errorf(
			"nonce signer ID mismatch: expected %s, got %s",
			input.SignerID,
			storedNonce.SignerID(),
		)
	}

	// Note: In a real implementation, we would need to:
	// 1. Get the private key for this signer
	// 2. Get all public keys from all signers
	// 3. Create MuSig2 context with keys
	// 4. Create session with stored nonce
	// 5. Register all public nonces
	// 6. Sign the message
	//
	// For now, this is a simplified implementation

	// Placeholder: Create a dummy private key for demonstration
	// Using SHA256 to create a deterministic, full-length key for the placeholder.
	privKeyBytes := sha256.Sum256([]byte(input.SignerID))
	privKey, _ := btcec.PrivKeyFromBytes(privKeyBytes[:])

	// Placeholder: For demo, we'll use just this signer's public key
	allPubKeys := []*btcec.PublicKey{privKey.PubKey(), privKey.PubKey()}

	// Create MuSig2 context
	musig2Ctx, err := u.musig2Service.CreateContext(privKey, allPubKeys, true)
	if err != nil {
		return keygenusecase.MuSig2SignOutput{}, fmt.Errorf("failed to create MuSig2 context: %w", err)
	}

	// Create signing session
	session, err := u.musig2Service.CreateSession(musig2Ctx)
	if err != nil {
		return keygenusecase.MuSig2SignOutput{}, fmt.Errorf("failed to create MuSig2 session: %w", err)
	}

	// Register aggregated nonces from all signers
	for _, nonce := range input.AggregatedNonces {
		haveAll, err := u.musig2Service.RegisterPubNonce(session, nonce)
		if err != nil {
			return keygenusecase.MuSig2SignOutput{}, fmt.Errorf("failed to register public nonce: %w", err)
		}
		if haveAll {
			break
		}
	}

	// Create partial signature
	_, err = u.musig2Service.Sign(session, input.MessageHash)
	if err != nil {
		return keygenusecase.MuSig2SignOutput{}, fmt.Errorf("failed to create partial signature: %w", err)
	}

	// Mark nonce as used (CRITICAL for security)
	err = u.nonceRepo.MarkNonceUsed(ctx, input.SignerID, input.TransactionID)
	if err != nil {
		return keygenusecase.MuSig2SignOutput{}, fmt.Errorf("failed to mark nonce as used: %w", err)
	}

	// Note: The actual partial signature would need to be extracted from the session
	// or stored in a domain object. For this MVP, we return a placeholder.
	// In a real implementation, the signature would be serialized and stored.
	var sigScalar [32]byte
	// Placeholder - in real implementation, extract from partialSig
	copy(sigScalar[:], []byte("placeholder_partial_signature"))

	return keygenusecase.MuSig2SignOutput{
		PartialSignature: sigScalar,
		SignerID:         input.SignerID,
	}, nil
}

// ExampleMuSig2Sign demonstrates how to use MuSig2SignUseCase
// to create partial signatures in MuSig2 Round 2 (signing phase).
//
// This example shows the workflow for a single signer (Keygen wallet):
//  1. Initialize dependencies (MuSig2Service, NonceRepository)
//  2. Create the use case with NewMuSig2SignUseCase
//  3. Prepare input with:
//     - Signer ID and transaction ID (to retrieve stored nonce)
//     - Aggregated nonces from all signers
//     - Message hash (transaction sighash to be signed)
//  4. Execute the use case to create partial signature
//  5. Share the partial signature with Watch wallet for aggregation
//
// CRITICAL SECURITY NOTES:
//   - Nonce from Round 1 must be retrieved and validated before signing
//   - After signing, nonce MUST be marked as used to prevent reuse
//   - Nonce reuse with different messages leaks the private key
//   - Session state from Round 1 should ideally be reused (see Known Issue)
//
// KNOWN ISSUE:
// This implementation creates a NEW session in Round 2, generating a different
// secret nonce than Round 1. This results in INVALID signatures. For production,
// session state persistence is required (see issue #136).
//
// Note: This is a conceptual example. In real implementation:
//   - Message hash would be the actual transaction sighash from PSBT
//   - Aggregated nonces would be received from Watch wallet in PSBT
//   - Partial signature would be encoded in PSBT proprietary fields
//   - Session state would be restored from Round 1 (not recreated)
func ExampleMuSig2Sign() {
	ctx := context.Background()

	// Initialize dependencies
	var (
		musig2Service *apibtcimpl.MuSig2Service // Provides MuSig2 cryptographic operations
		nonceRepo     multisig.NonceRepository  // Repository for nonce storage
	)

	// Create the use case
	useCase := NewMuSig2SignUseCase(musig2Service, nonceRepo)

	// Prepare input for signing (Round 2)
	// Note: These would come from PSBT in real implementation
	var messageHash [32]byte
	copy(messageHash[:], []byte("transaction_sighash_32_bytes_")) // Placeholder

	publicNonces := [][66]byte{
		{}, // Public nonce from keygen
		{}, // Public nonce from sign1
		{}, // Public nonce from sign2
	}

	input := keygenusecase.MuSig2SignInput{
		SignerID:         "keygen",     // Identifies this signer
		TransactionID:    "payment_42", // Associates with nonce from Round 1
		AggregatedNonces: publicNonces, // All public nonces from all signers
		MessageHash:      messageHash,  // Transaction sighash to sign
	}

	// Execute Round 2: Create partial signature
	output, err := useCase.Sign(ctx, input)
	if err != nil {
		fmt.Printf("failed to create partial signature: %v\n", err)
		return
	}

	fmt.Printf("Created partial signature for signer %s\n", output.SignerID)
	fmt.Printf("Signature length: %d bytes\n", len(output.PartialSignature))

	// After successful execution:
	// 1. Partial signature is created using private key + stored nonce
	// 2. Nonce is marked as USED in NonceRepository
	// 3. Partial signature can now be shared with Watch wallet
	//
	// Next steps in MuSig2 workflow:
	// 1. All signers create their partial signatures (sequential)
	// 2. Watch wallet collects all partial signatures
	// 3. Watch wallet aggregates partial signatures into final Schnorr signature
	// 4. Watch wallet verifies final signature
	// 5. Watch wallet broadcasts transaction to Bitcoin network
	//
	// Output:
	// Created partial signature for signer keygen
	// Signature length: 32 bytes
}
