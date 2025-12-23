package export

import (
	"github.com/spf13/cobra"

	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
)

// AddCommands adds all export subcommands
func AddCommands(parentCmd *cobra.Command, wallet *wallets.Signer) {
	// fullpubkey command
	fullpubkeyCmd := &cobra.Command{
		Use:   "fullpubkey",
		Short: "export full pubkey",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runFullPubkey(*wallet)
		},
	}
	parentCmd.AddCommand(fullpubkeyCmd)
}
