package btc

import (
	"errors"
	"fmt"
)

func runWalletPassphrase(btc btcKeygenAPICmds, passphrase string) error {
	fmt.Println("stores the wallet decryption key in memory for 'timeout' seconds")

	// validator
	if passphrase == "" {
		return errors.New("passphrase option [-passphrase] is required")
	}

	err := btc.WalletPassphrase(passphrase, 10)
	if err != nil {
		return fmt.Errorf("fail to call btc.WalletPassphrase() %w", err)
	}

	fmt.Println("wallet encryption is unlocked for 10s!")

	return nil
}
