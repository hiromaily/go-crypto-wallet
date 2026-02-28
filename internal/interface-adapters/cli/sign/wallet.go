package sign

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hiromaily/go-crypto-wallet/internal/di"
	"github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/cli/keygen/api/btc"
	"github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/cli/sign/create"
	"github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/cli/sign/export"
	"github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/cli/sign/imports"
	"github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/cli/sign/sign"
	wallets "github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/wallet"
	btcwallet "github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/wallet/btc"
)

// AddCommands adds all sign subcommands to the root command
func AddCommands(rootCmd *cobra.Command, wallet *wallets.Signer, containerGetter func() di.Container, version string) {
	// Create command
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "create resources",
	}
	rootCmd.AddCommand(createCmd)
	create.AddCommands(createCmd, wallet, containerGetter)

	// Export command
	exportCmd := &cobra.Command{
		Use:   "export",
		Short: "export resources",
	}
	rootCmd.AddCommand(exportCmd)
	export.AddCommands(exportCmd, wallet, containerGetter)

	// Import command
	importCmd := &cobra.Command{
		Use:   "import",
		Short: "import resources",
	}
	rootCmd.AddCommand(importCmd)
	imports.AddCommands(importCmd, wallet, containerGetter)

	// Sign command
	signCmd := &cobra.Command{
		Use:   "sign",
		Short: "sign unsigned transaction",
	}
	rootCmd.AddCommand(signCmd)
	sign.AddCommands(signCmd, wallet, containerGetter)

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
			default:
				fmt.Printf("[WARN] api command is not supported for this coin type\n")
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
