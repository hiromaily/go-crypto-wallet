package btc

import (
	"context"
	"encoding/json"
	"errors"
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
		return keygenusecase.ExportDescriptorOutput{}, errors.New("output path is required")
	}

	format := input.Format
	if format == "" {
		format = keygenusecase.DescriptorFormatBitcoinCore
	}

	addressTypes := supportedAddressTypes()
	descriptors := make([]string, 0, len(addressTypes)*2)
	var skippedTypes []string

	for _, addrType := range addressTypes {
		generate := func(isChange bool) error {
			output, err := u.generator.Generate(ctx, keygenusecase.GenerateDescriptorInput{
				AccountType: input.AccountType,
				AddressType: addrType,
				IsChange:    isChange,
			})
			if err != nil {
				// Skip address types that don't have keys instead of failing
				return err
			}
			descriptors = append(descriptors, output.Descriptor)
			return nil
		}

		if err := generate(false); err != nil {
			// Log and skip this address type (keys may not exist for all types)
			skippedTypes = append(skippedTypes, addrType.String())
			continue
		}

		if input.IncludeChange {
			if err := generate(true); err != nil {
				// If receive descriptor exists but change fails, it's an actual error
				return keygenusecase.ExportDescriptorOutput{},
					fmt.Errorf("generate change descriptor for %s/%s: %w",
						input.AccountType, addrType, err)
			}
		}
	}

	// If no descriptors were generated, return error
	if len(descriptors) == 0 {
		if len(skippedTypes) > 0 {
			return keygenusecase.ExportDescriptorOutput{},
				fmt.Errorf("no keys found for any address type (skipped: %v)", skippedTypes)
		}
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
