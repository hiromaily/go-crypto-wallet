package multisig

import (
	"testing"
)

func TestValidateSignerCount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		count       int
		wantErr     bool
		errContains string
	}{
		{
			name:    "valid count 2",
			count:   2,
			wantErr: false,
		},
		{
			name:    "valid count 5",
			count:   5,
			wantErr: false,
		},
		{
			name:    "valid count 15",
			count:   15,
			wantErr: false,
		},
		{
			name:        "invalid count 0",
			count:       0,
			wantErr:     true,
			errContains: "at least 2 signers",
		},
		{
			name:        "invalid count 1",
			count:       1,
			wantErr:     true,
			errContains: "at least 2 signers",
		},
		{
			name:        "invalid count 16",
			count:       16,
			wantErr:     true,
			errContains: "cannot exceed 15",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateSignerCount(tt.count)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateSignerCount() expected error but got none")
					return
				}
				if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("ValidateSignerCount() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("ValidateSignerCount() unexpected error = %v", err)
			}
		})
	}
}

func TestValidateNonceUniqueness(t *testing.T) {
	t.Parallel()

	nonce1, _ := NewNonceCommitment([66]byte{1, 2, 3}, "signer1")
	nonce2, _ := NewNonceCommitment([66]byte{4, 5, 6}, "signer2")
	nonce3, _ := NewNonceCommitment([66]byte{7, 8, 9}, "signer3")
	nonceDuplicate, _ := NewNonceCommitment([66]byte{1, 2, 3}, "signer4")

	tests := []struct {
		name        string
		nonces      []*NonceCommitment
		wantErr     bool
		errContains string
	}{
		{
			name:    "valid unique nonces",
			nonces:  []*NonceCommitment{nonce1, nonce2, nonce3},
			wantErr: false,
		},
		{
			name:        "empty nonces",
			nonces:      []*NonceCommitment{},
			wantErr:     true,
			errContains: "no nonces provided",
		},
		{
			name:        "nil nonce in list",
			nonces:      []*NonceCommitment{nonce1, nil, nonce2},
			wantErr:     true,
			errContains: "nonce at index",
		},
		{
			name:        "duplicate nonce",
			nonces:      []*NonceCommitment{nonce1, nonce2, nonceDuplicate},
			wantErr:     true,
			errContains: "duplicate nonce detected",
		},
		{
			name: "duplicate signer ID",
			nonces: []*NonceCommitment{
				nonce1,
				func() *NonceCommitment {
					n, _ := NewNonceCommitment([66]byte{10, 11, 12}, "signer1")
					return n
				}(),
			},
			wantErr:     true,
			errContains: "duplicate signer ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateNonceUniqueness(tt.nonces)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateNonceUniqueness() expected error but got none")
					return
				}
				if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("ValidateNonceUniqueness() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("ValidateNonceUniqueness() unexpected error = %v", err)
			}
		})
	}
}

func TestValidatePartialSignatures(t *testing.T) {
	t.Parallel()

	sig1, _ := NewPartialSignature([32]byte{1, 2, 3}, "signer1")
	sig2, _ := NewPartialSignature([32]byte{4, 5, 6}, "signer2")
	sig3, _ := NewPartialSignature([32]byte{7, 8, 9}, "signer3")

	tests := []struct {
		name        string
		signatures  []*PartialSignature
		wantErr     bool
		errContains string
	}{
		{
			name:       "valid signatures",
			signatures: []*PartialSignature{sig1, sig2, sig3},
			wantErr:    false,
		},
		{
			name:        "empty signatures",
			signatures:  []*PartialSignature{},
			wantErr:     true,
			errContains: "no partial signatures provided",
		},
		{
			name:        "nil signature in list",
			signatures:  []*PartialSignature{sig1, nil, sig2},
			wantErr:     true,
			errContains: "partial signature at index",
		},
		{
			name: "duplicate signer ID",
			signatures: []*PartialSignature{
				sig1,
				func() *PartialSignature {
					s, _ := NewPartialSignature([32]byte{10, 11, 12}, "signer1")
					return s
				}(),
			},
			wantErr:     true,
			errContains: "duplicate signer ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := ValidatePartialSignatures(tt.signatures)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidatePartialSignatures() expected error but got none")
					return
				}
				if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("ValidatePartialSignatures() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("ValidatePartialSignatures() unexpected error = %v", err)
			}
		})
	}
}

func TestValidatePublicKeysForMuSig2(t *testing.T) {
	t.Parallel()

	validKey33 := make([]byte, 33)
	validKey33[0] = 2
	for i := 1; i < 33; i++ {
		validKey33[i] = byte(i)
	}

	validKey65 := make([]byte, 65)
	validKey65[0] = 4
	for i := 1; i < 65; i++ {
		validKey65[i] = byte(i)
	}

	anotherKey33 := make([]byte, 33)
	if len(anotherKey33) > 0 {
		anotherKey33[0] = 3
	}
	for i := 1; i < len(anotherKey33); i++ {
		anotherKey33[i] = byte(i + 10)
	}

	tests := []struct {
		name        string
		publicKeys  [][]byte
		wantErr     bool
		errContains string
	}{
		{
			name:       "valid 2 keys (33 bytes)",
			publicKeys: [][]byte{validKey33, anotherKey33},
			wantErr:    false,
		},
		{
			name:       "valid 2 keys (mixed 33 and 65)",
			publicKeys: [][]byte{validKey33, validKey65},
			wantErr:    false,
		},
		{
			name:        "insufficient keys (1)",
			publicKeys:  [][]byte{validKey33},
			wantErr:     true,
			errContains: "at least 2 public keys required",
		},
		{
			name:        "empty key list",
			publicKeys:  [][]byte{},
			wantErr:     true,
			errContains: "at least 2 public keys required",
		},
		{
			name:        "too many keys (16)",
			publicKeys:  make([][]byte, 16),
			wantErr:     true,
			errContains: "cannot exceed 15",
		},
		{
			name:        "empty public key",
			publicKeys:  [][]byte{validKey33, {}},
			wantErr:     true,
			errContains: "public key at index 1 is empty",
		},
		{
			name:        "invalid key length",
			publicKeys:  [][]byte{validKey33, make([]byte, 32)},
			wantErr:     true,
			errContains: "invalid length",
		},
		{
			name:        "duplicate public key",
			publicKeys:  [][]byte{validKey33, validKey33},
			wantErr:     true,
			errContains: "duplicate public key",
		},
		{
			name:        "zero public key",
			publicKeys:  [][]byte{validKey33, make([]byte, 33)},
			wantErr:     true,
			errContains: "all zeros",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := ValidatePublicKeysForMuSig2(tt.publicKeys)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidatePublicKeysForMuSig2() expected error but got none")
					return
				}
				if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("ValidatePublicKeysForMuSig2() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("ValidatePublicKeysForMuSig2() unexpected error = %v", err)
			}
		})
	}
}

func TestValidateAggregatedPublicKey(t *testing.T) {
	t.Parallel()

	validKey := make([]byte, 33)
	validKey[0] = 2
	for i := 1; i < 33; i++ {
		validKey[i] = byte(i)
	}

	tests := []struct {
		name        string
		publicKey   []byte
		wantErr     bool
		errContains string
	}{
		{
			name:      "valid aggregated key",
			publicKey: validKey,
			wantErr:   false,
		},
		{
			name:        "empty key",
			publicKey:   []byte{},
			wantErr:     true,
			errContains: "aggregated public key is empty",
		},
		{
			name:        "invalid length (32 bytes)",
			publicKey:   make([]byte, 32),
			wantErr:     true,
			errContains: "invalid length",
		},
		{
			name:        "invalid length (65 bytes)",
			publicKey:   make([]byte, 65),
			wantErr:     true,
			errContains: "invalid length",
		},
		{
			name:        "zero key",
			publicKey:   make([]byte, 33),
			wantErr:     true,
			errContains: "all zeros",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateAggregatedPublicKey(tt.publicKey)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateAggregatedPublicKey() expected error but got none")
					return
				}
				if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("ValidateAggregatedPublicKey() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("ValidateAggregatedPublicKey() unexpected error = %v", err)
			}
		})
	}
}

func TestValidateSigningSessionComplete(t *testing.T) {
	t.Parallel()

	// Create a complete session
	completeSession, _ := NewSigningSession("complete", 2)
	nonce1, _ := NewNonceCommitment([66]byte{1, 2, 3}, "signer1")
	nonce2, _ := NewNonceCommitment([66]byte{4, 5, 6}, "signer2")
	sig1, _ := NewPartialSignature([32]byte{1, 2, 3}, "signer1")
	sig2, _ := NewPartialSignature([32]byte{4, 5, 6}, "signer2")
	_ = completeSession.AddNonceCommitment(nonce1)
	_ = completeSession.AddNonceCommitment(nonce2)
	_ = completeSession.AddPartialSignature(sig1)
	_ = completeSession.AddPartialSignature(sig2)

	// Create incomplete session (missing signatures)
	incompleteSession1, _ := NewSigningSession("incomplete1", 2)
	_ = incompleteSession1.AddNonceCommitment(nonce1)
	_ = incompleteSession1.AddNonceCommitment(nonce2)

	// Create incomplete session (missing nonces)
	incompleteSession2, _ := NewSigningSession("incomplete2", 2)
	_ = incompleteSession2.AddPartialSignature(sig1)
	_ = incompleteSession2.AddPartialSignature(sig2)

	tests := []struct {
		name        string
		session     *SigningSession
		wantErr     bool
		errContains string
	}{
		{
			name:    "complete session",
			session: completeSession,
			wantErr: false,
		},
		{
			name:        "nil session",
			session:     nil,
			wantErr:     true,
			errContains: "signing session is nil",
		},
		{
			name:        "missing signatures",
			session:     incompleteSession1,
			wantErr:     true,
			errContains: "missing partial signatures",
		},
		{
			name:        "missing nonces",
			session:     incompleteSession2,
			wantErr:     true,
			errContains: "missing nonce commitments",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateSigningSessionComplete(tt.session)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateSigningSessionComplete() expected error but got none")
					return
				}
				if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("ValidateSigningSessionComplete() error = %v, want error containing %v",
						err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("ValidateSigningSessionComplete() unexpected error = %v", err)
			}
		})
	}
}
