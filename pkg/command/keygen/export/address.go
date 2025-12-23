package export

import (
	"errors"
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
)

func runAddress(wallet wallets.Keygener, acnt string) error {
	fmt.Println("export generated PublicKey as csv file")

	// validator
	if !account.ValidateAccountType(acnt) {
		return errors.New("account option [-account] is invalid")
	}
	if !account.NotAllow(acnt, []account.AccountType{account.AccountTypeAuthorization}) {
		return fmt.Errorf("account: %s is not allowed", account.AccountTypeAuthorization)
	}

	// export generated PublicKey as csv file
	fileName, err := wallet.ExportAddress(account.AccountType(acnt))
	if err != nil {
		return fmt.Errorf("fail to call ExportAddress() %w", err)
	}
	fmt.Println("[fileName]: " + fileName)

	return nil
}
