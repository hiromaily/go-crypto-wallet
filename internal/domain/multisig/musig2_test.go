package multisig

import (
	"testing"
)

func TestNewNonceCommitment(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		publicNonce [66]byte
		signerID    string
		wantErr     bool
		errContains string
	}{
		{
			name: "valid nonce commitment",
			publicNonce: func() [66]byte {
				var nonce [66]byte
				for i := range nonce {
					nonce[i] = byte(i + 1)
				}
				return nonce
			}(),
			signerID: "signer1",
			wantErr:  false,
		},
		{
			name:        "empty signer ID",
			publicNonce: [66]byte{1, 2, 3},
			signerID:    "",
			wantErr:     true,
			errContains: "signer ID cannot be empty",
		},
		{
			name:        "zero nonce",
			publicNonce: [66]byte{},
			signerID:    "signer1",
			wantErr:     true,
			errContains: "public nonce cannot be zero",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			nonce, err := NewNonceCommitment(tt.publicNonce, tt.signerID)
			if tt.wantErr {
				if err == nil {
					t.Errorf("NewNonceCommitment() expected error but got none")
					return
				}
				if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("NewNonceCommitment() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("NewNonceCommitment() unexpected error = %v", err)
				return
			}

			if nonce.SignerID() != tt.signerID {
				t.Errorf("NewNonceCommitment() signer ID = %v, want %v", nonce.SignerID(), tt.signerID)
			}

			if nonce.PublicNonce() != tt.publicNonce {
				t.Errorf("NewNonceCommitment() public nonce mismatch")
			}
		})
	}
}

func TestNewPartialSignature(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		signatureScalar [32]byte
		signerID        string
		wantErr         bool
		errContains     string
	}{
		{
			name: "valid partial signature",
			signatureScalar: func() [32]byte {
				var scalar [32]byte
				for i := range scalar {
					scalar[i] = byte(i + 1)
				}
				return scalar
			}(),
			signerID: "signer1",
			wantErr:  false,
		},
		{
			name:            "empty signer ID",
			signatureScalar: [32]byte{1, 2, 3},
			signerID:        "",
			wantErr:         true,
			errContains:     "signer ID cannot be empty",
		},
		{
			name:            "zero signature scalar",
			signatureScalar: [32]byte{},
			signerID:        "signer1",
			wantErr:         true,
			errContains:     "signature scalar cannot be zero",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sig, err := NewPartialSignature(tt.signatureScalar, tt.signerID)
			if tt.wantErr {
				if err == nil {
					t.Errorf("NewPartialSignature() expected error but got none")
					return
				}
				if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("NewPartialSignature() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("NewPartialSignature() unexpected error = %v", err)
				return
			}

			if sig.SignerID() != tt.signerID {
				t.Errorf("NewPartialSignature() signer ID = %v, want %v", sig.SignerID(), tt.signerID)
			}

			if sig.SignatureScalar() != tt.signatureScalar {
				t.Errorf("NewPartialSignature() signature scalar mismatch")
			}
		})
	}
}

func TestNewAggregatedKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		publicKey   [33]byte
		tweaked     bool
		wantErr     bool
		errContains string
	}{
		{
			name: "valid aggregated key",
			publicKey: func() [33]byte {
				var key [33]byte
				key[0] = 2
				for i := 1; i < 33; i++ {
					key[i] = byte(i)
				}
				return key
			}(),
			tweaked: false,
			wantErr: false,
		},
		{
			name: "valid tweaked key",
			publicKey: func() [33]byte {
				var key [33]byte
				key[0] = 3
				for i := 1; i < 33; i++ {
					key[i] = byte(i)
				}
				return key
			}(),
			tweaked: true,
			wantErr: false,
		},
		{
			name:        "zero public key",
			publicKey:   [33]byte{},
			tweaked:     false,
			wantErr:     true,
			errContains: "aggregated public key cannot be zero",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			key, err := NewAggregatedKey(tt.publicKey, tt.tweaked)
			if tt.wantErr {
				if err == nil {
					t.Errorf("NewAggregatedKey() expected error but got none")
					return
				}
				if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("NewAggregatedKey() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("NewAggregatedKey() unexpected error = %v", err)
				return
			}

			if key.PublicKey() != tt.publicKey {
				t.Errorf("NewAggregatedKey() public key mismatch")
			}

			if key.IsTweaked() != tt.tweaked {
				t.Errorf("NewAggregatedKey() tweaked = %v, want %v", key.IsTweaked(), tt.tweaked)
			}
		})
	}
}

func TestNewSigningSession(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		sessionID        string
		participantCount int
		wantErr          bool
		errContains      string
	}{
		{
			name:             "valid session",
			sessionID:        "session1",
			participantCount: 2,
			wantErr:          false,
		},
		{
			name:             "empty session ID",
			sessionID:        "",
			participantCount: 2,
			wantErr:          true,
			errContains:      "session ID cannot be empty",
		},
		{
			name:             "invalid participant count (too low)",
			sessionID:        "session1",
			participantCount: 1,
			wantErr:          true,
			errContains:      "at least 2 signers",
		},
		{
			name:             "invalid participant count (too high)",
			sessionID:        "session1",
			participantCount: 20,
			wantErr:          true,
			errContains:      "cannot exceed 15",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			session, err := NewSigningSession(tt.sessionID, tt.participantCount)
			if tt.wantErr {
				if err == nil {
					t.Errorf("NewSigningSession() expected error but got none")
					return
				}
				if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("NewSigningSession() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("NewSigningSession() unexpected error = %v", err)
				return
			}

			if session.SessionID() != tt.sessionID {
				t.Errorf("NewSigningSession() session ID = %v, want %v", session.SessionID(), tt.sessionID)
			}

			if session.ParticipantCount() != tt.participantCount {
				t.Errorf("NewSigningSession() participant count = %v, want %v",
					session.ParticipantCount(), tt.participantCount)
			}

			if session.HasAllNonces() {
				t.Errorf("NewSigningSession() should not have all nonces initially")
			}

			if session.HasAllSignatures() {
				t.Errorf("NewSigningSession() should not have all signatures initially")
			}
		})
	}
}

func TestSigningSession_AddNonceCommitment(t *testing.T) {
	t.Parallel()

	session, err := NewSigningSession("test-session", 2)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	nonce1, err := NewNonceCommitment([66]byte{1, 2, 3}, "signer1")
	if err != nil {
		t.Fatalf("Failed to create nonce1: %v", err)
	}

	nonce2, err := NewNonceCommitment([66]byte{4, 5, 6}, "signer2")
	if err != nil {
		t.Fatalf("Failed to create nonce2: %v", err)
	}

	// Add first nonce
	if err := session.AddNonceCommitment(nonce1); err != nil {
		t.Errorf("AddNonceCommitment() unexpected error = %v", err)
	}

	if session.HasAllNonces() {
		t.Errorf("Session should not have all nonces yet")
	}

	// Add duplicate nonce should fail
	if err := session.AddNonceCommitment(nonce1); err == nil {
		t.Errorf("AddNonceCommitment() expected error for duplicate nonce")
	}

	// Add second nonce
	if err := session.AddNonceCommitment(nonce2); err != nil {
		t.Errorf("AddNonceCommitment() unexpected error = %v", err)
	}

	if !session.HasAllNonces() {
		t.Errorf("Session should have all nonces now")
	}

	// Test nil nonce
	if err := session.AddNonceCommitment(nil); err == nil {
		t.Errorf("AddNonceCommitment() expected error for nil nonce")
	}
}

func TestSigningSession_AddPartialSignature(t *testing.T) {
	t.Parallel()

	session, err := NewSigningSession("test-session", 2)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	sig1, err := NewPartialSignature([32]byte{1, 2, 3}, "signer1")
	if err != nil {
		t.Fatalf("Failed to create sig1: %v", err)
	}

	sig2, err := NewPartialSignature([32]byte{4, 5, 6}, "signer2")
	if err != nil {
		t.Fatalf("Failed to create sig2: %v", err)
	}

	// Add first signature
	if err := session.AddPartialSignature(sig1); err != nil {
		t.Errorf("AddPartialSignature() unexpected error = %v", err)
	}

	if session.HasAllSignatures() {
		t.Errorf("Session should not have all signatures yet")
	}

	// Add duplicate signature should fail
	if err := session.AddPartialSignature(sig1); err == nil {
		t.Errorf("AddPartialSignature() expected error for duplicate signature")
	}

	// Add second signature
	if err := session.AddPartialSignature(sig2); err != nil {
		t.Errorf("AddPartialSignature() unexpected error = %v", err)
	}

	if !session.HasAllSignatures() {
		t.Errorf("Session should have all signatures now")
	}

	// Test nil signature
	if err := session.AddPartialSignature(nil); err == nil {
		t.Errorf("AddPartialSignature() expected error for nil signature")
	}
}

func TestSigningSession_GetNoncesAndSignatures(t *testing.T) {
	t.Parallel()

	session, err := NewSigningSession("test-session", 2)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	nonce1, _ := NewNonceCommitment([66]byte{1, 2, 3}, "signer1")
	nonce2, _ := NewNonceCommitment([66]byte{4, 5, 6}, "signer2")
	sig1, _ := NewPartialSignature([32]byte{1, 2, 3}, "signer1")
	sig2, _ := NewPartialSignature([32]byte{4, 5, 6}, "signer2")

	_ = session.AddNonceCommitment(nonce1)
	_ = session.AddNonceCommitment(nonce2)
	_ = session.AddPartialSignature(sig1)
	_ = session.AddPartialSignature(sig2)

	nonces := session.GetNonceCommitments()
	if len(nonces) != 2 {
		t.Errorf("GetNonceCommitments() returned %d nonces, want 2", len(nonces))
	}

	signatures := session.GetPartialSignatures()
	if len(signatures) != 2 {
		t.Errorf("GetPartialSignatures() returned %d signatures, want 2", len(signatures))
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	if s == substr || len(substr) == 0 {
		return len(s) >= len(substr)
	}
	if len(s) > 0 && len(substr) > 0 {
		return containsHelper(s, substr)
	}
	return false
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
