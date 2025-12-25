package monitor

import (
	"github.com/spf13/cobra"

	"github.com/hiromaily/go-crypto-wallet/internal/di"
	wallets "github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/wallet"
)

// AddCommands adds all monitor subcommands
func AddCommands(parentCmd *cobra.Command, wallet *wallets.Watcher, container di.Container) {
	// senttx command
	var senttxAccount string
	senttxCmd := &cobra.Command{
		Use:   "senttx",
		Short: "monitor sent transactions",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSentTx(container, senttxAccount)
		},
	}
	senttxCmd.Flags().StringVar(&senttxAccount, "account", "", "account for monitoring")
	parentCmd.AddCommand(senttxCmd)

	// balance command
	var balanceConfirmationNum uint64
	balanceCmd := &cobra.Command{
		Use:   "balance",
		Short: "monitor balance",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runBalance(container, balanceConfirmationNum)
		},
	}
	balanceCmd.Flags().Uint64Var(&balanceConfirmationNum, "num", 6, "confirmation number")
	parentCmd.AddCommand(balanceCmd)
}
