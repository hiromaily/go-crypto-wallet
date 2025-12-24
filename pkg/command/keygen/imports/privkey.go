package imports

import (
	"errors"
	"fmt"

	domainAccount "github.com/hiromaily/go-crypto-wallet/pkg/domain/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
)

func runPrivKey(wallet wallets.Keygener, acnt string) error {
	fmt.Println("import generated private key in database to keygen wallet")

	// validator
	if !domainAccount.ValidateAccountType(acnt) {
		return errors.New("account option [-account] is invalid")
	}
	if !domainAccount.NotAllow(acnt, []domainAccount.AccountType{domainAccount.AccountTypeAuthorization}) {
		return fmt.Errorf("account: %s is not allowed", domainAccount.AccountTypeAuthorization)
	}

	// import generated private key to keygen wallet
	err := wallet.ImportPrivKey(domainAccount.AccountType(acnt))
	if err != nil {
		return fmt.Errorf("fail to call ImportPrivKey() %w", err)
	}
	fmt.Println("Done!")

	return nil
}
