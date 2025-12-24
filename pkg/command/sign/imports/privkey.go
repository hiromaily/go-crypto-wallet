package imports

import (
	"context"
	"fmt"

	signusecase "github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/sign"
	"github.com/hiromaily/go-crypto-wallet/pkg/di"
)

func runPrivKey(container di.Container) error {
	fmt.Println("import generated private key for Authorization account to database")

	// import generated private key to database
	// Note: AuthType is passed to the factory method, not the Input
	useCase := container.NewSignImportPrivateKeyUseCase(container.AuthType())
	err := useCase.Import(context.Background(), signusecase.ImportPrivateKeyInput{})
	if err != nil {
		return fmt.Errorf("fail to import private key: %w", err)
	}
	fmt.Println("Done!")

	return nil
}
