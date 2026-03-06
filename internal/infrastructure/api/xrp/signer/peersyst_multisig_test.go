package signer_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	dtoxrp "github.com/hiromaily/go-crypto-wallet/internal/application/dto/xrp"
	xrpsigner "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/xrp/signer"
)

// TestPeersystSigner_MultiSignature_EmptySigningPubKey tests that transactions
// with empty SigningPubKey are detected as multi-signature transactions.
func TestPeersystSigner_MultiSignature_EmptySigningPubKey(t *testing.T) {
	t.Parallel()

	signer := xrpsigner.NewPeersystSigner()
	ctx := context.Background()
	testSeed := "sEdTM1uX8pu2do5XvTnutH6HsouMaM2"

	// Create transaction with empty SigningPubKey (multi-sig indicator)
	txInput := &dtoxrp.TxInput{
		TransactionType:    "Payment",
		Account:            "rG31cLyErnqeVj2eomEjBZtq7PYaupGYzL",
		Destination:        "rPEPPER7kfTD9w2To4CQk6UCfuHM9c6GDY",
		Amount:             "1000000",
		Fee:                "36",
		Sequence:           1,
		SigningPubKey:      "", // Empty string indicates multi-sig
		LastLedgerSequence: 10000000,
	}

	// Execute
	txHash, signedBlob, err := signer.SignTransactionNative(ctx, txInput, testSeed, true, nil)

	// Assert - multi-sig signing should work
	require.NoError(t, err, "multi-sig signing should succeed")
	assert.NotEmpty(t, txHash, "transaction hash should not be empty")
	assert.NotEmpty(t, signedBlob, "signed blob should not be empty")

	// Multi-sig blobs should be longer than single-sig due to Signers array
	// This is a heuristic to verify Multisign() was actually called
	assert.Greater(t, len(signedBlob), 200, "multi-sig blob should be longer than typical single-sig")
}

// TestPeersystSigner_MultiSignature_VsSingleSignature tests that multi-sig
// transactions produce different output than single-sig transactions.
func TestPeersystSigner_MultiSignature_VsSingleSignature(t *testing.T) {
	t.Parallel()

	signer := xrpsigner.NewPeersystSigner()
	ctx := context.Background()
	testSeed := "sEdTM1uX8pu2do5XvTnutH6HsouMaM2"

	// Single-sig transaction
	// Note: For single-sig, we rely on wallet.Sign() to add SigningPubKey automatically
	// The absence of explicit SigningPubKey field setting means we'll skip adding it to the tx map
	singleSigTx := dtoxrp.TxInput{
		TransactionType:    "Payment",
		Account:            "rG31cLyErnqeVj2eomEjBZtq7PYaupGYzL",
		Destination:        "rPEPPER7kfTD9w2To4CQk6UCfuHM9c6GDY",
		Amount:             "1000000",
		Fee:                "12",
		Sequence:           1,
		LastLedgerSequence: 10000000,
		// SigningPubKey deliberately not set - will be added by wallet.Sign()
	}

	// Multi-sig transaction
	// CRITICAL: SigningPubKey must be explicitly set to "" for multi-sig
	multiSigTx := dtoxrp.TxInput{
		TransactionType:    "Payment",
		Account:            "rG31cLyErnqeVj2eomEjBZtq7PYaupGYzL",
		Destination:        "rPEPPER7kfTD9w2To4CQk6UCfuHM9c6GDY",
		Amount:             "1000000",
		Fee:                "12",
		Sequence:           1,
		LastLedgerSequence: 10000000,
		SigningPubKey:      "", // MUST be explicitly set for multi-sig
	}

	// Sign both
	_, singleBlob, err1 := signer.SignTransactionNative(ctx, &singleSigTx, testSeed, false, nil)
	_, multiBlob, err2 := signer.SignTransactionNative(ctx, &multiSigTx, testSeed, true, nil)

	// Both should succeed
	require.NoError(t, err1, "single-sig should succeed")
	require.NoError(t, err2, "multi-sig should succeed")

	// Blobs should be different
	assert.NotEqual(t, singleBlob, multiBlob, "single-sig and multi-sig blobs should differ")

	// Multi-sig blob should contain Signers array (longer)
	assert.Greater(t, len(multiBlob), len(singleBlob),
		"multi-sig blob should be longer due to Signers array structure")
}

// TestPeersystSigner_MultiSignature_FeeValidation tests that multi-signature
// transactions can be signed with various fee amounts.
// Note: Fee validation is done by XRP Ledger, not the signer.
func TestPeersystSigner_MultiSignature_FeeValidation(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		fee         string
		description string
	}{
		{
			name:        "multi-sig fee for 2 signers",
			fee:         "36", // (2+1) × 12
			description: "Standard multi-sig fee calculation",
		},
		{
			name:        "multi-sig fee for 3 signers",
			fee:         "48", // (3+1) × 12
			description: "Higher multi-sig fee is acceptable",
		},
		{
			name:        "base fee with multi-sig structure",
			fee:         "12", // Base fee
			description: "Low fee is accepted by signer but may be rejected by ledger",
		},
	}

	signer := xrpsigner.NewPeersystSigner()
	ctx := context.Background()
	testSeed := "sEdTM1uX8pu2do5XvTnutH6HsouMaM2"

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			txInput := &dtoxrp.TxInput{
				TransactionType:    "Payment",
				Account:            "rG31cLyErnqeVj2eomEjBZtq7PYaupGYzL",
				Destination:        "rPEPPER7kfTD9w2To4CQk6UCfuHM9c6GDY",
				Amount:             "1000000",
				Fee:                tc.fee,
				Sequence:           1,
				SigningPubKey:      "", // Multi-sig indicator
				LastLedgerSequence: 10000000,
			}

			// Execute
			_, _, err := signer.SignTransactionNative(ctx, txInput, testSeed, true, nil)

			// Signer doesn't validate fees - all should succeed
			assert.NoError(t, err, tc.description)
		})
	}
}

// TestPeersystSigner_MultiSignature_SequentialSigning tests that multiple
// signers can sign the same transaction sequentially (simulated).
func TestPeersystSigner_MultiSignature_SequentialSigning(t *testing.T) {
	t.Parallel()

	signer := xrpsigner.NewPeersystSigner()
	ctx := context.Background()

	// Two different seeds for two different signers
	seed1 := "sEdTM1uX8pu2do5XvTnutH6HsouMaM2"
	seed2 := "sEdSKaCy2JT7JaM7v95H9SxkhP9wS2r" // Different seed

	// Create multi-sig transaction
	txInput := &dtoxrp.TxInput{
		TransactionType:    "Payment",
		Account:            "rG31cLyErnqeVj2eomEjBZtq7PYaupGYzL",
		Destination:        "rPEPPER7kfTD9w2To4CQk6UCfuHM9c6GDY",
		Amount:             "1000000",
		Fee:                "36", // Multi-sig fee
		Sequence:           1,
		SigningPubKey:      "", // Multi-sig
		LastLedgerSequence: 10000000,
	}

	// First signer signs
	hash1, blob1, err1 := signer.SignTransactionNative(ctx, txInput, seed1, true, nil)
	require.NoError(t, err1, "first signer should succeed")
	assert.NotEmpty(t, hash1)
	assert.NotEmpty(t, blob1)

	// Second signer signs the same transaction
	hash2, blob2, err2 := signer.SignTransactionNative(ctx, txInput, seed2, true, nil)
	require.NoError(t, err2, "second signer should succeed")
	assert.NotEmpty(t, hash2)
	assert.NotEmpty(t, blob2)

	// Note: In real multi-sig flow, the second signer would receive the
	// partially signed transaction and add their signature.
	// For now, we're verifying that both can independently create signatures.
	// Full sequential signing support requires transaction file handling.
}

// TestPeersystSigner_MultiSignature_DeterministicSigning tests that
// multi-sig signing is deterministic.
func TestPeersystSigner_MultiSignature_DeterministicSigning(t *testing.T) {
	t.Parallel()

	signer := xrpsigner.NewPeersystSigner()
	ctx := context.Background()
	testSeed := "sEdTM1uX8pu2do5XvTnutH6HsouMaM2"

	txInput := &dtoxrp.TxInput{
		TransactionType:    "Payment",
		Account:            "rG31cLyErnqeVj2eomEjBZtq7PYaupGYzL",
		Destination:        "rPEPPER7kfTD9w2To4CQk6UCfuHM9c6GDY",
		Amount:             "1000000",
		Fee:                "36",
		Sequence:           1,
		SigningPubKey:      "", // Multi-sig
		LastLedgerSequence: 10000000,
	}

	// Sign twice
	hash1, blob1, err1 := signer.SignTransactionNative(ctx, txInput, testSeed, true, nil)
	hash2, blob2, err2 := signer.SignTransactionNative(ctx, txInput, testSeed, true, nil)

	// Both should succeed
	require.NoError(t, err1)
	require.NoError(t, err2)

	// Results should be identical (deterministic)
	assert.Equal(t, hash1, hash2, "multi-sig hash should be deterministic")
	assert.Equal(t, blob1, blob2, "multi-sig blob should be deterministic")
}
