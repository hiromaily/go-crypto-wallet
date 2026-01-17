//go:build integration

package sqlite_test

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	coldsqlite "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/cold/sqlite"
)

// setupSeedTable creates the seed table for testing
func setupSeedTable(t *testing.T, db *sql.DB) {
	t.Helper()

	createTableSQL := `
		CREATE TABLE IF NOT EXISTS seed (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			coin TEXT NOT NULL,
			seed TEXT NOT NULL,
			updated_at TEXT DEFAULT CURRENT_TIMESTAMP
		);
	`
	_, err := db.Exec(createTableSQL)
	require.NoError(t, err, "failed to create seed table")
}

// cleanupSeedTable drops the seed table
func cleanupSeedTable(t *testing.T, db *sql.DB) {
	t.Helper()
	_, _ = db.Exec("DROP TABLE IF EXISTS seed")
}

// TestSeedRepositorySqlc_Insert tests inserting a seed
func TestSeedRepositorySqlc_Insert(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupSeedTable(t, db)
	defer cleanupSeedTable(t, db)

	repo := coldsqlite.NewSeedRepositorySqlc(db, domainCoin.BTC)
	ctx := context.Background()

	testSeed := "test_seed_data_encrypted"

	err = repo.Insert(ctx, testSeed)
	require.NoError(t, err, "Insert should succeed")
}

// TestSeedRepositorySqlc_GetOne tests retrieving a seed
func TestSeedRepositorySqlc_GetOne(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupSeedTable(t, db)
	defer cleanupSeedTable(t, db)

	repo := coldsqlite.NewSeedRepositorySqlc(db, domainCoin.BTC)
	ctx := context.Background()

	testSeed := "test_seed_data_encrypted"

	// Insert first
	err = repo.Insert(ctx, testSeed)
	require.NoError(t, err)

	// Retrieve
	seed, err := repo.GetOne(ctx)
	require.NoError(t, err, "GetOne should succeed")
	require.NotNil(t, seed, "seed should not be nil")

	assert.Equal(t, domainCoin.BTC, seed.CoinTypeCode, "coin type should match")
	assert.Equal(t, testSeed, seed.Seed, "seed data should match")
}

// TestSeedRepositorySqlc_GetOne_NotFound tests retrieving when no seed exists
func TestSeedRepositorySqlc_GetOne_NotFound(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupSeedTable(t, db)
	defer cleanupSeedTable(t, db)

	repo := coldsqlite.NewSeedRepositorySqlc(db, domainCoin.BTC)
	ctx := context.Background()

	_, err = repo.GetOne(ctx)
	require.Error(t, err, "GetOne should fail when no seed exists")
}

// TestSeedRepositorySqlc_Constructor tests the repository constructor
func TestSeedRepositorySqlc_Constructor(t *testing.T) {
	t.Parallel()

	db := &sql.DB{}
	repo := coldsqlite.NewSeedRepositorySqlc(db, domainCoin.BTC)

	require.NotNil(t, repo, "NewSeedRepositorySqlc should not return nil")
}
