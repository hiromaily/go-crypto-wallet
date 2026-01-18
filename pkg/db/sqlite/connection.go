package sqlite

import (
	"database/sql"
	"fmt"

	// Pure Go SQLite driver (CGO-free)
	// Changed from mattn/go-sqlite3 to enable cross-compilation with CGO_ENABLED=0
	// Required for GoReleaser builds (see Issue #419)
	_ "modernc.org/sqlite"

	"github.com/hiromaily/go-crypto-wallet/pkg/config"
)

// NewSQLite connects to SQLite database and returns a database connection.
func NewSQLite(conf *config.SQLite) (*sql.DB, error) {
	db, err := sql.Open("sqlite", conf.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to open SQLite database: %w", err)
	}

	// Configure connection pool
	// IMPORTANT: Do NOT set MaxOpenConns to 1.
	// With MaxOpenConns=1, a read query holds the single connection,
	// causing Begin() to block indefinitely (deadlock).
	// See: config.SQLite.MaxOpenConns for details.
	maxOpenConns := conf.MaxOpenConns
	if maxOpenConns <= 0 {
		maxOpenConns = 2 // Default: minimum 2 for concurrent read + write
	}
	db.SetMaxIdleConns(maxOpenConns)
	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(0) // No max lifetime

	// Verify connection is established and working
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Enable foreign keys (disabled by default in SQLite)
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	// Enable WAL mode for better concurrency (allows concurrent reads with single writer)
	if _, err := db.Exec("PRAGMA journal_mode = WAL"); err != nil {
		return nil, fmt.Errorf("failed to enable WAL mode: %w", err)
	}

	// Set busy timeout to wait for locks instead of immediately failing (5 seconds)
	if _, err := db.Exec("PRAGMA busy_timeout = 5000"); err != nil {
		return nil, fmt.Errorf("failed to set busy timeout: %w", err)
	}

	return db, nil
}
