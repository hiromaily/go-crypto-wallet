package monitor

import (
	"context"
	"fmt"

	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	"github.com/hiromaily/go-crypto-wallet/pkg/di"
)

func runBalance(container di.Container, confirmationNum uint64) error {
	// Get use case from container
	useCase := container.NewWatchMonitorTransactionUseCase().(watchusecase.MonitorTransactionUseCase)

	if err := useCase.MonitorBalance(context.Background(), watchusecase.MonitorBalanceInput{
		ConfirmationNum: confirmationNum,
	}); err != nil {
		return fmt.Errorf("fail to monitor balance: %w", err)
	}

	return nil
}
