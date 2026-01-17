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
	db.SetMaxIdleConns(1)    // SQLite works best with minimal concurrent connections
	db.SetMaxOpenConns(1)    // Single connection for SQLite
	db.SetConnMaxLifetime(0) // No max lifetime

	// Verify connection is established and working
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Enable foreign keys (disabled by default in SQLite)
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	return db, nil
}
