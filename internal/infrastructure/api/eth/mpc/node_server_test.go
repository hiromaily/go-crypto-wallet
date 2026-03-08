package mpc_test

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apieth "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/eth"
	mpcinfra "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/eth/mpc"
)

// =============================================================================
// Minimal stub for MPCInboundTransport used by node server tests
// =============================================================================

// stubInboundTransport is a simple in-memory implementation of MPCInboundTransport
// for use in node server unit tests. It bypasses gRPC entirely.
type stubInboundTransport struct {
	inbound  chan []byte // coordinator → node messages (fed by tests)
	outbound chan []byte // node → coordinator messages (read by tests)
}

func newStubInboundTransport() *stubInboundTransport {
	return &stubInboundTransport{
		inbound:  make(chan []byte, 256),
		outbound: make(chan []byte, 256),
	}
}

func (*stubInboundTransport) Listen(_ context.Context, _ string) error { return nil }

func (s *stubInboundTransport) Receive(_ context.Context) (<-chan []byte, error) {
	return s.inbound, nil
}

func (s *stubInboundTransport) EnqueueOutbound(msg []byte) error {
	s.outbound <- msg
	return nil
}

func (*stubInboundTransport) Close() error { return nil }

// stubShardStorage is a minimal MPCKeyShardStorage that saves shards in memory.
type stubShardStorage struct {
	saved map[string][]byte
}

func newStubShardStorage() *stubShardStorage {
	return &stubShardStorage{saved: make(map[string][]byte)}
}

func (s *stubShardStorage) SaveShard(_ context.Context, path, _ string, data []byte) error {
	s.saved[path] = data
	return nil
}

func (s *stubShardStorage) LoadShard(_ context.Context, path, _ string) ([]byte, error) {
	data, ok := s.saved[path]
	if !ok {
		return nil, os.ErrNotExist
	}
	return data, nil
}

// =============================================================================
// GeneratePreParams tests
// =============================================================================

func TestMPCNodeServer_GeneratePreParams_ContextCancellation(t *testing.T) {
	t.Parallel()

	transport := newStubInboundTransport()
	storage := newStubShardStorage()
	server := mpcinfra.NewMPCNodeServer(transport, storage, slog.Default())

	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "pre_params.json")

	// Cancel immediately — GeneratePreParams must return an error.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := server.GeneratePreParams(ctx, outputPath)
	require.Error(t, err, "GeneratePreParams with cancelled context must return an error")
}

func TestMPCNodeServer_GeneratePreParams_CreatesJSONFile(t *testing.T) {
	if testing.Short() {
		t.Skip("GeneratePreParams takes several minutes; skipped in short mode")
	}
	// This test is excluded from normal CI; run with -run=TestMPCNodeServer_GeneratePreParams_CreatesJSONFile
	t.Parallel()

	transport := newStubInboundTransport()
	storage := newStubShardStorage()
	server := mpcinfra.NewMPCNodeServer(transport, storage, slog.Default())

	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "pre_params.json")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	err := server.GeneratePreParams(ctx, outputPath)
	require.NoError(t, err)

	data, err := os.ReadFile(outputPath)
	require.NoError(t, err)
	require.NotEmpty(t, data)

	// Must be valid JSON.
	var raw map[string]any
	require.NoError(t, json.Unmarshal(data, &raw), "pre_params.json must be valid JSON")
}

// =============================================================================
// RunDKG validation tests
// =============================================================================

func TestMPCNodeServer_RunDKG_MissingPreParamsFile(t *testing.T) {
	t.Parallel()

	transport := newStubInboundTransport()
	storage := newStubShardStorage()
	server := mpcinfra.NewMPCNodeServer(transport, storage, slog.Default())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := server.RunDKG(ctx, buildDKGParams(t, "/nonexistent/path/pre_params.json", "node-1"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "pre_params")
}

func TestMPCNodeServer_RunDKG_InvalidPreParamsJSON(t *testing.T) {
	t.Parallel()

	transport := newStubInboundTransport()
	storage := newStubShardStorage()
	server := mpcinfra.NewMPCNodeServer(transport, storage, slog.Default())

	tmpDir := t.TempDir()
	badPath := filepath.Join(tmpDir, "bad.json")
	require.NoError(t, os.WriteFile(badPath, []byte("not valid json {{{"), 0o600))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := server.RunDKG(ctx, buildDKGParams(t, badPath, "node-1"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unmarshal")
}

func TestMPCNodeServer_RunDKG_UnknownPartyID(t *testing.T) {
	t.Parallel()

	transport := newStubInboundTransport()
	storage := newStubShardStorage()
	server := mpcinfra.NewMPCNodeServer(transport, storage, slog.Default())

	// Params where PartyID is not in AllPartyIDs.
	prePath := writeFakePreParams(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	params := buildDKGParams(t, prePath, "unknown-party")

	_, err := server.RunDKG(ctx, params)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestMPCNodeServer_RunDKG_InvalidPreParams(t *testing.T) {
	t.Parallel()

	transport := newStubInboundTransport()
	storage := newStubShardStorage()
	server := mpcinfra.NewMPCNodeServer(transport, storage, slog.Default())

	// writeFakePreParams creates a file with null fields; ValidateWithProof() will return false.
	prePath := writeFakePreParams(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	params := buildDKGParams(t, prePath, "node-1")
	_, err := server.RunDKG(ctx, params)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "pre-params")
}

// =============================================================================
// Helpers
// =============================================================================

// buildDKGParams returns a minimal DKGParams for testing validation paths.
func buildDKGParams(t *testing.T, preParamsPath, partyID string) apieth.DKGParams {
	t.Helper()
	return apieth.DKGParams{
		PartyID:         partyID,
		AllPartyIDs:     []string{"node-1", "node-2"},
		Threshold:       2,
		PreParamsPath:   preParamsPath,
		ShardOutputPath: filepath.Join(t.TempDir(), "shard.enc"),
		Passphrase:      "test-passphrase",
	}
}

// writeFakePreParams writes minimal valid-looking pre-params JSON to a temp file.
// The key fields have zero values; tss-lib will reject them during DKG, which is fine
// for the context-cancellation and validation tests (they abort before using the params).
func writeFakePreParams(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "pre_params.json")
	// Minimal JSON that passes json.Unmarshal into LocalPreParams.
	raw := `{"PaillierSK":null,"NTildei":null,"H1i":null,"H2i":null,"Alpha":null,"Beta":null,"P":null,"Q":null}`
	require.NoError(t, os.WriteFile(path, []byte(raw), 0o600))
	return path
}
