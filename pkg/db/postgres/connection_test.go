package postgres_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	"github.com/hiromaily/go-crypto-wallet/pkg/db/postgres"
)

// TestNewPostgres_InvalidHost verifies that a connection attempt to an unreachable host
// returns a wrapped error with the expected "failed to ping postgres database" message.
func TestNewPostgres_InvalidHost(t *testing.T) {
	t.Parallel()

	conf := &config.PostgreSQL{
		Host:    "invalid-host-that-does-not-exist.example",
		Port:    5432,
		DB:      "testdb",
		User:    "testuser",
		Pass:    "testpass",
		SSLMode: "disable",
	}

	db, err := postgres.NewPostgres(conf)
	require.Error(t, err)
	assert.Nil(t, db)
	assert.Contains(t, err.Error(), "failed to ping postgres database")
}

// TestNewPostgres_DefaultPort verifies that Port=0 defaults to 5432 and does not
// produce a DSN parse error (only a ping error from the unreachable host).
func TestNewPostgres_DefaultPort(t *testing.T) {
	t.Parallel()

	conf := &config.PostgreSQL{
		Host:    "localhost",
		Port:    0, // Should default to 5432
		DB:      "nonexistent",
		User:    "testuser",
		Pass:    "testpass",
		SSLMode: "disable",
	}

	db, err := postgres.NewPostgres(conf)
	// Expect a ping error (connection refused/timeout), not a DSN parse error.
	// If we got "failed to open postgres connection" it would indicate a DSN issue.
	require.Error(t, err)
	assert.Nil(t, db)
	assert.Contains(t, err.Error(), "failed to ping postgres database")
}

// TestNewPostgres_DefaultSSLMode verifies that SSLMode="" defaults to "disable" and
// does not produce a DSN parse error.
func TestNewPostgres_DefaultSSLMode(t *testing.T) {
	t.Parallel()

	conf := &config.PostgreSQL{
		Host:    "localhost",
		Port:    5432,
		DB:      "nonexistent",
		User:    "testuser",
		Pass:    "testpass",
		SSLMode: "", // Should default to "disable"
	}

	db, err := postgres.NewPostgres(conf)
	// Expect a ping error, not a DSN error. An invalid SSL mode would produce
	// "failed to open postgres connection" instead.
	require.Error(t, err)
	assert.Nil(t, db)
	assert.Contains(t, err.Error(), "failed to ping postgres database")
}

// TestNewPostgres_ErrorReturnNilDB verifies that on error, a nil *sql.DB is returned
// (not a partially-initialized handle that could cause resource leaks).
func TestNewPostgres_ErrorReturnNilDB(t *testing.T) {
	t.Parallel()

	conf := &config.PostgreSQL{
		Host:    "invalid-host.example",
		Port:    5432,
		DB:      "testdb",
		User:    "testuser",
		Pass:    "testpass",
		SSLMode: "disable",
	}

	db, err := postgres.NewPostgres(conf)
	assert.Nil(t, db, "db must be nil when an error is returned")
	assert.Error(t, err, "error must be non-nil when connection fails")
}
