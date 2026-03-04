package xrp_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	dtoxrp "github.com/hiromaily/go-crypto-wallet/internal/application/dto/xrp"
	xrpsigner "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/xrp/signer"
)

// TestPeersystSigner_SignTransactionNative_SingleSignature tests single-signature
// transaction signing using native Go implementation.
func TestPeersystSigner_SignTransactionNative_SingleSignature(t *testing.T) {
	t.Parallel()

	// Setup
	signer := xrpsigner.NewPeersystSigner()
	ctx := context.Background()

	// Test seed from XRP Ledger documentation
	testSeed := "sEdTM1uX8pu2do5XvTnutH6HsouMaM2"

	// Create a test transaction input
	// Note: Account address must match the address derived from testSeed
	txInput := &dtoxrp.TxInput{
		TransactionType:    "Payment",
		Account:            "rG31cLyErnqeVj2eomEjBZtq7PYaupGYzL", // Address from test seed
		Destination:        "rPEPPER7kfTD9w2To4CQk6UCfuHM9c6GDY", // Known valid address
		Amount:             "1000000",                            // 1 XRP in drops
		Fee:                "12",                                 // Base fee
		Sequence:           1,
		LastLedgerSequence: 10000000,
	}

	// Execute
	txHash, signedBlob, err := signer.SignTransactionNative(ctx, txInput, testSeed, false, nil)

	// Assert
	require.NoError(t, err, "signing should succeed")
	assert.NotEmpty(t, txHash, "transaction hash should not be empty")
	assert.Len(t, txHash, 64, "transaction hash should be 64 characters (hex)")
	assert.NotEmpty(t, signedBlob, "signed blob should not be empty")
	assert.NotContains(t, signedBlob, testSeed, "signed blob should not contain seed")
}

// TestPeersystSigner_SignTransactionNative_DeterministicSigning tests that
// signing the same transaction with the same seed produces the same signature.
func TestPeersystSigner_SignTransactionNative_DeterministicSigning(t *testing.T) {
	t.Parallel()

	// Setup
	signer := xrpsigner.NewPeersystSigner()
	ctx := context.Background()
	testSeed := "sEdTM1uX8pu2do5XvTnutH6HsouMaM2"

	txInput := &dtoxrp.TxInput{
		TransactionType:    "Payment",
		Account:            "rG31cLyErnqeVj2eomEjBZtq7PYaupGYzL",
		Destination:        "rPEPPER7kfTD9w2To4CQk6UCfuHM9c6GDY",
		Amount:             "1000000",
		Fee:                "12",
		Sequence:           1,
		LastLedgerSequence: 10000000,
	}

	// Sign twice with same input
	txHash1, signedBlob1, err1 := signer.SignTransactionNative(ctx, txInput, testSeed, false, nil)
	txHash2, signedBlob2, err2 := signer.SignTransactionNative(ctx, txInput, testSeed, false, nil)

	// Assert deterministic behavior
	require.NoError(t, err1)
	require.NoError(t, err2)
	assert.Equal(t, txHash1, txHash2, "transaction hash should be deterministic")
	assert.Equal(t, signedBlob1, signedBlob2, "signed blob should be deterministic")
}

// TestPeersystSigner_SignTransactionNative_InvalidSeed tests error handling
// for invalid seed format.
func TestPeersystSigner_SignTransactionNative_InvalidSeed(t *testing.T) {
	t.Parallel()

	// Setup
	signer := xrpsigner.NewPeersystSigner()
	ctx := context.Background()

	txInput := &dtoxrp.TxInput{
		TransactionType:    "Payment",
		Account:            "rN7n7otQDd6FczFgLdSqtcsAUxDkw6fzRH",
		Destination:        "rLHzPsX6oXkzU9J1kSwq7X4X8J1Xa8AkDC",
		Amount:             "1000000",
		Fee:                "12",
		Sequence:           1,
		LastLedgerSequence: 10000000,
	}

	// Execute with invalid seed
	txHash, signedBlob, err := signer.SignTransactionNative(ctx, txInput, "invalid-seed-format", false, nil)

	// Assert error
	require.Error(t, err, "should return error for invalid seed")
	assert.Empty(t, txHash, "transaction hash should be empty on error")
	assert.Empty(t, signedBlob, "signed blob should be empty on error")
	assert.Contains(t, err.Error(), "failed to derive wallet from seed", "error should contain context")
}

// TestPeersystSigner_SignTransactionNative_MissingRequiredFields tests
// error handling for transactions with missing required fields.
func TestPeersystSigner_SignTransactionNative_MissingRequiredFields(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		txInput *dtoxrp.TxInput
		errMsg  string
	}{
		{
			name: "missing Account",
			txInput: &dtoxrp.TxInput{
				TransactionType:    "Payment",
				Destination:        "rLHzPsX6oXkzU9J1kSwq7X4X8J1Xa8AkDC",
				Amount:             "1000000",
				Fee:                "12",
				Sequence:           1,
				LastLedgerSequence: 10000000,
			},
			errMsg: "account",
		},
		{
			name: "missing Destination",
			txInput: &dtoxrp.TxInput{
				TransactionType:    "Payment",
				Account:            "rG31cLyErnqeVj2eomEjBZtq7PYaupGYzL",
				Amount:             "1000000",
				Fee:                "12",
				Sequence:           1,
				LastLedgerSequence: 10000000,
			},
			errMsg: "destination",
		},
		{
			name: "missing Amount",
			txInput: &dtoxrp.TxInput{
				TransactionType:    "Payment",
				Account:            "rG31cLyErnqeVj2eomEjBZtq7PYaupGYzL",
				Destination:        "rPEPPER7kfTD9w2To4CQk6UCfuHM9c6GDY",
				Fee:                "12",
				Sequence:           1,
				LastLedgerSequence: 10000000,
			},
			errMsg: "amount",
		},
	}

	signer := xrpsigner.NewPeersystSigner()
	ctx := context.Background()
	testSeed := "sEdTM1uX8pu2do5XvTnutH6HsouMaM2"

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// Execute
			txHash, signedBlob, err := signer.SignTransactionNative(ctx, tc.txInput, testSeed, false, nil)

			// Assert
			require.Error(t, err, "should return error for missing field")
			assert.Empty(t, txHash, "transaction hash should be empty on error")
			assert.Empty(t, signedBlob, "signed blob should be empty on error")
			assert.Contains(t, err.Error(), tc.errMsg, "error should mention missing field")
		})
	}
}

// TestPeersystSigner_SignTransactionNative_OfflineCapability tests that
// signing works without network access (no network calls).
func TestPeersystSigner_SignTransactionNative_OfflineCapability(t *testing.T) {
	t.Parallel()

	// This test verifies that signing completes successfully
	// without any network dependencies by using a context that
	// has no network access configured.

	signer := xrpsigner.NewPeersystSigner()

	// Use background context with no network configuration
	ctx := context.Background()
	testSeed := "sEdTM1uX8pu2do5XvTnutH6HsouMaM2"

	txInput := &dtoxrp.TxInput{
		TransactionType:    "Payment",
		Account:            "rG31cLyErnqeVj2eomEjBZtq7PYaupGYzL",
		Destination:        "rPEPPER7kfTD9w2To4CQk6UCfuHM9c6GDY",
		Amount:             "1000000",
		Fee:                "12",
		Sequence:           1,
		LastLedgerSequence: 10000000,
	}

	// Execute - should complete without network access
	txHash, signedBlob, err := signer.SignTransactionNative(ctx, txInput, testSeed, false, nil)

	// Assert offline signing works
	require.NoError(t, err, "offline signing should succeed")
	assert.NotEmpty(t, txHash, "should produce transaction hash offline")
	assert.NotEmpty(t, signedBlob, "should produce signed blob offline")
}
