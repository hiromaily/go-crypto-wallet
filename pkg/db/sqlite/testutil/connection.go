package testutil

import (
	"database/sql"
	"log"
	"os"

	// Pure Go SQLite driver (CGO-free)
	// Changed from mattn/go-sqlite3 to enable cross-compilation with CGO_ENABLED=0
	// Required for GoReleaser builds (see Issue #419)
	_ "modernc.org/sqlite"
)

// shared database connection
var dbConn *sql.DB

// GetDB returns shared database connection for tests
func GetDB() *sql.DB {
	if dbConn != nil {
		return dbConn
	}

	dbConn = GetDBWithMaxConns(2)
	return dbConn
}

// GetDBWithMaxConns returns a database connection with specified max open connections.
// This is useful for testing different connection pool configurations.
//
// IMPORTANT: maxOpenConns=1 will cause deadlock when a read query is followed by Begin().
// See TestSQLiteTransaction_MaxOpenConns1_Deadlock for the reproduction case.
func GetDBWithMaxConns(maxOpenConns int) *sql.DB {
	// Use in-memory database for tests
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatalf("fail to create db: %v", err)
	}

	// Configure connection pool
	db.SetMaxIdleConns(maxOpenConns)
	db.SetMaxOpenConns(maxOpenConns)

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		log.Fatalf("fail to enable foreign keys: %v", err)
	}

	// Enable WAL mode for better concurrency
	if _, err := db.Exec("PRAGMA journal_mode = WAL"); err != nil {
		log.Fatalf("fail to enable WAL mode: %v", err)
	}

	// Set busy timeout
	if _, err := db.Exec("PRAGMA busy_timeout = 5000"); err != nil {
		log.Fatalf("fail to set busy timeout: %v", err)
	}

	return db
}

// GetTestDBPath returns a temporary SQLite database file path for tests
func GetTestDBPath() string {
	tmpDir := os.TempDir()
	return tmpDir + "/go-crypto-wallet-test.db"
}

// CloseDB closes the shared database connection
func CloseDB() error {
	if dbConn != nil {
		err := dbConn.Close()
		dbConn = nil
		return err
	}
	return nil
}
