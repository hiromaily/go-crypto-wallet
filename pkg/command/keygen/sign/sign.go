package sign

import (
	"github.com/spf13/cobra"

	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
)

// AddCommands adds all sign subcommands
func AddCommands(parentCmd *cobra.Command, wallet *wallets.Keygener) {
	// signature command
	var signatureFile string
	signatureCmd := &cobra.Command{
		Use:   "signature",
		Short: "sign on unsigned transaction (account would be found from file name)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSignature(*wallet, signatureFile)
		},
	}
	signatureCmd.Flags().StringVar(&signatureFile, "file", "", "import file path for signed transactions")
	parentCmd.AddCommand(signatureCmd)
}
