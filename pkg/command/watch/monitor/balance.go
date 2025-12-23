package monitor

import (
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
)

func runBalance(wallet wallets.Watcher, confirmationNum uint64) error {
	if err := wallet.MonitorBalance(confirmationNum); err != nil {
		return fmt.Errorf("fail to call MonitorBalance() %w", err)
	}

	return nil
}
