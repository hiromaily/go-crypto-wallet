package btc

import (
	"errors"
	"fmt"

	apibtc "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/btc"
)

func runWalletPassphraseChange(btc apibtc.Bitcoiner, old, newPass string) error {
	fmt.Println("changes the wallet passphrase from 'oldpassphrase' to 'newpassphrase'")

	// validator
	if old == "" {
		return errors.New("old passphrase option [-old] is required")
	}
	if newPass == "" {
		return errors.New("new passphrase option [-new] is required")
	}

	err := btc.WalletPassphraseChange(old, newPass)
	if err != nil {
		return fmt.Errorf("fail to call btc.WalletPassphraseChange() %w", err)
	}

	fmt.Println("wallet passphrase was changed!")

	return nil
}
