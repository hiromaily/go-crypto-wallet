package shared

import (
	"context"
	"fmt"

	signusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/sign"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/cold"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/key"
)

type storeSeedUseCase struct {
	seedRepo cold.SeedRepositorier
}

// NewStoreSeedUseCase creates a new StoreSeedUseCase for sign wallet
func NewStoreSeedUseCase(seedRepo cold.SeedRepositorier) signusecase.StoreSeedUseCase {
	return &storeSeedUseCase{
		seedRepo: seedRepo,
	}
}

func (u *storeSeedUseCase) Store(
	ctx context.Context,
	input signusecase.StoreSeedInput,
) (signusecase.StoreSeedOutput, error) {
	// Convert seed string to bytes
	bSeed, err := key.SeedToByte(input.Seed)
	if err != nil {
		return signusecase.StoreSeedOutput{}, fmt.Errorf("fail to call key.SeedToByte(): %w", err)
	}

	// Insert seed in database
	err = u.seedRepo.Insert(input.Seed)
	if err != nil {
		return signusecase.StoreSeedOutput{}, fmt.Errorf("fail to call seedRepo.Insert(): %w", err)
	}

	return signusecase.StoreSeedOutput{
		Seed: bSeed,
	}, nil
}
