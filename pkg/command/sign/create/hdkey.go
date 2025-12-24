package create

import (
	"context"
	"fmt"

	"github.com/bookerzzz/grok"

	signusecase "github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/sign"
	"github.com/hiromaily/go-crypto-wallet/pkg/di"
)

func runHDKey(container di.Container) error {
	fmt.Println("create key for hd wallet for Authorization account")

	// create seed
	generateSeedUseCase := container.NewSignGenerateSeedUseCase()
	seedOutput, err := generateSeedUseCase.Generate(context.Background())
	if err != nil {
		return fmt.Errorf("fail to generate seed: %w", err)
	}

	// create key for hd wallet for Authorization account
	generateAuthKeyUseCase := container.NewSignGenerateAuthKeyUseCase()
	output, err := generateAuthKeyUseCase.Generate(context.Background(), signusecase.GenerateAuthKeyInput{
		AuthType: container.AuthType(),
		Seed:     seedOutput.Seed,
		Count:    1,
	})
	if err != nil {
		return fmt.Errorf("fail to generate auth key: %w", err)
	}
	grok.Value(output)

	return nil
}
