package xrp

import (
	"github.com/spf13/cobra"

	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/ripple"
	"github.com/hiromaily/go-crypto-wallet/pkg/config"
)

// AddCommands adds all Ripple API subcommands
func AddCommands(parentCmd *cobra.Command, xrp ripple.Rippler, txData *config.RippleTxData) {
	// sendcoin command
	var (
		sendcoinAddress string
		sendcoinAmount  float64
	)
	sendcoinCmd := &cobra.Command{
		Use:   "sendcoin",
		Short: "send coin from faucet coin",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSendCoin(xrp, txData, sendcoinAddress, sendcoinAmount)
		},
	}
	sendcoinCmd.Flags().StringVar(&sendcoinAddress, "address", "", "receiver address")
	sendcoinCmd.Flags().Float64Var(&sendcoinAmount, "amount", 0, "amount")
	parentCmd.AddCommand(sendcoinCmd)
}
