package sign

import (
	"github.com/spf13/cobra"

	"github.com/hiromaily/go-crypto-wallet/internal/di"
	wallets "github.com/hiromaily/go-crypto-wallet/internal/wallet"
)

// AddCommands adds all sign subcommands
func AddCommands(parentCmd *cobra.Command, wallet *wallets.Keygener, container di.Container) {
	// signature command
	var signatureFile string
	signatureCmd := &cobra.Command{
		Use:   "signature",
		Short: "sign on unsigned transaction (account would be found from file name)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSignature(container, signatureFile)
		},
	}
	signatureCmd.Flags().StringVar(&signatureFile, "file", "", "import file path for signed transactions")
	parentCmd.AddCommand(signatureCmd)
}
