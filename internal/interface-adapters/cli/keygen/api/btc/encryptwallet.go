package btc

import (
	"errors"
	"fmt"

	apibtc "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/btc"
)

func runEncryptWallet(btc apibtc.PKGRPCProvider, passphrase string) error {
	fmt.Println("encrypts the wallet with 'passphrase'")

	// validator
	if passphrase == "" {
		return errors.New("passphrase option [-passphrase] is required")
	}

	err := btc.GetPkgRPC().EncryptWallet(passphrase)
	if err != nil {
		return fmt.Errorf("fail to call btc.GetPkgRPC().EncryptWallet(): %w", err)
	}

	fmt.Println("wallet is encrypted!")

	return nil
}
