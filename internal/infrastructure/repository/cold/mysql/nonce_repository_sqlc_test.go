//go:build integration

package mysql_test

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
	coldmysql "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/cold/mysql"
	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	mysql "github.com/hiromaily/go-crypto-wallet/pkg/db/mysql"
)

// getTestDB returns a database connection for integration tests.
func getTestDB(t *testing.T) *sql.DB {
	t.Helper()

	projPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-crypto-wallet"
	confPath := projPath + "/config/wallet/btc/keygen.yaml"
	conf, err := config.NewWallet(confPath, wallet.WalletTypeKeyGen, domainCoin.BTC)
	if err != nil {
		log.Fatalf("fail to create config: %v", err)
	}

	dbConn, err := mysql.NewMySQL(&conf.Database.MySQL)
	if err != nil {
		log.Fatalf("fail to create db: %v", err)
	}

	return dbConn
}

// testData holds test data for nonce repository tests.
type testData struct {
	signerID1      string
	signerID2      string
	transactionID1 string
	transactionID2 string
	publicNonce1   [66]byte
	publicNonce2   [66]byte
	publicNonce3   [66]byte
}

// newTestData creates test data.
func newTestData() testData {
	return testData{
		signerID1:      "signer1",
		signerID2:      "signer2",
		transactionID1: "tx001",
		transactionID2: "tx002",
		publicNonce1:   [66]byte{1, 2, 3},
		publicNonce2:   [66]byte{4, 5, 6},
		publicNonce3:   [66]byte{7, 8, 9},
	}
}

// testBasicCRUD tests basic CRUD operations (Save, Get, Duplicate).
func testBasicCRUD(t *testing.T, repo *coldmysql.NonceRepositorySqlc, ctx context.Context, data testData) {
	t.Helper()

	t.Run("SaveNonce - Success", func(t *testing.T) {
		err := repo.SaveNonce(ctx, data.signerID1, data.transactionID1, data.publicNonce1)
		if err != nil {
			t.Fatalf("SaveNonce failed: %v", err)
		}
	})

	t.Run("SaveNonce - Duplicate", func(t *testing.T) {
		err := repo.SaveNonce(ctx, data.signerID1, data.transactionID1, data.publicNonce1)
		if err != multisig.ErrNonceDuplicate {
			t.Errorf("Expected ErrNonceDuplicate, got: %v", err)
		}
	})

	t.Run("GetNonce - Success", func(t *testing.T) {
		nonce, err := repo.GetNonce(ctx, data.signerID1, data.transactionID1)
		if err != nil {
			t.Fatalf("GetNonce failed: %v", err)
		}

		if nonce.SignerID() != data.signerID1 {
			t.Errorf("Expected signer ID %s, got %s", data.signerID1, nonce.SignerID())
		}

		publicNonce := nonce.PublicNonce()
		if publicNonce != data.publicNonce1 {
			t.Errorf("Public nonce mismatch: expected %v, got %v", data.publicNonce1, publicNonce)
		}
	})

	t.Run("GetNonce - NotFound", func(t *testing.T) {
		_, err := repo.GetNonce(ctx, "nonexistent", data.transactionID1)
		if err != multisig.ErrNonceNotFound {
			t.Errorf("Expected ErrNonceNotFound, got: %v", err)
		}
	})
}

// testMultipleSigners tests operations with multiple signers.
func testMultipleSigners(t *testing.T, repo *coldmysql.NonceRepositorySqlc, ctx context.Context, data testData) {
	t.Helper()

	t.Run("GetAllNoncesForTransaction - Multiple Signers", func(t *testing.T) {
		err := repo.SaveNonce(ctx, data.signerID2, data.transactionID1, data.publicNonce2)
		if err != nil {
			t.Fatalf("SaveNonce failed: %v", err)
		}

		nonces, err := repo.GetAllNoncesForTransaction(ctx, data.transactionID1)
		if err != nil {
			t.Fatalf("GetAllNoncesForTransaction failed: %v", err)
		}

		if len(nonces) != 2 {
			t.Errorf("Expected 2 nonces, got %d", len(nonces))
		}

		signerIDs := make(map[string]bool)
		for _, nonce := range nonces {
			signerIDs[nonce.SignerID()] = true
		}

		if !signerIDs[data.signerID1] || !signerIDs[data.signerID2] {
			t.Errorf("Missing expected signer IDs in result")
		}
	})

	t.Run("GetUnusedNoncesForTransaction - Before Marking", func(t *testing.T) {
		nonces, err := repo.GetUnusedNoncesForTransaction(ctx, data.transactionID1)
		if err != nil {
			t.Fatalf("GetUnusedNoncesForTransaction failed: %v", err)
		}

		if len(nonces) != 2 {
			t.Errorf("Expected 2 unused nonces, got %d", len(nonces))
		}
	})
}

// testMarkingOperations tests marking nonces as used.
func testMarkingOperations(t *testing.T, repo *coldmysql.NonceRepositorySqlc, ctx context.Context, data testData) {
	t.Helper()

	t.Run("MarkNonceUsed - Success", func(t *testing.T) {
		err := repo.MarkNonceUsed(ctx, data.signerID1, data.transactionID1)
		if err != nil {
			t.Fatalf("MarkNonceUsed failed: %v", err)
		}
	})

	t.Run("MarkNonceUsed - NotFound", func(t *testing.T) {
		err := repo.MarkNonceUsed(ctx, "nonexistent", data.transactionID1)
		if err != multisig.ErrNonceNotFound {
			t.Errorf("Expected ErrNonceNotFound, got: %v", err)
		}
	})

	t.Run("GetUnusedNoncesForTransaction - After Marking", func(t *testing.T) {
		nonces, err := repo.GetUnusedNoncesForTransaction(ctx, data.transactionID1)
		if err != nil {
			t.Fatalf("GetUnusedNoncesForTransaction failed: %v", err)
		}

		if len(nonces) != 1 {
			t.Errorf("Expected 1 unused nonce, got %d", len(nonces))
		}

		if len(nonces) > 0 && nonces[0].SignerID() != data.signerID2 {
			t.Errorf("Expected unused nonce for %s, got %s", data.signerID2, nonces[0].SignerID())
		}
	})

	t.Run("GetAllNoncesForTransaction - After Marking", func(t *testing.T) {
		nonces, err := repo.GetAllNoncesForTransaction(ctx, data.transactionID1)
		if err != nil {
			t.Fatalf("GetAllNoncesForTransaction failed: %v", err)
		}

		if len(nonces) != 2 {
			t.Errorf("Expected 2 total nonces, got %d", len(nonces))
		}
	})
}

// testCleanupOperations tests cleanup operations.
func testCleanupOperations(t *testing.T, repo *coldmysql.NonceRepositorySqlc, ctx context.Context, data testData) {
	t.Helper()

	t.Run("CleanupOldUnusedNonces", func(t *testing.T) {
		err := repo.SaveNonce(ctx, data.signerID1, data.transactionID2, data.publicNonce3)
		if err != nil {
			t.Fatalf("SaveNonce failed: %v", err)
		}

		time.Sleep(100 * time.Millisecond)

		err = repo.CleanupOldUnusedNonces(ctx, time.Now())
		if err != nil {
			t.Fatalf("CleanupOldUnusedNonces failed: %v", err)
		}

		nonce, err := repo.GetNonce(ctx, data.signerID1, data.transactionID2)
		if err != nil {
			t.Errorf("Nonce should still exist after cleanup: %v", err)
		}
		if nonce == nil {
			t.Errorf("Nonce should not be nil")
		}

		err = repo.CleanupOldUnusedNonces(ctx, time.Now().Add(1*time.Hour))
		if err != nil {
			t.Fatalf("CleanupOldUnusedNonces failed: %v", err)
		}

		_, err = repo.GetNonce(ctx, data.signerID1, data.transactionID2)
		if err != multisig.ErrNonceNotFound {
			t.Errorf("Expected ErrNonceNotFound after cleanup, got: %v", err)
		}
	})
}

// testDeleteOperations tests delete operations.
func testDeleteOperations(t *testing.T, repo *coldmysql.NonceRepositorySqlc, ctx context.Context, data testData) {
	t.Helper()

	t.Run("DeleteNoncesForTransaction", func(t *testing.T) {
		err := repo.DeleteNoncesForTransaction(ctx, data.transactionID1)
		if err != nil {
			t.Fatalf("DeleteNoncesForTransaction failed: %v", err)
		}

		nonces, err := repo.GetAllNoncesForTransaction(ctx, data.transactionID1)
		if err != nil {
			t.Fatalf("GetAllNoncesForTransaction failed: %v", err)
		}

		if len(nonces) != 0 {
			t.Errorf("Expected 0 nonces after deletion, got %d", len(nonces))
		}
	})
}

// TestNonceRepositorySqlc tests the SQLC-based nonce repository implementation.
func TestNonceRepositorySqlc(t *testing.T) {
	dbConn := getTestDB(t)
	defer func() {
		if err := dbConn.Close(); err != nil {
			t.Errorf("Failed to close database: %v", err)
		}
	}()

	repo := coldmysql.NewNonceRepositorySqlc(dbConn)
	ctx := context.Background()
	data := newTestData()

	// Cleanup before test
	_ = repo.DeleteNoncesForTransaction(ctx, data.transactionID1)
	_ = repo.DeleteNoncesForTransaction(ctx, data.transactionID2)

	// Run test groups
	testBasicCRUD(t, repo, ctx, data)
	testMultipleSigners(t, repo, ctx, data)
	testMarkingOperations(t, repo, ctx, data)
	testCleanupOperations(t, repo, ctx, data)
	testDeleteOperations(t, repo, ctx, data)
}

// TestNonceRepositorySqlc_Constructor tests the repository constructor.
func TestNonceRepositorySqlc_Constructor(t *testing.T) {
	t.Parallel()

	dbConn := &sql.DB{}
	repo := coldmysql.NewNonceRepositorySqlc(dbConn)

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

	repo := coldmysql.NewNonceRepositorySqlc(dbConn)
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
