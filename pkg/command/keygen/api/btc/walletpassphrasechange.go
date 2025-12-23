package btc

import (
	"errors"
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/btcgrp"
)

func runWalletPassphraseChange(btc btcgrp.Bitcoiner, old, newPass string) error {
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
