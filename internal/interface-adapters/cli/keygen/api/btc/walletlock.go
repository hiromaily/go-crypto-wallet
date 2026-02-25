package btc

import (
	"fmt"

	apibtc "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/btc"
)

func runWalletLock(btc apibtc.WalletSecurityManager) error {
	fmt.Println("removes the wallet encryption key from memory, locking the wallet")

	err := btc.WalletLock()
	if err != nil {
		return fmt.Errorf("fail to call WalletLock() %w", err)
	}

	fmt.Println("wallet is locked!")

	return nil
}
