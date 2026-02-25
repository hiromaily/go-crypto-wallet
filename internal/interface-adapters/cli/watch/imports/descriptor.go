package imports

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	watchusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/watch"
	"github.com/hiromaily/go-crypto-wallet/internal/di"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	wallets "github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/wallet"
	btcwallet "github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/wallet/btc"
)

func addDescriptorCommands(
	parentCmd *cobra.Command, wallet *wallets.Watcher, containerGetter func() di.Container,
) {
	descCmd := &cobra.Command{
		Use:   "descriptor",
		Short: "Descriptor operations for BTC watch wallet (BTC only)",
	}

	descCmd.AddCommand(
		newDescriptorImportCommand(wallet, containerGetter),
		newDescriptorValidateCommand(wallet, containerGetter),
	)

	parentCmd.AddCommand(descCmd)
}

func newDescriptorImportCommand(wallet *wallets.Watcher, containerGetter func() di.Container) *cobra.Command {
	var (
		file       string
		account    string
		startIndex uint32
		count      uint32
	)

	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import descriptors from file",
		RunE: func(cmd *cobra.Command, args []string) error {
			if wallet == nil || *wallet == nil {
				return errors.New("wallet not initialized, check --coin flag")
			}
			if _, ok := (*wallet).(*btcwallet.BTCWatch); !ok {
				fmt.Println("this command is only supported for BTC")
				return nil
			}
			return runDescriptorImport(cmd.Context(), containerGetter(), watchusecase.ImportDescriptorInput{
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

func newDescriptorValidateCommand(wallet *wallets.Watcher, containerGetter func() di.Container) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate descriptors without importing",
		RunE: func(cmd *cobra.Command, args []string) error {
			if wallet == nil || *wallet == nil {
				return errors.New("wallet not initialized, check --coin flag")
			}
			if _, ok := (*wallet).(*btcwallet.BTCWatch); !ok {
				fmt.Println("this command is only supported for BTC")
				return nil
			}
			return runDescriptorImport(cmd.Context(), containerGetter(), watchusecase.ImportDescriptorInput{
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
