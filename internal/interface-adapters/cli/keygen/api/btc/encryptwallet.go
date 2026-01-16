package btc

import (
	"errors"
	"fmt"

	portsBitcoin "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/bitcoin"
)

func runEncryptWallet(btc portsBitcoin.Bitcoiner, passphrase string) error {
	fmt.Println("encrypts the wallet with 'passphrase'")

	// validator
	if passphrase == "" {
		return errors.New("passphrase option [-passphrase] is required")
	}

	err := btc.EncryptWallet(passphrase)
	if err != nil {
		return fmt.Errorf("fail to call btc.EncryptWallet() %w", err)
	}

	fmt.Println("wallet is encrypted!")

	return nil
}
