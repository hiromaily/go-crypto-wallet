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
	var (
		file         string
		account      string
		startIndex   uint32
		count        uint32
		validateOnly bool
	)

	descCmd := &cobra.Command{
		Use:   "descriptor",
		Short: "Import descriptors into Bitcoin Core (BTC only)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if wallet == nil || *wallet == nil {
				return errors.New("wallet not initialized, check --coin flag")
			}
			if _, ok := (*wallet).(*btcwallet.BTCWatch); !ok {
				fmt.Println("[WARN] descriptor command is not supported for this coin type")
				return nil
			}
			return runDescriptorImport(cmd.Context(), containerGetter(), watchusecase.ImportDescriptorInput{
				FilePath:     file,
				AccountType:  domainAccount.AccountType(account),
				StartIndex:   startIndex,
				Count:        count,
				ValidateOnly: validateOnly,
			})
		},
	}

	descCmd.Flags().StringVar(&file, "file", "", "descriptor file path (omit to read from stdin)")
	descCmd.Flags().StringVar(&account, "account", "", "target account (e.g. deposit, payment)")
	descCmd.Flags().Uint32Var(&startIndex, "start", 0, "start index for address derivation")
	descCmd.Flags().Uint32Var(&count, "count", 100, "number of addresses to derive per descriptor")
	descCmd.Flags().BoolVar(&validateOnly, "validate", false, "validate descriptors without importing")
	_ = descCmd.MarkFlagRequired("account")

	parentCmd.AddCommand(descCmd)
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
