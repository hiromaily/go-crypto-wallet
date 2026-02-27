package postgres

import (
	"database/sql"
	"fmt"
	"net/url"

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
		// Default to "prefer": use TLS when the server supports it, but allow
		// unencrypted fallback. See pkg/config/wallet.go PostgreSQL.SSLMode docs.
		sslMode = "prefer"
	}

	// Use URL-format DSN with url.URL to safely encode user, password, and
	// database name, preventing DSN injection from values containing spaces
	// or special characters.
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(conf.User, conf.Pass),
		Host:   fmt.Sprintf("%s:%d", conf.Host, port),
		Path:   conf.DB,
	}
	q := u.Query()
	q.Set("sslmode", sslMode)
	u.RawQuery = q.Encode()

	db, err := sql.Open("postgres", u.String())
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
