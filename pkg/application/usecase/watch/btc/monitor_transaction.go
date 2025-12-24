package btc

import (
	"context"
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/watch"
	btcwatchsrv "github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/watch/btc"
)

type monitorTransactionUseCase struct {
	txMonitor *btcwatchsrv.TxMonitor
}

// NewMonitorTransactionUseCase creates a new MonitorTransactionUseCase
func NewMonitorTransactionUseCase(txMonitor *btcwatchsrv.TxMonitor) watch.MonitorTransactionUseCase {
	return &monitorTransactionUseCase{
		txMonitor: txMonitor,
	}
}

func (u *monitorTransactionUseCase) UpdateTxStatus(ctx context.Context) error {
	if err := u.txMonitor.UpdateTxStatus(); err != nil {
		return fmt.Errorf("failed to update tx status: %w", err)
	}
	return nil
}

func (u *monitorTransactionUseCase) MonitorBalance(ctx context.Context, input watch.MonitorBalanceInput) error {
	if err := u.txMonitor.MonitorBalance(input.ConfirmationNum); err != nil {
		return fmt.Errorf("failed to monitor balance: %w", err)
	}
	return nil
}
