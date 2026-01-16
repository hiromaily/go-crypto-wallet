package btc

import (
	"fmt"

	portsBitcoin "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/bitcoin"
)

func runUnlockTx(btc portsBitcoin.Bitcoiner) error {
	// unlock locked transaction for unspent transaction
	err := btc.UnlockUnspent()
	if err != nil {
		return fmt.Errorf("fail to call BTC.UnlockUnspent() %w", err)
	}

	return nil
}
