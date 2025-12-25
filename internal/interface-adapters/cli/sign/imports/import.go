package imports

import (
	"github.com/spf13/cobra"

	"github.com/hiromaily/go-crypto-wallet/internal/di"
	wallets "github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/wallet"
)

// AddCommands adds all import subcommands
func AddCommands(parentCmd *cobra.Command, wallet *wallets.Signer, container di.Container) {
	// privkey command
	privkeyCmd := &cobra.Command{
		Use:   "privkey",
		Short: "import generated private key for Authorization account to database",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPrivKey(container)
		},
	}
	parentCmd.AddCommand(privkeyCmd)
}
