package testutil

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite" // Pure Go SQLite driver (CGO-free)
)

// shared database connection
var dbConn *sql.DB

// GetDB returns shared database connection for tests
func GetDB() *sql.DB {
	if dbConn != nil {
		return dbConn
	}

	// Use in-memory database for tests
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatalf("fail to create db: %v", err)
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		log.Fatalf("fail to enable foreign keys: %v", err)
	}

	dbConn = db
	return dbConn
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
