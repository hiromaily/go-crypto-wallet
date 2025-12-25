package imports

import (
	"context"
	"errors"
	"fmt"

	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	"github.com/hiromaily/go-crypto-wallet/internal/di"
)

func runAddress(container di.Container, filePath string, isRescan bool) error {
	fmt.Println("-file: " + filePath)

	// validator
	if filePath == "" {
		return errors.New("file path option [-file] is required")
	}

	// Get use case from container
	useCase := container.NewWatchImportAddressUseCase()

	// import public address
	err := useCase.Execute(context.Background(), watchusecase.ImportAddressInput{
		FileName: filePath,
		Rescan:   isRescan,
	})
	if err != nil {
		return fmt.Errorf("fail to import address: %w", err)
	}
	fmt.Println("Done!")

	return nil
}
