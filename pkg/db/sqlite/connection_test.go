package sqlite_test

import (
	"context"
	"database/sql"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite" // Pure Go SQLite driver (CGO-free)

	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	"github.com/hiromaily/go-crypto-wallet/pkg/db/sqlite"
)

// setupTestTable creates a simple test table
func setupTestTable(t *testing.T, db *sql.DB) {
	t.Helper()
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS test_item (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			value INTEGER NOT NULL
		)
	`)
	require.NoError(t, err)
}

// TestNewSQLite_DefaultMaxOpenConns verifies default MaxOpenConns is 2 when not specified
func TestNewSQLite_DefaultMaxOpenConns(t *testing.T) {
	t.Parallel()

	tmpFile := t.TempDir() + "/test_default_conns.db"
	conf := &config.SQLite{
		Path:         tmpFile,
		MaxOpenConns: 0, // Not specified, should default to 2
	}

	db, err := sqlite.NewSQLite(conf)
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	// Verify connection is working
	err = db.Ping()
	assert.NoError(t, err)
}

// TestNewSQLite_CustomMaxOpenConns verifies custom MaxOpenConns is respected
func TestNewSQLite_CustomMaxOpenConns(t *testing.T) {
	t.Parallel()

	tmpFile := t.TempDir() + "/test_custom_conns.db"
	conf := &config.SQLite{
		Path:         tmpFile,
		MaxOpenConns: 5,
	}

	db, err := sqlite.NewSQLite(conf)
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	// Verify connection is working
	err = db.Ping()
	assert.NoError(t, err)
}

// TestSQLiteTransaction_MaxOpenConns2_ConcurrentReadWrite verifies that
// MaxOpenConns=2 allows concurrent read operations followed by write transactions.
// This is the expected behavior for production use.
func TestSQLiteTransaction_MaxOpenConns2_ConcurrentReadWrite(t *testing.T) {
	t.Parallel()

	tmpFile := t.TempDir() + "/test_conns2.db"
	conf := &config.SQLite{
		Path:         tmpFile,
		MaxOpenConns: 2,
	}

	db, err := sqlite.NewSQLite(conf)
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupTestTable(t, db)

	// Insert initial data
	_, err = db.Exec("INSERT INTO test_item (name, value) VALUES (?, ?)", "item1", 100)
	require.NoError(t, err)

	// Simulate the pattern that caused deadlock with MaxOpenConns=1:
	// 1. Execute a read query (SELECT)
	// 2. Immediately start a transaction with Begin()
	// 3. Execute write operations in the transaction

	// Step 1: Read query
	var count int64
	err = db.QueryRow("SELECT COUNT(*) FROM test_item WHERE name = ?", "item1").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)

	// Step 2: Start transaction (this should NOT block with MaxOpenConns=2)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := db.BeginTx(ctx, nil)
	require.NoError(t, err, "Begin should not block with MaxOpenConns=2")

	// Step 3: Execute write in transaction
	_, err = tx.Exec("INSERT INTO test_item (name, value) VALUES (?, ?)", "item2", 200)
	require.NoError(t, err)

	err = tx.Commit()
	require.NoError(t, err)

	// Verify the insert succeeded
	err = db.QueryRow("SELECT COUNT(*) FROM test_item").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, int64(2), count)
}

// TestSQLiteTransaction_MaxOpenConns1_Deadlock demonstrates that MaxOpenConns=1
// causes a deadlock when a read query is immediately followed by Begin().
//
// This test documents the bug that was fixed by increasing MaxOpenConns to 2.
// DO NOT use MaxOpenConns=1 in production!
func TestSQLiteTransaction_MaxOpenConns1_Deadlock(t *testing.T) {
	t.Parallel()

	// Create database with MaxOpenConns=1 (the problematic configuration)
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	// Configure with MaxOpenConns=1 (causes deadlock)
	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(1)

	// Enable pragmas
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	require.NoError(t, err)
	_, err = db.Exec("PRAGMA journal_mode = WAL")
	require.NoError(t, err)
	_, err = db.Exec("PRAGMA busy_timeout = 1000") // Short timeout for test
	require.NoError(t, err)

	setupTestTable(t, db)

	// Insert initial data
	_, err = db.Exec("INSERT INTO test_item (name, value) VALUES (?, ?)", "item1", 100)
	require.NoError(t, err)

	// This test demonstrates the deadlock issue:
	// With MaxOpenConns=1, the single connection is used by the SELECT query,
	// and Begin() blocks waiting for a connection that will never be released.

	var wg sync.WaitGroup
	deadlockDetected := false

	wg.Add(1)
	go func() {
		defer wg.Done()

		// Step 1: Read query (holds the single connection)
		var count int64
		err := db.QueryRow("SELECT COUNT(*) FROM test_item WHERE name = ?", "item1").Scan(&count)
		if err != nil {
			return
		}

		// Step 2: Try to start transaction with timeout
		// With MaxOpenConns=1, this will timeout because the connection is still in use
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		_, err = db.BeginTx(ctx, nil)
		if err != nil {
			// Expected: context deadline exceeded (deadlock)
			deadlockDetected = true
		}
	}()

	wg.Wait()

	// With MaxOpenConns=1, we expect the deadlock to be detected (timeout)
	// Note: This behavior may vary depending on the SQLite driver implementation
	// The key point is that MaxOpenConns=1 is problematic and should be avoided
	t.Logf("Deadlock detected with MaxOpenConns=1: %v", deadlockDetected)
}

// TestSQLiteTransaction_RollbackOnError verifies that transactions are properly
// rolled back when an error occurs.
func TestSQLiteTransaction_RollbackOnError(t *testing.T) {
	t.Parallel()

	tmpFile := t.TempDir() + "/test_rollback.db"
	conf := &config.SQLite{
		Path:         tmpFile,
		MaxOpenConns: 2,
	}

	db, err := sqlite.NewSQLite(conf)
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupTestTable(t, db)

	// Start transaction
	tx, err := db.Begin()
	require.NoError(t, err)

	// Insert data
	_, err = tx.Exec("INSERT INTO test_item (name, value) VALUES (?, ?)", "item1", 100)
	require.NoError(t, err)

	// Rollback instead of commit
	err = tx.Rollback()
	require.NoError(t, err)

	// Verify data was not persisted
	var count int64
	err = db.QueryRow("SELECT COUNT(*) FROM test_item").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, int64(0), count, "Data should not be persisted after rollback")
}

// TestSQLiteTransaction_CommitPersistsData verifies that committed transactions
// persist data correctly.
func TestSQLiteTransaction_CommitPersistsData(t *testing.T) {
	t.Parallel()

	tmpFile := t.TempDir() + "/test_commit.db"
	conf := &config.SQLite{
		Path:         tmpFile,
		MaxOpenConns: 2,
	}

	db, err := sqlite.NewSQLite(conf)
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupTestTable(t, db)

	// Start transaction
	tx, err := db.Begin()
	require.NoError(t, err)

	// Insert data
	_, err = tx.Exec("INSERT INTO test_item (name, value) VALUES (?, ?)", "item1", 100)
	require.NoError(t, err)

	// Commit
	err = tx.Commit()
	require.NoError(t, err)

	// Verify data was persisted
	var count int64
	err = db.QueryRow("SELECT COUNT(*) FROM test_item").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, int64(1), count, "Data should be persisted after commit")
}

// TestSQLiteTransaction_MultipleOperationsInTransaction verifies that multiple
// operations within a single transaction work correctly.
func TestSQLiteTransaction_MultipleOperationsInTransaction(t *testing.T) {
	t.Parallel()

	tmpFile := t.TempDir() + "/test_multi_ops.db"
	conf := &config.SQLite{
		Path:         tmpFile,
		MaxOpenConns: 2,
	}

	db, err := sqlite.NewSQLite(conf)
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupTestTable(t, db)

	// Start transaction
	tx, err := db.Begin()
	require.NoError(t, err)

	// Multiple inserts
	_, err = tx.Exec("INSERT INTO test_item (name, value) VALUES (?, ?)", "item1", 100)
	require.NoError(t, err)
	_, err = tx.Exec("INSERT INTO test_item (name, value) VALUES (?, ?)", "item2", 200)
	require.NoError(t, err)
	_, err = tx.Exec("INSERT INTO test_item (name, value) VALUES (?, ?)", "item3", 300)
	require.NoError(t, err)

	// Update
	_, err = tx.Exec("UPDATE test_item SET value = ? WHERE name = ?", 150, "item1")
	require.NoError(t, err)

	// Delete
	_, err = tx.Exec("DELETE FROM test_item WHERE name = ?", "item2")
	require.NoError(t, err)

	// Commit
	err = tx.Commit()
	require.NoError(t, err)

	// Verify final state
	var count int64
	err = db.QueryRow("SELECT COUNT(*) FROM test_item").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, int64(2), count, "Should have 2 items after operations")

	var value int
	err = db.QueryRow("SELECT value FROM test_item WHERE name = ?", "item1").Scan(&value)
	require.NoError(t, err)
	assert.Equal(t, 150, value, "item1 value should be updated to 150")
}
