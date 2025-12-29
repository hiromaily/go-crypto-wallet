package btc

import (
	"context"
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
	privKeyBytes := make([]byte, 32)
	copy(privKeyBytes, []byte(input.SignerID)) // Simplified for demo
	privKey, _ := btcec.PrivKeyFromBytes(privKeyBytes)

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
