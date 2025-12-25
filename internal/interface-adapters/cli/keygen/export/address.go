package export

import (
	"context"
	"errors"
	"fmt"

	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	"github.com/hiromaily/go-crypto-wallet/internal/di"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
)

func runAddress(container di.Container, acnt string) error {
	fmt.Println("export generated PublicKey as csv file")

	// validator
	if !domainAccount.ValidateAccountType(acnt) {
		return errors.New("account option [-account] is invalid")
	}
	if !domainAccount.NotAllow(acnt, []domainAccount.AccountType{domainAccount.AccountTypeAuthorization}) {
		return fmt.Errorf("account: %s is not allowed", domainAccount.AccountTypeAuthorization)
	}

	// export generated PublicKey as csv file
	useCase := container.NewKeygenExportAddressUseCase()
	output, err := useCase.Export(context.Background(), keygenusecase.ExportAddressInput{
		AccountType: domainAccount.AccountType(acnt),
	})
	if err != nil {
		return fmt.Errorf("fail to export address: %w", err)
	}
	fmt.Println("[fileName]: " + output.FileName)

	return nil
}
