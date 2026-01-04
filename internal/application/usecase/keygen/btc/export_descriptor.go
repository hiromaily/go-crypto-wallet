package btc

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	portsStorage "github.com/hiromaily/go-crypto-wallet/internal/application/ports/storage"
	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
)

type exportDescriptorUseCase struct {
	generator  keygenusecase.GenerateDescriptorUseCase
	fileWriter portsStorage.DescriptorFileWriter
}

// NewExportDescriptorUseCase creates a descriptor export use case.
func NewExportDescriptorUseCase(
	generator keygenusecase.GenerateDescriptorUseCase,
	fileWriter portsStorage.DescriptorFileWriter,
) keygenusecase.ExportDescriptorUseCase {
	return &exportDescriptorUseCase{
		generator:  generator,
		fileWriter: fileWriter,
	}
}

func (u *exportDescriptorUseCase) Export(
	ctx context.Context,
	input keygenusecase.ExportDescriptorInput,
) (keygenusecase.ExportDescriptorOutput, error) {
	if input.OutputPath == "" {
		return keygenusecase.ExportDescriptorOutput{}, fmt.Errorf("output path is required")
	}

	format := input.Format
	if format == "" {
		format = keygenusecase.DescriptorFormatBitcoinCore
	}

	addressTypes := supportedAddressTypes()
	descriptors := make([]string, 0, len(addressTypes)*2)

	for _, addrType := range addressTypes {
		output, err := u.generator.Generate(ctx, keygenusecase.GenerateDescriptorInput{
			AccountType: input.AccountType,
			AddressType: addrType,
			IsChange:    false,
		})
		if err != nil {
			return keygenusecase.ExportDescriptorOutput{}, fmt.Errorf("generate receive descriptor for %s/%s: %w", input.AccountType, addrType, err)
		}
		descriptors = append(descriptors, output.Descriptor)

		if input.IncludeChange {
			changeOutput, err := u.generator.Generate(ctx, keygenusecase.GenerateDescriptorInput{
				AccountType: input.AccountType,
				AddressType: addrType,
				IsChange:    true,
			})
			if err != nil {
				return keygenusecase.ExportDescriptorOutput{}, fmt.Errorf("generate change descriptor for %s/%s: %w", input.AccountType, addrType, err)
			}
			descriptors = append(descriptors, changeOutput.Descriptor)
		}
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

func (u *exportDescriptorUseCase) formatDescriptors(
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

func supportedAddressTypes() []domainAddress.AddrType {
	return []domainAddress.AddrType{
		domainAddress.AddrTypeTaproot,
		domainAddress.AddrTypeBech32,
		domainAddress.AddrTypeP2shSegwit,
		domainAddress.AddrTypeLegacy,
	}
}

type bitcoinCoreDescriptor struct {
	Descriptor string `json:"desc"`
	Timestamp  string `json:"timestamp"`
	Range      [2]int `json:"range"`
	WatchOnly  bool   `json:"watchonly"`
}
