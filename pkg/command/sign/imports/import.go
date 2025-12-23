package imports

import (
	"github.com/spf13/cobra"

	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
)

// AddCommands adds all import subcommands
func AddCommands(parentCmd *cobra.Command, wallet *wallets.Signer) {
	// privkey command
	privkeyCmd := &cobra.Command{
		Use:   "privkey",
		Short: "import generated private key for Authorization account to database",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPrivKey(*wallet)
		},
	}
	parentCmd.AddCommand(privkeyCmd)
}
