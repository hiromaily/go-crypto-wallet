package create

import (
	"errors"
	"fmt"

	domainAccount "github.com/hiromaily/go-crypto-wallet/pkg/domain/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
)

// runMultisigWithAccount is the actual implementation that accepts parsed flags
func runMultisigWithAccount(wallet wallets.Keygener, acnt string) error {
	fmt.Println("create multisig address")

	// validator
	if !domainAccount.ValidateAccountType(acnt) {
		return errors.New("account option [-account] is invalid")
	}

	// create multisig address
	err := wallet.CreateMultisigAddress(domainAccount.AccountType(acnt))
	if err != nil {
		return fmt.Errorf("fail to call CreateMultisigAddress() %w", err)
	}

	return nil
}
