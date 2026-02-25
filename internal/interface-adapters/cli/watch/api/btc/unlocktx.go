package btc

import (
	"fmt"
)

func runUnlockTx(btc btcWatchAPICmds) error {
	// unlock locked transaction for unspent transaction
	err := btc.UnlockUnspent()
	if err != nil {
		return fmt.Errorf("fail to call BTC.UnlockUnspent() %w", err)
	}

	return nil
}
