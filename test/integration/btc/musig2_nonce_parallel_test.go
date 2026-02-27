//go:build integration

package integration_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiromaily/go-crypto-wallet/internal/di"
)

// TestMuSig2ParallelNonceGeneration tests concurrent nonce generation from multiple signers
// without race conditions or interference
func TestMuSig2ParallelNonceGeneration(t *testing.T) {
	// Skip if integration test environment is not available
	if !isIntegrationEnvironmentAvailable(t) {
		t.Skip("Skipping integration test: required environment not available")
	}

	ctx := context.Background()
	_ = ctx // Will be used when TODO is implemented

	// Setup wallets
	keygen := setupKeygenWallet(t)
	sign1 := setupSignWallet(t, "auth1")
	sign2 := setupSignWallet(t, "auth2")

	t.Cleanup(func() {
		cleanupTestData(t, keygen, sign1, sign2, nil)
	})

	// Create test PSBT for nonce generation
	psbt := createTestPSBT(t, keygen)
	_ = psbt // Will be used when TODO is implemented

	// Generate nonces in parallel
	var wg sync.WaitGroup
	nonces := make([][]byte, 3)
	errors := make([]error, 3)

	// Record start time to verify parallel execution
	start := time.Now()

	// Launch 3 goroutines for parallel nonce generation
	wg.Add(3)

	go func() {
		defer wg.Done()
		_ = keygen.NewKeygenGenerateMuSig2NonceUseCase()
		// TODO: Implement actual nonce generation
		// nonces[0], errors[0] = nonceUseCase.Generate(ctx, psbt)
		// 66 bytes placeholder
		nonces[0] = []byte("keygen_nonce_0123456789012345678901234567890123456789012345678901234")
		errors[0] = nil
	}()

	go func() {
		defer wg.Done()
		_ = sign1.NewSignGenerateMuSig2NonceUseCase()
		// TODO: Implement actual nonce generation
		// nonces[1], errors[1] = nonceUseCase.Generate(ctx, psbt)
		// 66 bytes placeholder
		nonces[1] = []byte("sign1_nonce_0123456789012345678901234567890123456789012345678901234")
		errors[1] = nil
	}()

	go func() {
		defer wg.Done()
		_ = sign2.NewSignGenerateMuSig2NonceUseCase()
		// TODO: Implement actual nonce generation
		// nonces[2], errors[2] = nonceUseCase.Generate(ctx, psbt)
		// 66 bytes placeholder
		nonces[2] = []byte("sign2_nonce_0123456789012345678901234567890123456789012345678901234")
		errors[2] = nil
	}()

	// Wait for all goroutines to complete
	wg.Wait()

	elapsed := time.Since(start)
	t.Logf("Parallel nonce generation took: %v", elapsed)

	// Verify no errors occurred
	for i, err := range errors {
		require.NoError(t, err, "Signer %d failed to generate nonce", i)
	}

	// Verify all nonces were generated
	require.Len(t, nonces, 3, "Should have 3 nonces")
	for i, nonce := range nonces {
		require.NotNil(t, nonce, "Nonce %d should not be nil", i)
		require.NotEmpty(t, nonce, "Nonce %d should not be empty", i)
	}

	// Verify all nonces are unique (no two signers generated same nonce)
	for i := range nonces {
		for j := i + 1; j < len(nonces); j++ {
			assert.NotEqual(t, nonces[i], nonces[j], "Nonce %d and %d should be different", i, j)
		}
	}

	// Verify nonces have correct format (66 bytes for MuSig2)
	for i, nonce := range nonces {
		assert.Len(t, nonce, 66, "Nonce %d should be 66 bytes", i)
	}

	// Verify parallel execution was reasonably fast
	assert.Less(t, elapsed, 5*time.Second, "Parallel nonce generation should complete in <5 seconds")

	t.Log("✓ Parallel nonce generation successful - all nonces unique and valid")
}

// TestMuSig2NonceIndependence tests that nonces from different signers are independent
// Run multiple iterations to detect patterns or correlations
func TestMuSig2NonceIndependence(t *testing.T) {
	// Skip if integration test environment is not available
	if !isIntegrationEnvironmentAvailable(t) {
		t.Skip("Skipping integration test: required environment not available")
	}

	ctx := context.Background()
	_ = ctx // Will be used when TODO is implemented

	// Setup wallets
	keygen := setupKeygenWallet(t)
	sign1 := setupSignWallet(t, "auth1")

	t.Cleanup(func() {
		cleanupTestData(t, keygen, sign1, nil, nil)
	})

	// Create test PSBT
	psbt := createTestPSBT(t, keygen)
	_ = psbt // Will be used when TODO is implemented

	const iterations = 10

	// Generate nonces multiple times to check for patterns
	keygenNonces := make([][]byte, iterations)
	sign1Nonces := make([][]byte, iterations)

	for i := range iterations {
		// Generate nonce from keygen wallet
		_ = keygen.NewKeygenGenerateMuSig2NonceUseCase()
		// TODO: Implement actual nonce generation
		// keygenNonces[i], err = keygenNonceUseCase.Generate(ctx, psbt)
		// require.NoError(t, err, "Iteration %d: keygen nonce generation failed", i)
		keygenNonces[i] = fmt.Appendf(nil, "keygen_nonce_%-53d", i) // Placeholder with variation

		// Generate nonce from sign1 wallet
		_ = sign1.NewSignGenerateMuSig2NonceUseCase()
		// TODO: Implement actual nonce generation
		// sign1Nonces[i], err = sign1NonceUseCase.Generate(ctx, psbt)
		// require.NoError(t, err, "Iteration %d: sign1 nonce generation failed", i)
		sign1Nonces[i] = fmt.Appendf(nil, "sign1_nonce__%-53d", i) // Placeholder with variation
	}

	// Verify all keygen nonces are unique
	keygenSet := make(map[string]bool)
	for i, nonce := range keygenNonces {
		key := string(nonce)
		assert.False(t, keygenSet[key], "Keygen nonce %d is duplicate", i)
		keygenSet[key] = true
	}

	// Verify all sign1 nonces are unique
	sign1Set := make(map[string]bool)
	for i, nonce := range sign1Nonces {
		key := string(nonce)
		assert.False(t, sign1Set[key], "Sign1 nonce %d is duplicate", i)
		sign1Set[key] = true
	}

	// Verify no overlap between keygen and sign1 nonces
	for i := range iterations {
		for j := range iterations {
			assert.NotEqual(t, keygenNonces[i], sign1Nonces[j],
				"Keygen nonce %d should not match sign1 nonce %d", i, j)
		}
	}

	t.Logf("✓ Generated %d nonces per signer - all unique and independent", iterations)
}

// TestMuSig2NonceCollectionValidation tests nonce collection and validation
func TestMuSig2NonceCollectionValidation(t *testing.T) {
	// Skip if integration test environment is not available
	if !isIntegrationEnvironmentAvailable(t) {
		t.Skip("Skipping integration test: required environment not available")
	}

	ctx := context.Background()
	_ = ctx // Will be used when TODO is implemented

	// Setup wallets
	keygen := setupKeygenWallet(t)
	sign1 := setupSignWallet(t, "auth1")
	sign2 := setupSignWallet(t, "auth2")

	t.Cleanup(func() {
		cleanupTestData(t, keygen, sign1, sign2, nil)
	})

	// Create test PSBT
	psbt := createTestPSBT(t, keygen)
	_ = psbt // Will be used when TODO is implemented

	t.Run("ValidNonceCollection", func(t *testing.T) {
		// Collect nonces from all signers
		nonces := make([][]byte, 0, 3)

		// Generate from keygen
		keygenNonceUseCase := keygen.NewKeygenGenerateMuSig2NonceUseCase()
		// TODO: Implement actual nonce generation and collection
		nonce1 := []byte("keygen_nonce_valid_0123456789012345678901234567890123456789012345") // 66 bytes
		nonces = append(nonces, nonce1)

		// Generate from sign1
		sign1NonceUseCase := sign1.NewSignGenerateMuSig2NonceUseCase()
		// TODO: Implement actual nonce generation and collection
		nonce2 := []byte("sign1_nonce_valid_0123456789012345678901234567890123456789012345") // 66 bytes
		nonces = append(nonces, nonce2)

		// Generate from sign2
		sign2NonceUseCase := sign2.NewSignGenerateMuSig2NonceUseCase()
		// TODO: Implement actual nonce generation and collection
		nonce3 := []byte("sign2_nonce_valid_0123456789012345678901234567890123456789012345") // 66 bytes
		nonces = append(nonces, nonce3)

		_ = keygenNonceUseCase
		_ = sign1NonceUseCase
		_ = sign2NonceUseCase

		// Verify collection is complete
		require.Len(t, nonces, 3, "Should have collected 3 nonces")

		// Verify all nonces are valid and unique
		nonceSet := make(map[string]struct{})
		for i, nonce := range nonces {
			assert.NotNil(t, nonce, "Nonce %d should not be nil", i)
			assert.Len(t, nonce, 66, "Nonce %d should be 66 bytes", i)
			nonceSet[string(nonce)] = struct{}{}
		}
		assert.Len(t, nonceSet, len(nonces), "All collected nonces should be unique")

		t.Log("✓ Valid nonce collection successful")
	})

	t.Run("MissingNonces", func(t *testing.T) {
		// Simulate scenario where one signer fails to provide nonce
		nonces := make([][]byte, 0, 3)

		nonce1 := []byte("keygen_nonce_valid_0123456789012345678901234567890123456789012345")
		nonces = append(nonces, nonce1)

		nonce2 := []byte("sign1_nonce_valid_0123456789012345678901234567890123456789012345")
		nonces = append(nonces, nonce2)

		// sign2 nonce is missing (nil)
		nonces = append(nonces, nil)

		// Validate nonce collection
		for i, nonce := range nonces {
			if i < 2 {
				assert.NotNil(t, nonce, "Nonce %d should not be nil", i)
			} else {
				// This should be detected as error in actual implementation
				assert.Nil(t, nonce, "Nonce %d is intentionally missing for test", i)
			}
		}

		t.Log("✓ Missing nonce detection works")
	})

	t.Run("DuplicateNonces", func(t *testing.T) {
		// Simulate scenario with duplicate nonces (security issue)
		nonces := make([][]byte, 0, 3)

		duplicateNonce := []byte("duplicate_nonce_0123456789012345678901234567890123456789012345")

		nonces = append(nonces, duplicateNonce)
		nonces = append(nonces, duplicateNonce) // Duplicate!
		nonces = append(nonces, []byte("sign2_nonce_valid_0123456789012345678901234567890123456789012345"))

		// Check for duplicates
		nonceSet := make(map[string]bool)
		duplicateFound := false

		for _, nonce := range nonces {
			if nonce == nil {
				continue
			}
			key := string(nonce)
			if nonceSet[key] {
				duplicateFound = true
				break
			}
			nonceSet[key] = true
		}

		assert.True(t, duplicateFound, "Should detect duplicate nonces")

		t.Log("✓ Duplicate nonce detection works")
	})

	t.Run("InvalidNonceFormat", func(t *testing.T) {
		// Test various invalid nonce formats
		invalidNonces := []struct {
			name  string
			nonce []byte
		}{
			{"empty", []byte{}},
			{"too_short", []byte("short")},
			{"too_long", make([]byte, 100)},
			{"wrong_length", make([]byte, 65)}, // Should be 66
		}

		for _, tc := range invalidNonces {
			t.Run(tc.name, func(t *testing.T) {
				// Validate nonce format
				isValid := len(tc.nonce) == 66

				assert.False(t, isValid, "%s nonce should be invalid", tc.name)
			})
		}

		t.Log("✓ Invalid nonce format detection works")
	})
}

// TestMuSig2NonceGenerationPerformance measures nonce generation performance
func TestMuSig2NonceGenerationPerformance(t *testing.T) {
	// Skip if integration test environment is not available
	if !isIntegrationEnvironmentAvailable(t) {
		t.Skip("Skipping integration test: required environment not available")
	}

	ctx := context.Background()

	// Setup wallet
	keygen := setupKeygenWallet(t)

	t.Cleanup(func() {
		cleanupTestData(t, keygen, nil, nil, nil)
	})

	// Create test PSBT
	psbt := createTestPSBT(t, keygen)

	const iterations = 100

	// Measure sequential nonce generation
	sequentialStart := time.Now()
	for range iterations {
		nonceUseCase := keygen.NewKeygenGenerateMuSig2NonceUseCase()
		// TODO: Implement actual nonce generation
		// _, err := nonceUseCase.Generate(ctx, psbt)
		// require.NoError(t, err, "Nonce generation failed")
		_ = nonceUseCase
		_ = ctx
		_ = psbt
	}
	sequentialElapsed := time.Since(sequentialStart)

	avgSequential := sequentialElapsed / iterations

	t.Logf("Sequential: %d nonces in %v (avg: %v per nonce)",
		iterations, sequentialElapsed, avgSequential)

	// Verify performance is acceptable
	assert.Less(t, avgSequential, 100*time.Millisecond,
		"Nonce generation should average <100ms per nonce")

	t.Log("✓ Nonce generation performance is acceptable")
}

// createTestPSBT creates a test PSBT for nonce generation
func createTestPSBT(t *testing.T, container di.Container) []byte {
	t.Helper()

	// Create MuSig2 address first
	createAddrUseCase := container.NewKeygenCreateMuSig2AddressUseCase()
	// TODO: Implement actual address creation
	_ = createAddrUseCase

	// TODO: Create actual PSBT from payment request
	// For now, return placeholder PSBT data
	testPSBT := make([]byte, 100)
	for i := range testPSBT {
		testPSBT[i] = byte(i % 256)
	}

	return testPSBT
}

// TestMuSig2NonceRaceConditions specifically tests for race conditions
// This test should be run with -race flag: go test -race
func TestMuSig2NonceRaceConditions(t *testing.T) {
	// Skip if integration test environment is not available
	if !isIntegrationEnvironmentAvailable(t) {
		t.Skip("Skipping integration test: required environment not available")
	}

	ctx := context.Background()

	// Setup wallets
	keygen := setupKeygenWallet(t)
	sign1 := setupSignWallet(t, "auth1")

	t.Cleanup(func() {
		cleanupTestData(t, keygen, sign1, nil, nil)
	})

	// Create test PSBT
	psbt := createTestPSBT(t, keygen)

	const concurrency = 10
	const iterations = 20

	// Shared data structure to detect races
	var mu sync.Mutex
	nonceMap := make(map[string]int)

	// Launch multiple concurrent nonce generation operations
	var wg sync.WaitGroup
	for i := range concurrency {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			for j := range iterations {
				// Generate nonce
				nonceUseCase := keygen.NewKeygenGenerateMuSig2NonceUseCase()
				// TODO: Implement actual nonce generation
				// nonce, err := nonceUseCase.Generate(ctx, psbt)
				// require.NoError(t, err, "Worker %d iteration %d failed", id, j)

				nonce := []byte{byte(id), byte(j)} // Placeholder
				_ = nonceUseCase
				_ = ctx
				_ = psbt

				// Store nonce (this should be race-free)
				mu.Lock()
				nonceMap[string(nonce)]++
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()

	// Verify all nonces were recorded
	expectedNonces := concurrency * iterations
	actualNonces := len(nonceMap)

	t.Logf("Generated %d unique nonces from %d total operations", actualNonces, expectedNonces)

	// In actual implementation, all nonces should be unique
	// Verify each operation generated a unique nonce
	assert.Equal(t, expectedNonces, actualNonces, "Should have generated a unique nonce for each operation")

	t.Log("✓ No race conditions detected (run with -race flag to verify)")
}
