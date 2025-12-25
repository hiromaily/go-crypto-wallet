package create

import (
	"context"
	"fmt"
	"os"

	signusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/sign"
	"github.com/hiromaily/go-crypto-wallet/internal/di"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/wallet/key"
)

func runSeed(container di.Container, seed string) error {
	fmt.Println("create seed")

	var bSeed []byte

	if seed == "" {
		seed = os.Getenv("SIGN_SEED")
		if seed != "" {
			fmt.Println("seed is found from environment variable")
		}
	}

	if seed != "" {
		// store seed into database, not generate seed
		storeSeedUseCase := container.NewSignStoreSeedUseCase()
		output, err := storeSeedUseCase.Store(context.Background(), signusecase.StoreSeedInput{
			Seed: seed,
		})
		if err != nil {
			return fmt.Errorf("fail to store seed: %w", err)
		}
		bSeed = output.Seed
	} else {
		// create seed
		generateSeedUseCase := container.NewSignGenerateSeedUseCase()
		output, err := generateSeedUseCase.Generate(context.Background())
		if err != nil {
			return fmt.Errorf("fail to generate seed: %w", err)
		}
		bSeed = output.Seed
	}
	fmt.Println("seed: " + key.SeedToString(bSeed))

	return nil
}
