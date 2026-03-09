package eth_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
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

// testTxData holds a valid raw transaction hex and the corresponding signing hash,
// used to satisfy the hash verification introduced to prevent blind signing.
type testTxData struct {
	rawTxHex string
	hash     []byte
}

// makeTestTxData builds an EIP-1559 transaction and returns its encoded hex
// together with the signing hash (signer.Hash) for chainID 1337.
func makeTestTxData(t *testing.T) testTxData {
	t.Helper()

	chainID := big.NewInt(1337)
	to := common.HexToAddress("0x72cCC7a7C3fa28C79aaC4f834168767A5762a7D0")
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     7,
		To:        &to,
		Value:     big.NewInt(500_000_000_000_000_000),
		Gas:       21000,
		GasTipCap: big.NewInt(2_000_000_000),
		GasFeeCap: big.NewInt(30_000_000_000),
	})

	txHex, err := pkgeth.EncodeTx(tx)
	require.NoError(t, err)
	require.NotNil(t, txHex)

	signer := types.LatestSignerForChainID(chainID)
	hash := signer.Hash(tx)
	return testTxData{rawTxHex: *txHex, hash: hash.Bytes()}
}

func makeSessionInfo(sessionID string, txd testTxData) apieth.MPCSigningSessionInfo {
	return apieth.MPCSigningSessionInfo{
		SessionID: sessionID,
		Hash:      txd.hash,
		PartyIDs:  []string{"node-1", "node-2"},
		Threshold: 2,
		RawTxHex:  txd.rawTxHex,
	}
}

func TestServeMPC_SigningCompletes(t *testing.T) {
	t.Parallel()

	txd := makeTestTxData(t)
	deps := newServeMPCDeps(t)

	deps.shardStorage.EXPECT().
		LoadShard(mock.Anything, dkgShardPath, dkgPassphrase).
		Return([]byte(`{"dummy":"shard"}`), nil).Once()

	deps.transport.EXPECT().
		Listen(mock.Anything, "localhost:9001").
		Return(nil).Once()

	deps.transport.EXPECT().
		AwaitSessionInfo(mock.Anything).
		Return(makeSessionInfo("session-ok", txd), nil).Once()

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

	txd := makeTestTxData(t)
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
		Return(makeSessionInfo("session-err", txd), nil).Once()

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

func TestServeMPC_HashMismatch(t *testing.T) {
	t.Parallel()

	txd := makeTestTxData(t)
	deps := newServeMPCDeps(t)

	deps.shardStorage.EXPECT().
		LoadShard(mock.Anything, dkgShardPath, dkgPassphrase).
		Return([]byte(`{"dummy":"shard"}`), nil).Once()

	deps.transport.EXPECT().
		Listen(mock.Anything, "localhost:9001").
		Return(nil).Once()

	// Provide a tampered hash that does not match the raw tx.
	tamperedHash := make([]byte, 32)
	tamperedHash[0] = 0xff
	sessionInfo := apieth.MPCSigningSessionInfo{
		SessionID: "session-tampered",
		Hash:      tamperedHash,
		PartyIDs:  []string{"node-1", "node-2"},
		Threshold: 2,
		RawTxHex:  txd.rawTxHex,
	}

	deps.transport.EXPECT().
		AwaitSessionInfo(mock.Anything).
		Return(sessionInfo, nil).Once()

	deps.transport.EXPECT().
		Close().
		Return(nil).Once()

	uc := newServeMPCUseCase(deps)
	err := uc.Serve(context.Background(), newServeMPCInput())

	require.Error(t, err)
	require.Contains(t, err.Error(), "tx hash verification failed")
}
