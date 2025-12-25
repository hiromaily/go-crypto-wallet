package create

import (
	"context"
	"errors"
	"fmt"

	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/di"
)

// runMultisigWithAccount is the actual implementation that accepts parsed flags
func runMultisigWithAccount(container di.Container, acnt string) error {
	fmt.Println("create multisig address")

	// validator
	if !domainAccount.ValidateAccountType(acnt) {
		return errors.New("account option [-account] is invalid")
	}

	// create multisig address
	useCase := container.NewKeygenCreateMultisigAddressUseCase()
	err := useCase.Create(context.Background(), keygenusecase.CreateMultisigAddressInput{
		AccountType: domainAccount.AccountType(acnt),
		AddressType: container.AddressType(),
	})
	if err != nil {
		return fmt.Errorf("fail to create multisig address: %w", err)
	}

	return nil
}
