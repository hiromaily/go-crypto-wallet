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

type muSig2SignUseCase struct {
	musig2Service *btc.MuSig2Service
	nonceRepo     multisig.NonceRepository
}

// NewMuSig2SignUseCase creates a new MuSig2SignUseCase for BTC keygen.
func NewMuSig2SignUseCase(
	musig2Service *btc.MuSig2Service,
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
