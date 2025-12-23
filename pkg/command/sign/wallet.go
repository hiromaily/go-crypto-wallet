package sign

import (
	"github.com/spf13/cobra"

	"github.com/hiromaily/go-crypto-wallet/pkg/command/keygen/api/btc"
	"github.com/hiromaily/go-crypto-wallet/pkg/command/sign/create"
	"github.com/hiromaily/go-crypto-wallet/pkg/command/sign/export"
	"github.com/hiromaily/go-crypto-wallet/pkg/command/sign/imports"
	"github.com/hiromaily/go-crypto-wallet/pkg/command/sign/sign"
	ethapi "github.com/hiromaily/go-crypto-wallet/pkg/command/watch/api/eth"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets/btcwallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets/ethwallet"
)

// AddCommands adds all sign subcommands to the root command
func AddCommands(rootCmd *cobra.Command, wallet *wallets.Signer, version string) {
	// Create command
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "create resources",
	}
	rootCmd.AddCommand(createCmd)
	create.AddCommands(createCmd, wallet)

	// Export command
	exportCmd := &cobra.Command{
		Use:   "export",
		Short: "export resources",
	}
	rootCmd.AddCommand(exportCmd)
	export.AddCommands(exportCmd, wallet)

	// Import command
	importCmd := &cobra.Command{
		Use:   "import",
		Short: "import resources",
	}
	rootCmd.AddCommand(importCmd)
	imports.AddCommands(importCmd, wallet)

	// Sign command
	signCmd := &cobra.Command{
		Use:   "sign",
		Short: "sign unsigned transaction",
	}
	rootCmd.AddCommand(signCmd)
	sign.AddCommands(signCmd, wallet)

	// API commands - wallet-type specific
	rootCmd.AddCommand(&cobra.Command{
		Use:   "api",
		Short: "api commands",
		RunE: func(cmd *cobra.Command, args []string) error {
			if *wallet == nil {
				return nil
			}
			switch v := (*wallet).(type) {
			case *btcwallet.BTCSign:
				apiCmd := &cobra.Command{
					Use:   "api",
					Short: "Bitcoin API commands",
				}
				btc.AddCommands(apiCmd, v.BTC)
				return apiCmd.Execute()
			case *ethwallet.ETHSign:
				apiCmd := &cobra.Command{
					Use:   "api",
					Short: "Ethereum API commands",
				}
				ethapi.AddCommands(apiCmd, v.ETH)
				return apiCmd.Execute()
			}
			return nil
		},
	})
}
