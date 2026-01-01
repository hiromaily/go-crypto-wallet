package btc

import (
	"fmt"

	portsBitcoin "github.com/hiromaily/go-crypto-wallet/internal/application/ports/bitcoin"
)

func runWalletLock(btc portsBitcoin.Bitcoiner) error {
	fmt.Println("removes the wallet encryption key from memory, locking the wallet")

	err := btc.WalletLock()
	if err != nil {
		return fmt.Errorf("fail to call WalletLock() %w", err)
	}

	fmt.Println("wallet is locked!")

	return nil
}
