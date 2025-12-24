package shared

import (
	"context"
	"fmt"

	signusecase "github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/sign"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/service"
)

type storeSeedUseCase struct {
	seeder service.Seeder
}

// NewStoreSeedUseCase creates a new StoreSeedUseCase for sign wallet
func NewStoreSeedUseCase(seeder service.Seeder) signusecase.StoreSeedUseCase {
	return &storeSeedUseCase{
		seeder: seeder,
	}
}

func (u *storeSeedUseCase) Store(
	ctx context.Context, input signusecase.StoreSeedInput,
) (signusecase.StoreSeedOutput, error) {
	seed, err := u.seeder.Store(input.Seed)
	if err != nil {
		return signusecase.StoreSeedOutput{}, fmt.Errorf("failed to store seed: %w", err)
	}

	return signusecase.StoreSeedOutput{
		Seed: seed,
	}, nil
}
