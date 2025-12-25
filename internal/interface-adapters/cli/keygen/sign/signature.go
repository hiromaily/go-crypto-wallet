package sign

import (
	"context"
	"errors"
	"fmt"

	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	"github.com/hiromaily/go-crypto-wallet/internal/di"
)

func runSignature(container di.Container, filePath string) error {
	fmt.Println("sign on unsigned transaction (account would be found from file name)")

	// validator
	if filePath == "" {
		return errors.New("file path option [-file] is required")
	}

	// sign on unsigned transactions, action(deposit/payment) could be found from file name
	useCase := container.NewKeygenSignTransactionUseCase()
	output, err := useCase.Sign(context.Background(), keygenusecase.SignTransactionInput{
		FilePath: filePath,
	})
	if err != nil {
		return fmt.Errorf("fail to sign transaction: %w", err)
	}

	// TODO: output should be json if json option is true
	fmt.Printf("[isCompleted]: %t\n[fileName]: %s\n[signedCount]: %d\n[unsignedCount]: %d\n",
		output.IsDone, output.FilePath, output.SignedCount, output.UnsignedCount)

	return nil
}
