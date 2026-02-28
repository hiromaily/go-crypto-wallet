package create

import (
	"github.com/spf13/cobra"

	"github.com/hiromaily/go-crypto-wallet/internal/di"
	wallets "github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/wallet"
)

// AddCommands adds all create subcommands
func AddCommands(parentCmd *cobra.Command, wallet *wallets.Watcher, containerGetter func() di.Container) {
	// deposit command
	var depositFee float64
	depositCmd := &cobra.Command{
		Use:   "deposit",
		Short: "create a deposit unsigned transaction file for client account",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDeposit(containerGetter(), depositFee)
		},
	}
	depositCmd.Flags().Float64Var(&depositFee, "fee", 0, "adjustment fee")
	parentCmd.AddCommand(depositCmd)

	// payment command
	var paymentFee float64
	paymentCmd := &cobra.Command{
		Use:   "payment",
		Short: "create a payment unsigned transaction file for payment account",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPayment(containerGetter(), paymentFee)
		},
	}
	paymentCmd.Flags().Float64Var(&paymentFee, "fee", 0, "adjustment fee")
	parentCmd.AddCommand(paymentCmd)

	// transfer command
	var (
		transferAccount1 string
		transferAccount2 string
		transferAmount   float64
		transferFee      float64
	)
	transferCmd := &cobra.Command{
		Use:   "transfer",
		Short: "create unsigned transaction for transfer among accounts",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTransfer(containerGetter(), transferAccount1, transferAccount2, transferAmount, transferFee)
		},
	}
	transferCmd.Flags().StringVar(&transferAccount1, "account1", "", "sender account")
	transferCmd.Flags().StringVar(&transferAccount2, "account2", "", "receiver account")
	transferCmd.Flags().Float64Var(
		&transferAmount, "amount", 0, "amount to send coin. if amount=0, all coin is sent")
	transferCmd.Flags().Float64Var(&transferFee, "fee", 0, "adjustment fee")
	parentCmd.AddCommand(transferCmd)
}
