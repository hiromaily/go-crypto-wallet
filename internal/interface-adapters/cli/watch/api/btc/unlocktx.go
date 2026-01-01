package btc

import (
	"fmt"

	portsBtc "github.com/hiromaily/go-crypto-wallet/internal/application/ports/btc"
)

func runUnlockTx(btc portsBtc.Bitcoiner) error {
	// unlock locked transaction for unspent transaction
	err := btc.UnlockUnspent()
	if err != nil {
		return fmt.Errorf("fail to call BTC.UnlockUnspent() %w", err)
	}

	return nil
}
