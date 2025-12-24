package shared

import (
	"context"
	"fmt"

	signusecase "github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/sign"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/service"
)

type generateSeedUseCase struct {
	seeder service.Seeder
}

// NewGenerateSeedUseCase creates a new GenerateSeedUseCase for sign wallet
func NewGenerateSeedUseCase(seeder service.Seeder) signusecase.GenerateSeedUseCase {
	return &generateSeedUseCase{
		seeder: seeder,
	}
}

func (u *generateSeedUseCase) Generate(ctx context.Context) (signusecase.GenerateSeedOutput, error) {
	seed, err := u.seeder.Generate()
	if err != nil {
		return signusecase.GenerateSeedOutput{}, fmt.Errorf("failed to generate seed: %w", err)
	}

	return signusecase.GenerateSeedOutput{
		Seed: seed,
	}, nil
}
