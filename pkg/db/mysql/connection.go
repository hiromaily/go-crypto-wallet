package mysql

import (
	"database/sql"
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/pkg/config"
)

// NewMySQL connects to MySQL server and returns a database connection.
// TODO: retry functionality and retry count should be configured in config file
func NewMySQL(conf *config.MySQL) (*sql.DB, error) {
	db, err := sql.Open("mysql",
		fmt.Sprintf(
			"%s:%s@tcp(%s)/%s?parseTime=true&charset=utf8mb4",
			conf.User,
			conf.Pass,
			conf.Host,
			conf.DB))
	if err != nil {
		return nil, fmt.Errorf("Connection(): error: %v", err)
	}

	// Configure connection pool to prevent stale connections
	// Limit idle connections to force fresh connections for each process
	db.SetMaxIdleConns(0) // No idle connections - always create fresh
	db.SetMaxOpenConns(10) // Reasonable limit for concurrent operations
	db.SetConnMaxLifetime(0) // No max lifetime - connections live until closed

	// Verify connection is established and working
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return db, nil
}
