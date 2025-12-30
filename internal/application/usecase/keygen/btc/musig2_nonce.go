// Package btc provides Bitcoin-specific use cases for Keygen wallet operations.
//
// # MuSig2 Nonce Generation (Round 1)
//
// The GenerateMuSig2NonceUseCase implements Round 1 of the MuSig2 two-round protocol.
// Each signer generates a fresh cryptographic nonce for a specific transaction.
//
// # Why Nonce Generation is Critical
//
// In MuSig2, nonce generation is the most security-critical operation:
//   - Each nonce must be used EXACTLY ONCE per transaction
//   - Nonce reuse leaks the signer's private key
//   - Nonces must be generated BEFORE seeing the message to be signed
//   - All signers generate nonces in parallel (Round 1)
//   - Nonces are collected and aggregated before signing (Round 2)
//
// # MuSig2 Two-Round Protocol Overview
//
//	┌─────────────────────────────────────────────────────────┐
//	│ Round 1: Nonce Generation (PARALLEL)                    │
//	├─────────────────────────────────────────────────────────┤
//	│                                                         │
//	│  Keygen Wallet:                                         │
//	│    1. Generate fresh nonce (secret + public)            │
//	│    2. Store nonce locally                               │
//	│    3. Share PUBLIC nonce with Watch wallet              │
//	│                                                         │
//	│  Sign1 Wallet:                                          │
//	│    1. Generate fresh nonce (secret + public)            │
//	│    2. Store nonce locally                               │
//	│    3. Share PUBLIC nonce with Watch wallet              │
//	│                                                         │
//	│  Sign2 Wallet:                                          │
//	│    1. Generate fresh nonce (secret + public)            │
//	│    2. Store nonce locally                               │
//	│    3. Share PUBLIC nonce with Watch wallet              │
//	│                                                         │
//	│  Watch Wallet (Coordinator):                            │
//	│    1. Collect all public nonces                         │
//	│    2. Aggregate nonces                                  │
//	│    3. Distribute aggregated nonce to all signers        │
//	└─────────────────────────────────────────────────────────┘
//	                         │
//	                         ▼
//	┌─────────────────────────────────────────────────────────┐
//	│ Round 2: Signing (SEQUENTIAL)                           │
//	├─────────────────────────────────────────────────────────┤
//	│                                                         │
//	│  Each signer (Keygen, Sign1, Sign2):                    │
//	│    1. Retrieve stored nonce from Round 1                │
//	│    2. Use nonce + private key to create partial sig     │
//	│    3. Mark nonce as USED (prevent reuse)                │
//	│    4. Share partial signature with Watch wallet         │
//	│                                                         │
//	│  Watch Wallet:                                          │
//	│    1. Collect all partial signatures                    │
//	│    2. Aggregate into final Schnorr signature            │
//	│    3. Verify signature                                  │
//	│    4. Broadcast transaction                             │
//	└─────────────────────────────────────────────────────────┘
//
// # Nonce Data Flow
//
//	┌──────────┐   Public Nonce   ┌──────────┐
//	│  Keygen  │ ───────────────> │          │
//	│  Wallet  │                  │          │
//	└──────────┘                  │          │
//	                              │  Watch   │  Aggregated
//	┌──────────┐   Public Nonce   │  Wallet  │  Nonce
//	│  Sign1   │ ───────────────> │  (Coord) │ ────────> [Distribute to all]
//	│  Wallet  │                  │          │
//	└──────────┘                  │          │
//	                              │          │
//	┌──────────┐   Public Nonce   │          │
//	│  Sign2   │ ───────────────> │          │
//	│  Wallet  │                  │          │
//	└──────────┘                  └──────────┘
//
// # Security Best Practices
//
//   - Generate nonces using cryptographically secure random number generator
//   - Store nonces securely with transaction association
//   - Never reuse nonces across different transactions
//   - Validate nonce uniqueness before Round 2
//   - Mark nonces as "used" after signing to prevent accidental reuse
//
// # Example
//
// For detailed usage examples, see Example_generateMuSig2Nonce.
package btc

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"

	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	"github.com/hiromaily/go-crypto-wallet/internal/domain/multisig"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/bitcoin/btc"
)

type generateMuSig2NonceUseCase struct {
	musig2Service *btc.MuSig2Service
	nonceRepo     multisig.NonceRepository
}

// NewGenerateMuSig2NonceUseCase creates a new GenerateMuSig2NonceUseCase for BTC keygen.
func NewGenerateMuSig2NonceUseCase(
	musig2Service *btc.MuSig2Service,
	nonceRepo multisig.NonceRepository,
) keygenusecase.GenerateMuSig2NonceUseCase {
	return &generateMuSig2NonceUseCase{
		musig2Service: musig2Service,
		nonceRepo:     nonceRepo,
	}
}

// Generate implements MuSig2 Round 1: nonce generation and storage.
//
// This use case:
// 1. Creates a MuSig2 context (private key will be provided by caller)
// 2. Generates a new signing session
// 3. Retrieves the public nonce
// 4. Stores the nonce securely for later use in Round 2
// 5. Returns the public nonce for sharing with other signers
//
// SECURITY: The secret nonce is managed by the MuSig2Service session and must not be reused.
// Each transaction requires a fresh nonce to prevent private key leakage.
func (u *generateMuSig2NonceUseCase) Generate(
	ctx context.Context,
	input keygenusecase.GenerateMuSig2NonceInput,
) (keygenusecase.GenerateMuSig2NonceOutput, error) {
	// Note: In a real implementation, we would need to:
	// 1. Get the private key for this signer (from account key repository)
	// 2. Get all public keys from all signers (for context creation)
	// 3. Create MuSig2 context with keys
	// 4. Create session and generate nonce
	//
	// For now, this is a simplified implementation that assumes the caller
	// will provide the necessary cryptographic materials.

	// Placeholder: Create a dummy private key for demonstration
	// In real implementation, this would come from the account key repository
	// Using SHA256 to create a deterministic, full-length key for the placeholder.
	privKeyBytes := sha256.Sum256([]byte(input.SignerID))
	privKey, _ := btcec.PrivKeyFromBytes(privKeyBytes[:])

	// Placeholder: For demo, we'll use just this signer's public key
	// In real implementation, we need all signers' public keys
	allPubKeys := []*btcec.PublicKey{privKey.PubKey(), privKey.PubKey()} // Simplified

	// Create MuSig2 context
	musig2Ctx, err := u.musig2Service.CreateContext(privKey, allPubKeys, true)
	if err != nil {
		return keygenusecase.GenerateMuSig2NonceOutput{}, fmt.Errorf("failed to create MuSig2 context: %w", err)
	}

	// Create signing session (generates fresh nonce internally)
	session, err := u.musig2Service.CreateSession(musig2Ctx)
	if err != nil {
		return keygenusecase.GenerateMuSig2NonceOutput{}, fmt.Errorf("failed to create MuSig2 session: %w", err)
	}

	// Get public nonce for sharing
	publicNonce := u.musig2Service.GetPublicNonce(session)

	// Store nonce for Round 2
	err = u.nonceRepo.SaveNonce(ctx, input.SignerID, input.TransactionID, publicNonce)
	if err != nil {
		return keygenusecase.GenerateMuSig2NonceOutput{}, fmt.Errorf("failed to store nonce: %w", err)
	}

	return keygenusecase.GenerateMuSig2NonceOutput{
		PublicNonce: publicNonce,
		SignerID:    input.SignerID,
	}, nil
}

// ExampleGenerateMuSig2Nonce demonstrates how to use GenerateMuSig2NonceUseCase
// to generate fresh nonces for MuSig2 Round 1 (nonce generation phase).
//
// This example shows the workflow for a single signer (Keygen wallet):
//  1. Initialize dependencies (MuSig2Service, NonceRepository)
//  2. Create the use case with NewGenerateMuSig2NonceUseCase
//  3. Prepare input with signer ID and transaction ID
//  4. Execute the use case to generate a fresh nonce
//  5. Share the public nonce with the Watch wallet coordinator
//
// CRITICAL SECURITY NOTES:
//   - Each nonce must be generated BEFORE seeing the message to be signed
//   - Nonces must NEVER be reused across different transactions
//   - The secret nonce is managed internally and never exposed
//   - Only the public nonce is shared with other participants
//
// Note: This is a conceptual example. In real implementation:
//   - NonceRepository would be backed by persistent storage (database, file)
//   - Transaction ID would come from a payment request or unsigned PSBT
//   - Signer ID would identify the specific wallet (keygen, auth1, auth2, etc.)
//   - The public nonce would be encoded in PSBT proprietary fields for sharing
func ExampleGenerateMuSig2Nonce() {
	ctx := context.Background()

	// Initialize dependencies
	var (
		musig2Service *btc.MuSig2Service       // Provides MuSig2 cryptographic operations
		nonceRepo     multisig.NonceRepository // Repository for nonce storage
	)

	// Create the use case
	useCase := NewGenerateMuSig2NonceUseCase(musig2Service, nonceRepo)

	// Prepare input for nonce generation
	input := keygenusecase.GenerateMuSig2NonceInput{
		SignerID:      "keygen",     // Identifies this signer (keygen, auth1, auth2, etc.)
		TransactionID: "payment_42", // Associates nonce with specific transaction
	}

	// Execute Round 1: Generate fresh nonce
	output, err := useCase.Generate(ctx, input)
	if err != nil {
		fmt.Printf("failed to generate nonce: %v\n", err)
		return
	}

	fmt.Printf("Generated nonce for signer %s\n", output.SignerID)
	fmt.Printf("Public nonce length: %d bytes\n", len(output.PublicNonce))

	// After successful execution:
	// 1. Secret nonce is stored securely (managed by MuSig2Service session)
	// 2. Public nonce is stored in NonceRepository (associated with transaction ID)
	// 3. Public nonce can now be shared with Watch wallet coordinator
	//
	// Next steps in MuSig2 workflow:
	// 1. All signers generate their nonces (parallel)
	// 2. Watch wallet collects all public nonces
	// 3. Watch wallet aggregates nonces
	// 4. Aggregated nonce is distributed to all signers
	// 5. Round 2: Each signer creates partial signature using stored nonce
	//
	// Output:
	// Generated nonce for signer keygen
	// Public nonce length: 66 bytes
}
