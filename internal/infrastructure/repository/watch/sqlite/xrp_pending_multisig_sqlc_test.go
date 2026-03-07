//go:build integration

package sqlite_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"

	domainXRP "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/xrp"
	watchsqlite "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/watch/sqlite"
)

func setupXRPPendingMultisigTable(t *testing.T, db *sql.DB) {
	t.Helper()
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS xrp_pending_multisig (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			tx_uuid TEXT NOT NULL UNIQUE,
			account_id TEXT NOT NULL,
			unsigned_tx_json TEXT NOT NULL,
			xrp_tx_type TEXT NOT NULL,
			required_quorum INTEGER NOT NULL,
			current_weight INTEGER NOT NULL DEFAULT 0,
			status TEXT NOT NULL DEFAULT 'pending',
			combined_tx_blob TEXT,
			submitted_tx_hash TEXT,
			expires_at TEXT,
			created_at TEXT,
			updated_at TEXT
		)
	`)
	require.NoError(t, err)
}

func cleanupXRPPendingMultisigTable(t *testing.T, db *sql.DB) {
	t.Helper()
	_, _ = db.Exec("DROP TABLE IF EXISTS xrp_pending_multisig")
}

func insertTestPendingMultisig(
	t *testing.T, repo *watchsqlite.XRPPendingMultisigRepositorySqlc,
	txUUID, account string,
) int64 {
	t.Helper()
	id, err := repo.Insert(context.Background(), &domainXRP.XRPPendingMultisig{
		TxUUID:         txUUID,
		AccountID:      account,
		UnsignedTxJSON: `{"TransactionType":"Payment"}`,
		XRPTxType:      "Payment",
		RequiredQuorum: 2,
		CurrentWeight:  0,
		Status:         domainXRP.MultisigStatusPending,
	})
	require.NoError(t, err)
	return id
}

func TestXRPPendingMultisigRepositorySqlc_Constructor(t *testing.T) {
	t.Parallel()

	db := &sql.DB{}
	repo := watchsqlite.NewXRPPendingMultisigRepositorySqlc(db)
	require.NotNil(t, repo)
}

func TestXRPPendingMultisigRepositorySqlc_InsertAndGetByID(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupXRPPendingMultisigTable(t, db)
	defer cleanupXRPPendingMultisigTable(t, db)

	repo := watchsqlite.NewXRPPendingMultisigRepositorySqlc(db)
	id := insertTestPendingMultisig(t, repo, "uuid-001", "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh")
	assert.Positive(t, id)

	got, err := repo.GetByID(context.Background(), id)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "uuid-001", got.TxUUID)
	assert.Equal(t, uint32(2), got.RequiredQuorum)
	assert.Equal(t, domainXRP.MultisigStatusPending, got.Status)
}

func TestXRPPendingMultisigRepositorySqlc_GetByID_NotFound(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupXRPPendingMultisigTable(t, db)
	defer cleanupXRPPendingMultisigTable(t, db)

	repo := watchsqlite.NewXRPPendingMultisigRepositorySqlc(db)

	got, err := repo.GetByID(context.Background(), 9999)
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestXRPPendingMultisigRepositorySqlc_GetByUUID(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupXRPPendingMultisigTable(t, db)
	defer cleanupXRPPendingMultisigTable(t, db)

	repo := watchsqlite.NewXRPPendingMultisigRepositorySqlc(db)
	insertTestPendingMultisig(t, repo, "uuid-abc", "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh")

	got, err := repo.GetByUUID(context.Background(), "uuid-abc")
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "uuid-abc", got.TxUUID)

	notFound, err := repo.GetByUUID(context.Background(), "uuid-nonexistent")
	require.NoError(t, err)
	assert.Nil(t, notFound)
}

func TestXRPPendingMultisigRepositorySqlc_GetByAccountIDAndStatus(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupXRPPendingMultisigTable(t, db)
	defer cleanupXRPPendingMultisigTable(t, db)

	repo := watchsqlite.NewXRPPendingMultisigRepositorySqlc(db)
	account := "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh"
	ctx := context.Background()

	insertTestPendingMultisig(t, repo, "uuid-p1", account)
	insertTestPendingMultisig(t, repo, "uuid-p2", account)

	pending, err := repo.GetByAccountIDAndStatus(ctx, account, domainXRP.MultisigStatusPending)
	require.NoError(t, err)
	assert.Len(t, pending, 2)

	ready, err := repo.GetByAccountIDAndStatus(ctx, account, domainXRP.MultisigStatusReady)
	require.NoError(t, err)
	assert.Empty(t, ready)
}

func TestXRPPendingMultisigRepositorySqlc_GetByStatus(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupXRPPendingMultisigTable(t, db)
	defer cleanupXRPPendingMultisigTable(t, db)

	repo := watchsqlite.NewXRPPendingMultisigRepositorySqlc(db)
	ctx := context.Background()

	insertTestPendingMultisig(t, repo, "uuid-s1", "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh")
	insertTestPendingMultisig(t, repo, "uuid-s2", "rPEPPER7kfTD9w2To4CQk6UCfuHM9c6GDY")

	pending, err := repo.GetByStatus(ctx, domainXRP.MultisigStatusPending)
	require.NoError(t, err)
	assert.Len(t, pending, 2)
}

func TestXRPPendingMultisigRepositorySqlc_GetExpired(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupXRPPendingMultisigTable(t, db)
	defer cleanupXRPPendingMultisigTable(t, db)

	repo := watchsqlite.NewXRPPendingMultisigRepositorySqlc(db)
	ctx := context.Background()
	past := time.Now().Add(-1 * time.Hour)

	_, err = repo.Insert(ctx, &domainXRP.XRPPendingMultisig{
		TxUUID:         "uuid-expired",
		AccountID:      "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh",
		UnsignedTxJSON: `{"TransactionType":"Payment"}`,
		XRPTxType:      "Payment",
		RequiredQuorum: 2,
		Status:         domainXRP.MultisigStatusPending,
		ExpiresAt:      &past,
	})
	require.NoError(t, err)

	expired, err := repo.GetExpired(ctx)
	require.NoError(t, err)
	assert.Len(t, expired, 1)
	assert.Equal(t, "uuid-expired", expired[0].TxUUID)
}

func TestXRPPendingMultisigRepositorySqlc_UpdateWeight(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupXRPPendingMultisigTable(t, db)
	defer cleanupXRPPendingMultisigTable(t, db)

	repo := watchsqlite.NewXRPPendingMultisigRepositorySqlc(db)
	ctx := context.Background()
	id := insertTestPendingMultisig(t, repo, "uuid-w1", "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh")

	err = repo.UpdateWeight(ctx, id, 1)
	require.NoError(t, err)

	got, err := repo.GetByID(ctx, id)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, uint32(1), got.CurrentWeight)
}

func TestXRPPendingMultisigRepositorySqlc_UpdateStatus(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupXRPPendingMultisigTable(t, db)
	defer cleanupXRPPendingMultisigTable(t, db)

	repo := watchsqlite.NewXRPPendingMultisigRepositorySqlc(db)
	ctx := context.Background()
	id := insertTestPendingMultisig(t, repo, "uuid-us1", "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh")

	err = repo.UpdateStatus(ctx, id, domainXRP.MultisigStatusReady)
	require.NoError(t, err)

	got, err := repo.GetByID(ctx, id)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, domainXRP.MultisigStatusReady, got.Status)
}

func TestXRPPendingMultisigRepositorySqlc_SetReady(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupXRPPendingMultisigTable(t, db)
	defer cleanupXRPPendingMultisigTable(t, db)

	repo := watchsqlite.NewXRPPendingMultisigRepositorySqlc(db)
	ctx := context.Background()
	id := insertTestPendingMultisig(t, repo, "uuid-r1", "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh")

	err = repo.SetReady(ctx, id, "ABCDEF1234567890COMBINED")
	require.NoError(t, err)

	got, err := repo.GetByID(ctx, id)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.NotNil(t, got.CombinedTxBlob)
	assert.Equal(t, "ABCDEF1234567890COMBINED", *got.CombinedTxBlob)
}

func TestXRPPendingMultisigRepositorySqlc_SetSubmitted(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupXRPPendingMultisigTable(t, db)
	defer cleanupXRPPendingMultisigTable(t, db)

	repo := watchsqlite.NewXRPPendingMultisigRepositorySqlc(db)
	ctx := context.Background()
	id := insertTestPendingMultisig(t, repo, "uuid-sub1", "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh")

	err = repo.SetSubmitted(ctx, id, "TXHASH0000000001")
	require.NoError(t, err)

	got, err := repo.GetByID(ctx, id)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.NotNil(t, got.SubmittedTxHash)
	assert.Equal(t, "TXHASH0000000001", *got.SubmittedTxHash)
}

func TestXRPPendingMultisigRepositorySqlc_SetConfirmed(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupXRPPendingMultisigTable(t, db)
	defer cleanupXRPPendingMultisigTable(t, db)

	repo := watchsqlite.NewXRPPendingMultisigRepositorySqlc(db)
	ctx := context.Background()
	id := insertTestPendingMultisig(t, repo, "uuid-conf1", "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh")

	err = repo.SetConfirmed(ctx, id)
	require.NoError(t, err)

	got, err := repo.GetByID(ctx, id)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, domainXRP.MultisigStatusConfirmed, got.Status)
}

func TestXRPPendingMultisigRepositorySqlc_SetFailed(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupXRPPendingMultisigTable(t, db)
	defer cleanupXRPPendingMultisigTable(t, db)

	repo := watchsqlite.NewXRPPendingMultisigRepositorySqlc(db)
	ctx := context.Background()
	id := insertTestPendingMultisig(t, repo, "uuid-fail1", "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh")

	err = repo.SetFailed(ctx, id)
	require.NoError(t, err)

	got, err := repo.GetByID(ctx, id)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, domainXRP.MultisigStatusFailed, got.Status)
}

func TestXRPPendingMultisigRepositorySqlc_ExpireAll(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	setupXRPPendingMultisigTable(t, db)
	defer cleanupXRPPendingMultisigTable(t, db)

	repo := watchsqlite.NewXRPPendingMultisigRepositorySqlc(db)
	ctx := context.Background()
	past := time.Now().Add(-1 * time.Hour)

	_, err = repo.Insert(ctx, &domainXRP.XRPPendingMultisig{
		TxUUID:         "uuid-exp-all",
		AccountID:      "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh",
		UnsignedTxJSON: `{"TransactionType":"Payment"}`,
		XRPTxType:      "Payment",
		RequiredQuorum: 2,
		Status:         domainXRP.MultisigStatusPending,
		ExpiresAt:      &past,
	})
	require.NoError(t, err)

	count, err := repo.ExpireAll(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)

	got, err := repo.GetByUUID(ctx, "uuid-exp-all")
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, domainXRP.MultisigStatusExpired, got.Status)
}
