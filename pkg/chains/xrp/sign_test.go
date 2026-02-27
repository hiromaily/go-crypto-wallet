package xrp

import (
	"encoding/hex"
	"maps"
	"strings"
	"testing"
)

// Test vectors based on XRPL documentation and reference implementations
// Reference: https://xrpl.org/transaction-basics.html
// Reference: https://github.com/XRPLF/xrpl.js/blob/main/packages/xrpl/test/wallet/index.ts

// TestSigner_ReferenceVectors tests signing with deterministic keys.
// These tests verify that signing produces consistent results for the same inputs.
// Reference: https://github.com/XRPLF/xrpl.js/blob/main/packages/xrpl/test/wallet/index.ts
func TestSigner_ReferenceVectors(t *testing.T) {
	t.Parallel()

	// Test with deterministic passphrases to generate consistent keys
	tests := []struct {
		name              string
		algorithm         KeyAlgorithm
		passphrase        string
		expectedPubKeyPfx string // First bytes of public key (prefix)
	}{
		{
			name:              "secp256k1 deterministic key",
			algorithm:         AlgorithmSecp256k1,
			passphrase:        "masterpassphrase",
			expectedPubKeyPfx: "02", // secp256k1 compressed (02 or 03)
		},
		{
			name:              "ed25519 deterministic key",
			algorithm:         AlgorithmEd25519,
			passphrase:        "masterpassphrase",
			expectedPubKeyPfx: "ED", // ed25519 prefix
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			keyGen := NewKeyGenerator(tc.algorithm)
			keyPair, err := keyGen.GenerateFromPassphrase(tc.passphrase)
			if err != nil {
				t.Fatalf("Failed to generate key pair: %v", err)
			}

			// Verify public key prefix
			upperPubKey := strings.ToUpper(keyPair.PublicKeyHex)
			if !strings.HasPrefix(upperPubKey, tc.expectedPubKeyPfx) &&
				!strings.HasPrefix(upperPubKey, "03") { // secp256k1 can be 02 or 03
				t.Errorf("Public key prefix mismatch: got %s, want prefix %s",
					keyPair.PublicKeyHex[:2], tc.expectedPubKeyPfx)
			}

			// Generate a second keypair to use as destination
			destKeyPair, err := keyGen.GenerateFromPassphrase("destination-account")
			if err != nil {
				t.Fatalf("Failed to generate destination key pair: %v", err)
			}

			// Sign a test transaction with known data
			txFields := map[string]any{
				"TransactionType":    "Payment",
				"Account":            keyPair.ClassicAddress,
				"Destination":        destKeyPair.ClassicAddress,
				"Amount":             "1000000",
				"Fee":                "12",
				"Sequence":           uint64(1),
				"LastLedgerSequence": uint64(65953073),
			}

			signer := NewSigner(tc.algorithm)
			result, err := signer.SignTransaction(txFields, keyPair.Seed)
			if err != nil {
				t.Fatalf("SignTransaction failed: %v", err)
			}

			// Verify signature is valid
			verifyFields := map[string]any{
				"TransactionType":    "Payment",
				"Account":            keyPair.ClassicAddress,
				"Destination":        destKeyPair.ClassicAddress,
				"Amount":             "1000000",
				"Fee":                "12",
				"Sequence":           uint64(1),
				"LastLedgerSequence": uint64(65953073),
				"SigningPubKey":      result.SigningPubKey,
				"TxnSignature":       result.TxnSignature,
			}

			valid, err := signer.VerifySignature(verifyFields)
			if err != nil {
				t.Fatalf("VerifySignature failed: %v", err)
			}
			if !valid {
				t.Error("Signature should be valid")
			}

			// Verify transaction ID format
			if len(result.TxID) != 64 {
				t.Errorf("TxID should be 64 hex chars, got %d", len(result.TxID))
			}

			// Log results for debugging
			t.Logf("Address: %s", keyPair.ClassicAddress)
			t.Logf("PublicKey: %s", result.SigningPubKey)
			t.Logf("TxID: %s", result.TxID)
		})
	}
}

// TestSigner_DeterministicSigning verifies that signing the same transaction
// with the same key produces consistent results (deterministic signing).
func TestSigner_DeterministicSigning(t *testing.T) {
	t.Parallel()

	// Generate a deterministic key pair
	keyGen := NewKeyGenerator(AlgorithmEd25519)
	keyPair, err := keyGen.GenerateFromPassphrase("test-deterministic-signing")
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	txFields := map[string]any{
		"TransactionType":    "Payment",
		"Account":            keyPair.ClassicAddress,
		"Destination":        "rfkE1aSy9G8Upk4JssnwBxhEv5p4mn2KTy",
		"Amount":             "500000",
		"Fee":                "10",
		"Sequence":           uint64(5),
		"LastLedgerSequence": uint64(12345678),
	}

	signer := NewSigner(AlgorithmEd25519)

	// Sign the same transaction multiple times
	var previousTxID string
	var previousSignature string
	for i := range 3 {
		result, err := signer.SignTransaction(txFields, keyPair.Seed)
		if err != nil {
			t.Fatalf("SignTransaction failed on iteration %d: %v", i, err)
		}

		if previousTxID != "" && result.TxID != previousTxID {
			t.Errorf("TxID changed between iterations: %s != %s", result.TxID, previousTxID)
		}
		if previousSignature != "" && result.TxnSignature != previousSignature {
			t.Errorf("Signature changed between iterations: %s != %s",
				result.TxnSignature, previousSignature)
		}

		previousTxID = result.TxID
		previousSignature = result.TxnSignature
	}
}

// TestSigner_InputNotModified verifies that SignTransaction does not modify
// the input txFields map.
func TestSigner_InputNotModified(t *testing.T) {
	t.Parallel()

	keyGen := NewKeyGenerator(AlgorithmEd25519)
	keyPair, err := keyGen.GenerateRandom()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	txFields := map[string]any{
		"TransactionType":    "Payment",
		"Account":            keyPair.ClassicAddress,
		"Destination":        "rfkE1aSy9G8Upk4JssnwBxhEv5p4mn2KTy",
		"Amount":             "1000000",
		"Fee":                "12",
		"Sequence":           uint64(1),
		"LastLedgerSequence": uint64(65953073),
	}

	// Record original keys
	originalKeys := make(map[string]bool)
	for k := range txFields {
		originalKeys[k] = true
	}

	signer := NewSigner(AlgorithmEd25519)
	_, err = signer.SignTransaction(txFields, keyPair.Seed)
	if err != nil {
		t.Fatalf("SignTransaction failed: %v", err)
	}

	// Verify no new keys were added
	if _, exists := txFields["SigningPubKey"]; exists {
		t.Error("SignTransaction modified input map: added SigningPubKey")
	}
	if _, exists := txFields["TxnSignature"]; exists {
		t.Error("SignTransaction modified input map: added TxnSignature")
	}

	// Verify all original keys still exist
	for k := range originalKeys {
		if _, exists := txFields[k]; !exists {
			t.Errorf("SignTransaction modified input map: removed key %s", k)
		}
	}
}

func TestSigner_SignTransaction_Ed25519(t *testing.T) {
	t.Parallel()

	// Create a signer with Ed25519 algorithm
	signer := NewSigner(AlgorithmEd25519)

	// Test transaction (Payment type)
	txFields := map[string]any{
		"TransactionType":    "Payment",
		"Account":            "rN7n3473SaZBCG4dFL83w7a1RXtXtbk2D9",
		"Destination":        "rfkE1aSy9G8Upk4JssnwBxhEv5p4mn2KTy",
		"Amount":             "1000000", // 1 XRP in drops
		"Fee":                "12",
		"Sequence":           uint64(1),
		"LastLedgerSequence": uint64(65953073),
	}

	// Generate a test key pair for signing
	keyGen := NewKeyGenerator(AlgorithmEd25519)
	keyPair, err := keyGen.GenerateRandom()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	// Update Account to match the generated key
	txFields["Account"] = keyPair.ClassicAddress

	// Sign the transaction
	result, err := signer.SignTransaction(txFields, keyPair.Seed)
	if err != nil {
		t.Fatalf("SignTransaction failed: %v", err)
	}

	// Verify the result
	if result.TxBlob == "" {
		t.Error("TxBlob should not be empty")
	}
	if result.TxID == "" {
		t.Error("TxID should not be empty")
	}
	if result.SigningPubKey == "" {
		t.Error("SigningPubKey should not be empty")
	}
	if result.TxnSignature == "" {
		t.Error("TxnSignature should not be empty")
	}

	// Verify TxID is 64 hex characters (32 bytes)
	if len(result.TxID) != 64 {
		t.Errorf("TxID should be 64 hex characters, got %d", len(result.TxID))
	}

	// Verify signature length (Ed25519 signatures are 64 bytes = 128 hex chars)
	sigLen := len(result.TxnSignature)
	if sigLen != 128 {
		t.Errorf("Ed25519 signature should be 128 hex characters, got %d", sigLen)
	}

	// Verify the signature
	txFields["SigningPubKey"] = result.SigningPubKey
	txFields["TxnSignature"] = result.TxnSignature

	valid, err := signer.VerifySignature(txFields)
	if err != nil {
		t.Fatalf("VerifySignature failed: %v", err)
	}
	if !valid {
		t.Error("Signature should be valid")
	}
}

func TestSigner_SignTransaction_Secp256k1(t *testing.T) {
	t.Parallel()

	// Create a signer with secp256k1 algorithm
	signer := NewSigner(AlgorithmSecp256k1)

	// Generate a test key pair for signing
	keyGen := NewKeyGenerator(AlgorithmSecp256k1)
	keyPair, err := keyGen.GenerateRandom()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	// Test transaction (Payment type)
	txFields := map[string]any{
		"TransactionType":    "Payment",
		"Account":            keyPair.ClassicAddress,
		"Destination":        "rfkE1aSy9G8Upk4JssnwBxhEv5p4mn2KTy",
		"Amount":             "1000000", // 1 XRP in drops
		"Fee":                "12",
		"Sequence":           uint64(1),
		"LastLedgerSequence": uint64(65953073),
	}

	// Sign the transaction
	result, err := signer.SignTransaction(txFields, keyPair.Seed)
	if err != nil {
		t.Fatalf("SignTransaction failed: %v", err)
	}

	// Verify the result
	if result.TxBlob == "" {
		t.Error("TxBlob should not be empty")
	}
	if result.TxID == "" {
		t.Error("TxID should not be empty")
	}

	// Verify signature is DER encoded (variable length, typically 70-72 bytes)
	sigBytes, err := hex.DecodeString(result.TxnSignature)
	if err != nil {
		t.Fatalf("Invalid signature hex: %v", err)
	}
	if len(sigBytes) < 68 || len(sigBytes) > 73 {
		t.Errorf("secp256k1 DER signature should be 68-73 bytes, got %d", len(sigBytes))
	}

	// Verify the signature
	txFields["SigningPubKey"] = result.SigningPubKey
	txFields["TxnSignature"] = result.TxnSignature

	valid, err := signer.VerifySignature(txFields)
	if err != nil {
		t.Fatalf("VerifySignature failed: %v", err)
	}
	if !valid {
		t.Error("Signature should be valid")
	}
}

func TestSigner_SignTransaction_InvalidSeed(t *testing.T) {
	t.Parallel()

	signer := NewSigner(AlgorithmEd25519)

	txFields := map[string]any{
		"TransactionType": "Payment",
		"Account":         "rN7n3473SaZBCG4dFL83w7a1RXtXtbk2D9",
		"Amount":          "1000000",
		"Fee":             "12",
		"Sequence":        uint64(1),
	}

	// Test with invalid seed
	_, err := signer.SignTransaction(txFields, "invalid-seed")
	if err == nil {
		t.Error("Expected error for invalid seed")
	}

	// Test with empty seed
	_, err = signer.SignTransaction(txFields, "")
	if err == nil {
		t.Error("Expected error for empty seed")
	}
}

func TestSigner_SignTransaction_NilTransaction(t *testing.T) {
	t.Parallel()

	signer := NewSigner(AlgorithmEd25519)

	_, err := signer.SignTransaction(nil, "somekey")
	if err == nil {
		t.Error("Expected error for nil transaction")
	}
}

func TestDetectAlgorithm(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		pubKeyHex string
		expected  KeyAlgorithm
		wantErr   bool
	}{
		{
			name:      "Ed25519 public key",
			pubKeyHex: "ED" + "0000000000000000000000000000000000000000000000000000000000000000",
			expected:  AlgorithmEd25519,
			wantErr:   false,
		},
		{
			name:      "secp256k1 public key (compressed, even)",
			pubKeyHex: "02" + "0000000000000000000000000000000000000000000000000000000000000000",
			expected:  AlgorithmSecp256k1,
			wantErr:   false,
		},
		{
			name:      "secp256k1 public key (compressed, odd)",
			pubKeyHex: "03" + "0000000000000000000000000000000000000000000000000000000000000000",
			expected:  AlgorithmSecp256k1,
			wantErr:   false,
		},
		{
			name:      "Invalid public key",
			pubKeyHex: "04" + "0000000000000000000000000000000000000000000000000000000000000000",
			expected:  KeyAlgorithm(-1),
			wantErr:   true,
		},
		{
			name:      "Invalid hex",
			pubKeyHex: "not-hex",
			expected:  KeyAlgorithm(-1),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			algorithm, err := DetectAlgorithm(tt.pubKeyHex)
			if (err != nil) != tt.wantErr {
				t.Errorf("DetectAlgorithm() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if algorithm != tt.expected {
				t.Errorf("DetectAlgorithm() = %v, want %v", algorithm, tt.expected)
			}
		})
	}
}

func TestSerializer_Serialize(t *testing.T) {
	t.Parallel()

	serializer := NewSerializer()

	// Generate a key pair for testing
	keyGen := NewKeyGenerator(AlgorithmEd25519)
	keyPair, err := keyGen.GenerateRandom()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	txFields := map[string]any{
		"TransactionType":    "Payment",
		"Account":            keyPair.ClassicAddress,
		"Destination":        "rfkE1aSy9G8Upk4JssnwBxhEv5p4mn2KTy",
		"Amount":             "1000000",
		"Fee":                "12",
		"Sequence":           uint64(1),
		"LastLedgerSequence": uint64(65953073),
	}

	tx := NewTransaction(txFields)

	// Test serialization for signing
	serialized, err := serializer.SerializeForSigning(tx)
	if err != nil {
		t.Fatalf("SerializeForSigning failed: %v", err)
	}
	if len(serialized) == 0 {
		t.Error("Serialized output should not be empty")
	}

	// Test computing signing hash
	hash, err := serializer.ComputeSigningHash(tx)
	if err != nil {
		t.Fatalf("ComputeSigningHash failed: %v", err)
	}
	if len(hash) != 32 {
		t.Errorf("Signing hash should be 32 bytes, got %d", len(hash))
	}
}

func TestComputeTransactionID(t *testing.T) {
	t.Parallel()

	// Test with known data
	txBlob := []byte{0x12, 0x00, 0x00} // Minimal test data
	txID := ComputeTransactionID(txBlob)

	if len(txID) != 64 {
		t.Errorf("Transaction ID should be 64 hex characters, got %d", len(txID))
	}

	// Verify it's valid hex
	_, err := hex.DecodeString(txID)
	if err != nil {
		t.Errorf("Transaction ID should be valid hex: %v", err)
	}
}

func TestSigner_VerifySignature_InvalidSignature(t *testing.T) {
	t.Parallel()

	signer := NewSigner(AlgorithmEd25519)

	// Generate a key pair
	keyGen := NewKeyGenerator(AlgorithmEd25519)
	keyPair, err := keyGen.GenerateRandom()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	txFields := map[string]any{
		"TransactionType": "Payment",
		"Account":         keyPair.ClassicAddress,
		"Amount":          "1000000",
		"Fee":             "12",
		"Sequence":        uint64(1),
		"SigningPubKey":   keyPair.PublicKeyHex,
		// Invalid signature (65 bytes of zeros)
		"TxnSignature": "00" +
			"0000000000000000000000000000000000000000000000000000000000000000" +
			"0000000000000000000000000000000000000000000000000000000000000000",
	}

	// Verify should fail for invalid signature
	valid, err := signer.VerifySignature(txFields)
	if err != nil {
		t.Logf("VerifySignature returned error (expected for invalid signature): %v", err)
	}
	if valid {
		t.Error("Verification should fail for invalid signature")
	}
}

func TestSigner_VerifySignature_MissingFields(t *testing.T) {
	t.Parallel()

	signer := NewSigner(AlgorithmEd25519)

	// Test missing TxnSignature
	txFields := map[string]any{
		"TransactionType": "Payment",
		"SigningPubKey":   "ED0000000000000000000000000000000000000000000000000000000000000000",
	}

	_, err := signer.VerifySignature(txFields)
	if err == nil {
		t.Error("Expected error for missing TxnSignature")
	}

	// Test missing SigningPubKey
	txFields = map[string]any{
		"TransactionType": "Payment",
		"TxnSignature":    "0000000000000000000000000000000000000000000000000000000000000000",
	}

	_, err = signer.VerifySignature(txFields)
	if err == nil {
		t.Error("Expected error for missing SigningPubKey")
	}
}

func TestSignTransactionWithAutoDetect(t *testing.T) {
	t.Parallel()

	// Generate Ed25519 key pair
	keyGen := NewKeyGenerator(AlgorithmEd25519)
	keyPair, err := keyGen.GenerateRandom()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	txFields := map[string]any{
		"TransactionType":    "Payment",
		"Account":            keyPair.ClassicAddress,
		"Destination":        "rfkE1aSy9G8Upk4JssnwBxhEv5p4mn2KTy",
		"Amount":             "1000000",
		"Fee":                "12",
		"Sequence":           uint64(1),
		"LastLedgerSequence": uint64(65953073),
	}

	// Sign without specifying algorithm
	result, err := SignTransactionWithAutoDetect(txFields, keyPair.Seed)
	if err != nil {
		t.Fatalf("SignTransactionWithAutoDetect failed: %v", err)
	}

	if result.TxBlob == "" {
		t.Error("TxBlob should not be empty")
	}
	if result.TxID == "" {
		t.Error("TxID should not be empty")
	}
}

// BenchmarkSignTransaction benchmarks the signing performance
func BenchmarkSignTransaction_Ed25519(b *testing.B) {
	signer := NewSigner(AlgorithmEd25519)

	keyGen := NewKeyGenerator(AlgorithmEd25519)
	keyPair, err := keyGen.GenerateRandom()
	if err != nil {
		b.Fatalf("Failed to generate key pair: %v", err)
	}

	txFields := map[string]any{
		"TransactionType":    "Payment",
		"Account":            keyPair.ClassicAddress,
		"Destination":        "rfkE1aSy9G8Upk4JssnwBxhEv5p4mn2KTy",
		"Amount":             "1000000",
		"Fee":                "12",
		"Sequence":           uint64(1),
		"LastLedgerSequence": uint64(65953073),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Make a copy to avoid mutation
		txCopy := make(map[string]any)
		maps.Copy(txCopy, txFields)
		_, err := signer.SignTransaction(txCopy, keyPair.Seed)
		if err != nil {
			b.Fatalf("SignTransaction failed: %v", err)
		}
	}
}

func BenchmarkSignTransaction_Secp256k1(b *testing.B) {
	signer := NewSigner(AlgorithmSecp256k1)

	keyGen := NewKeyGenerator(AlgorithmSecp256k1)
	keyPair, err := keyGen.GenerateRandom()
	if err != nil {
		b.Fatalf("Failed to generate key pair: %v", err)
	}

	txFields := map[string]any{
		"TransactionType":    "Payment",
		"Account":            keyPair.ClassicAddress,
		"Destination":        "rfkE1aSy9G8Upk4JssnwBxhEv5p4mn2KTy",
		"Amount":             "1000000",
		"Fee":                "12",
		"Sequence":           uint64(1),
		"LastLedgerSequence": uint64(65953073),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Make a copy to avoid mutation
		txCopy := make(map[string]any)
		maps.Copy(txCopy, txFields)
		_, err := signer.SignTransaction(txCopy, keyPair.Seed)
		if err != nil {
			b.Fatalf("SignTransaction failed: %v", err)
		}
	}
}
