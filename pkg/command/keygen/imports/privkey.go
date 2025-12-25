package imports

import (
	"context"
	"errors"
	"fmt"

	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/di"
)

func runPrivKey(container di.Container, acnt string) error {
	fmt.Println("import generated private key in database to keygen wallet")

	// validator
	if !domainAccount.ValidateAccountType(acnt) {
		return errors.New("account option [-account] is invalid")
	}
	if !domainAccount.NotAllow(acnt, []domainAccount.AccountType{domainAccount.AccountTypeAuthorization}) {
		return fmt.Errorf("account: %s is not allowed", domainAccount.AccountTypeAuthorization)
	}

	// import generated private key to keygen wallet
	useCase := container.NewKeygenImportPrivateKeyUseCase()
	err := useCase.Import(context.Background(), keygenusecase.ImportPrivateKeyInput{
		AccountType: domainAccount.AccountType(acnt),
	})
	if err != nil {
		return fmt.Errorf("fail to import private key: %w", err)
	}
	fmt.Println("Done!")

	return nil
}
