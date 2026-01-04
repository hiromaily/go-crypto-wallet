package btc

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	"github.com/hiromaily/go-crypto-wallet/internal/di"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
)

// AddDescriptorCommands adds descriptor-related subcommands for BTC.
func AddDescriptorCommands(parentCmd *cobra.Command, container di.Container) {
	descriptorCmd := &cobra.Command{
		Use:   "descriptor",
		Short: "Descriptor operations for BTC keygen wallet",
	}

	descriptorCmd.AddCommand(
		newDescriptorGenerateCommand(container),
		newDescriptorExportCommand(container),
		newDescriptorExportAllCommand(container),
	)

	parentCmd.AddCommand(descriptorCmd)
}

func newDescriptorGenerateCommand(container di.Container) *cobra.Command {
	var (
		account       string
		addressType   string
		includeChange bool
		requiredSigs  int
	)

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a descriptor for the specified account and address type",
		Example: "  keygen --coin btc descriptor generate --account deposit --address-type taproot\n" +
			"  keygen --coin btc descriptor generate --account deposit --address-type bech32 --change",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDescriptorGenerate(cmd.Context(), container, account, addressType, includeChange, requiredSigs)
		},
	}

	cmd.Flags().StringVar(&account, "account", "", "target account (e.g. deposit, payment)")
	cmd.Flags().StringVar(&addressType, "address-type", "", "address type (taproot|bech32|p2sh-segwit|legacy)")
	cmd.Flags().BoolVar(&includeChange, "change", false, "generate change descriptor instead of receive")
	cmd.Flags().IntVar(&requiredSigs, "required-sigs", 0, "required signatures for multisig accounts (optional)")
	_ = cmd.MarkFlagRequired("account")
	_ = cmd.MarkFlagRequired("address-type")

	return cmd
}

func newDescriptorExportCommand(container di.Container) *cobra.Command {
	var (
		account       string
		outputPath    string
		format        string
		includeChange bool
	)

	cmd := &cobra.Command{
		Use:     "export",
		Short:   "Export descriptors for the account to a file",
		Example: "  keygen --coin btc descriptor export --account deposit --output /tmp/descriptors.txt --format bitcoin-core --include-change",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDescriptorExport(cmd.Context(), container, account, outputPath, format, includeChange)
		},
	}

	cmd.Flags().StringVar(&account, "account", "", "target account (e.g. deposit, payment)")
	cmd.Flags().StringVar(&outputPath, "output", "", "output file path")
	cmd.Flags().StringVar(&format, "format", string(keygenusecase.DescriptorFormatBitcoinCore), "output format (text|json|bitcoin-core)")
	cmd.Flags().BoolVar(&includeChange, "include-change", false, "include change descriptors")
	_ = cmd.MarkFlagRequired("account")
	_ = cmd.MarkFlagRequired("output")

	return cmd
}

func newDescriptorExportAllCommand(container di.Container) *cobra.Command {
	var (
		account    string
		outputPath string
		format     string
	)

	cmd := &cobra.Command{
		Use:     "export-all",
		Short:   "Export all descriptors (receive and change) for the account",
		Example: "  keygen --coin btc descriptor export-all --account deposit --output /tmp/deposit_descriptors.json",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDescriptorExport(cmd.Context(), container, account, outputPath, format, true)
		},
	}

	cmd.Flags().StringVar(&account, "account", "", "target account (e.g. deposit, payment)")
	cmd.Flags().StringVar(&outputPath, "output", "", "output file path")
	cmd.Flags().StringVar(&format, "format", string(keygenusecase.DescriptorFormatJSON), "output format (text|json|bitcoin-core)")
	_ = cmd.MarkFlagRequired("account")
	_ = cmd.MarkFlagRequired("output")

	return cmd
}

func runDescriptorGenerate(
	ctx context.Context,
	container di.Container,
	account string,
	addressType string,
	isChange bool,
	requiredSigs int,
) error {
	acnt, err := parseAccountType(account)
	if err != nil {
		return err
	}
	addrType, err := parseAddressType(addressType)
	if err != nil {
		return err
	}
	if container == nil {
		return fmt.Errorf("container is not initialized")
	}

	useCase := container.NewKeygenGenerateDescriptorUseCase()
	output, err := useCase.Generate(ctx, keygenusecase.GenerateDescriptorInput{
		AccountType:  acnt,
		AddressType:  addrType,
		IsChange:     isChange,
		RequiredSigs: requiredSigs,
	})
	if err != nil {
		return fmt.Errorf("failed to generate descriptor: %w", err)
	}

	fmt.Printf("%s\n", output.Descriptor)
	return nil
}

func runDescriptorExport(
	ctx context.Context,
	container di.Container,
	account string,
	outputPath string,
	format string,
	includeChange bool,
) error {
	acnt, err := parseAccountType(account)
	if err != nil {
		return err
	}
	formatType, err := parseDescriptorFormat(format)
	if err != nil {
		return err
	}
	if container == nil {
		return fmt.Errorf("container is not initialized")
	}

	useCase := container.NewKeygenExportDescriptorUseCase()
	output, err := useCase.Export(ctx, keygenusecase.ExportDescriptorInput{
		AccountType:   acnt,
		OutputPath:    outputPath,
		Format:        formatType,
		IncludeChange: includeChange,
	})
	if err != nil {
		return fmt.Errorf("failed to export descriptors: %w", err)
	}

	fmt.Printf("descriptors exported to %s\n", output.FilePath)
	return nil
}

func parseAccountType(account string) (domainAccount.AccountType, error) {
	if !domainAccount.ValidateAccountType(account) {
		return "", fmt.Errorf("account option [--account] is invalid")
	}
	return domainAccount.AccountType(account), nil
}

func parseAddressType(addressType string) (domainAddress.AddrType, error) {
	switch strings.ToLower(addressType) {
	case domainAddress.AddrTypeTaproot.String():
		return domainAddress.AddrTypeTaproot, nil
	case domainAddress.AddrTypeBech32.String():
		return domainAddress.AddrTypeBech32, nil
	case domainAddress.AddrTypeP2shSegwit.String():
		return domainAddress.AddrTypeP2shSegwit, nil
	case domainAddress.AddrTypeLegacy.String():
		return domainAddress.AddrTypeLegacy, nil
	default:
		return "", fmt.Errorf("address-type is invalid: %s", addressType)
	}
}

func parseDescriptorFormat(format string) (keygenusecase.DescriptorFormat, error) {
	switch strings.ToLower(format) {
	case string(keygenusecase.DescriptorFormatText):
		return keygenusecase.DescriptorFormatText, nil
	case string(keygenusecase.DescriptorFormatJSON):
		return keygenusecase.DescriptorFormatJSON, nil
	case string(keygenusecase.DescriptorFormatBitcoinCore):
		return keygenusecase.DescriptorFormatBitcoinCore, nil
	default:
		return "", fmt.Errorf("format is invalid: %s", format)
	}
}
