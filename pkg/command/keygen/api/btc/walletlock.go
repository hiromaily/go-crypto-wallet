package btc

import (
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/btcgrp"
)

func runWalletLock(btc btcgrp.Bitcoiner) error {
	fmt.Println("removes the wallet encryption key from memory, locking the wallet")

	err := btc.WalletLock()
	if err != nil {
		return fmt.Errorf("fail to call WalletLock() %w", err)
	}

	fmt.Println("wallet is locked!")

	return nil
}
