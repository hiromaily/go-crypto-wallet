//go:build integration

package sqlite_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "modernc.org/sqlite" // Pure Go SQLite driver (CGO-free)

	coldsqlite "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/cold/sqlite"
	"github.com/hiromaily/go-crypto-wallet/pkg/chains/btc/multisig"
)

// setupNonceTable creates the musig2_nonces table for testing
func setupNonceTable(t *testing.T, db *sql.DB) {
	t.Helper()

	createTableSQL := `
		CREATE TABLE IF NOT EXISTS musig2_nonces (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			signer_id TEXT NOT NULL,
			transaction_id TEXT NOT NULL,
			public_nonce BLOB NOT NULL,
			is_used INTEGER NOT NULL DEFAULT 0,
			created_at TEXT DEFAULT CURRENT_TIMESTAMP,
			used_at TEXT,
			UNIQUE(signer_id, transaction_id)
		);
		CREATE INDEX IF NOT EXISTS idx_transaction_id ON musig2_nonces(transaction_id);
		CREATE INDEX IF NOT EXISTS idx_is_used ON musig2_nonces(is_used);
		CREATE INDEX IF NOT EXISTS idx_created_at ON musig2_nonces(created_at);
	`
	_, err := db.Exec(createTableSQL)
	if err != nil {
		t.Fatalf("failed to create musig2_nonces table: %v", err)
	}
}

// cleanupNonceTable drops the musig2_nonces table
func cleanupNonceTable(t *testing.T, db *sql.DB) {
	t.Helper()
	_, _ = db.Exec("DROP TABLE IF EXISTS musig2_nonces")
}

// testData holds test data for nonce repository tests
type testData struct {
	signerID1      string
	signerID2      string
	transactionID1 string
	transactionID2 string
	publicNonce1   [66]byte
	publicNonce2   [66]byte
	publicNonce3   [66]byte
}

// newTestData creates test data
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

// testBasicCRUD tests basic CRUD operations (Save, Get, Duplicate)
func testBasicCRUD(t *testing.T, repo *coldsqlite.NonceRepositorySqlc, ctx context.Context, data testData) {
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

// testMultipleSigners tests operations with multiple signers
func testMultipleSigners(t *testing.T, repo *coldsqlite.NonceRepositorySqlc, ctx context.Context, data testData) {
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

// testMarkingOperations tests marking nonces as used
func testMarkingOperations(t *testing.T, repo *coldsqlite.NonceRepositorySqlc, ctx context.Context, data testData) {
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

// testCleanupOperations tests cleanup operations
func testCleanupOperations(t *testing.T, repo *coldsqlite.NonceRepositorySqlc, ctx context.Context, data testData) {
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

// testDeleteOperations tests delete operations
func testDeleteOperations(t *testing.T, repo *coldsqlite.NonceRepositorySqlc, ctx context.Context, data testData) {
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

// TestNonceRepositorySqlc tests the SQLC-based nonce repository implementation
func TestNonceRepositorySqlc(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("Failed to close database: %v", err)
		}
	}()

	setupNonceTable(t, db)
	defer cleanupNonceTable(t, db)

	repo := coldsqlite.NewNonceRepositorySqlc(db)
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

// TestNonceRepositorySqlc_Constructor tests the repository constructor
func TestNonceRepositorySqlc_Constructor(t *testing.T) {
	t.Parallel()

	dbConn := &sql.DB{}
	repo := coldsqlite.NewNonceRepositorySqlc(dbConn)

	if repo == nil {
		t.Error("NewNonceRepositorySqlc returned nil")
	}
}
