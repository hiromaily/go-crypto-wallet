package btc

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	"github.com/hiromaily/go-crypto-wallet/internal/di"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
)

// AddDescriptorCommands adds descriptor-related commands for watch wallet (BTC).
func AddDescriptorCommands(parentCmd *cobra.Command, container di.Container) {
	descCmd := &cobra.Command{
		Use:   "descriptor",
		Short: "Descriptor operations for BTC watch wallet",
	}

	descCmd.AddCommand(
		newDescriptorImportCommand(container),
		newDescriptorValidateCommand(container),
	)

	parentCmd.AddCommand(descCmd)
}

func newDescriptorImportCommand(container di.Container) *cobra.Command {
	var (
		file       string
		account    string
		startIndex uint32
		count      uint32
	)

	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import descriptors from file or stdin",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDescriptorImport(cmd.Context(), container, watchusecase.ImportDescriptorInput{
				FilePath:    file,
				AccountType: domainAccount.AccountType(account),
				StartIndex:  startIndex,
				Count:       count,
			})
		},
	}

	cmd.Flags().StringVar(&file, "file", "", "descriptor file path (omit to read from stdin)")
	cmd.Flags().StringVar(&account, "account", "", "target account (e.g. deposit, payment)")
	cmd.Flags().Uint32Var(&startIndex, "start", 0, "start index for address derivation")
	cmd.Flags().Uint32Var(&count, "count", 100, "number of addresses to derive per descriptor")
	_ = cmd.MarkFlagRequired("account")

	return cmd
}

func newDescriptorValidateCommand(container di.Container) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate descriptors without importing",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDescriptorImport(cmd.Context(), container, watchusecase.ImportDescriptorInput{
				FilePath:     file,
				ValidateOnly: true,
			})
		},
	}

	cmd.Flags().StringVar(&file, "file", "", "descriptor file path (omit to read from stdin)")

	return cmd
}

func runDescriptorImport(ctx context.Context, container di.Container, input watchusecase.ImportDescriptorInput) error {
	if container == nil {
		return errors.New("container is not initialized")
	}

	if input.AccountType != "" && !domainAccount.ValidateAccountType(input.AccountType.String()) {
		return errors.New("account option [--account] is invalid")
	}

	useCase := container.NewWatchImportDescriptorUseCase()
	output, err := useCase.Import(ctx, input)
	if err != nil {
		return err
	}

	if !input.ValidateOnly {
		fmt.Printf(
			"Descriptors imported: %d, addresses generated: %d\n",
			output.DescriptorsImported,
			output.AddressesGenerated,
		)
	} else {
		fmt.Println("Descriptors validated successfully")
	}

	if len(output.Errors) > 0 {
		fmt.Println("Some descriptors reported errors:")
		for _, e := range output.Errors {
			fmt.Printf("- %s\n", e)
		}
	}

	return nil
}
