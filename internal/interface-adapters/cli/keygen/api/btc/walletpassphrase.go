package btc

import (
	"errors"
	"fmt"

	apibtc "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/btc"
)

func runWalletPassphrase(btc apibtc.PKGRPCProvider, passphrase string) error {
	fmt.Println("stores the wallet decryption key in memory for 'timeout' seconds")

	// validator
	if passphrase == "" {
		return errors.New("passphrase option [-passphrase] is required")
	}

	err := btc.GetPkgRPC().WalletPassphrase(passphrase, 10)
	if err != nil {
		return fmt.Errorf("fail to call btc.GetPkgRPC().WalletPassphrase(): %w", err)
	}

	fmt.Println("wallet encryption is unlocked for 10s!")

	return nil
}
