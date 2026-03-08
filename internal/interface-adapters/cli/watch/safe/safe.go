// Package safe provides CLI commands for Safe contract operations on the Watch wallet.
package safe

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hiromaily/go-crypto-wallet/internal/di"
	wallets "github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/wallet"
	ethwallet "github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/wallet/eth"
)

// AddCommands adds all safe subcommands to parentCmd.
func AddCommands(parentCmd *cobra.Command, wallet *wallets.Watcher, _ func() di.Container) {
	parentCmd.AddCommand(newInfoCommand(wallet))
}

// newInfoCommand creates the `watch safe info` command.
func newInfoCommand(wallet *wallets.Watcher) *cobra.Command {
	var safeAddress string

	cmd := &cobra.Command{
		Use:   "info",
		Short: "Retrieve on-chain Safe contract state (ETH only)",
		Long: `Query the on-chain state of a Safe proxy contract:
owners list, threshold, current nonce, and balance in Wei.`,
		Example: "  watch --coin eth safe info --safe 0xYourSafeAddress",
		RunE: func(cmd *cobra.Command, args []string) error {
			if wallet == nil || *wallet == nil {
				return errors.New("wallet not initialized, check --coin flag")
			}
			ethWatch, ok := (*wallet).(*ethwallet.ETHWatch)
			if !ok {
				fmt.Println("this command is only supported for ETH")
				return nil
			}
			return runSafeInfo(ethWatch, safeAddress)
		},
	}

	cmd.Flags().StringVar(&safeAddress, "safe", "", "Safe proxy contract address (required)")
	cobra.CheckErr(cmd.MarkFlagRequired("safe"))

	return cmd
}

func runSafeInfo(ethWatch *ethwallet.ETHWatch, safeAddress string) error {
	output, err := ethWatch.GetSafeInfo(safeAddress)
	if err != nil {
		return fmt.Errorf("failed to get Safe info: %w", err)
	}

	fmt.Printf("Safe address:  %s\n", safeAddress)
	fmt.Printf("Threshold:     %d\n", output.Threshold)
	fmt.Printf("Nonce:         %d\n", output.Nonce)
	fmt.Printf("Balance (Wei): %s\n", output.BalanceWei)
	fmt.Printf("Owners (%d):\n", len(output.Owners))
	for i, owner := range output.Owners {
		fmt.Printf("  [%d] %s\n", i+1, owner)
	}

	return nil
}
