//go:build integration

package sqlite_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"

	domainXRP "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/xrp"
	coldsqlite "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/cold/sqlite"
)

// setupXRPSignerEntryTable creates the xrp_signer_entry table for testing.
func setupXRPSignerEntryTable(t *testing.T, db *sql.DB) {
	t.Helper()
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS xrp_signer_entry (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			signer_list_id INTEGER NOT NULL,
			signer_account TEXT NOT NULL,
			signer_weight INTEGER NOT NULL,
			created_at TEXT
		)
	`)
	require.NoError(t, err)
}

// cleanupXRPSignerEntryTable drops the table after the test.
func cleanupXRPSignerEntryTable(t *testing.T, db *sql.DB) {
	t.Helper()
	_, _ = db.Exec("DROP TABLE IF EXISTS xrp_signer_entry")
}

func TestXRPSignerEntryRepositorySqlc_Constructor(t *testing.T) {
	t.Parallel()

	db := &sql.DB{}
	repo := coldsqlite.NewXRPSignerEntryRepositorySqlc(db)
	require.NotNil(t, repo)
}

func TestXRPSignerEntryRepositorySqlc_InsertAndGetByListID(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupXRPSignerEntryTable(t, db)
	defer cleanupXRPSignerEntryTable(t, db)

	repo := coldsqlite.NewXRPSignerEntryRepositorySqlc(db)
	const listID int64 = 1

	id, err := repo.Insert(context.Background(), &domainXRP.XRPSignerEntry{
		SignerListID:  listID,
		SignerAccount: "rPEPPER7kfTD9w2To4CQk6UCfuHM9c6GDY",
		SignerWeight:  1,
	})
	require.NoError(t, err)
	assert.Positive(t, id)

	entries, err := repo.GetByListID(context.Background(), listID)
	require.NoError(t, err)
	require.Len(t, entries, 1)
	assert.Equal(t, "rPEPPER7kfTD9w2To4CQk6UCfuHM9c6GDY", entries[0].SignerAccount)
	assert.Equal(t, uint32(1), entries[0].SignerWeight)
	assert.Equal(t, listID, entries[0].SignerListID)
}

func TestXRPSignerEntryRepositorySqlc_GetByID(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupXRPSignerEntryTable(t, db)
	defer cleanupXRPSignerEntryTable(t, db)

	repo := coldsqlite.NewXRPSignerEntryRepositorySqlc(db)

	id, err := repo.Insert(context.Background(), &domainXRP.XRPSignerEntry{
		SignerListID:  1,
		SignerAccount: "rPEPPER7kfTD9w2To4CQk6UCfuHM9c6GDY",
		SignerWeight:  2,
	})
	require.NoError(t, err)

	got, err := repo.GetByID(context.Background(), id)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, id, got.ID)
	assert.Equal(t, uint32(2), got.SignerWeight)
}

func TestXRPSignerEntryRepositorySqlc_GetByID_NotFound(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupXRPSignerEntryTable(t, db)
	defer cleanupXRPSignerEntryTable(t, db)

	repo := coldsqlite.NewXRPSignerEntryRepositorySqlc(db)

	got, err := repo.GetByID(context.Background(), 9999)
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestXRPSignerEntryRepositorySqlc_GetByListAndAccount(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupXRPSignerEntryTable(t, db)
	defer cleanupXRPSignerEntryTable(t, db)

	repo := coldsqlite.NewXRPSignerEntryRepositorySqlc(db)
	const listID int64 = 1

	_, err = repo.Insert(context.Background(), &domainXRP.XRPSignerEntry{
		SignerListID:  listID,
		SignerAccount: "rPEPPER7kfTD9w2To4CQk6UCfuHM9c6GDY",
		SignerWeight:  1,
	})
	require.NoError(t, err)

	got, err := repo.GetByListAndAccount(context.Background(), listID, "rPEPPER7kfTD9w2To4CQk6UCfuHM9c6GDY")
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "rPEPPER7kfTD9w2To4CQk6UCfuHM9c6GDY", got.SignerAccount)

	// Not found case
	notFound, err := repo.GetByListAndAccount(context.Background(), listID, "rNonExistent")
	require.NoError(t, err)
	assert.Nil(t, notFound)
}

func TestXRPSignerEntryRepositorySqlc_GetTotalWeight(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupXRPSignerEntryTable(t, db)
	defer cleanupXRPSignerEntryTable(t, db)

	repo := coldsqlite.NewXRPSignerEntryRepositorySqlc(db)
	const listID int64 = 1

	_, err = repo.Insert(context.Background(), &domainXRP.XRPSignerEntry{
		SignerListID:  listID,
		SignerAccount: "rPEPPER7kfTD9w2To4CQk6UCfuHM9c6GDY",
		SignerWeight:  2,
	})
	require.NoError(t, err)
	_, err = repo.Insert(context.Background(), &domainXRP.XRPSignerEntry{
		SignerListID:  listID,
		SignerAccount: "rG31cLyErnqeVj2eomEjBZtq7PYaupGYzL",
		SignerWeight:  3,
	})
	require.NoError(t, err)

	total, err := repo.GetTotalWeight(context.Background(), listID)
	require.NoError(t, err)
	assert.Equal(t, uint32(5), total)
}

func TestXRPSignerEntryRepositorySqlc_InsertBulk(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupXRPSignerEntryTable(t, db)
	defer cleanupXRPSignerEntryTable(t, db)

	repo := coldsqlite.NewXRPSignerEntryRepositorySqlc(db)
	const listID int64 = 1

	entries := []*domainXRP.XRPSignerEntry{
		{SignerListID: listID, SignerAccount: "rPEPPER7kfTD9w2To4CQk6UCfuHM9c6GDY", SignerWeight: 1},
		{SignerListID: listID, SignerAccount: "rG31cLyErnqeVj2eomEjBZtq7PYaupGYzL", SignerWeight: 1},
	}

	err = repo.InsertBulk(context.Background(), entries)
	require.NoError(t, err)

	got, err := repo.GetByListID(context.Background(), listID)
	require.NoError(t, err)
	assert.Len(t, got, 2)
}

func TestXRPSignerEntryRepositorySqlc_DeleteByListID(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupXRPSignerEntryTable(t, db)
	defer cleanupXRPSignerEntryTable(t, db)

	repo := coldsqlite.NewXRPSignerEntryRepositorySqlc(db)
	const listID int64 = 1

	entries := []*domainXRP.XRPSignerEntry{
		{SignerListID: listID, SignerAccount: "rPEPPER7kfTD9w2To4CQk6UCfuHM9c6GDY", SignerWeight: 1},
		{SignerListID: listID, SignerAccount: "rG31cLyErnqeVj2eomEjBZtq7PYaupGYzL", SignerWeight: 1},
	}
	err = repo.InsertBulk(context.Background(), entries)
	require.NoError(t, err)

	err = repo.DeleteByListID(context.Background(), listID)
	require.NoError(t, err)

	got, err := repo.GetByListID(context.Background(), listID)
	require.NoError(t, err)
	assert.Empty(t, got)
}
