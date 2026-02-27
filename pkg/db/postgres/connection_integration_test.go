//go:build integration

package postgres_test

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	"github.com/hiromaily/go-crypto-wallet/pkg/db/postgres"
)

// postgresTestConf builds a PostgreSQL config from environment variables.
// Defaults match the Docker Compose wallet-postgres service configuration.
func postgresTestConf() *config.PostgreSQL {
	port := 5432
	if p := os.Getenv("POSTGRES_PORT"); p != "" {
		if v, err := strconv.Atoi(p); err == nil {
			port = v
		}
	}
	return &config.PostgreSQL{
		Host:    envOrDefault("POSTGRES_HOST", "localhost"),
		Port:    port,
		DB:      envOrDefault("POSTGRES_DB", "watch"),
		User:    envOrDefault("POSTGRES_USER", "postgres"),
		Pass:    envOrDefault("POSTGRES_PASS", "postgres"),
		SSLMode: "disable",
	}
}

func envOrDefault(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

// TestNewPostgres_ValidConnection verifies that a valid configuration produces a
// working *sql.DB that can execute a Ping.
func TestNewPostgres_ValidConnection(t *testing.T) {
	conf := postgresTestConf()

	db, err := postgres.NewPostgres(conf)
	require.NoError(t, err)
	require.NotNil(t, db)
	defer func() { _ = db.Close() }()

	err = db.Ping()
	assert.NoError(t, err, "Ping should succeed after a successful NewPostgres call")
}

// TestNewPostgres_InvalidCredentials verifies that wrong credentials produce a
// "failed to ping postgres database" error.
func TestNewPostgres_InvalidCredentials(t *testing.T) {
	conf := postgresTestConf()
	conf.User = "invalid_user_xyzzy"
	conf.Pass = "wrong_password_xyzzy"

	db, err := postgres.NewPostgres(conf)
	require.Error(t, err)
	assert.Nil(t, db)
	assert.Contains(t, err.Error(), "failed to ping postgres database")
}

// TestNewPostgres_PoolConfiguration verifies the connection pool settings applied
// by NewPostgres: MaxOpenConns=10, MaxIdleConns=0, ConnMaxLifetime=0 (unlimited).
func TestNewPostgres_PoolConfiguration(t *testing.T) {
	conf := postgresTestConf()

	db, err := postgres.NewPostgres(conf)
	require.NoError(t, err)
	require.NotNil(t, db)
	defer func() { _ = db.Close() }()

	stats := db.Stats()
	assert.Equal(t, 10, stats.MaxOpenConnections, "MaxOpenConns should be 10")
	// MaxIdleConns=0 means idle connections are closed immediately; not surfaced in Stats.
	// ConnMaxLifetime=0 means no expiry; also not surfaced in Stats.
	// Verifying MaxOpenConnections is the primary observable pool configuration.
}

// TestNewPostgres_SSLModeDisable verifies that SSLMode="disable" produces a working
// connection (no TLS negotiation required).
func TestNewPostgres_SSLModeDisable(t *testing.T) {
	conf := postgresTestConf() // SSLMode defaults to "disable" in postgresTestConf

	db, err := postgres.NewPostgres(conf)
	require.NoError(t, err)
	require.NotNil(t, db)
	defer func() { _ = db.Close() }()

	assert.NoError(t, db.Ping())
}

// TestNewPostgres_SSLModeRequire verifies that SSLMode="require" produces the expected
// error when the server does not support SSL (common in local Docker Compose setup).
func TestNewPostgres_SSLModeRequire(t *testing.T) {
	conf := postgresTestConf()
	conf.SSLMode = "require"

	db, err := postgres.NewPostgres(conf)
	if err != nil {
		// Expected: server does not support SSL in default Docker Compose setup
		assert.Nil(t, db)
		assert.Contains(t, err.Error(), "failed to ping postgres database")
		return
	}
	// If the server does support SSL, the connection should work
	defer func() { _ = db.Close() }()
	assert.NoError(t, db.Ping())
}
