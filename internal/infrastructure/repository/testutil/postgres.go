// Package testutil provides shared helpers for PostgreSQL repository integration tests.
package testutil

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	pkgpostgres "github.com/hiromaily/go-crypto-wallet/pkg/db/postgres"
)

// EnvOrDefault returns the value of the environment variable key, or defaultVal if unset.
func EnvOrDefault(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

// NewTestPostgresConf builds a PostgreSQL config from environment variables.
// dbEnvKey is the env var for the DB name (e.g. "POSTGRES_DB"); defaultDB is the fallback.
// Host, port, user, and password are read from the standard POSTGRES_* env vars.
func NewTestPostgresConf(dbEnvKey, defaultDB string) *config.PostgreSQL {
	port := 5432
	if p := os.Getenv("POSTGRES_PORT"); p != "" {
		if v, err := strconv.Atoi(p); err == nil {
			port = v
		}
	}
	return &config.PostgreSQL{
		Host:    EnvOrDefault("POSTGRES_HOST", "localhost"),
		Port:    port,
		DB:      EnvOrDefault(dbEnvKey, defaultDB),
		User:    EnvOrDefault("POSTGRES_USER", "postgres"),
		Pass:    EnvOrDefault("POSTGRES_PASS", "postgres"),
		SSLMode: "disable",
	}
}

// OpenTestDB opens a PostgreSQL connection for testing and registers db.Close with t.Cleanup.
func OpenTestDB(t *testing.T, conf *config.PostgreSQL) *sql.DB {
	t.Helper()
	db, err := pkgpostgres.NewPostgres(conf)
	require.NoError(t, err, "failed to connect to PostgreSQL database")
	t.Cleanup(func() { _ = db.Close() })
	return db
}

// UniqueStr generates a unique string with a nano-timestamp suffix to avoid test collisions.
func UniqueStr(prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
}
