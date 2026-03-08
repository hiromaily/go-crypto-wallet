package mpc_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mpcstore "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/mpc"
)

func TestMPCShardStore_SaveAndLoad_RoundTrip(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	store := mpcstore.NewMPCShardStore()

	tempDir := t.TempDir()
	shardPath := filepath.Join(tempDir, "shard.enc")
	passphrase := "test-passphrase-correct"
	// Sample shard JSON payload (mimics LocalPartySaveData wrapper format).
	originalData := []byte(
		`{"version":1,"party_id":"node-1","all_party_ids":["node-1","node-2"],` +
			`"threshold":2,"eth_address":"0x0742D35CC6634c0532925A3b844bc9E7595f0Beb"}`,
	)

	err := store.SaveShard(ctx, shardPath, passphrase, originalData)
	require.NoError(t, err)

	// File must be created.
	info, err := os.Stat(shardPath)
	require.NoError(t, err)
	assert.Greater(t, info.Size(), int64(0))

	// Round-trip: loaded data must match original.
	loaded, err := store.LoadShard(ctx, shardPath, passphrase)
	require.NoError(t, err)
	assert.Equal(t, originalData, loaded)
}

func TestMPCShardStore_SaveAndLoad_DifferentData(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	store := mpcstore.NewMPCShardStore()

	tempDir := t.TempDir()
	shardPath := filepath.Join(tempDir, "shard2.enc")
	passphrase := "another-passphrase"
	data := []byte(`{"version":1,"party_id":"node-2","threshold":3}`)

	require.NoError(t, store.SaveShard(ctx, shardPath, passphrase, data))
	loaded, err := store.LoadShard(ctx, shardPath, passphrase)
	require.NoError(t, err)
	assert.Equal(t, data, loaded)
}

func TestMPCShardStore_Load_WrongPassphrase(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	store := mpcstore.NewMPCShardStore()

	tempDir := t.TempDir()
	shardPath := filepath.Join(tempDir, "shard.enc")
	data := []byte(`{"version":1}`)

	require.NoError(t, store.SaveShard(ctx, shardPath, "correct-passphrase", data))

	_, err := store.LoadShard(ctx, shardPath, "wrong-passphrase")
	require.Error(t, err)
}

func TestMPCShardStore_Load_MissingFile(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	store := mpcstore.NewMPCShardStore()

	_, err := store.LoadShard(ctx, "/nonexistent/path/shard.enc", "passphrase")
	require.Error(t, err)
}

func TestMPCShardStore_Save_CreatesParentDirectory(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	store := mpcstore.NewMPCShardStore()

	tempDir := t.TempDir()
	shardPath := filepath.Join(tempDir, "subdir", "nested", "shard.enc")
	data := []byte(`{"version":1}`)

	err := store.SaveShard(ctx, shardPath, "passphrase", data)
	require.NoError(t, err)

	_, err = os.Stat(shardPath)
	require.NoError(t, err)
}

func TestMPCShardStore_Load_CorruptedFile(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	store := mpcstore.NewMPCShardStore()

	tempDir := t.TempDir()
	shardPath := filepath.Join(tempDir, "corrupt.enc")

	// Write garbage bytes that are long enough but not a valid ciphertext.
	err := os.WriteFile(shardPath, []byte("not-a-valid-encrypted-shard-file-at-all-xxxxxxxxxxxx"), 0o600)
	require.NoError(t, err)

	_, err = store.LoadShard(ctx, shardPath, "passphrase")
	require.Error(t, err)
}
