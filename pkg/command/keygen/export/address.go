package export

import (
	"errors"
	"fmt"

	domainAccount "github.com/hiromaily/go-crypto-wallet/pkg/domain/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
)

func runAddress(wallet wallets.Keygener, acnt string) error {
	fmt.Println("export generated PublicKey as csv file")

	// validator
	if !domainAccount.ValidateAccountType(acnt) {
		return errors.New("account option [-account] is invalid")
	}
	if !domainAccount.NotAllow(acnt, []domainAccount.AccountType{domainAccount.AccountTypeAuthorization}) {
		return fmt.Errorf("account: %s is not allowed", domainAccount.AccountTypeAuthorization)
	}

	// export generated PublicKey as csv file
	fileName, err := wallet.ExportAddress(domainAccount.AccountType(acnt))
	if err != nil {
		return fmt.Errorf("fail to call ExportAddress() %w", err)
	}
	fmt.Println("[fileName]: " + fileName)

	return nil
}
