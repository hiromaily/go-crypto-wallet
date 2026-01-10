package imports

import (
	"context"
	"errors"
	"fmt"

	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	"github.com/hiromaily/go-crypto-wallet/internal/di"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
)

func runDescriptor(container di.Container, filePath string, acnt string) error {
	fmt.Println("import descriptors from file")

	// Validate file path
	if filePath == "" {
		return errors.New("file path [-file] is required")
	}

	// Validate account
	if !domainAccount.ValidateAccountType(acnt) {
		return errors.New("account option [-account] is invalid")
	}
	if !domainAccount.NotAllow(acnt, []domainAccount.AccountType{domainAccount.AccountTypeAuthorization}) {
		return fmt.Errorf("account: %s is not allowed", domainAccount.AccountTypeAuthorization)
	}

	// Import descriptors
	useCase := container.NewWatchImportDescriptorUseCase()
	output, err := useCase.Import(context.Background(), watchusecase.ImportDescriptorInput{
		FilePath:    filePath,
		AccountType: domainAccount.AccountType(acnt),
		StartIndex:  0,
		Count:       1000, // Default range for descriptor import
	})
	if err != nil {
		return fmt.Errorf("fail to import descriptor: %w", err)
	}

	fmt.Printf("[descriptors_imported]: %d\n", output.DescriptorsImported)
	fmt.Printf("[addresses_generated]: %d\n", output.AddressesGenerated)
	if len(output.Errors) > 0 {
		fmt.Println("[errors]:")
		for _, errMsg := range output.Errors {
			fmt.Printf("  - %s\n", errMsg)
		}
	}

	return nil
}
