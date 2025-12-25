package shared

import (
	"context"
	"errors"
	"fmt"

	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/cold"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/wallet/key"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

type generateSeedUseCase struct {
	seedRepo cold.SeedRepositorier
}

// NewGenerateSeedUseCase creates a new GenerateSeedUseCase
func NewGenerateSeedUseCase(seedRepo cold.SeedRepositorier) keygenusecase.GenerateSeedUseCase {
	return &generateSeedUseCase{
		seedRepo: seedRepo,
	}
}

func (u *generateSeedUseCase) Generate(ctx context.Context) (keygenusecase.GenerateSeedOutput, error) {
	// Try to retrieve existing seed from database
	bSeed, err := u.retrieveSeed()
	if err == nil {
		return keygenusecase.GenerateSeedOutput{
			Seed: bSeed,
		}, nil
	}

	// Generate new seed if not found
	bSeed, err = key.GenerateSeed()
	if err != nil {
		return keygenusecase.GenerateSeedOutput{}, fmt.Errorf("fail to call key.GenerateSeed(): %w", err)
	}
	strSeed := key.SeedToString(bSeed)

	// Insert seed in database
	err = u.seedRepo.Insert(strSeed)
	if err != nil {
		return keygenusecase.GenerateSeedOutput{}, fmt.Errorf("fail to call seedRepo.Insert(): %w", err)
	}

	return keygenusecase.GenerateSeedOutput{
		Seed: bSeed,
	}, nil
}

func (u *generateSeedUseCase) Store(
	ctx context.Context,
	input keygenusecase.StoreSeedInput,
) (keygenusecase.StoreSeedOutput, error) {
	// Convert seed string to bytes
	bSeed, err := key.SeedToByte(input.Seed)
	if err != nil {
		return keygenusecase.StoreSeedOutput{}, fmt.Errorf("fail to call key.SeedToByte(): %w", err)
	}

	// Insert seed in database
	err = u.seedRepo.Insert(input.Seed)
	if err != nil {
		return keygenusecase.StoreSeedOutput{}, fmt.Errorf("fail to call seedRepo.Insert(): %w", err)
	}

	return keygenusecase.StoreSeedOutput{
		Seed: bSeed,
	}, nil
}

// retrieveSeed retrieves seed from database
func (u *generateSeedUseCase) retrieveSeed() ([]byte, error) {
	// Get seed from database, seed is expected to have only one record
	seed, err := u.seedRepo.GetOne()
	if err == nil && seed.Seed != "" {
		logger.Info("seed have already been generated")
		return key.SeedToByte(seed.Seed)
	}
	if err != nil {
		return nil, fmt.Errorf("fail to call seedRepo.GetOne(): %w", err)
	}
	// In this case, though err didn't happen, but seed is blank
	return nil, errors.New("somehow seed retrieved from database is blank")
}
