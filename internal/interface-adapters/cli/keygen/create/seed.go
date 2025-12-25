package create

import (
	"context"
	"fmt"
	"os"

	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	"github.com/hiromaily/go-crypto-wallet/internal/di"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/wallet/key"
)

func runSeed(container di.Container, seed string) error {
	fmt.Println("create seed")

	var bSeed []byte

	if seed == "" {
		seed = os.Getenv("KEYGEN_SEED")
		if seed != "" {
			fmt.Println("seed is found from environment variable")
		}
	}

	seedUseCase := container.NewKeygenGenerateSeedUseCase()

	if seed != "" {
		// store seed into database, not generate seed
		output, err := seedUseCase.Store(context.Background(), keygenusecase.StoreSeedInput{
			Seed: seed,
		})
		if err != nil {
			return fmt.Errorf("fail to store seed: %w", err)
		}
		bSeed = output.Seed
	} else {
		// create seed
		output, err := seedUseCase.Generate(context.Background())
		if err != nil {
			return fmt.Errorf("fail to generate seed: %w", err)
		}
		bSeed = output.Seed
	}
	fmt.Println("seed: " + key.SeedToString(bSeed))

	return nil
}
