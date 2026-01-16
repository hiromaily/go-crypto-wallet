package btc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	file "github.com/hiromaily/go-crypto-wallet/internal/application/ports/file"
	repocold "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository/cold"
	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
)

type exportDescriptorUseCase struct {
	generator      keygenusecase.GenerateDescriptorUseCase
	fileWriter     file.DescriptorFileWriter
	accountKeyRepo repocold.BTCAccountKeyRepositorier
}

// NewExportDescriptorUseCase creates a descriptor export use case.
func NewExportDescriptorUseCase(
	generator keygenusecase.GenerateDescriptorUseCase,
	fileWriter file.DescriptorFileWriter,
	accountKeyRepo repocold.BTCAccountKeyRepositorier,
) keygenusecase.ExportDescriptorUseCase {
	return &exportDescriptorUseCase{
		generator:      generator,
		fileWriter:     fileWriter,
		accountKeyRepo: accountKeyRepo,
	}
}

func (u *exportDescriptorUseCase) Export(
	ctx context.Context,
	input keygenusecase.ExportDescriptorInput,
) (keygenusecase.ExportDescriptorOutput, error) {
	if input.OutputPath == "" {
		return keygenusecase.ExportDescriptorOutput{}, errors.New("output path is required")
	}

	format := input.Format
	if format == "" {
		format = keygenusecase.DescriptorFormatBitcoinCore
	}

	// Determine which address types have keys in the database by checking the account key
	addressTypes, err := u.getAvailableAddressTypes(ctx, input.AccountType)
	if err != nil {
		return keygenusecase.ExportDescriptorOutput{},
			fmt.Errorf("determine available address types: %w", err)
	}

	if len(addressTypes) == 0 {
		return keygenusecase.ExportDescriptorOutput{},
			fmt.Errorf("no keys found for account %s - run 'keygen create hdkey' first", input.AccountType)
	}

	descriptors := make([]string, 0, len(addressTypes)*2)

	for _, addrType := range addressTypes {
		generate := func(isChange bool) error {
			output, err := u.generator.Generate(ctx, keygenusecase.GenerateDescriptorInput{
				AccountType: input.AccountType,
				AddressType: addrType,
				IsChange:    isChange,
			})
			if err != nil {
				return err
			}
			descriptors = append(descriptors, output.Descriptor)
			return nil
		}

		if err := generate(false); err != nil {
			return keygenusecase.ExportDescriptorOutput{},
				fmt.Errorf("generate descriptor for %s/%s: %w",
					input.AccountType, addrType, err)
		}

		if input.IncludeChange {
			if err := generate(true); err != nil {
				return keygenusecase.ExportDescriptorOutput{},
					fmt.Errorf("generate change descriptor for %s/%s: %w",
						input.AccountType, addrType, err)
			}
		}
	}

	// If no descriptors were generated, return error
	if len(descriptors) == 0 {
		return keygenusecase.ExportDescriptorOutput{}, errors.New("no descriptors generated")
	}

	content, err := u.formatDescriptors(descriptors, format)
	if err != nil {
		return keygenusecase.ExportDescriptorOutput{}, fmt.Errorf("format descriptors: %w", err)
	}

	if err := u.fileWriter.WriteFile(input.OutputPath, content); err != nil {
		return keygenusecase.ExportDescriptorOutput{}, fmt.Errorf("write descriptor file: %w", err)
	}

	return keygenusecase.ExportDescriptorOutput{
		FilePath: input.OutputPath,
	}, nil
}

func (*exportDescriptorUseCase) formatDescriptors(
	descriptors []string,
	format keygenusecase.DescriptorFormat,
) ([]byte, error) {
	switch format {
	case keygenusecase.DescriptorFormatText:
		return []byte(strings.Join(descriptors, "\n")), nil
	case keygenusecase.DescriptorFormatJSON:
		return json.MarshalIndent(descriptors, "", "  ")
	case keygenusecase.DescriptorFormatBitcoinCore:
		items := make([]bitcoinCoreDescriptor, 0, len(descriptors))
		for _, desc := range descriptors {
			items = append(items, bitcoinCoreDescriptor{
				Descriptor: desc,
				Timestamp:  "now",
				Range:      [2]int{0, 1000},
				WatchOnly:  true,
			})
		}
		return json.MarshalIndent(items, "", "  ")
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

// getAvailableAddressTypes queries the database to find which key types exist
// for the given account, and returns the corresponding address types.
// This prevents trying to export descriptors for key types that don't exist.
func (u *exportDescriptorUseCase) getAvailableAddressTypes(
	_ context.Context,
	accountType domainAccount.AccountType,
) ([]domainAddress.AddrType, error) {
	// Get the account key to determine which key_type exists
	accountKey, err := u.accountKeyRepo.GetOneMaxID(accountType)
	if err != nil {
		return nil, fmt.Errorf("get account key: %w", err)
	}
	if accountKey == nil {
		return nil, nil // No keys found
	}

	// Convert key_type string to KeyType
	keyType := domainKey.KeyType(accountKey.KeyType)

	// Determine the address type from the key type
	var addrType domainAddress.AddrType
	switch keyType {
	case domainKey.KeyTypeBIP44:
		addrType = domainAddress.AddrTypeLegacy
	case domainKey.KeyTypeBIP49:
		addrType = domainAddress.AddrTypeP2shSegwit
	case domainKey.KeyTypeBIP84:
		addrType = domainAddress.AddrTypeBech32
	case domainKey.KeyTypeBIP86:
		addrType = domainAddress.AddrTypeTaproot
	case domainKey.KeyTypeMuSig2:
		addrType = domainAddress.AddrTypeTaproot
	default:
		return nil, fmt.Errorf("unsupported key type: %s", keyType)
	}

	return []domainAddress.AddrType{addrType}, nil
}

type bitcoinCoreDescriptor struct {
	Descriptor string `json:"desc"`
	Timestamp  string `json:"timestamp"`
	Range      [2]int `json:"range"`
	WatchOnly  bool   `json:"watchonly"`
}
