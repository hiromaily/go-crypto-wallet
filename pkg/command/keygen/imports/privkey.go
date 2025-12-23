package imports

import (
	"errors"
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
)

func runPrivKey(wallet wallets.Keygener, acnt string) error {
	fmt.Println("import generated private key in database to keygen wallet")

	// validator
	if !account.ValidateAccountType(acnt) {
		return errors.New("account option [-account] is invalid")
	}
	if !account.NotAllow(acnt, []account.AccountType{account.AccountTypeAuthorization}) {
		return fmt.Errorf("account: %s is not allowed", account.AccountTypeAuthorization)
	}

	// import generated private key to keygen wallet
	err := wallet.ImportPrivKey(account.AccountType(acnt))
	if err != nil {
		return fmt.Errorf("fail to call ImportPrivKey() %w", err)
	}
	fmt.Println("Done!")

	return nil
}
