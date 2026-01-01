package btc

import (
	"fmt"

	portsBtc "github.com/hiromaily/go-crypto-wallet/internal/application/ports/btc"
)

func runWalletLock(btc portsBtc.Bitcoiner) error {
	fmt.Println("removes the wallet encryption key from memory, locking the wallet")

	err := btc.WalletLock()
	if err != nil {
		return fmt.Errorf("fail to call WalletLock() %w", err)
	}

	fmt.Println("wallet is locked!")

	return nil
}
