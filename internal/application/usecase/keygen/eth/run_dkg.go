package eth

import (
	"context"
	"fmt"

	apieth "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/eth"
	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

type runDKGUseCase struct {
	keyGenerator apieth.MPCKeyGeneratorPort
}

// NewRunDKGUseCase constructs the use case for running the DKG ceremony on this node.
func NewRunDKGUseCase(keyGenerator apieth.MPCKeyGeneratorPort) keygenusecase.RunDKGUseCase {
	return &runDKGUseCase{keyGenerator: keyGenerator}
}

// Execute validates the input, delegates DKG execution to the port, and returns the joint
// Ethereum address. No shard is persisted if the port returns an error.
func (u *runDKGUseCase) Execute(
	ctx context.Context,
	input keygenusecase.RunDKGInput,
) (keygenusecase.RunDKGOutput, error) {
	// 1. Validate threshold (fail fast on zero before comparing against party count).
	if input.Threshold < 1 {
		return keygenusecase.RunDKGOutput{}, fmt.Errorf("threshold must be >= 1, got %d", input.Threshold)
	}
	if len(input.AllPartyIDs) < input.Threshold {
		return keygenusecase.RunDKGOutput{},
			fmt.Errorf("invalid DKG configuration: len(AllPartyIDs)=%d is less than Threshold=%d",
				len(input.AllPartyIDs), input.Threshold)
	}

	// 2. Build DKGParams from input.
	params := apieth.DKGParams{
		PartyID:         input.PartyID,
		AllPartyIDs:     input.AllPartyIDs,
		Threshold:       input.Threshold,
		PeerAddrs:       input.PeerAddrs,
		PreParamsPath:   input.PreParamsPath,
		ShardOutputPath: input.ShardOutputPath,
		Passphrase:      input.Passphrase,
	}

	// 3. Run DKG ceremony via the port. No shard is saved on error.
	result, err := u.keyGenerator.RunDKG(ctx, params)
	if err != nil {
		return keygenusecase.RunDKGOutput{}, fmt.Errorf("DKG ceremony failed: %w", err)
	}

	logger.Info("DKG ceremony completed",
		"ethAddress", result.EthAddress,
		"shardPath", result.ShardPath,
	)

	return keygenusecase.RunDKGOutput{
		EthAddress: result.EthAddress,
		ShardPath:  result.ShardPath,
	}, nil
}
