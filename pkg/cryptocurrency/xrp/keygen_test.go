package xrp

import (
	"encoding/hex"
	"strings"
	"testing"
)

func TestKeyAlgorithm_String(t *testing.T) {
	tests := []struct {
		name      string
		algorithm KeyAlgorithm
		want      string
	}{
		{
			name:      "secp256k1",
			algorithm: AlgorithmSecp256k1,
			want:      "secp256k1",
		},
		{
			name:      "ed25519",
			algorithm: AlgorithmEd25519,
			want:      "ed25519",
		},
		{
			name:      "unknown",
			algorithm: KeyAlgorithm(99),
			want:      "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.algorithm.String(); got != tt.want {
				t.Errorf("KeyAlgorithm.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseKeyAlgorithm(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    KeyAlgorithm
		wantErr bool
	}{
		{
			name:    "secp256k1",
			input:   "secp256k1",
			want:    AlgorithmSecp256k1,
			wantErr: false,
		},
		{
			name:    "ed25519",
			input:   "ed25519",
			want:    AlgorithmEd25519,
			wantErr: false,
		},
		{
			name:    "unknown algorithm",
			input:   "unknown",
			want:    KeyAlgorithm(-1),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseKeyAlgorithm(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseKeyAlgorithm() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseKeyAlgorithm() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeyGenerator_GenerateRandom_Secp256k1(t *testing.T) {
	gen := NewKeyGenerator(AlgorithmSecp256k1)

	keyPair, err := gen.GenerateRandom()
	if err != nil {
		t.Fatalf("GenerateRandom() error = %v", err)
	}

	// Verify the key pair
	validateKeyPair(t, keyPair, AlgorithmSecp256k1)
}

func TestKeyGenerator_GenerateRandom_Ed25519(t *testing.T) {
	gen := NewKeyGenerator(AlgorithmEd25519)

	keyPair, err := gen.GenerateRandom()
	if err != nil {
		t.Fatalf("GenerateRandom() error = %v", err)
	}

	// Verify the key pair
	validateKeyPair(t, keyPair, AlgorithmEd25519)
}

func TestKeyGenerator_GenerateFromEntropy_Secp256k1(t *testing.T) {
	gen := NewKeyGenerator(AlgorithmSecp256k1)

	// Use fixed entropy for reproducible results
	entropy, _ := hex.DecodeString("0123456789abcdef0123456789abcdef")

	keyPair1, err := gen.GenerateFromEntropy(entropy)
	if err != nil {
		t.Fatalf("GenerateFromEntropy() error = %v", err)
	}

	// Generate again with same entropy - should produce same result
	keyPair2, err := gen.GenerateFromEntropy(entropy)
	if err != nil {
		t.Fatalf("GenerateFromEntropy() second call error = %v", err)
	}

	// Verify deterministic generation
	if keyPair1.Seed != keyPair2.Seed {
		t.Errorf("Seeds don't match: %s != %s", keyPair1.Seed, keyPair2.Seed)
	}
	if keyPair1.ClassicAddress != keyPair2.ClassicAddress {
		t.Errorf("Addresses don't match: %s != %s", keyPair1.ClassicAddress, keyPair2.ClassicAddress)
	}
	if keyPair1.PublicKey != keyPair2.PublicKey {
		t.Errorf("PublicKeys don't match: %s != %s", keyPair1.PublicKey, keyPair2.PublicKey)
	}

	validateKeyPair(t, keyPair1, AlgorithmSecp256k1)
}

func TestKeyGenerator_GenerateFromEntropy_Ed25519(t *testing.T) {
	gen := NewKeyGenerator(AlgorithmEd25519)

	// Use fixed entropy for reproducible results
	entropy, _ := hex.DecodeString("0123456789abcdef0123456789abcdef")

	keyPair1, err := gen.GenerateFromEntropy(entropy)
	if err != nil {
		t.Fatalf("GenerateFromEntropy() error = %v", err)
	}

	// Generate again with same entropy - should produce same result
	keyPair2, err := gen.GenerateFromEntropy(entropy)
	if err != nil {
		t.Fatalf("GenerateFromEntropy() second call error = %v", err)
	}

	// Verify deterministic generation
	if keyPair1.Seed != keyPair2.Seed {
		t.Errorf("Seeds don't match: %s != %s", keyPair1.Seed, keyPair2.Seed)
	}
	if keyPair1.ClassicAddress != keyPair2.ClassicAddress {
		t.Errorf("Addresses don't match: %s != %s", keyPair1.ClassicAddress, keyPair2.ClassicAddress)
	}
	if keyPair1.PublicKey != keyPair2.PublicKey {
		t.Errorf("PublicKeys don't match: %s != %s", keyPair1.PublicKey, keyPair2.PublicKey)
	}

	validateKeyPair(t, keyPair1, AlgorithmEd25519)
}

func TestKeyGenerator_GenerateFromEntropy_ShortEntropy(t *testing.T) {
	gen := NewKeyGenerator(AlgorithmSecp256k1)

	// Use entropy that's too short
	shortEntropy := []byte{0x01, 0x02, 0x03}

	_, err := gen.GenerateFromEntropy(shortEntropy)
	if err == nil {
		t.Error("GenerateFromEntropy() should fail with short entropy")
	}
}

func TestKeyGenerator_GenerateFromPassphrase_Secp256k1(t *testing.T) {
	gen := NewKeyGenerator(AlgorithmSecp256k1)

	passphrase := "test passphrase for XRP key generation"

	keyPair1, err := gen.GenerateFromPassphrase(passphrase)
	if err != nil {
		t.Fatalf("GenerateFromPassphrase() error = %v", err)
	}

	// Generate again with same passphrase - should produce same result
	keyPair2, err := gen.GenerateFromPassphrase(passphrase)
	if err != nil {
		t.Fatalf("GenerateFromPassphrase() second call error = %v", err)
	}

	// Verify deterministic generation
	if keyPair1.Seed != keyPair2.Seed {
		t.Errorf("Seeds don't match: %s != %s", keyPair1.Seed, keyPair2.Seed)
	}
	if keyPair1.ClassicAddress != keyPair2.ClassicAddress {
		t.Errorf("Addresses don't match: %s != %s", keyPair1.ClassicAddress, keyPair2.ClassicAddress)
	}

	validateKeyPair(t, keyPair1, AlgorithmSecp256k1)
}

func TestKeyGenerator_GenerateFromPassphrase_Ed25519(t *testing.T) {
	gen := NewKeyGenerator(AlgorithmEd25519)

	passphrase := "test passphrase for XRP key generation"

	keyPair1, err := gen.GenerateFromPassphrase(passphrase)
	if err != nil {
		t.Fatalf("GenerateFromPassphrase() error = %v", err)
	}

	// Generate again with same passphrase - should produce same result
	keyPair2, err := gen.GenerateFromPassphrase(passphrase)
	if err != nil {
		t.Fatalf("GenerateFromPassphrase() second call error = %v", err)
	}

	// Verify deterministic generation
	if keyPair1.Seed != keyPair2.Seed {
		t.Errorf("Seeds don't match: %s != %s", keyPair1.Seed, keyPair2.Seed)
	}
	if keyPair1.ClassicAddress != keyPair2.ClassicAddress {
		t.Errorf("Addresses don't match: %s != %s", keyPair1.ClassicAddress, keyPair2.ClassicAddress)
	}

	validateKeyPair(t, keyPair1, AlgorithmEd25519)
}

func TestKeyGenerator_GenerateFromPassphrase_Empty(t *testing.T) {
	gen := NewKeyGenerator(AlgorithmSecp256k1)

	_, err := gen.GenerateFromPassphrase("")
	if err == nil {
		t.Error("GenerateFromPassphrase() should fail with empty passphrase")
	}
}

func TestKeyGenerator_DeriveKeyFromHDKey_Secp256k1(t *testing.T) {
	gen := NewKeyGenerator(AlgorithmSecp256k1)

	// Simulate a 32-byte HD private key
	hdPrivateKey, _ := hex.DecodeString(
		"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
	)

	keyPair1, err := gen.DeriveKeyFromHDKey(hdPrivateKey)
	if err != nil {
		t.Fatalf("DeriveKeyFromHDKey() error = %v", err)
	}

	// Generate again with same HD key - should produce same result
	keyPair2, err := gen.DeriveKeyFromHDKey(hdPrivateKey)
	if err != nil {
		t.Fatalf("DeriveKeyFromHDKey() second call error = %v", err)
	}

	// Verify deterministic generation
	if keyPair1.Seed != keyPair2.Seed {
		t.Errorf("Seeds don't match: %s != %s", keyPair1.Seed, keyPair2.Seed)
	}
	if keyPair1.ClassicAddress != keyPair2.ClassicAddress {
		t.Errorf("Addresses don't match: %s != %s", keyPair1.ClassicAddress, keyPair2.ClassicAddress)
	}

	validateKeyPair(t, keyPair1, AlgorithmSecp256k1)
}

func TestKeyGenerator_DeriveKeyFromHDKey_Ed25519(t *testing.T) {
	gen := NewKeyGenerator(AlgorithmEd25519)

	// Simulate a 32-byte HD private key
	hdPrivateKey, _ := hex.DecodeString(
		"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
	)

	keyPair1, err := gen.DeriveKeyFromHDKey(hdPrivateKey)
	if err != nil {
		t.Fatalf("DeriveKeyFromHDKey() error = %v", err)
	}

	// Generate again with same HD key - should produce same result
	keyPair2, err := gen.DeriveKeyFromHDKey(hdPrivateKey)
	if err != nil {
		t.Fatalf("DeriveKeyFromHDKey() second call error = %v", err)
	}

	// Verify deterministic generation
	if keyPair1.Seed != keyPair2.Seed {
		t.Errorf("Seeds don't match: %s != %s", keyPair1.Seed, keyPair2.Seed)
	}
	if keyPair1.ClassicAddress != keyPair2.ClassicAddress {
		t.Errorf("Addresses don't match: %s != %s", keyPair1.ClassicAddress, keyPair2.ClassicAddress)
	}

	validateKeyPair(t, keyPair1, AlgorithmEd25519)
}

func TestKeyGenerator_DeriveKeyFromHDKey_ShortKey(t *testing.T) {
	gen := NewKeyGenerator(AlgorithmSecp256k1)

	// Use HD key that's too short
	shortKey := []byte{0x01, 0x02, 0x03}

	_, err := gen.DeriveKeyFromHDKey(shortKey)
	if err == nil {
		t.Error("DeriveKeyFromHDKey() should fail with short key")
	}
}

func TestDifferentAlgorithms_ProduceDifferentAddresses(t *testing.T) {
	// Use same entropy for both algorithms
	entropy, _ := hex.DecodeString("0123456789abcdef0123456789abcdef")

	secpGen := NewKeyGenerator(AlgorithmSecp256k1)
	edGen := NewKeyGenerator(AlgorithmEd25519)

	secpKeyPair, err := secpGen.GenerateFromEntropy(entropy)
	if err != nil {
		t.Fatalf("secp256k1 GenerateFromEntropy() error = %v", err)
	}

	edKeyPair, err := edGen.GenerateFromEntropy(entropy)
	if err != nil {
		t.Fatalf("ed25519 GenerateFromEntropy() error = %v", err)
	}

	// Seeds should be identical since they're both derived from the same entropy
	// using NewFamilySeed(seedBytes) which doesn't depend on the algorithm
	if secpKeyPair.Seed != edKeyPair.Seed {
		t.Errorf("Seeds should match for same entropy: secp=%s, ed=%s",
			secpKeyPair.Seed, edKeyPair.Seed)
	}

	// Addresses should be different due to different key derivation
	if secpKeyPair.ClassicAddress == edKeyPair.ClassicAddress {
		t.Error("Different algorithms should produce different addresses from same entropy")
	}

	// Verify Ed25519 public key has 0xED prefix
	if !strings.HasPrefix(secpKeyPair.PublicKeyHex, "02") &&
		!strings.HasPrefix(secpKeyPair.PublicKeyHex, "03") {
		t.Errorf("secp256k1 public key should start with 02 or 03, got: %s",
			secpKeyPair.PublicKeyHex[:4])
	}

	pubKeyBytes, _ := hex.DecodeString(edKeyPair.PublicKeyHex)
	if len(pubKeyBytes) > 0 && pubKeyBytes[0] != 0xED {
		t.Errorf("ed25519 public key should start with 0xED, got: 0x%02X", pubKeyBytes[0])
	}
}

// validateKeyPair validates that a generated key pair has the expected format.
func validateKeyPair(t *testing.T, kp *XRPKeyPair, expectedAlg KeyAlgorithm) {
	t.Helper()

	// Check algorithm
	if kp.Algorithm != expectedAlg {
		t.Errorf("Algorithm = %v, want %v", kp.Algorithm, expectedAlg)
	}

	// Check seed format (should start with 's')
	if !strings.HasPrefix(kp.Seed, "s") {
		t.Errorf("Seed should start with 's', got: %s", kp.Seed[:min(5, len(kp.Seed))])
	}

	// Check seed hex is valid hex
	if _, err := hex.DecodeString(kp.SeedHex); err != nil {
		t.Errorf("SeedHex is not valid hex: %v", err)
	}

	// Check seed hex length (should be 32 chars = 16 bytes)
	if len(kp.SeedHex) != 32 {
		t.Errorf("SeedHex length = %d, want 32", len(kp.SeedHex))
	}

	// Check classic address format (should start with 'r')
	if !strings.HasPrefix(kp.ClassicAddress, "r") {
		t.Errorf("ClassicAddress should start with 'r', got: %s",
			kp.ClassicAddress[:min(5, len(kp.ClassicAddress))])
	}

	// Check public key format (should start with 'a' for account public key)
	if !strings.HasPrefix(kp.PublicKey, "a") {
		t.Errorf("PublicKey should start with 'a', got: %s",
			kp.PublicKey[:min(5, len(kp.PublicKey))])
	}

	// Check public key hex is valid hex
	if _, err := hex.DecodeString(kp.PublicKeyHex); err != nil {
		t.Errorf("PublicKeyHex is not valid hex: %v", err)
	}
}
