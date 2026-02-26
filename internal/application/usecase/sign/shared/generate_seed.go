package shared

import (
	"context"
	"errors"
	"fmt"

	repocold "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository/cold"
	signusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/sign"
	btcpkg "github.com/hiromaily/go-crypto-wallet/pkg/chains/btc"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

type generateSeedUseCase struct {
	seedRepo repocold.SeedRepositorier
}

// NewGenerateSeedUseCase creates a new GenerateSeedUseCase for sign wallet
func NewGenerateSeedUseCase(seedRepo repocold.SeedRepositorier) signusecase.GenerateSeedUseCase {
	return &generateSeedUseCase{
		seedRepo: seedRepo,
	}
}

func (u *generateSeedUseCase) Generate(ctx context.Context) (signusecase.GenerateSeedOutput, error) {
	// Try to retrieve existing seed from database
	bSeed, err := u.retrieveSeed(ctx)
	if err == nil {
		return signusecase.GenerateSeedOutput{
			Seed: bSeed,
		}, nil
	}

	// Generate new seed if not found
	bSeed, err = btcpkg.GenerateSeed()
	if err != nil {
		return signusecase.GenerateSeedOutput{}, fmt.Errorf("fail to call btcpkg.GenerateSeed(): %w", err)
	}
	strSeed := btcpkg.SeedToString(bSeed)

	// Insert seed in database
	err = u.seedRepo.Insert(ctx, strSeed)
	if err != nil {
		return signusecase.GenerateSeedOutput{}, fmt.Errorf("fail to call seedRepo.Insert(): %w", err)
	}

	return signusecase.GenerateSeedOutput{
		Seed: bSeed,
	}, nil
}

// retrieveSeed retrieves seed from database
func (u *generateSeedUseCase) retrieveSeed(ctx context.Context) ([]byte, error) {
	// Get seed from database, seed is expected to have only one record
	seed, err := u.seedRepo.GetOne(ctx)
	if err == nil && seed.Seed != "" {
		logger.Info("seed have already been generated")
		return btcpkg.SeedToByte(seed.Seed)
	}
	if err != nil {
		return nil, fmt.Errorf("fail to call seedRepo.GetOne(): %w", err)
	}
	// In this case, though err didn't happen, but seed is blank
	return nil, errors.New("somehow seed retrieved from database is blank")
}
