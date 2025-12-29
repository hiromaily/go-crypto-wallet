package multisig

import (
	"bytes"
	"errors"
	"fmt"
)

// NonceCommitment represents a MuSig2 public nonce commitment from a signer.
// This is safe to share with all participants in the signing protocol.
//
// A nonce commitment consists of a 66-byte public nonce (two 33-byte compressed elliptic curve points)
// and a signer identifier to track which signer provided the nonce.
type NonceCommitment struct {
	publicNonce [66]byte // Two 33-byte compressed EC points (R1 || R2)
	signerID    string   // Identifier for the signer who generated this nonce
}

// NewNonceCommitment creates a new nonce commitment with validation.
func NewNonceCommitment(publicNonce [66]byte, signerID string) (*NonceCommitment, error) {
	if signerID == "" {
		return nil, errors.New("signer ID cannot be empty")
	}

	// Validate that nonce is not all zeros
	var zero [66]byte
	if bytes.Equal(publicNonce[:], zero[:]) {
		return nil, errors.New("public nonce cannot be zero")
	}

	return &NonceCommitment{
		publicNonce: publicNonce,
		signerID:    signerID,
	}, nil
}

// PublicNonce returns a copy of the public nonce.
func (n *NonceCommitment) PublicNonce() [66]byte {
	return n.publicNonce
}

// SignerID returns the signer identifier.
func (n *NonceCommitment) SignerID() string {
	return n.signerID
}

// PartialSignature represents a MuSig2 partial signature from a single signer.
// Partial signatures from all signers are aggregated to produce the final signature.
//
// A partial signature contains the signature scalar (32 bytes) and the signer identifier
// to track which signer produced this partial signature.
type PartialSignature struct {
	signatureScalar [32]byte // The 's' value in the Schnorr signature
	signerID        string   // Identifier for the signer who generated this signature
}

// NewPartialSignature creates a new partial signature with validation.
func NewPartialSignature(signatureScalar [32]byte, signerID string) (*PartialSignature, error) {
	if signerID == "" {
		return nil, errors.New("signer ID cannot be empty")
	}

	// Validate that signature scalar is not all zeros
	var zero [32]byte
	if bytes.Equal(signatureScalar[:], zero[:]) {
		return nil, errors.New("signature scalar cannot be zero")
	}

	return &PartialSignature{
		signatureScalar: signatureScalar,
		signerID:        signerID,
	}, nil
}

// SignatureScalar returns a copy of the signature scalar.
func (p *PartialSignature) SignatureScalar() [32]byte {
	return p.signatureScalar
}

// SignerID returns the signer identifier.
func (p *PartialSignature) SignerID() string {
	return p.signerID
}

// AggregatedKey represents a MuSig2 aggregated public key.
// This is the combined public key from all participants in the multisig.
type AggregatedKey struct {
	publicKey [33]byte // Compressed 33-byte aggregated public key
	tweaked   bool     // Whether the key has been tweaked (e.g., for Taproot)
}

// NewAggregatedKey creates a new aggregated key with validation.
func NewAggregatedKey(publicKey [33]byte, tweaked bool) (*AggregatedKey, error) {
	// Validate that public key is not all zeros
	var zero [33]byte
	if bytes.Equal(publicKey[:], zero[:]) {
		return nil, errors.New("aggregated public key cannot be zero")
	}

	return &AggregatedKey{
		publicKey: publicKey,
		tweaked:   tweaked,
	}, nil
}

// PublicKey returns a copy of the aggregated public key.
func (a *AggregatedKey) PublicKey() [33]byte {
	return a.publicKey
}

// IsTweaked returns whether the key has been tweaked.
func (a *AggregatedKey) IsTweaked() bool {
	return a.tweaked
}

// SigningSession represents the state of a MuSig2 signing session.
// A signing session tracks the nonce commitments and partial signatures
// collected during the two-round signing protocol.
type SigningSession struct {
	sessionID         string                       // Unique identifier for this session
	participantCount  int                          // Total number of participants
	nonceCommitments  map[string]*NonceCommitment  // Nonces indexed by signer ID
	partialSignatures map[string]*PartialSignature // Partial signatures indexed by signer ID
}

// NewSigningSession creates a new signing session with validation.
func NewSigningSession(sessionID string, participantCount int) (*SigningSession, error) {
	if sessionID == "" {
		return nil, errors.New("session ID cannot be empty")
	}

	if err := ValidateSignerCount(participantCount); err != nil {
		return nil, fmt.Errorf("invalid participant count: %w", err)
	}

	return &SigningSession{
		sessionID:         sessionID,
		participantCount:  participantCount,
		nonceCommitments:  make(map[string]*NonceCommitment),
		partialSignatures: make(map[string]*PartialSignature),
	}, nil
}

// SessionID returns the session identifier.
func (s *SigningSession) SessionID() string {
	return s.sessionID
}

// ParticipantCount returns the total number of participants.
func (s *SigningSession) ParticipantCount() int {
	return s.participantCount
}

// AddNonceCommitment adds a nonce commitment to the session.
func (s *SigningSession) AddNonceCommitment(nonce *NonceCommitment) error {
	if nonce == nil {
		return errors.New("nonce commitment cannot be nil")
	}

	signerID := nonce.SignerID()
	if _, exists := s.nonceCommitments[signerID]; exists {
		return fmt.Errorf("nonce commitment already exists for signer %s", signerID)
	}

	s.nonceCommitments[signerID] = nonce
	return nil
}

// AddPartialSignature adds a partial signature to the session.
func (s *SigningSession) AddPartialSignature(sig *PartialSignature) error {
	if sig == nil {
		return errors.New("partial signature cannot be nil")
	}

	signerID := sig.SignerID()
	if _, exists := s.partialSignatures[signerID]; exists {
		return fmt.Errorf("partial signature already exists for signer %s", signerID)
	}

	s.partialSignatures[signerID] = sig
	return nil
}

// HasAllNonces returns true if all expected nonce commitments have been received.
func (s *SigningSession) HasAllNonces() bool {
	return len(s.nonceCommitments) == s.participantCount
}

// HasAllSignatures returns true if all expected partial signatures have been received.
func (s *SigningSession) HasAllSignatures() bool {
	return len(s.partialSignatures) == s.participantCount
}

// GetNonceCommitments returns all nonce commitments.
func (s *SigningSession) GetNonceCommitments() []*NonceCommitment {
	nonces := make([]*NonceCommitment, 0, len(s.nonceCommitments))
	for _, nonce := range s.nonceCommitments {
		nonces = append(nonces, nonce)
	}
	return nonces
}

// GetPartialSignatures returns all partial signatures.
func (s *SigningSession) GetPartialSignatures() []*PartialSignature {
	signatures := make([]*PartialSignature, 0, len(s.partialSignatures))
	for _, sig := range s.partialSignatures {
		signatures = append(signatures, sig)
	}
	return signatures
}
