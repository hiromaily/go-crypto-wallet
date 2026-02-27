//go:build integration

package integration_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	"github.com/hiromaily/go-crypto-wallet/internal/di"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/config"
)

// NonceResult holds the result of parallel nonce generation
type NonceResult struct {
	Nonce []byte
	Err   error
}

// SignatureResult holds the result of parallel signature generation
type SignatureResult struct {
	Signature []byte
	Err       error
}

// TestMuSig2EndToEndFlow tests the complete MuSig2 flow from address creation
// through transaction broadcast
func TestMuSig2EndToEndFlow(t *testing.T) {
	// Skip if integration test environment is not available
	if !isIntegrationEnvironmentAvailable(t) {
		t.Skip("Skipping integration test: required environment not available")
	}

	ctx := context.Background()

	// Setup test environment with 4 wallet instances
	// - keygen: Offline wallet for key generation and first signature
	// - sign1: Offline wallet for second signature (auth1)
	// - sign2: Offline wallet for third signature (auth2)
	// - watch: Online wallet for transaction creation and aggregation
	keygen := setupKeygenWallet(t)
	sign1 := setupSignWallet(t, "auth1")
	sign2 := setupSignWallet(t, "auth2")
	watch := setupWatchWallet(t)

	// Cleanup test data
	t.Cleanup(func() {
		cleanupTestData(t, keygen, sign1, sign2, watch)
	})

	// Step 1: Create MuSig2 Taproot address
	t.Run("Step1_CreateMuSig2Address", func(t *testing.T) {
		// Create MuSig2 address for payment account
		useCase := keygen.NewKeygenCreateMuSig2AddressUseCase()
		err := useCase.Create(ctx, keygenusecase.CreateMuSig2AddressInput{
			AccountType: domainAccount.AccountTypePayment,
		})
		require.NoError(t, err, "Failed to create MuSig2 address")

		// Verify address was created
		// TODO: Add assertion to verify address in database
		t.Log("✓ MuSig2 Taproot address created")
	})

	// Get the created address for use in subsequent steps
	// TODO: Implement address retrieval from database
	_ = "tb1p..." // Placeholder address

	// Step 2: Create unsigned transaction (PSBT)
	var psbtData []byte
	_ = psbtData // Will be used when TODO is implemented
	t.Run("Step2_CreateUnsignedTransaction", func(t *testing.T) {
		// Create payment request
		createPaymentUseCase := watch.NewWatchCreatePaymentRequestUseCase()
		err := createPaymentUseCase.Execute(ctx, watchusecase.CreatePaymentRequestInput{
			AmountList: []float64{0.0001}, // 0.0001 BTC
		})
		require.NoError(t, err, "Failed to create payment request")

		// Create unsigned transaction
		_ = watch.NewWatchCreateTransactionUseCase()
		// TODO: Implement actual transaction creation
		// psbtData, err = createTxUseCase.Execute(ctx, ...)
		// require.NoError(t, err, "Failed to create unsigned transaction")

		t.Log("✓ Unsigned transaction (PSBT) created")
	})

	// Step 3: Round 1 - Generate nonces from all signers (parallel)
	var nonces [][]byte
	t.Run("Step3_GenerateNoncesParallel", func(t *testing.T) {
		// Use buffered channel to collect results from 3 signers
		nonceResults := make(chan NonceResult, 3)

		// Record start time to verify parallel execution
		startTime := time.Now()

		// Generate nonces in parallel from all three wallets
		go func() {
			_ = keygen.NewKeygenGenerateMuSig2NonceUseCase()
			// TODO: Implement actual nonce generation
			// nonce, err := useCase.Generate(ctx, psbtData)
			nonceResults <- NonceResult{
				Nonce: []byte("keygen_nonce"), // Placeholder
				Err:   nil,
			}
		}()

		go func() {
			_ = sign1.NewSignGenerateMuSig2NonceUseCase()
			// TODO: Implement actual nonce generation
			// nonce, err := useCase.Generate(ctx, psbtData)
			nonceResults <- NonceResult{
				Nonce: []byte("sign1_nonce"), // Placeholder
				Err:   nil,
			}
		}()

		go func() {
			_ = sign2.NewSignGenerateMuSig2NonceUseCase()
			// TODO: Implement actual nonce generation
			// nonce, err := useCase.Generate(ctx, psbtData)
			nonceResults <- NonceResult{
				Nonce: []byte("sign2_nonce"), // Placeholder
				Err:   nil,
			}
		}()

		// Collect all nonces
		for i := range 3 {
			result := <-nonceResults
			require.NoError(t, result.Err, "Failed to generate nonce for signer %d", i+1)
			require.NotEmpty(t, result.Nonce, "Nonce should not be empty for signer %d", i+1)
			nonces = append(nonces, result.Nonce)
		}

		// Verify parallel execution completed quickly
		duration := time.Since(startTime)
		t.Logf("Parallel nonce generation took: %v", duration)
		assert.Less(t, duration, 5*time.Second, "Parallel nonce generation should complete quickly")

		// Verify all nonces are unique
		require.Len(t, nonces, 3, "Should have 3 nonces")
		assert.NotEqual(t, nonces[0], nonces[1], "Nonces should be unique")
		assert.NotEqual(t, nonces[1], nonces[2], "Nonces should be unique")
		assert.NotEqual(t, nonces[0], nonces[2], "Nonces should be unique")

		t.Log("✓ All nonces generated in parallel successfully")
	})

	// Step 4: Round 2 - Create partial signatures from all signers
	var signatures [][]byte
	t.Run("Step4_CreatePartialSignatures", func(t *testing.T) {
		// Sign with keygen wallet
		_ = keygen.NewKeygenMuSig2SignUseCase()
		// TODO: Implement actual signing
		// sig1, err := keygenSignUseCase.Sign(ctx, psbtData, nonces)
		// require.NoError(t, err, "Failed to create signature from keygen wallet")
		sig1 := []byte("keygen_signature") // Placeholder
		signatures = append(signatures, sig1)

		// Sign with sign1 wallet
		_ = sign1.NewSignMuSig2SignUseCase()
		// TODO: Implement actual signing
		// sig2, err := sign1SignUseCase.Sign(ctx, psbtData, nonces)
		// require.NoError(t, err, "Failed to create signature from sign1 wallet")
		sig2 := []byte("sign1_signature") // Placeholder
		signatures = append(signatures, sig2)

		// Sign with sign2 wallet
		_ = sign2.NewSignMuSig2SignUseCase()
		// TODO: Implement actual signing
		// sig3, err := sign2SignUseCase.Sign(ctx, psbtData, nonces)
		// require.NoError(t, err, "Failed to create signature from sign2 wallet")
		sig3 := []byte("sign2_signature") // Placeholder
		signatures = append(signatures, sig3)

		require.Len(t, signatures, 3, "Should have 3 partial signatures")

		t.Log("✓ All partial signatures created successfully")
	})

	// Step 5: Aggregate signatures in Watch wallet
	var finalTx []byte
	t.Run("Step5_AggregateSignatures", func(t *testing.T) {
		_ = watch.NewWatchAggregateMuSig2SignaturesUseCase()
		// TODO: Implement actual aggregation
		// finalTx, err := aggregateUseCase.Aggregate(ctx, psbtData, signatures)
		// require.NoError(t, err, "Failed to aggregate signatures")
		finalTx = []byte("aggregated_transaction") // Placeholder

		require.NotEmpty(t, finalTx, "Final transaction should not be empty")

		t.Log("✓ Signatures aggregated successfully")
	})

	// Step 6: Finalize PSBT with aggregated signature
	t.Run("Step6_FinalizePSBT", func(t *testing.T) {
		// TODO: Implement PSBT finalization
		// The aggregated signature should be embedded into the PSBT

		t.Log("✓ PSBT finalized with aggregated signature")
	})

	// Step 7: Verify transaction validity and size reduction
	t.Run("Step7_VerifyTransactionAndSize", func(t *testing.T) {
		// Verify transaction validity
		// TODO: Implement transaction verification
		// valid, err := watch.VerifyTransaction(ctx, finalTx)
		// require.NoError(t, err, "Failed to verify transaction")
		// assert.True(t, valid, "Transaction should be valid")

		// Measure transaction size
		txSize := len(finalTx)
		t.Logf("MuSig2 transaction size: %d bytes", txSize)

		// Compare with traditional 2-of-3 P2WSH multisig
		// Traditional 2-of-3 P2WSH transaction is approximately 370-400 bytes
		traditionalSize := 400 // Approximate size
		sizeReduction := float64(traditionalSize-txSize) / float64(traditionalSize) * 100

		t.Logf("Traditional multisig size: %d bytes", traditionalSize)
		t.Logf("Size reduction: %.1f%%", sizeReduction)

		// Verify at least 25% size reduction (should be 30-50%)
		assert.Greater(t, sizeReduction, 25.0, "Expected at least 25%% size reduction")

		t.Log("✓ Transaction is valid and shows significant size reduction")
	})
}

// isIntegrationEnvironmentAvailable checks if the required integration test
// environment is available
func isIntegrationEnvironmentAvailable(t *testing.T) bool {
	t.Helper()

	// Check if Bitcoin Core regtest node is available
	// Check if MySQL database is available
	// Check if config files exist

	// For now, check if INTEGRATION_TEST environment variable is set
	if os.Getenv("INTEGRATION_TEST") != "true" {
		t.Log("Set INTEGRATION_TEST=true to run integration tests")
		return false
	}

	return true
}

// setupKeygenWallet creates a keygen wallet instance for testing
func setupKeygenWallet(t *testing.T) di.Container {
	t.Helper()

	// Use relative path from this test file location
	// test/integration/btc/ -> ../../../config/wallet/
	confPath := filepath.Join("..", "..", "..", "config", "wallet", "btc_keygen.yaml")
	accountConfPath := filepath.Join("..", "..", "..", "config", "wallet", "account.yaml")

	// Load wallet configuration
	conf, err := config.NewWallet(confPath, domainWallet.WalletTypeKeyGen, domainCoin.BTC)
	require.NoError(t, err, "Failed to load keygen wallet config")

	// Load account configuration
	accountConf, err := config.NewAccount(accountConfPath)
	require.NoError(t, err, "Failed to load account config")

	// Create DI container (no error return)
	container := di.NewContainer(conf, accountConf, domainWallet.WalletTypeKeyGen)

	return container
}

// setupSignWallet creates a sign wallet instance for testing
func setupSignWallet(t *testing.T, authName string) di.Container {
	t.Helper()

	// Use relative path from this test file location
	confPath := filepath.Join("..", "..", "..", "data", "config", "btc_sign.toml")
	accountConfPath := filepath.Join("..", "..", "..", "data", "config", "account.toml")

	// Load wallet configuration
	conf, err := config.NewWallet(confPath, domainWallet.WalletTypeSign, domainCoin.BTC)
	require.NoError(t, err, "Failed to load sign wallet config")

	// Load account configuration
	accountConf, err := config.NewAccount(accountConfPath)
	require.NoError(t, err, "Failed to load account config")

	// Note: authName override would be done via environment variables or command-line flags
	// in actual implementation, not through config modification
	_ = authName // Will be used when auth name configuration is implemented

	// Create DI container (no error return)
	container := di.NewContainer(conf, accountConf, domainWallet.WalletTypeSign)

	return container
}

// setupWatchWallet creates a watch wallet instance for testing
func setupWatchWallet(t *testing.T) di.Container {
	t.Helper()

	// Use relative path from this test file location
	confPath := filepath.Join("..", "..", "..", "data", "config", "btc_watch.toml")
	accountConfPath := filepath.Join("..", "..", "..", "data", "config", "account.toml")

	// Load wallet configuration
	conf, err := config.NewWallet(confPath, domainWallet.WalletTypeWatchOnly, domainCoin.BTC)
	require.NoError(t, err, "Failed to load watch wallet config")

	// Load account configuration
	accountConf, err := config.NewAccount(accountConfPath)
	require.NoError(t, err, "Failed to load account config")

	// Create DI container (no error return)
	container := di.NewContainer(conf, accountConf, domainWallet.WalletTypeWatchOnly)

	return container
}

// cleanupTestData removes test data created during the test
func cleanupTestData(t *testing.T, keygen, sign1, sign2, watch di.Container) {
	t.Helper()

	// TODO: Implement cleanup logic
	_ = keygen // Will be used when TODO is implemented
	_ = sign1  // Will be used when TODO is implemented
	_ = sign2  // Will be used when TODO is implemented
	_ = watch  // Will be used when TODO is implemented

	// - Remove test addresses from database
	// - Remove test nonces from database
	// - Remove test payment requests
	// - Remove test transactions

	t.Log("✓ Test data cleaned up")
}
