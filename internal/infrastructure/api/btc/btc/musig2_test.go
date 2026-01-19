package btc

import (
	"crypto/sha256"
	"strings"
	"testing"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcec/v2/schnorr/musig2"
)

func TestNewMuSig2Service(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		wantNil bool
	}{
		{
			name:    "happy path",
			wantNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service := NewMuSig2Service()
			if (service == nil) != tt.wantNil {
				t.Errorf("NewMuSig2Service() = %v, wantNil %v", service, tt.wantNil)
			}
		})
	}
}

func TestMuSig2Service_CreateContext(t *testing.T) {
	t.Parallel()

	service := NewMuSig2Service()

	// Generate test keys
	privKey1, err := btcec.NewPrivateKey()
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	privKey2, err := btcec.NewPrivateKey()
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	pubKey1 := privKey1.PubKey()
	pubKey2 := privKey2.PubKey()

	tests := []struct {
		name        string
		privKey     *btcec.PrivateKey
		allPubKeys  []*btcec.PublicKey
		sortKeys    bool
		wantErr     bool
		errContains string
	}{
		{
			name:       "valid 2-of-2",
			privKey:    privKey1,
			allPubKeys: []*btcec.PublicKey{pubKey1, pubKey2},
			sortKeys:   true,
			wantErr:    false,
		},
		{
			name:       "valid unsorted keys",
			privKey:    privKey1,
			allPubKeys: []*btcec.PublicKey{pubKey1, pubKey2},
			sortKeys:   false,
			wantErr:    false,
		},
		{
			name:        "nil private key",
			privKey:     nil,
			allPubKeys:  []*btcec.PublicKey{pubKey1, pubKey2},
			sortKeys:    true,
			wantErr:     true,
			errContains: "private key cannot be nil",
		},
		{
			name:        "insufficient public keys",
			privKey:     privKey1,
			allPubKeys:  []*btcec.PublicKey{pubKey1},
			sortKeys:    true,
			wantErr:     true,
			errContains: "at least 2 public keys required",
		},
		{
			name:        "empty public keys",
			privKey:     privKey1,
			allPubKeys:  []*btcec.PublicKey{},
			sortKeys:    true,
			wantErr:     true,
			errContains: "at least 2 public keys required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, err := service.CreateContext(tt.privKey, tt.allPubKeys, tt.sortKeys)
			if tt.wantErr {
				if err == nil {
					t.Errorf("CreateContext() expected error but got none")
					return
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("CreateContext() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("CreateContext() unexpected error = %v", err)
				return
			}

			if ctx == nil {
				t.Errorf("CreateContext() returned nil context")
			}
		})
	}
}

func TestMuSig2Service_CreateContextWithTaproot(t *testing.T) {
	t.Parallel()

	service := NewMuSig2Service()

	// Generate test keys
	privKey1, err := btcec.NewPrivateKey()
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	privKey2, err := btcec.NewPrivateKey()
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	pubKey1 := privKey1.PubKey()
	pubKey2 := privKey2.PubKey()

	tests := []struct {
		name        string
		privKey     *btcec.PrivateKey
		allPubKeys  []*btcec.PublicKey
		sortKeys    bool
		wantErr     bool
		errContains string
	}{
		{
			name:       "valid 2-of-2 with Taproot",
			privKey:    privKey1,
			allPubKeys: []*btcec.PublicKey{pubKey1, pubKey2},
			sortKeys:   true,
			wantErr:    false,
		},
		{
			name:        "nil private key",
			privKey:     nil,
			allPubKeys:  []*btcec.PublicKey{pubKey1, pubKey2},
			sortKeys:    true,
			wantErr:     true,
			errContains: "private key cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, err := service.CreateContextWithTaproot(tt.privKey, tt.allPubKeys, tt.sortKeys)
			if tt.wantErr {
				if err == nil {
					t.Errorf("CreateContextWithTaproot() expected error but got none")
					return
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("CreateContextWithTaproot() error = %v, want error containing %v",
						err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("CreateContextWithTaproot() unexpected error = %v", err)
				return
			}

			if ctx == nil {
				t.Errorf("CreateContextWithTaproot() returned nil context")
			}
		})
	}
}

func TestMuSig2Service_CreateSession(t *testing.T) {
	t.Parallel()

	service := NewMuSig2Service()

	// Generate test keys
	privKey1, _ := btcec.NewPrivateKey()
	privKey2, _ := btcec.NewPrivateKey()
	pubKey1 := privKey1.PubKey()
	pubKey2 := privKey2.PubKey()

	// Create a valid context
	ctx, err := service.CreateContext(privKey1, []*btcec.PublicKey{pubKey1, pubKey2}, true)
	if err != nil {
		t.Fatalf("Failed to create context: %v", err)
	}

	tests := []struct {
		name        string
		ctx         *musig2.Context
		wantErr     bool
		errContains string
	}{
		{
			name:    "valid context",
			ctx:     ctx,
			wantErr: false,
		},
		{
			name:        "nil context",
			ctx:         nil,
			wantErr:     true,
			errContains: "context cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			session, err := service.CreateSession(tt.ctx)
			if tt.wantErr {
				if err == nil {
					t.Errorf("CreateSession() expected error but got none")
					return
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("CreateSession() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("CreateSession() unexpected error = %v", err)
				return
			}

			if session == nil {
				t.Errorf("CreateSession() returned nil session")
			}
		})
	}
}

// testSetupSigner creates a context and session for a signer.
func testSetupSigner(
	t *testing.T,
	service *MuSig2Service,
	privKey *btcec.PrivateKey,
	allPubKeys []*btcec.PublicKey,
	signerNum int,
) (*musig2.Context, *musig2.Session) {
	t.Helper()

	ctx, err := service.CreateContext(privKey, allPubKeys, true)
	if err != nil {
		t.Fatalf("Failed to create context for signer %d: %v", signerNum, err)
	}

	session, err := service.CreateSession(ctx)
	if err != nil {
		t.Fatalf("Failed to create session for signer %d: %v", signerNum, err)
	}

	return ctx, session
}

// testExchangeNonces performs Round 1: nonce generation and exchange between two sessions.
func testExchangeNonces(t *testing.T, service *MuSig2Service, session1, session2 *musig2.Session) {
	t.Helper()

	nonce1 := service.GetPublicNonce(session1)
	nonce2 := service.GetPublicNonce(session2)

	// Session 1 registers session 2's nonce
	haveAll, err := service.RegisterPubNonce(session1, nonce2)
	if err != nil {
		t.Fatalf("Failed to register nonce in session 1: %v", err)
	}
	if !haveAll {
		t.Errorf("Session 1 should have all nonces after registration")
	}

	// Session 2 registers session 1's nonce
	haveAll, err = service.RegisterPubNonce(session2, nonce1)
	if err != nil {
		t.Fatalf("Failed to register nonce in session 2: %v", err)
	}
	if !haveAll {
		t.Errorf("Session 2 should have all nonces after registration")
	}
}

// testSignAndAggregate performs Round 2: signing and aggregation.
func testSignAndAggregate(
	t *testing.T,
	service *MuSig2Service,
	session1, session2 *musig2.Session,
	message [32]byte,
) *schnorr.Signature {
	t.Helper()

	// Both signers create partial signatures
	_, err := service.Sign(session1, message)
	if err != nil {
		t.Fatalf("Failed to sign with session 1: %v", err)
	}

	partialSig2, err := service.Sign(session2, message)
	if err != nil {
		t.Fatalf("Failed to sign with session 2: %v", err)
	}

	// Aggregate in session 1
	haveAll, err := service.CombineSig(session1, partialSig2)
	if err != nil {
		t.Fatalf("Failed to combine signature: %v", err)
	}
	if !haveAll {
		t.Errorf("Session 1 should have all signatures after combination")
	}

	// Get final signature
	finalSig, err := service.FinalSig(session1)
	if err != nil {
		t.Fatalf("Failed to get final signature: %v", err)
	}
	if finalSig == nil {
		t.Fatalf("Final signature is nil")
	}

	return finalSig
}

func TestMuSig2Service_CompleteTwoRoundProtocol(t *testing.T) {
	t.Parallel()

	service := NewMuSig2Service()

	// Generate test keys for 2-of-2 multisig
	privKey1, _ := btcec.NewPrivateKey()
	privKey2, _ := btcec.NewPrivateKey()
	pubKey1 := privKey1.PubKey()
	pubKey2 := privKey2.PubKey()
	allPubKeys := []*btcec.PublicKey{pubKey1, pubKey2}

	// Message to sign
	message := sha256.Sum256([]byte("test transaction"))

	// Setup signers
	ctx1, session1 := testSetupSigner(t, service, privKey1, allPubKeys, 1)
	_, session2 := testSetupSigner(t, service, privKey2, allPubKeys, 2)

	// Round 1: Nonce generation and exchange
	testExchangeNonces(t, service, session1, session2)

	// Round 2: Signing and aggregation
	finalSig := testSignAndAggregate(t, service, session1, session2, message)

	// Verify signature
	aggregatedKey, err := service.GetCombinedKey(ctx1)
	if err != nil {
		t.Fatalf("Failed to get combined key: %v", err)
	}

	valid := service.VerifySignature(finalSig, message, aggregatedKey)
	if !valid {
		t.Errorf("Signature verification failed")
	}

	// Verify signature size (Schnorr signatures are 64 bytes)
	sigBytes := finalSig.Serialize()
	if len(sigBytes) != 64 {
		t.Errorf("Final signature size = %d bytes, want 64 bytes", len(sigBytes))
	}
}

func TestMuSig2Service_RegisterPubNonce(t *testing.T) {
	// Note: Cannot use t.Parallel() because tests mutate session state

	service := NewMuSig2Service()

	// Generate test keys
	privKey1, _ := btcec.NewPrivateKey()
	privKey2, _ := btcec.NewPrivateKey()
	pubKey1 := privKey1.PubKey()
	pubKey2 := privKey2.PubKey()

	// Create context and session
	ctx, _ := service.CreateContext(privKey1, []*btcec.PublicKey{pubKey1, pubKey2}, true)
	session, _ := service.CreateSession(ctx)

	// Get a valid nonce from another session
	ctx2, _ := service.CreateContext(privKey2, []*btcec.PublicKey{pubKey1, pubKey2}, true)
	session2, _ := service.CreateSession(ctx2)
	validNonce := service.GetPublicNonce(session2)

	tests := []struct {
		name        string
		session     *musig2.Session
		nonce       [66]byte
		wantErr     bool
		errContains string
	}{
		{
			name:    "valid nonce",
			session: session,
			nonce:   validNonce,
			wantErr: false,
		},
		{
			name:        "nil session",
			session:     nil,
			nonce:       validNonce,
			wantErr:     true,
			errContains: "session cannot be nil",
		},
		{
			name:        "zero nonce",
			session:     session,
			nonce:       [66]byte{},
			wantErr:     true,
			errContains: "public nonce cannot be zero",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: Cannot use t.Parallel() here because session is mutated

			_, err := service.RegisterPubNonce(tt.session, tt.nonce)
			if tt.wantErr {
				if err == nil {
					t.Errorf("RegisterPubNonce() expected error but got none")
					return
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("RegisterPubNonce() error = %v, want error containing %v",
						err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("RegisterPubNonce() unexpected error = %v", err)
			}
		})
	}
}

func TestMuSig2Service_Sign(t *testing.T) {
	// Note: Cannot use t.Parallel() because tests mutate session state

	service := NewMuSig2Service()

	// Generate test keys
	privKey1, _ := btcec.NewPrivateKey()
	privKey2, _ := btcec.NewPrivateKey()
	pubKey1 := privKey1.PubKey()
	pubKey2 := privKey2.PubKey()

	// Create context and sessions
	ctx1, _ := service.CreateContext(privKey1, []*btcec.PublicKey{pubKey1, pubKey2}, true)
	session1, _ := service.CreateSession(ctx1)

	ctx2, _ := service.CreateContext(privKey2, []*btcec.PublicKey{pubKey1, pubKey2}, true)
	session2, _ := service.CreateSession(ctx2)

	// Exchange nonces
	nonce1 := service.GetPublicNonce(session1)
	nonce2 := service.GetPublicNonce(session2)
	_, _ = service.RegisterPubNonce(session1, nonce2)
	_, _ = service.RegisterPubNonce(session2, nonce1)

	message := sha256.Sum256([]byte("test message"))

	tests := []struct {
		name        string
		session     *musig2.Session
		messageHash [32]byte
		wantErr     bool
		errContains string
	}{
		{
			name:        "valid signing",
			session:     session1,
			messageHash: message,
			wantErr:     false,
		},
		{
			name:        "nil session",
			session:     nil,
			messageHash: message,
			wantErr:     true,
			errContains: "session cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: Cannot use t.Parallel() here because session is mutated

			sig, err := service.Sign(tt.session, tt.messageHash)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Sign() expected error but got none")
					return
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Sign() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("Sign() unexpected error = %v", err)
				return
			}

			if sig == nil {
				t.Errorf("Sign() returned nil signature")
			}
		})
	}
}

func TestMuSig2Service_GetCombinedKey(t *testing.T) {
	t.Parallel()

	service := NewMuSig2Service()

	// Generate test keys
	privKey1, _ := btcec.NewPrivateKey()
	privKey2, _ := btcec.NewPrivateKey()
	pubKey1 := privKey1.PubKey()
	pubKey2 := privKey2.PubKey()

	ctx, _ := service.CreateContext(privKey1, []*btcec.PublicKey{pubKey1, pubKey2}, true)

	tests := []struct {
		name        string
		ctx         *musig2.Context
		wantErr     bool
		errContains string
	}{
		{
			name:    "valid context",
			ctx:     ctx,
			wantErr: false,
		},
		{
			name:        "nil context",
			ctx:         nil,
			wantErr:     true,
			errContains: "context cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			key, err := service.GetCombinedKey(tt.ctx)
			if tt.wantErr {
				if err == nil {
					t.Errorf("GetCombinedKey() expected error but got none")
					return
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("GetCombinedKey() error = %v, want error containing %v",
						err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("GetCombinedKey() unexpected error = %v", err)
				return
			}

			if key == nil {
				t.Errorf("GetCombinedKey() returned nil key")
			}
		})
	}
}

func TestMuSig2Service_VerifySignature(t *testing.T) {
	t.Parallel()

	service := NewMuSig2Service()

	// Create a complete signature for testing
	privKey1, _ := btcec.NewPrivateKey()
	privKey2, _ := btcec.NewPrivateKey()
	pubKey1 := privKey1.PubKey()
	pubKey2 := privKey2.PubKey()
	allPubKeys := []*btcec.PublicKey{pubKey1, pubKey2}

	message := sha256.Sum256([]byte("test message"))

	// Complete protocol to get a valid signature
	ctx1, _ := service.CreateContext(privKey1, allPubKeys, true)
	session1, _ := service.CreateSession(ctx1)

	ctx2, _ := service.CreateContext(privKey2, allPubKeys, true)
	session2, _ := service.CreateSession(ctx2)

	nonce1 := service.GetPublicNonce(session1)
	nonce2 := service.GetPublicNonce(session2)
	_, _ = service.RegisterPubNonce(session1, nonce2)
	_, _ = service.RegisterPubNonce(session2, nonce1)

	_, _ = service.Sign(session1, message)
	partialSig2, _ := service.Sign(session2, message)

	// Combine only OTHER signer's partial signature (our own is already in session)
	_, _ = service.CombineSig(session1, partialSig2)
	finalSig, _ := service.FinalSig(session1)

	aggregatedKey, _ := service.GetCombinedKey(ctx1)

	wrongMessage := sha256.Sum256([]byte("wrong message"))

	tests := []struct {
		name          string
		signature     *schnorr.Signature
		messageHash   [32]byte
		aggregatedKey *btcec.PublicKey
		wantValid     bool
	}{
		{
			name:          "valid signature",
			signature:     finalSig,
			messageHash:   message,
			aggregatedKey: aggregatedKey,
			wantValid:     true,
		},
		{
			name:          "wrong message",
			signature:     finalSig,
			messageHash:   wrongMessage,
			aggregatedKey: aggregatedKey,
			wantValid:     false,
		},
		{
			name:          "nil signature",
			signature:     nil,
			messageHash:   message,
			aggregatedKey: aggregatedKey,
			wantValid:     false,
		},
		{
			name:          "nil key",
			signature:     finalSig,
			messageHash:   message,
			aggregatedKey: nil,
			wantValid:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			valid := service.VerifySignature(tt.signature, tt.messageHash, tt.aggregatedKey)
			if valid != tt.wantValid {
				t.Errorf("VerifySignature() = %v, want %v", valid, tt.wantValid)
			}
		})
	}
}
