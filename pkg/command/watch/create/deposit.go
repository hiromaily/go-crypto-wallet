package create

import (
	"context"
	"fmt"

	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	"github.com/hiromaily/go-crypto-wallet/pkg/di"
)

func runDeposit(container di.Container, fee float64) error {
	// Detect transaction for clients from blockchain network and create deposit unsigned transaction
	// It would be run manually on the daily basis because signature is manual task

	// Get use case from container
	useCase := container.NewWatchCreateTransactionUseCase().(watchusecase.CreateTransactionUseCase)

	output, err := useCase.Execute(context.Background(), watchusecase.CreateTransactionInput{
		ActionType:    domainTx.ActionTypeDeposit.String(),
		AdjustmentFee: fee,
	})
	if err != nil {
		return fmt.Errorf("fail to create deposit transaction: %w", err)
	}

	if output.TransactionHex == "" {
		fmt.Println("No utxo")
		return nil
	}

	// TODO: output should be json if json option is true
	fmt.Printf("[hex]: %s\n[fileName]: %s\n", output.TransactionHex, output.FileName)

	return nil
}
