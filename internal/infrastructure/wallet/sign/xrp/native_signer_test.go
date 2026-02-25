package xrp

import (
	"context"
	"testing"

	dtoxrp "github.com/hiromaily/go-crypto-wallet/internal/application/dto/xrp"
	xrpcrypto "github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp"
)

func TestNativeSigner_SignTransaction_Ed25519(t *testing.T) {
	t.Parallel()

	// Create native signer with Ed25519
	signer := NewNativeSignerEd25519()

	// Generate a test key pair
	keyGen := xrpcrypto.NewKeyGenerator(xrpcrypto.AlgorithmEd25519)
	keyPair, err := keyGen.GenerateRandom()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	// Create test transaction
	txInput := &dtoxrp.TxInput{
		TransactionType:    "Payment",
		Account:            keyPair.ClassicAddress,
		Destination:        "rfkE1aSy9G8Upk4JssnwBxhEv5p4mn2KTy",
		Amount:             "1000000", // 1 XRP in drops
		Fee:                "12",
		Sequence:           1,
		LastLedgerSequence: 65953073,
	}

	// Sign the transaction
	ctx := context.Background()
	txID, txBlob, err := signer.SignTransaction(ctx, txInput, keyPair.Seed)
	if err != nil {
		t.Fatalf("SignTransaction failed: %v", err)
	}

	// Verify results
	if txID == "" {
		t.Error("txID should not be empty")
	}
	if txBlob == "" {
		t.Error("txBlob should not be empty")
	}
	if len(txID) != 64 {
		t.Errorf("txID should be 64 hex characters, got %d", len(txID))
	}

	t.Logf("Transaction ID: %s", txID)
	t.Logf("Transaction Blob length: %d", len(txBlob))
}

func TestNativeSigner_SignTransaction_Secp256k1(t *testing.T) {
	t.Parallel()

	// Create native signer with secp256k1
	signer := NewNativeSignerSecp256k1()

	// Generate a test key pair
	keyGen := xrpcrypto.NewKeyGenerator(xrpcrypto.AlgorithmSecp256k1)
	keyPair, err := keyGen.GenerateRandom()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	// Create test transaction
	txInput := &dtoxrp.TxInput{
		TransactionType:    "Payment",
		Account:            keyPair.ClassicAddress,
		Destination:        "rfkE1aSy9G8Upk4JssnwBxhEv5p4mn2KTy",
		Amount:             "1000000",
		Fee:                "12",
		Sequence:           1,
		LastLedgerSequence: 65953073,
	}

	// Sign the transaction
	ctx := context.Background()
	txID, txBlob, err := signer.SignTransaction(ctx, txInput, keyPair.Seed)
	if err != nil {
		t.Fatalf("SignTransaction failed: %v", err)
	}

	// Verify results
	if txID == "" {
		t.Error("txID should not be empty")
	}
	if txBlob == "" {
		t.Error("txBlob should not be empty")
	}

	t.Logf("Transaction ID: %s", txID)
}

func TestNativeSigner_SignTransactionWithDetails(t *testing.T) {
	t.Parallel()

	signer := NewNativeSignerEd25519()

	// Generate a test key pair
	keyGen := xrpcrypto.NewKeyGenerator(xrpcrypto.AlgorithmEd25519)
	keyPair, err := keyGen.GenerateRandom()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	txInput := &dtoxrp.TxInput{
		TransactionType:    "Payment",
		Account:            keyPair.ClassicAddress,
		Destination:        "rfkE1aSy9G8Upk4JssnwBxhEv5p4mn2KTy",
		Amount:             "1000000",
		Fee:                "12",
		Sequence:           1,
		LastLedgerSequence: 65953073,
	}

	ctx := context.Background()
	result, err := signer.SignTransactionWithDetails(ctx, txInput, keyPair.Seed)
	if err != nil {
		t.Fatalf("SignTransactionWithDetails failed: %v", err)
	}

	if result.TxID == "" {
		t.Error("TxID should not be empty")
	}
	if result.TxBlob == "" {
		t.Error("TxBlob should not be empty")
	}
	if result.SigningPubKey == "" {
		t.Error("SigningPubKey should not be empty")
	}
	if result.TxnSignature == "" {
		t.Error("TxnSignature should not be empty")
	}

	// Ed25519 signature should be 128 hex characters (64 bytes)
	if len(result.TxnSignature) != 128 {
		t.Errorf("Ed25519 signature should be 128 hex characters, got %d", len(result.TxnSignature))
	}
}

func TestNativeSigner_SignTransaction_NilInput(t *testing.T) {
	t.Parallel()

	signer := NewNativeSignerEd25519()
	ctx := context.Background()

	_, _, err := signer.SignTransaction(ctx, nil, "someseed")
	if err == nil {
		t.Error("Expected error for nil input")
	}
}

func TestNativeSigner_SignTransaction_EmptySecret(t *testing.T) {
	t.Parallel()

	signer := NewNativeSignerEd25519()
	ctx := context.Background()

	txInput := &dtoxrp.TxInput{
		TransactionType: "Payment",
		Account:         "rN7n3473SaZBCG4dFL83w7a1RXtXtbk2D9",
		Amount:          "1000000",
		Fee:             "12",
		Sequence:        1,
	}

	_, _, err := signer.SignTransaction(ctx, txInput, "")
	if err == nil {
		t.Error("Expected error for empty secret")
	}
}

func TestNativeSigner_VerifySignature(t *testing.T) {
	t.Parallel()

	signer := NewNativeSignerEd25519()

	// Generate a test key pair
	keyGen := xrpcrypto.NewKeyGenerator(xrpcrypto.AlgorithmEd25519)
	keyPair, err := keyGen.GenerateRandom()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	txInput := &dtoxrp.TxInput{
		TransactionType:    "Payment",
		Account:            keyPair.ClassicAddress,
		Destination:        "rfkE1aSy9G8Upk4JssnwBxhEv5p4mn2KTy",
		Amount:             "1000000",
		Fee:                "12",
		Sequence:           1,
		LastLedgerSequence: 65953073,
	}

	// Sign the transaction
	ctx := context.Background()
	result, err := signer.SignTransactionWithDetails(ctx, txInput, keyPair.Seed)
	if err != nil {
		t.Fatalf("SignTransactionWithDetails failed: %v", err)
	}

	// Create signed txInput for verification
	signedTxInput := &dtoxrp.TxInput{
		TransactionType:    "Payment",
		Account:            keyPair.ClassicAddress,
		Destination:        "rfkE1aSy9G8Upk4JssnwBxhEv5p4mn2KTy",
		Amount:             "1000000",
		Fee:                "12",
		Sequence:           1,
		LastLedgerSequence: 65953073,
		SigningPubKey:      result.SigningPubKey,
		TxnSignature:       result.TxnSignature,
	}

	// Verify the signature
	valid, err := signer.VerifySignature(signedTxInput)
	if err != nil {
		t.Fatalf("VerifySignature failed: %v", err)
	}
	if !valid {
		t.Error("Signature should be valid")
	}
}

func TestNativeSigner_VerifySignature_MissingFields(t *testing.T) {
	t.Parallel()

	signer := NewNativeSignerEd25519()

	// Test with nil input
	_, err := signer.VerifySignature(nil)
	if err == nil {
		t.Error("Expected error for nil input")
	}

	// Test with missing SigningPubKey
	txInput := &dtoxrp.TxInput{
		TransactionType: "Payment",
		TxnSignature:    "someSignature",
	}
	_, err = signer.VerifySignature(txInput)
	if err == nil {
		t.Error("Expected error for missing SigningPubKey")
	}

	// Test with missing TxnSignature
	txInput = &dtoxrp.TxInput{
		TransactionType: "Payment",
		SigningPubKey:   "someKey",
	}
	_, err = signer.VerifySignature(txInput)
	if err == nil {
		t.Error("Expected error for missing TxnSignature")
	}
}

func TestNativeSigner_Algorithm(t *testing.T) {
	t.Parallel()

	ed25519Signer := NewNativeSignerEd25519()
	if ed25519Signer.Algorithm() != xrpcrypto.AlgorithmEd25519 {
		t.Errorf("Expected AlgorithmEd25519, got %v", ed25519Signer.Algorithm())
	}

	secp256k1Signer := NewNativeSignerSecp256k1()
	if secp256k1Signer.Algorithm() != xrpcrypto.AlgorithmSecp256k1 {
		t.Errorf("Expected AlgorithmSecp256k1, got %v", secp256k1Signer.Algorithm())
	}
}

// Benchmark tests
func BenchmarkNativeSigner_SignTransaction_Ed25519(b *testing.B) {
	signer := NewNativeSignerEd25519()

	keyGen := xrpcrypto.NewKeyGenerator(xrpcrypto.AlgorithmEd25519)
	keyPair, err := keyGen.GenerateRandom()
	if err != nil {
		b.Fatalf("Failed to generate key pair: %v", err)
	}

	txInput := &dtoxrp.TxInput{
		TransactionType:    "Payment",
		Account:            keyPair.ClassicAddress,
		Destination:        "rfkE1aSy9G8Upk4JssnwBxhEv5p4mn2KTy",
		Amount:             "1000000",
		Fee:                "12",
		Sequence:           1,
		LastLedgerSequence: 65953073,
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := signer.SignTransaction(ctx, txInput, keyPair.Seed)
		if err != nil {
			b.Fatalf("SignTransaction failed: %v", err)
		}
	}
}

func BenchmarkNativeSigner_SignTransaction_Secp256k1(b *testing.B) {
	signer := NewNativeSignerSecp256k1()

	keyGen := xrpcrypto.NewKeyGenerator(xrpcrypto.AlgorithmSecp256k1)
	keyPair, err := keyGen.GenerateRandom()
	if err != nil {
		b.Fatalf("Failed to generate key pair: %v", err)
	}

	txInput := &dtoxrp.TxInput{
		TransactionType:    "Payment",
		Account:            keyPair.ClassicAddress,
		Destination:        "rfkE1aSy9G8Upk4JssnwBxhEv5p4mn2KTy",
		Amount:             "1000000",
		Fee:                "12",
		Sequence:           1,
		LastLedgerSequence: 65953073,
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := signer.SignTransaction(ctx, txInput, keyPair.Seed)
		if err != nil {
			b.Fatalf("SignTransaction failed: %v", err)
		}
	}
}
