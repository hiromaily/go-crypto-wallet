package btc

import (
	"errors"
	"fmt"

	apibtc "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/btc"
)

var timeoutSec int64 = 10

func runWalletPassphrase(btc apibtc.WalletSecurityManager, passphrase string) error {
	fmt.Println("stores the wallet decryption key in memory for 'timeout' seconds")

	// validator
	if passphrase == "" {
		return errors.New("passphrase option [-passphrase] is required")
	}

	err := btc.WalletPassphrase(passphrase, timeoutSec)
	if err != nil {
		return fmt.Errorf("fail to call btc.WalletPassphrase(): %w", err)
	}

	fmt.Println("wallet encryption is unlocked for 10s!")

	return nil
}
