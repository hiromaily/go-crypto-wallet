package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/hiromaily/go-crypto-wallet/pkg/config"
)

// NewPostgres connects to PostgreSQL server and returns a database connection.
func NewPostgres(conf *config.PostgreSQL) (*sql.DB, error) {
	port := conf.Port
	if port == 0 {
		port = 5432
	}

	sslMode := conf.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		conf.Host, port, conf.User, conf.Pass, conf.DB, sslMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open postgres connection: %w", err)
	}

	// Configure connection pool to prevent stale connections
	// Limit idle connections to force fresh connections for each process
	db.SetMaxIdleConns(0)    // No idle connections - always create fresh
	db.SetMaxOpenConns(10)   // Reasonable limit for concurrent operations
	db.SetConnMaxLifetime(0) // No max lifetime - connections live until closed

	// Verify connection is established and working
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping postgres database: %w", err)
	}

	return db, nil
}
