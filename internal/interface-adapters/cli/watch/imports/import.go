package imports

import (
	"github.com/spf13/cobra"

	"github.com/hiromaily/go-crypto-wallet/internal/di"
	wallets "github.com/hiromaily/go-crypto-wallet/internal/wallet"
)

// AddCommands adds all import subcommands
func AddCommands(parentCmd *cobra.Command, wallet *wallets.Watcher, container di.Container) {
	// address command
	var (
		addressFilePath string
		addressIsRescan bool
	)
	addressCmd := &cobra.Command{
		Use:   "address",
		Short: "import generated addresses by keygen wallet",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAddress(container, addressFilePath, addressIsRescan)
		},
	}
	addressCmd.Flags().StringVar(&addressFilePath, "file", "", "import file path for generated addresses")
	addressCmd.Flags().BoolVar(&addressIsRescan, "rescan", false, "run rescan when importing addresses or not")
	parentCmd.AddCommand(addressCmd)
}
