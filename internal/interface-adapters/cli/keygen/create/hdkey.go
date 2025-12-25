package create

import (
	"context"
	"errors"
	"fmt"

	"github.com/bookerzzz/grok"

	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	"github.com/hiromaily/go-crypto-wallet/internal/di"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
)

// runHDKeyWithFlags is the actual implementation that accepts parsed flags
func runHDKeyWithFlags(container di.Container, keyNum uint64, acnt string, _ bool) error {
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
	seedUseCase := container.NewKeygenGenerateSeedUseCase()
	seedOutput, err := seedUseCase.Generate(context.Background())
	if err != nil {
		return fmt.Errorf("fail to generate seed: %w", err)
	}

	// generate key for hd wallet
	hdWalletUseCase := container.NewKeygenGenerateHDWalletUseCase()
	hdWalletOutput, err := hdWalletUseCase.Generate(context.Background(), keygenusecase.GenerateHDWalletInput{
		AccountType: domainAccount.AccountType(acnt),
		Seed:        seedOutput.Seed,
		Count:       uint32(keyNum),
	})
	if err != nil {
		return fmt.Errorf("fail to generate HD wallet keys: %w", err)
	}
	grok.Value(hdWalletOutput)

	return nil
}
