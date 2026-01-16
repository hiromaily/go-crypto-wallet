package btc

import (
	"fmt"

	apibtc "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/btc"
)

func runUnlockTx(btc apibtc.Bitcoiner) error {
	// unlock locked transaction for unspent transaction
	err := btc.UnlockUnspent()
	if err != nil {
		return fmt.Errorf("fail to call BTC.UnlockUnspent() %w", err)
	}

	return nil
}
