package monitor

import (
	"context"
	"fmt"

	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	"github.com/hiromaily/go-crypto-wallet/pkg/di"
)

func runSentTx(container di.Container, _ string) error {
	// Get use case from container
	useCase := container.NewWatchMonitorTransactionUseCase().(watchusecase.MonitorTransactionUseCase)

	// monitor sent transactions
	err := useCase.UpdateTxStatus(context.Background())
	if err != nil {
		return fmt.Errorf("fail to update transaction status: %w", err)
	}

	return nil
}
