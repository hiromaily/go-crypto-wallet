package eth_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	apieth "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/eth"
	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	keygenusecaseeth "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen/eth"
	ethapiamocks "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/eth/mocks"
)

const (
	dkgPartyID    = "node-1"
	dkgShardPath  = "/tmp/shard.json"
	dkgPassphrase = "test-passphrase"
)

type runDKGDeps struct {
	keyGenerator *ethapiamocks.MockMPCKeyGeneratorPort
}

func newRunDKGDeps(t *testing.T) *runDKGDeps {
	t.Helper()
	return &runDKGDeps{
		keyGenerator: ethapiamocks.NewMockMPCKeyGeneratorPort(t),
	}
}

func newRunDKGUseCase(deps *runDKGDeps) keygenusecase.RunDKGUseCase {
	return keygenusecaseeth.NewRunDKGUseCase(deps.keyGenerator)
}

func TestRunDKG_Success(t *testing.T) {
	t.Parallel()

	deps := newRunDKGDeps(t)
	deps.keyGenerator.EXPECT().
		RunDKG(mock.Anything, mock.MatchedBy(func(p apieth.DKGParams) bool {
			return p.PartyID == dkgPartyID && p.Threshold == 2 && len(p.AllPartyIDs) == 3
		})).
		Return(&apieth.DKGResult{
			EthAddress: "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
			ShardPath:  dkgShardPath,
		}, nil).Once()

	uc := newRunDKGUseCase(deps)
	out, err := uc.Execute(context.Background(), keygenusecase.RunDKGInput{
		PartyID:         dkgPartyID,
		AllPartyIDs:     []string{"node-1", "node-2", "node-3"},
		Threshold:       2,
		PeerAddrs:       []string{"localhost:9001", "localhost:9002"},
		PreParamsPath:   "/tmp/pre_params.json",
		ShardOutputPath: dkgShardPath,
		Passphrase:      dkgPassphrase,
	})

	require.NoError(t, err)
	assert.Equal(t, "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266", out.EthAddress)
	assert.Equal(t, dkgShardPath, out.ShardPath)
}

func TestRunDKG_InvalidThreshold_BelowPartyCount(t *testing.T) {
	t.Parallel()

	deps := newRunDKGDeps(t)
	// Port should NOT be called when validation fails.
	uc := newRunDKGUseCase(deps)

	_, err := uc.Execute(context.Background(), keygenusecase.RunDKGInput{
		PartyID:     dkgPartyID,
		AllPartyIDs: []string{"node-1"}, // only 1 party
		Threshold:   3,                  // but threshold is 3
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "len(AllPartyIDs)=1 is less than Threshold=3")
}

func TestRunDKG_ZeroThreshold(t *testing.T) {
	t.Parallel()

	deps := newRunDKGDeps(t)
	uc := newRunDKGUseCase(deps)

	_, err := uc.Execute(context.Background(), keygenusecase.RunDKGInput{
		PartyID:     dkgPartyID,
		AllPartyIDs: []string{"node-1", "node-2"},
		Threshold:   0,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "threshold must be >= 1")
}

func TestRunDKG_PortError_NoShardSaved(t *testing.T) {
	t.Parallel()

	deps := newRunDKGDeps(t)
	portErr := errors.New("DKG round 3 timeout")
	deps.keyGenerator.EXPECT().
		RunDKG(mock.Anything, mock.Anything).
		Return(nil, portErr).Once()

	uc := newRunDKGUseCase(deps)
	_, err := uc.Execute(context.Background(), keygenusecase.RunDKGInput{
		PartyID:     dkgPartyID,
		AllPartyIDs: []string{"node-1", "node-2"},
		Threshold:   2,
	})

	require.Error(t, err)
	assert.ErrorIs(t, err, portErr)
}
