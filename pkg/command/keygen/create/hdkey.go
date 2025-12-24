package create

import (
	"errors"
	"fmt"

	"github.com/bookerzzz/grok"

	domainAccount "github.com/hiromaily/go-crypto-wallet/pkg/domain/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
)

// runHDKeyWithFlags is the actual implementation that accepts parsed flags
func runHDKeyWithFlags(wallet wallets.Keygener, keyNum uint64, acnt string, isKeyPair bool) error {
	fmt.Println("create HD key")

	// validator
	if keyNum == 0 {
		return errors.New("number of key option [-keynum] is required")
	}
	if !domainAccount.ValidateAccountType(acnt) {
		return errors.New("account option [-account] is invalid")
	}
	if !domainAccount.NotAllow(acnt, []domainAccount.AccountType{domainAccount.AccountTypeAuthorization}) {
		return fmt.Errorf("account: %s is not allowed", domainAccount.AccountTypeAuthorization)
	}

	// create seed
	bSeed, err := wallet.GenerateSeed()
	if err != nil {
		return fmt.Errorf("fail to call GenerateSeed() %w", err)
	}

	// generate key for hd wallet
	keys, err := wallet.GenerateAccountKey(domainAccount.AccountType(acnt), bSeed, uint32(keyNum), isKeyPair)
	if err != nil {
		return fmt.Errorf("fail to call GenerateAccountKey() %w", err)
	}
	grok.Value(keys)

	return nil
}
