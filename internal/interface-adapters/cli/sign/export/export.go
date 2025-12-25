package export

import (
	"github.com/spf13/cobra"

	"github.com/hiromaily/go-crypto-wallet/pkg/di"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
)

// AddCommands adds all export subcommands
func AddCommands(parentCmd *cobra.Command, wallet *wallets.Signer, container di.Container) {
	// fullpubkey command
	fullpubkeyCmd := &cobra.Command{
		Use:   "fullpubkey",
		Short: "export full pubkey",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runFullPubkey(container)
		},
	}
	parentCmd.AddCommand(fullpubkeyCmd)
}
