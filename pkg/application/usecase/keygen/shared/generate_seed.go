package shared

import (
	"context"
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/keygen"
	sharedkeygensrv "github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/keygen/shared"
)

type generateSeedUseCase struct {
	seed *sharedkeygensrv.Seed
}

// NewGenerateSeedUseCase creates a new GenerateSeedUseCase
func NewGenerateSeedUseCase(seed *sharedkeygensrv.Seed) keygen.GenerateSeedUseCase {
	return &generateSeedUseCase{
		seed: seed,
	}
}

func (u *generateSeedUseCase) Generate(ctx context.Context) (keygen.GenerateSeedOutput, error) {
	seed, err := u.seed.Generate()
	if err != nil {
		return keygen.GenerateSeedOutput{}, fmt.Errorf("failed to generate seed: %w", err)
	}

	return keygen.GenerateSeedOutput{
		Seed: seed,
	}, nil
}

func (u *generateSeedUseCase) Store(ctx context.Context, input keygen.StoreSeedInput) (keygen.StoreSeedOutput, error) {
	seed, err := u.seed.Store(input.Seed)
	if err != nil {
		return keygen.StoreSeedOutput{}, fmt.Errorf("failed to store seed: %w", err)
	}

	return keygen.StoreSeedOutput{
		Seed: seed,
	}, nil
}
