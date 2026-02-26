package shared

import (
	"context"
	"fmt"

	repocold "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository/cold"
	signusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/sign"
	btcpkg "github.com/hiromaily/go-crypto-wallet/pkg/chains/btc"
)

type storeSeedUseCase struct {
	seedRepo repocold.SeedRepositorier
}

// NewStoreSeedUseCase creates a new StoreSeedUseCase for sign wallet
func NewStoreSeedUseCase(seedRepo repocold.SeedRepositorier) signusecase.StoreSeedUseCase {
	return &storeSeedUseCase{
		seedRepo: seedRepo,
	}
}

func (u *storeSeedUseCase) Store(
	ctx context.Context,
	input signusecase.StoreSeedInput,
) (signusecase.StoreSeedOutput, error) {
	// Convert seed string to bytes
	bSeed, err := btcpkg.SeedToByte(input.Seed)
	if err != nil {
		return signusecase.StoreSeedOutput{}, fmt.Errorf("fail to call btcpkg.SeedToByte(): %w", err)
	}

	// Insert seed in database
	err = u.seedRepo.Insert(ctx, input.Seed)
	if err != nil {
		return signusecase.StoreSeedOutput{}, fmt.Errorf("fail to call seedRepo.Insert(): %w", err)
	}

	return signusecase.StoreSeedOutput{
		Seed: bSeed,
	}, nil
}
