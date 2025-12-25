package btc

import (
	"errors"
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/bitcoin"
)

func runEncryptWallet(btc bitcoin.Bitcoiner, passphrase string) error {
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
