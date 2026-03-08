package eth_test

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	apieth "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/eth"
	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	keygenusecaseeth "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen/eth"
	ethapiamocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/eth/mocks"
	pkgeth "github.com/hiromaily/go-crypto-wallet/pkg/chains/eth"
)

type serveMPCDeps struct {
	shardStorage *ethapiamocks.MockMPCKeyShardStorage
	transport    *ethapiamocks.MockMPCInboundTransport
}

func newServeMPCDeps(t *testing.T) *serveMPCDeps {
	t.Helper()
	return &serveMPCDeps{
		shardStorage: ethapiamocks.NewMockMPCKeyShardStorage(t),
		transport:    ethapiamocks.NewMockMPCInboundTransport(t),
	}
}

func newServeMPCUseCase(deps *serveMPCDeps) keygenusecase.ServeMPCUseCase {
	return keygenusecaseeth.NewServeMPCUseCase(deps.shardStorage, deps.transport)
}

func newServeMPCInput() keygenusecase.ServeMPCInput {
	return keygenusecase.ServeMPCInput{
		ListenAddr:  "localhost:9001",
		ShardPath:   dkgShardPath,
		Passphrase:  dkgPassphrase,
		PartyID:     dkgPartyID,
		AllPartyIDs: []string{"node-1", "node-2"},
	}
}

// buildValidSessionRequest creates a well-formed session request with consistent tx hash.
func buildValidSessionRequest(t *testing.T, sessionID string) ([]byte, string) {
	t.Helper()
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   big.NewInt(31337),
		Nonce:     1,
		GasTipCap: big.NewInt(1e9),
		GasFeeCap: big.NewInt(2e9),
		Gas:       21000,
		To:        nil,
		Value:     big.NewInt(1e18),
	})
	txHex, err := pkgeth.EncodeTx(tx)
	require.NoError(t, err)

	req := apieth.MPCSessionRequest{
		SessionID: sessionID,
		TxHash:    tx.Hash().Hex(),
		RawTxHex:  *txHex,
		PartyIDs:  []string{"node-1", "node-2"},
		Threshold: 2,
	}
	data, err := json.Marshal(req)
	require.NoError(t, err)
	return data, *txHex
}

// buildMismatchedSessionRequest creates a session request where TxHash doesn't match RawTxHex.
func buildMismatchedSessionRequest(t *testing.T, sessionID string) []byte {
	t.Helper()
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID: big.NewInt(31337),
		Nonce:   1,
		Gas:     21000,
		Value:   big.NewInt(1e18),
	})
	txHex, err := pkgeth.EncodeTx(tx)
	require.NoError(t, err)

	req := apieth.MPCSessionRequest{
		SessionID: sessionID,
		TxHash:    "0x" + "ab" + "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890", // wrong hash
		RawTxHex:  *txHex,
		PartyIDs:  []string{"node-1", "node-2"},
		Threshold: 2,
	}
	data, err := json.Marshal(req)
	require.NoError(t, err)
	return data
}

func TestServeMPC_TxHashMismatch(t *testing.T) {
	t.Parallel()

	deps := newServeMPCDeps(t)

	deps.shardStorage.EXPECT().
		LoadShard(mock.Anything, dkgShardPath, dkgPassphrase).
		Return([]byte(`{"dummy":"shard"}`), nil).Once()

	deps.transport.EXPECT().
		Listen(mock.Anything, "localhost:9001").
		Return(nil).Once()

	ch := make(chan []byte, 1)
	msgBytes := buildMismatchedSessionRequest(t, "session-bad-hash")
	ch <- msgBytes
	deps.transport.EXPECT().
		Receive(mock.Anything).
		Return((<-chan []byte)(ch), nil).Once()

	deps.transport.EXPECT().
		Close().
		Return(nil).Once()

	uc := newServeMPCUseCase(deps)
	err := uc.Serve(context.Background(), newServeMPCInput())

	require.Error(t, err)
	require.ErrorIs(t, err, keygenusecaseeth.ErrTxHashMismatch)
}

func TestServeMPC_ValidHash_BlocksUntilCancel(t *testing.T) {
	t.Parallel()

	deps := newServeMPCDeps(t)

	deps.shardStorage.EXPECT().
		LoadShard(mock.Anything, dkgShardPath, dkgPassphrase).
		Return([]byte(`{"dummy":"shard"}`), nil).Once()

	deps.transport.EXPECT().
		Listen(mock.Anything, "localhost:9001").
		Return(nil).Once()

	ch := make(chan []byte, 1)
	msgBytes, _ := buildValidSessionRequest(t, "session-valid")
	ch <- msgBytes
	deps.transport.EXPECT().
		Receive(mock.Anything).
		Return((<-chan []byte)(ch), nil).Once()

	deps.transport.EXPECT().
		Close().
		Return(nil).Once()

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan error, 1)
	go func() {
		done <- newServeMPCUseCase(deps).Serve(ctx, newServeMPCInput())
	}()

	// Cancel after a short time to simulate graceful shutdown.
	cancel()

	err := <-done
	// Expect context cancellation, not ErrTxHashMismatch.
	require.Error(t, err)
	require.ErrorIs(t, err, context.Canceled)
}

func TestServeMPC_ShardLoadFailure(t *testing.T) {
	t.Parallel()

	deps := newServeMPCDeps(t)

	shardErr := errors.New("shard decryption failed: wrong passphrase")
	deps.shardStorage.EXPECT().
		LoadShard(mock.Anything, dkgShardPath, dkgPassphrase).
		Return(nil, shardErr).Once()

	uc := newServeMPCUseCase(deps)
	err := uc.Serve(context.Background(), newServeMPCInput())

	require.Error(t, err)
	require.ErrorIs(t, err, shardErr)
}
