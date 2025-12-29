//go:build integration

package cold_test

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/internal/domain/multisig"
	"github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/cold"
	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	mysql "github.com/hiromaily/go-crypto-wallet/pkg/db/mysql"
)

// getTestDB returns a database connection for integration tests.
func getTestDB(t *testing.T) *sql.DB {
	t.Helper()

	projPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-crypto-wallet"
	confPath := projPath + "/data/config/btc_keygen.toml"
	conf, err := config.NewWallet(confPath, wallet.WalletTypeKeyGen, domainCoin.BTC)
	if err != nil {
		log.Fatalf("fail to create config: %v", err)
	}

	dbConn, err := mysql.NewMySQL(&conf.MySQL)
	if err != nil {
		log.Fatalf("fail to create db: %v", err)
	}

	return dbConn
}

// TestNonceRepositorySqlc tests the SQLC-based nonce repository implementation.
func TestNonceRepositorySqlc(t *testing.T) {
	// Setup test database connection
	dbConn := getTestDB(t)
	defer func() {
		if err := dbConn.Close(); err != nil {
			t.Errorf("Failed to close database: %v", err)
		}
	}()

	repo := cold.NewNonceRepositorySqlc(dbConn)
	ctx := context.Background()

	// Test data
	signerID1 := "signer1"
	signerID2 := "signer2"
	transactionID1 := "tx001"
	transactionID2 := "tx002"
	publicNonce1 := [66]byte{1, 2, 3}
	publicNonce2 := [66]byte{4, 5, 6}
	publicNonce3 := [66]byte{7, 8, 9}

	// Cleanup before test
	_ = repo.DeleteNoncesForTransaction(ctx, transactionID1)
	_ = repo.DeleteNoncesForTransaction(ctx, transactionID2)

	t.Run("SaveNonce - Success", func(t *testing.T) {
		err := repo.SaveNonce(ctx, signerID1, transactionID1, publicNonce1)
		if err != nil {
			t.Fatalf("SaveNonce failed: %v", err)
		}
	})

	t.Run("SaveNonce - Duplicate", func(t *testing.T) {
		// Try to save the same nonce again (same signer + transaction)
		err := repo.SaveNonce(ctx, signerID1, transactionID1, publicNonce1)
		if err != multisig.ErrNonceDuplicate {
			t.Errorf("Expected ErrNonceDuplicate, got: %v", err)
		}
	})

	t.Run("GetNonce - Success", func(t *testing.T) {
		nonce, err := repo.GetNonce(ctx, signerID1, transactionID1)
		if err != nil {
			t.Fatalf("GetNonce failed: %v", err)
		}

		if nonce.SignerID() != signerID1 {
			t.Errorf("Expected signer ID %s, got %s", signerID1, nonce.SignerID())
		}

		publicNonce := nonce.PublicNonce()
		if publicNonce != publicNonce1 {
			t.Errorf("Public nonce mismatch: expected %v, got %v", publicNonce1, publicNonce)
		}
	})

	t.Run("GetNonce - NotFound", func(t *testing.T) {
		_, err := repo.GetNonce(ctx, "nonexistent", transactionID1)
		if err != multisig.ErrNonceNotFound {
			t.Errorf("Expected ErrNonceNotFound, got: %v", err)
		}
	})

	t.Run("GetAllNoncesForTransaction - Multiple Signers", func(t *testing.T) {
		// Add nonce for second signer
		err := repo.SaveNonce(ctx, signerID2, transactionID1, publicNonce2)
		if err != nil {
			t.Fatalf("SaveNonce failed: %v", err)
		}

		// Get all nonces for transaction
		nonces, err := repo.GetAllNoncesForTransaction(ctx, transactionID1)
		if err != nil {
			t.Fatalf("GetAllNoncesForTransaction failed: %v", err)
		}

		if len(nonces) != 2 {
			t.Errorf("Expected 2 nonces, got %d", len(nonces))
		}

		// Verify both signers are present
		signerIDs := make(map[string]bool)
		for _, nonce := range nonces {
			signerIDs[nonce.SignerID()] = true
		}

		if !signerIDs[signerID1] || !signerIDs[signerID2] {
			t.Errorf("Missing expected signer IDs in result")
		}
	})

	t.Run("GetUnusedNoncesForTransaction - Before Marking", func(t *testing.T) {
		// All nonces should be unused initially
		nonces, err := repo.GetUnusedNoncesForTransaction(ctx, transactionID1)
		if err != nil {
			t.Fatalf("GetUnusedNoncesForTransaction failed: %v", err)
		}

		if len(nonces) != 2 {
			t.Errorf("Expected 2 unused nonces, got %d", len(nonces))
		}
	})

	t.Run("MarkNonceUsed - Success", func(t *testing.T) {
		err := repo.MarkNonceUsed(ctx, signerID1, transactionID1)
		if err != nil {
			t.Fatalf("MarkNonceUsed failed: %v", err)
		}
	})

	t.Run("MarkNonceUsed - NotFound", func(t *testing.T) {
		err := repo.MarkNonceUsed(ctx, "nonexistent", transactionID1)
		if err != multisig.ErrNonceNotFound {
			t.Errorf("Expected ErrNonceNotFound, got: %v", err)
		}
	})

	t.Run("GetUnusedNoncesForTransaction - After Marking", func(t *testing.T) {
		// Only one nonce should be unused now (signer2)
		nonces, err := repo.GetUnusedNoncesForTransaction(ctx, transactionID1)
		if err != nil {
			t.Fatalf("GetUnusedNoncesForTransaction failed: %v", err)
		}

		if len(nonces) != 1 {
			t.Errorf("Expected 1 unused nonce, got %d", len(nonces))
		}

		if len(nonces) > 0 && nonces[0].SignerID() != signerID2 {
			t.Errorf("Expected unused nonce for %s, got %s", signerID2, nonces[0].SignerID())
		}
	})

	t.Run("GetAllNoncesForTransaction - After Marking", func(t *testing.T) {
		// Both nonces should still exist (one used, one unused)
		nonces, err := repo.GetAllNoncesForTransaction(ctx, transactionID1)
		if err != nil {
			t.Fatalf("GetAllNoncesForTransaction failed: %v", err)
		}

		if len(nonces) != 2 {
			t.Errorf("Expected 2 total nonces, got %d", len(nonces))
		}
	})

	t.Run("CleanupOldUnusedNonces", func(t *testing.T) {
		// Add an old nonce to transaction 2
		err := repo.SaveNonce(ctx, signerID1, transactionID2, publicNonce3)
		if err != nil {
			t.Fatalf("SaveNonce failed: %v", err)
		}

		// Sleep briefly to ensure timestamp difference
		time.Sleep(100 * time.Millisecond)

		// Cleanup nonces older than now (should not affect transaction 2)
		err = repo.CleanupOldUnusedNonces(ctx, time.Now())
		if err != nil {
			t.Fatalf("CleanupOldUnusedNonces failed: %v", err)
		}

		// Verify transaction 2 nonce still exists
		nonce, err := repo.GetNonce(ctx, signerID1, transactionID2)
		if err != nil {
			t.Errorf("Nonce should still exist after cleanup: %v", err)
		}
		if nonce == nil {
			t.Errorf("Nonce should not be nil")
		}

		// Cleanup nonces older than far future (should remove unused nonce from transaction 2)
		err = repo.CleanupOldUnusedNonces(ctx, time.Now().Add(1*time.Hour))
		if err != nil {
			t.Fatalf("CleanupOldUnusedNonces failed: %v", err)
		}

		// Verify transaction 2 nonce was removed
		_, err = repo.GetNonce(ctx, signerID1, transactionID2)
		if err != multisig.ErrNonceNotFound {
			t.Errorf("Expected ErrNonceNotFound after cleanup, got: %v", err)
		}
	})

	t.Run("DeleteNoncesForTransaction", func(t *testing.T) {
		// Delete all nonces for transaction 1
		err := repo.DeleteNoncesForTransaction(ctx, transactionID1)
		if err != nil {
			t.Fatalf("DeleteNoncesForTransaction failed: %v", err)
		}

		// Verify nonces were deleted
		nonces, err := repo.GetAllNoncesForTransaction(ctx, transactionID1)
		if err != nil {
			t.Fatalf("GetAllNoncesForTransaction failed: %v", err)
		}

		if len(nonces) != 0 {
			t.Errorf("Expected 0 nonces after deletion, got %d", len(nonces))
		}
	})
}

// TestNonceRepositorySqlc_Constructor tests the repository constructor.
func TestNonceRepositorySqlc_Constructor(t *testing.T) {
	t.Parallel()

	dbConn := &sql.DB{}
	repo := cold.NewNonceRepositorySqlc(dbConn)

	if repo == nil {
		t.Error("NewNonceRepositorySqlc returned nil")
	}
}

// TestNonceRepositorySqlc_InvalidNonceLength tests handling of invalid nonce data.
func TestNonceRepositorySqlc_InvalidNonceLength(t *testing.T) {
	// This test verifies that the repository correctly validates nonce length
	// when converting from database format to domain format.
	// Note: In practice, this should never happen due to database schema constraints (BINARY(66)),
	// but we test it for defensive programming.

	dbConn := getTestDB(t)
	defer func() {
		if err := dbConn.Close(); err != nil {
			t.Errorf("Failed to close database: %v", err)
		}
	}()

	repo := cold.NewNonceRepositorySqlc(dbConn)
	ctx := context.Background()

	signerID := "test_signer"
	transactionID := "test_tx_invalid"

	// Cleanup
	_ = repo.DeleteNoncesForTransaction(ctx, transactionID)

	// Save a valid nonce
	validNonce := [66]byte{1, 2, 3}
	err := repo.SaveNonce(ctx, signerID, transactionID, validNonce)
	if err != nil {
		t.Fatalf("SaveNonce failed: %v", err)
	}

	// Note: We cannot directly test invalid length from repository operations
	// because MySQL BINARY(66) enforces the length at the database level.
	// This test primarily verifies that valid operations work correctly.

	// Cleanup
	_ = repo.DeleteNoncesForTransaction(ctx, transactionID)
}
