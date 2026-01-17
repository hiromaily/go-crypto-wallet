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

type muSig2SignUseCase struct {
	musig2Service *apibtcimpl.MuSig2Service
	nonceRepo     multisig.NonceRepository
	authKeyRepo   repocold.AuthAccountKeyRepositorier
}

// NewMuSig2SignUseCase creates a new MuSig2SignUseCase for BTC sign wallet.
func NewMuSig2SignUseCase(
	musig2Service *apibtcimpl.MuSig2Service,
	nonceRepo multisig.NonceRepository,
	authKeyRepo repocold.AuthAccountKeyRepositorier,
) signusecase.MuSig2SignUseCase {
	return &muSig2SignUseCase{
		musig2Service: musig2Service,
		nonceRepo:     nonceRepo,
		authKeyRepo:   authKeyRepo,
	}
}

// Sign implements MuSig2 Round 2: partial signature creation for Sign wallet.
//
// This use case:
// 1. Retrieves the auth key from auth_account_key table (Sign wallet's key)
// 2. Retrieves the stored nonce from Round 1
// 3. Creates a MuSig2 context and session
// 4. Registers aggregated nonces from all signers
// 5. Creates a partial signature for the transaction
// 6. Marks the nonce as used (critical for security)
// 7. Returns the partial signature for aggregation
//
// SECURITY: After creating the partial signature, the nonce MUST be marked as used
// to prevent accidental reuse, which would leak the private key.
//
// CRITICAL ISSUE - Session Reuse Required:
// This implementation creates a NEW session in Round 2, which generates a different secret
// nonce than what was stored in Round 1. This is a FUNDAMENTAL CORRECTNESS BUG that will
// result in INVALID signatures when aggregated.
//
// Why this is critical:
//   - MuSig2 protocol requires using the SAME session object from Round 1 for signing
//   - Each session generates a unique secret nonce internally
//   - The public nonce stored in Round 1 corresponds to Round 1's secret nonce
//   - Using a different session in Round 2 means a different secret nonce
//   - This nonce mismatch causes signature verification to fail
//
// Required fix:
//   - Implement session state persistence (in-memory or serialized storage)
//   - Store the entire session object after Round 1 nonce generation
//   - Retrieve and reuse the same session object in Round 2 for signing
//   - This is a prerequisite for valid MuSig2 signatures
//
// See: https://github.com/hiromaily/go-crypto-wallet/issues/137#issuecomment
//
// Sign Wallet Specifics:
//   - Uses auth key from auth_account_key table (not account_key)
//   - Works offline (no Bitcoin Core RPC required)
//   - Runs after Keygen wallet completes Round 2 (adds second signature)
func (u *muSig2SignUseCase) Sign(
	ctx context.Context,
	input signusecase.MuSig2SignInput,
) (signusecase.MuSig2SignOutput, error) {
	// Get auth key from auth_account_key table (Sign wallet's key)
	authKey, err := u.authKeyRepo.GetOne(ctx, input.AuthType)
	if err != nil {
		return signusecase.MuSig2SignOutput{}, fmt.Errorf(
			"failed to get auth key for authType %s: %w",
			input.AuthType,
			err,
		)
	}

	// Use auth type as signer ID for Sign wallet
	signerID := string(input.AuthType)

	// Retrieve stored nonce from Round 1
	storedNonce, err := u.nonceRepo.GetNonce(ctx, signerID, input.TransactionID)
	if err != nil {
		return signusecase.MuSig2SignOutput{}, fmt.Errorf("failed to retrieve stored nonce: %w", err)
	}

	// Validate stored nonce matches signer
	if storedNonce.SignerID() != signerID {
		return signusecase.MuSig2SignOutput{}, fmt.Errorf(
			"nonce signer ID mismatch: expected %s, got %s",
			signerID,
			storedNonce.SignerID(),
		)
	}

	// Note: In a real implementation, we would need to:
	// 1. Parse the private key from WIF format
	// 2. Get all public keys from all signers
	// 3. Create MuSig2 context with keys
	// 4. Create session with stored nonce
	// 5. Register all public nonces
	// 6. Sign the message
	//
	// For now, this is a simplified implementation

	// Placeholder: Create a dummy private key for demonstration
	// In real implementation, this would come from authKey.WalletImportFormat
	// Using SHA256 to create a deterministic, full-length key for the placeholder.
	privKeyBytes := sha256.Sum256([]byte(authKey.WalletImportFormat))
	privKey, pubKey := btcec.PrivKeyFromBytes(privKeyBytes[:])

	// Placeholder: For demo, we'll use just this signer's public key
	allPubKeys := []*btcec.PublicKey{pubKey, pubKey}

	// Create MuSig2 context
	musig2Ctx, err := u.musig2Service.CreateContext(privKey, allPubKeys, true)
	if err != nil {
		return signusecase.MuSig2SignOutput{}, fmt.Errorf("failed to create MuSig2 context: %w", err)
	}

	// Create signing session
	session, err := u.musig2Service.CreateSession(musig2Ctx)
	if err != nil {
		return signusecase.MuSig2SignOutput{}, fmt.Errorf("failed to create MuSig2 session: %w", err)
	}

	// Register aggregated nonces from all signers
	for _, nonce := range input.AggregatedNonces {
		haveAll, err := u.musig2Service.RegisterPubNonce(session, nonce)
		if err != nil {
			return signusecase.MuSig2SignOutput{}, fmt.Errorf("failed to register public nonce: %w", err)
		}
		if haveAll {
			break
		}
	}

	// Create partial signature
	partialSig, err := u.musig2Service.Sign(session, input.MessageHash)
	if err != nil {
		return signusecase.MuSig2SignOutput{}, fmt.Errorf("failed to create partial signature: %w", err)
	}

	// Mark nonce as used (CRITICAL for security)
	err = u.nonceRepo.MarkNonceUsed(ctx, signerID, input.TransactionID)
	if err != nil {
		return signusecase.MuSig2SignOutput{}, fmt.Errorf("failed to mark nonce as used: %w", err)
	}

	// LIMITATION - Partial Signature Serialization:
	// The btcd musig2.PartialSignature type does not expose a Serialize() method or any way
	// to extract the underlying scalar value. The signature data is encapsulated within the
	// type with no public accessor methods.
	//
	// Current workaround: Return a placeholder until we either:
	//   1. Contribute a Serialize() method to the btcd library
	//   2. Use reflection to extract the internal field (not recommended)
	//   3. Find an alternative way to represent partial signatures
	//
	// Note: The partialSig is captured above but cannot be serialized yet.
	// This prevents the use case from producing a valid partial signature for aggregation.
	_ = partialSig // Captured but cannot serialize yet
	var sigScalar [32]byte
	copy(sigScalar[:], []byte("placeholder_partial_signature"))

	return signusecase.MuSig2SignOutput{
		PartialSignature: sigScalar,
		SignerID:         signerID,
	}, nil
}
