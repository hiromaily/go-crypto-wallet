package create

import (
	"github.com/spf13/cobra"

	"github.com/hiromaily/go-crypto-wallet/internal/di"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
)

// AddCommands adds all create subcommands
func AddCommands(parentCmd *cobra.Command, wallet *wallets.Watcher, container di.Container) {
	// deposit command
	var depositFee float64
	depositCmd := &cobra.Command{
		Use:   "deposit",
		Short: "create a deposit unsigned transaction file for client account",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDeposit(container, depositFee)
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
			return runPayment(container, paymentFee)
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
			return runTransfer(container, transferAccount1, transferAccount2, transferAmount, transferFee)
		},
	}
	transferCmd.Flags().StringVar(&transferAccount1, "account1", "", "sender account")
	transferCmd.Flags().StringVar(&transferAccount2, "account2", "", "receiver account")
	transferCmd.Flags().Float64Var(
		&transferAmount, "amount", 0, "amount to send coin. if amount=0, all coin is sent")
	transferCmd.Flags().Float64Var(&transferFee, "fee", 0, "adjustment fee")
	parentCmd.AddCommand(transferCmd)

	// db command
	var dbTable string
	dbCmd := &cobra.Command{
		Use:        "db",
		Short:      "create payment_request table with dummy data for development use",
		Deprecated: "Use query with shell script instead of go code",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDB(container, dbTable)
		},
	}
	dbCmd.Flags().StringVar(&dbTable, "table", "", "target table name")
	parentCmd.AddCommand(dbCmd)
}
