package btc

import (
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcec/v2/schnorr/musig2"

	"github.com/hiromaily/go-crypto-wallet/internal/domain/multisig"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// MuSig2Service provides MuSig2 multi-signature functionality using btcd/btcec/v2.
// It implements the two-round MuSig2 protocol:
//   - Round 1: Nonce generation and exchange
//   - Round 2: Partial signature creation and aggregation
//
// The service uses the high-level Session API from btcd which provides built-in
// protection against nonce reuse and other common pitfalls.
type MuSig2Service struct {
	logger logger.Logger
}

// NewMuSig2Service creates a new MuSig2 service instance.
func NewMuSig2Service(log logger.Logger) *MuSig2Service {
	if log == nil {
		log = logger.NewNoopLogger()
	}
	return &MuSig2Service{
		logger: log,
	}
}

// validateContextInputs validates common inputs for context creation.
func validateContextInputs(privKey *btcec.PrivateKey, allPubKeys []*btcec.PublicKey) error {
	if privKey == nil {
		return errors.New("private key cannot be nil")
	}

	if len(allPubKeys) < 2 {
		return fmt.Errorf("at least 2 public keys required for MuSig2: got %d", len(allPubKeys))
	}

	// Validate using domain layer
	pubKeysBytes := make([][]byte, len(allPubKeys))
	for i, pk := range allPubKeys {
		pubKeysBytes[i] = pk.SerializeCompressed()
	}
	if err := multisig.ValidatePublicKeysForMuSig2(pubKeysBytes); err != nil {
		return fmt.Errorf("invalid public keys: %w", err)
	}

	return nil
}

// CreateContext creates a new MuSig2 context for a signer.
// The context manages key aggregation and must be created before starting a signing session.
//
// Parameters:
//   - privKey: The signer's private key
//   - allPubKeys: All participants' public keys (including this signer's)
//   - sortKeys: Whether to sort keys for deterministic aggregation (should be true)
//
// Returns:
//   - *musig2.Context: The created context
//   - error: Any error that occurred during context creation
func (m *MuSig2Service) CreateContext(
	privKey *btcec.PrivateKey,
	allPubKeys []*btcec.PublicKey,
	sortKeys bool,
) (*musig2.Context, error) {
	if err := validateContextInputs(privKey, allPubKeys); err != nil {
		return nil, err
	}

	// Create context with known signers
	ctx, err := musig2.NewContext(
		privKey,
		sortKeys,
		musig2.WithKnownSigners(allPubKeys),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create MuSig2 context: %w", err)
	}

	m.logger.Debug("MuSig2 context created", map[string]any{
		"num_signers": len(allPubKeys),
		"sort_keys":   sortKeys,
	})

	return ctx, nil
}

// CreateContextWithTaproot creates a new MuSig2 context with Taproot tweaking.
// This is used for Taproot key-path spending (BIP 86).
//
// Parameters:
//   - privKey: The signer's private key
//   - allPubKeys: All participants' public keys (including this signer's)
//   - sortKeys: Whether to sort keys for deterministic aggregation (should be true)
//
// Returns:
//   - *musig2.Context: The created context with Taproot tweak
//   - error: Any error that occurred during context creation
func (m *MuSig2Service) CreateContextWithTaproot(
	privKey *btcec.PrivateKey,
	allPubKeys []*btcec.PublicKey,
	sortKeys bool,
) (*musig2.Context, error) {
	if err := validateContextInputs(privKey, allPubKeys); err != nil {
		return nil, err
	}

	// Create context with Taproot BIP 86 tweak (key-path spending only)
	ctx, err := musig2.NewContext(
		privKey,
		sortKeys,
		musig2.WithKnownSigners(allPubKeys),
		musig2.WithBip86TweakCtx(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create MuSig2 context with Taproot: %w", err)
	}

	m.logger.Debug("MuSig2 context created with Taproot tweak", map[string]any{
		"num_signers": len(allPubKeys),
		"sort_keys":   sortKeys,
	})

	return ctx, nil
}

// CreateSession creates a new signing session from a context.
// Each transaction signing requires a new session to ensure fresh nonces.
//
// SECURITY: Never reuse sessions or nonces across different transactions.
// Nonce reuse will leak the private key!
//
// Parameters:
//   - ctx: The MuSig2 context
//
// Returns:
//   - *musig2.Session: The created session
//   - error: Any error that occurred during session creation
func (m *MuSig2Service) CreateSession(ctx *musig2.Context) (*musig2.Session, error) {
	if ctx == nil {
		return nil, errors.New("context cannot be nil")
	}

	session, err := ctx.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create MuSig2 session: %w", err)
	}

	m.logger.Debug("MuSig2 session created")

	return session, nil
}

// GetPublicNonce retrieves the public nonce from a session for Round 1.
// This nonce should be shared with all other signers.
//
// SECURITY: Only share the PUBLIC nonce, never the secret nonce.
//
// Parameters:
//   - session: The MuSig2 session
//
// Returns:
//   - [66]byte: The public nonce (two 33-byte compressed EC points)
func (*MuSig2Service) GetPublicNonce(session *musig2.Session) [66]byte {
	return session.PublicNonce()
}

// RegisterPubNonce registers a public nonce from another signer.
// This must be called for each other signer's nonce before signing.
//
// Parameters:
//   - session: The MuSig2 session
//   - pubNonce: The public nonce from another signer
//
// Returns:
//   - bool: Whether all nonces have been received
//   - error: Any error that occurred during nonce registration
func (m *MuSig2Service) RegisterPubNonce(
	session *musig2.Session,
	pubNonce [66]byte,
) (bool, error) {
	if session == nil {
		return false, errors.New("session cannot be nil")
	}

	// Validate nonce using domain layer
	var zero [66]byte
	if pubNonce == zero {
		return false, multisig.ErrPublicNonceZero
	}

	haveAllNonces, err := session.RegisterPubNonce(pubNonce)
	if err != nil {
		return false, fmt.Errorf("failed to register public nonce: %w", err)
	}

	m.logger.Debug("Public nonce registered", map[string]any{
		"have_all_nonces": haveAllNonces,
	})

	return haveAllNonces, nil
}

// Sign creates a partial signature for Round 2.
// This should only be called after all nonces have been registered.
//
// Parameters:
//   - session: The MuSig2 session (must have all nonces registered)
//   - messageHash: The 32-byte message hash to sign (e.g., transaction hash)
//
// Returns:
//   - *musig2.PartialSignature: The partial signature
//   - error: Any error that occurred during signing
func (m *MuSig2Service) Sign(
	session *musig2.Session,
	messageHash [32]byte,
) (*musig2.PartialSignature, error) {
	if session == nil {
		return nil, errors.New("session cannot be nil")
	}

	partialSig, err := session.Sign(messageHash)
	if err != nil {
		return nil, fmt.Errorf("failed to create partial signature: %w", err)
	}

	m.logger.Debug("Partial signature created")

	return partialSig, nil
}

// CombineSig adds a partial signature to the session.
// This must be called for each signer's partial signature.
//
// Parameters:
//   - session: The MuSig2 session
//   - partialSig: A partial signature from a signer
//
// Returns:
//   - bool: Whether all partial signatures have been received
//   - error: Any error that occurred during signature combination
func (m *MuSig2Service) CombineSig(
	session *musig2.Session,
	partialSig *musig2.PartialSignature,
) (bool, error) {
	if session == nil {
		return false, errors.New("session cannot be nil")
	}

	if partialSig == nil {
		return false, multisig.ErrPartialSignatureNil
	}

	haveAllSigs, err := session.CombineSig(partialSig)
	if err != nil {
		return false, fmt.Errorf("failed to combine partial signature: %w", err)
	}

	m.logger.Debug("Partial signature combined", map[string]any{
		"have_all_sigs": haveAllSigs,
	})

	return haveAllSigs, nil
}

// FinalSig retrieves the final aggregated signature.
// This should only be called after all partial signatures have been combined.
//
// Parameters:
//   - session: The MuSig2 session (must have all partial signatures)
//
// Returns:
//   - *schnorr.Signature: The final aggregated signature (64 bytes)
//   - error: Any error that occurred
func (m *MuSig2Service) FinalSig(session *musig2.Session) (*schnorr.Signature, error) {
	if session == nil {
		return nil, errors.New("session cannot be nil")
	}

	finalSig := session.FinalSig()
	if finalSig == nil {
		return nil, errors.New("final signature not available")
	}

	m.logger.Debug("Final signature retrieved")

	return finalSig, nil
}

// GetCombinedKey retrieves the aggregated public key from the context.
// This key can be used to verify the final signature.
//
// Parameters:
//   - ctx: The MuSig2 context
//
// Returns:
//   - *btcec.PublicKey: The aggregated public key
//   - error: Any error that occurred
func (m *MuSig2Service) GetCombinedKey(ctx *musig2.Context) (*btcec.PublicKey, error) {
	if ctx == nil {
		return nil, errors.New("context cannot be nil")
	}

	aggregatedKey, err := ctx.CombinedKey()
	if err != nil {
		return nil, fmt.Errorf("failed to get combined key: %w", err)
	}

	// Validate using domain layer
	keyBytes := aggregatedKey.SerializeCompressed()
	if err := multisig.ValidateAggregatedPublicKey(keyBytes); err != nil {
		return nil, fmt.Errorf("invalid aggregated key: %w", err)
	}

	m.logger.Debug("Combined key retrieved")

	return aggregatedKey, nil
}

// VerifySignature verifies a MuSig2 signature against a message and aggregated key.
//
// Parameters:
//   - signature: The aggregated signature to verify
//   - messageHash: The 32-byte message hash
//   - aggregatedKey: The aggregated public key
//
// Returns:
//   - bool: Whether the signature is valid
func (m *MuSig2Service) VerifySignature(
	signature *schnorr.Signature,
	messageHash [32]byte,
	aggregatedKey *btcec.PublicKey,
) bool {
	if signature == nil || aggregatedKey == nil {
		m.logger.Warn("Invalid input for signature verification", map[string]any{
			"signature_nil": signature == nil,
			"key_nil":       aggregatedKey == nil,
		})
		return false
	}

	valid := signature.Verify(messageHash[:], aggregatedKey)

	m.logger.Debug("Signature verification completed", map[string]any{
		"valid": valid,
	})

	return valid
}
