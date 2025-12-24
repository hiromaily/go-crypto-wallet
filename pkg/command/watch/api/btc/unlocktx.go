package btc

import (
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/bitcoin"
)

func runUnlockTx(btc bitcoin.Bitcoiner) error {
	// unlock locked transaction for unspent transaction
	err := btc.UnlockUnspent()
	if err != nil {
		return fmt.Errorf("fail to call BTC.UnlockUnspent() %w", err)
	}

	return nil
}
