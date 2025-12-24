package sign

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/hiromaily/go-crypto-wallet/pkg/command/keygen/api/btc"
	"github.com/hiromaily/go-crypto-wallet/pkg/command/sign/create"
	"github.com/hiromaily/go-crypto-wallet/pkg/command/sign/export"
	"github.com/hiromaily/go-crypto-wallet/pkg/command/sign/imports"
	"github.com/hiromaily/go-crypto-wallet/pkg/command/sign/sign"
	ethapi "github.com/hiromaily/go-crypto-wallet/pkg/command/watch/api/eth"
	"github.com/hiromaily/go-crypto-wallet/pkg/di"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets/btcwallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets/ethwallet"
)

// AddCommands adds all sign subcommands to the root command
func AddCommands(rootCmd *cobra.Command, wallet *wallets.Signer, container di.Container, version string) {
	// Create command
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "create resources",
	}
	rootCmd.AddCommand(createCmd)
	create.AddCommands(createCmd, wallet, container)

	// Export command
	exportCmd := &cobra.Command{
		Use:   "export",
		Short: "export resources",
	}
	rootCmd.AddCommand(exportCmd)
	export.AddCommands(exportCmd, wallet, container)

	// Import command
	importCmd := &cobra.Command{
		Use:   "import",
		Short: "import resources",
	}
	rootCmd.AddCommand(importCmd)
	imports.AddCommands(importCmd, wallet, container)

	// Sign command
	signCmd := &cobra.Command{
		Use:   "sign",
		Short: "sign unsigned transaction",
	}
	rootCmd.AddCommand(signCmd)
	sign.AddCommands(signCmd, wallet, container)

	// API command - wallet-type specific, dynamically configured
	apiCmd := &cobra.Command{
		Use:   "api",
		Short: "API commands for the selected coin",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if *wallet == nil {
				return errors.New("wallet not initialized, check --coin flag")
			}
			// Clear existing subcommands to handle multiple runs in tests
			cmd.ResetCommands()
			switch v := (*wallet).(type) {
			case *btcwallet.BTCSign:
				btc.AddCommands(cmd, v.BTC)
			case *ethwallet.ETHSign:
				ethapi.AddCommands(cmd, v.ETH)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// This will run if no subcommand is given, e.g., `sign api`
			return cmd.Help()
		},
	}
	rootCmd.AddCommand(apiCmd)
}
