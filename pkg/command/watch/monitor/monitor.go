package monitor

import (
	"github.com/spf13/cobra"

	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
)

// AddCommands adds all monitor subcommands
func AddCommands(parentCmd *cobra.Command, wallet *wallets.Watcher) {
	// senttx command
	var senttxAccount string
	senttxCmd := &cobra.Command{
		Use:   "senttx",
		Short: "monitor sent transactions",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSentTx(*wallet, senttxAccount)
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
			return runBalance(*wallet, balanceConfirmationNum)
		},
	}
	balanceCmd.Flags().Uint64Var(&balanceConfirmationNum, "num", 6, "confirmation number")
	parentCmd.AddCommand(balanceCmd)
}
