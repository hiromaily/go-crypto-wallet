package export

import (
	"context"
	"errors"
	"fmt"

	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	"github.com/hiromaily/go-crypto-wallet/internal/di"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
)

func runDescriptor(container di.Container, acnt string, outputPath string, format string, includeChange bool) error {
	fmt.Println("export output descriptors to file")

	// Validate account
	if !domainAccount.ValidateAccountType(acnt) {
		return errors.New("account option [-account] is invalid")
	}
	if !domainAccount.NotAllow(acnt, []domainAccount.AccountType{domainAccount.AccountTypeAuthorization}) {
		return fmt.Errorf("account: %s is not allowed", domainAccount.AccountTypeAuthorization)
	}

	// Validate output path
	if outputPath == "" {
		return errors.New("output path [-output] is required")
	}

	// Validate format
	var descriptorFormat keygenusecase.DescriptorFormat
	switch format {
	case "text":
		descriptorFormat = keygenusecase.DescriptorFormatText
	case "json":
		descriptorFormat = keygenusecase.DescriptorFormatJSON
	case "bitcoin-core", "":
		descriptorFormat = keygenusecase.DescriptorFormatBitcoinCore
	default:
		return fmt.Errorf("invalid format: %s (valid: text, json, bitcoin-core)", format)
	}

	// Export descriptors
	useCase := container.NewKeygenExportDescriptorUseCase()
	output, err := useCase.Export(context.Background(), keygenusecase.ExportDescriptorInput{
		AccountType:   domainAccount.AccountType(acnt),
		OutputPath:    outputPath,
		Format:        descriptorFormat,
		IncludeChange: includeChange,
	})
	if err != nil {
		return fmt.Errorf("fail to export descriptor: %w", err)
	}
	fmt.Printf("[exported]: %s\n", output.FilePath)

	return nil
}
