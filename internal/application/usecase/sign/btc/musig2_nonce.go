package btc

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"

	repocold "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository/cold"
	signusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/sign"
	"github.com/hiromaily/go-crypto-wallet/internal/domain/multisig"
	apibtcimpl "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/btc/btc"
)

type generateMuSig2NonceUseCase struct {
	musig2Service *apibtcimpl.MuSig2Service
	nonceRepo     multisig.NonceRepository
	authKeyRepo   repocold.AuthAccountKeyRepositorier
}

// NewGenerateMuSig2NonceUseCase creates a new GenerateMuSig2NonceUseCase for BTC sign wallet.
func NewGenerateMuSig2NonceUseCase(
	musig2Service *apibtcimpl.MuSig2Service,
	nonceRepo multisig.NonceRepository,
	authKeyRepo repocold.AuthAccountKeyRepositorier,
) signusecase.GenerateMuSig2NonceUseCase {
	return &generateMuSig2NonceUseCase{
		musig2Service: musig2Service,
		nonceRepo:     nonceRepo,
		authKeyRepo:   authKeyRepo,
	}
}

// Generate implements MuSig2 Round 1: nonce generation and storage for Sign wallet.
//
// This use case:
// 1. Retrieves the auth key from auth_account_key table (Sign wallet's key)
// 2. Creates a MuSig2 context
// 3. Generates a new signing session
// 4. Retrieves the public nonce
// 5. Stores the nonce securely for later use in Round 2
// 6. Returns the public nonce for sharing with other signers
//
// SECURITY: The secret nonce is managed by the MuSig2Service session and must not be reused.
// Each transaction requires a fresh nonce to prevent private key leakage.
//
// Sign Wallet Specifics:
//   - Uses auth key from auth_account_key table (not account_key)
//   - Works offline (no Bitcoin Core RPC required)
//   - Can run in parallel with Keygen wallet's nonce generation
func (u *generateMuSig2NonceUseCase) Generate(
	ctx context.Context,
	input signusecase.GenerateMuSig2NonceInput,
) (signusecase.GenerateMuSig2NonceOutput, error) {
	// Get auth key from auth_account_key table (Sign wallet's key)
	authKey, err := u.authKeyRepo.GetOne(input.AuthType)
	if err != nil {
		return signusecase.GenerateMuSig2NonceOutput{}, fmt.Errorf(
			"failed to get auth key for authType %s: %w",
			input.AuthType,
			err,
		)
	}

	// Note: In a real implementation, we would need to:
	// 1. Parse the private key from WIF format
	// 2. Get all public keys from all signers (for context creation)
	// 3. Create MuSig2 context with keys
	// 4. Create session and generate nonce
	//
	// For now, this is a simplified implementation that assumes the caller
	// will provide the necessary cryptographic materials.

	// Placeholder: Create a dummy private key for demonstration
	// In real implementation, this would come from authKey.WalletImportFormat
	// Using SHA256 to create a deterministic, full-length key for the placeholder.
	privKeyBytes := sha256.Sum256([]byte(authKey.WalletImportFormat))
	privKey, pubKey := btcec.PrivKeyFromBytes(privKeyBytes[:])

	// Placeholder: For demo, we'll use just this signer's public key
	// In real implementation, we need all signers' public keys
	allPubKeys := []*btcec.PublicKey{pubKey, pubKey} // Simplified

	// Create MuSig2 context
	musig2Ctx, err := u.musig2Service.CreateContext(privKey, allPubKeys, true)
	if err != nil {
		return signusecase.GenerateMuSig2NonceOutput{}, fmt.Errorf("failed to create MuSig2 context: %w", err)
	}

	// Create signing session (generates fresh nonce internally)
	session, err := u.musig2Service.CreateSession(musig2Ctx)
	if err != nil {
		return signusecase.GenerateMuSig2NonceOutput{}, fmt.Errorf("failed to create MuSig2 session: %w", err)
	}

	// Get public nonce for sharing
	publicNonce := u.musig2Service.GetPublicNonce(session)

	// Use auth type as signer ID for Sign wallet
	signerID := string(input.AuthType)

	// Store nonce for Round 2
	err = u.nonceRepo.SaveNonce(ctx, signerID, input.TransactionID, publicNonce)
	if err != nil {
		return signusecase.GenerateMuSig2NonceOutput{}, fmt.Errorf("failed to store nonce: %w", err)
	}

	return signusecase.GenerateMuSig2NonceOutput{
		PublicNonce: publicNonce,
		SignerID:    signerID,
	}, nil
}
