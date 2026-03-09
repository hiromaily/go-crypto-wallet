package eth_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	apieth "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/eth"
	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	keygenusecaseeth "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen/eth"
	ethapiamocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/eth/mocks"
)

type serveMPCDeps struct {
	shardStorage *ethapiamocks.MockMPCKeyShardStorage
	transport    *ethapiamocks.MockMPCInboundTransport
	signingNode  *ethapiamocks.MockMPCSigningNodePort
}

func newServeMPCDeps(t *testing.T) *serveMPCDeps {
	t.Helper()
	return &serveMPCDeps{
		shardStorage: ethapiamocks.NewMockMPCKeyShardStorage(t),
		transport:    ethapiamocks.NewMockMPCInboundTransport(t),
		signingNode:  ethapiamocks.NewMockMPCSigningNodePort(t),
	}
}

func newServeMPCUseCase(deps *serveMPCDeps) keygenusecase.ServeMPCUseCase {
	return keygenusecaseeth.NewServeMPCUseCase(deps.shardStorage, deps.transport, deps.signingNode)
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

func makeSessionInfo(sessionID string) apieth.MPCSigningSessionInfo {
	return apieth.MPCSigningSessionInfo{
		SessionID: sessionID,
		Hash:      make([]byte, 32),
		PartyIDs:  []string{"node-1", "node-2"},
		Threshold: 2,
	}
}

func TestServeMPC_SigningCompletes(t *testing.T) {
	t.Parallel()

	deps := newServeMPCDeps(t)

	deps.shardStorage.EXPECT().
		LoadShard(mock.Anything, dkgShardPath, dkgPassphrase).
		Return([]byte(`{"dummy":"shard"}`), nil).Once()

	deps.transport.EXPECT().
		Listen(mock.Anything, "localhost:9001").
		Return(nil).Once()

	deps.transport.EXPECT().
		AwaitSessionInfo(mock.Anything).
		Return(makeSessionInfo("session-ok"), nil).Once()

	deps.transport.EXPECT().
		Close().
		Return(nil).Once()

	deps.signingNode.EXPECT().
		RunSigning(mock.Anything, mock.MatchedBy(func(p apieth.MPCSigningNodeParams) bool {
			return p.SessionID == "session-ok" && p.PartyID == dkgPartyID
		})).
		Return(nil).Once()

	uc := newServeMPCUseCase(deps)
	err := uc.Serve(context.Background(), newServeMPCInput())
	require.NoError(t, err)
}

func TestServeMPC_AwaitSessionInfoCancelled(t *testing.T) {
	t.Parallel()

	deps := newServeMPCDeps(t)

	deps.shardStorage.EXPECT().
		LoadShard(mock.Anything, dkgShardPath, dkgPassphrase).
		Return([]byte(`{"dummy":"shard"}`), nil).Once()

	deps.transport.EXPECT().
		Listen(mock.Anything, "localhost:9001").
		Return(nil).Once()

	deps.transport.EXPECT().
		AwaitSessionInfo(mock.Anything).
		Return(apieth.MPCSigningSessionInfo{}, context.Canceled).Once()

	deps.transport.EXPECT().
		Close().
		Return(nil).Once()

	uc := newServeMPCUseCase(deps)
	err := uc.Serve(context.Background(), newServeMPCInput())

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

func TestServeMPC_SigningNodeError(t *testing.T) {
	t.Parallel()

	deps := newServeMPCDeps(t)
	signingErr := errors.New("tss round 2 failed")

	deps.shardStorage.EXPECT().
		LoadShard(mock.Anything, dkgShardPath, dkgPassphrase).
		Return([]byte(`{"dummy":"shard"}`), nil).Once()

	deps.transport.EXPECT().
		Listen(mock.Anything, "localhost:9001").
		Return(nil).Once()

	deps.transport.EXPECT().
		AwaitSessionInfo(mock.Anything).
		Return(makeSessionInfo("session-err"), nil).Once()

	deps.transport.EXPECT().
		Close().
		Return(nil).Once()

	deps.signingNode.EXPECT().
		RunSigning(mock.Anything, mock.Anything).
		Return(signingErr).Once()

	uc := newServeMPCUseCase(deps)
	err := uc.Serve(context.Background(), newServeMPCInput())

	require.Error(t, err)
	require.ErrorIs(t, err, signingErr)
}
