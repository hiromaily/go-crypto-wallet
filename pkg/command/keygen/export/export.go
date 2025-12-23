package export

import (
	"github.com/spf13/cobra"

	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
)

// AddCommands adds all export subcommands
func AddCommands(parentCmd *cobra.Command, wallet *wallets.Keygener) {
	// address command
	var addressAccount string
	addressCmd := &cobra.Command{
		Use:   "address",
		Short: "export generated PublicKey as csv file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAddress(*wallet, addressAccount)
		},
	}
	addressCmd.Flags().StringVar(&addressAccount, "account", "", "target account")
	parentCmd.AddCommand(addressCmd)
}
