package monitor

import (
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
)

func runSentTx(wallet wallets.Watcher, _ string) error {
	// monitor sent transactions
	err := wallet.UpdateTxStatus()
	if err != nil {
		return fmt.Errorf("fail to call UpdateTxStatus() %w", err)
	}

	return nil
}
