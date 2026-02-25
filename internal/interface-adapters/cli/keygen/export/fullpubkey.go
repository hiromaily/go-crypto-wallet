package export

import (
	"context"
	"errors"
	"fmt"

	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	"github.com/hiromaily/go-crypto-wallet/internal/di"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
)

func runFullPubkey(container di.Container, acnt string) error {
	fmt.Println("export account-level extended public key (xpub) as csv file")

	// validator
	if !domainAccount.ValidateAccountType(acnt) {
		return errors.New("account option [-account] is invalid")
	}
	if domainAccount.Allow(acnt, []domainAccount.AccountType{domainAccount.AccountTypeAuthorization}) {
		return fmt.Errorf("account: %s is not allowed", domainAccount.AccountTypeAuthorization)
	}

	// export accountXpub to file for Watch wallet address derivation
	useCase := container.NewKeygenExportFullPubkeyUseCase()
	output, err := useCase.Export(context.Background(), keygenusecase.ExportFullPubkeyInput{
		AccountType: domainAccount.AccountType(acnt),
	})
	if err != nil {
		return fmt.Errorf("fail to export full pubkey: %w", err)
	}
	fmt.Println("[fileName]: " + output.FileName)

	return nil
}
