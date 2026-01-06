package keygen

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/hiromaily/go-crypto-wallet/internal/di"
	btcapi "github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/cli/keygen/api/btc"
	"github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/cli/keygen/api/eth"
	btckeygen "github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/cli/keygen/btc"
	"github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/cli/keygen/create"
	"github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/cli/keygen/export"
	"github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/cli/keygen/imports"
	"github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/cli/keygen/sign"
	wallets "github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/wallet"
	btcwallet "github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/wallet/btc"
	ethwallet "github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/wallet/eth"
)

// AddCommands adds all keygen subcommands to the root command
func AddCommands(rootCmd *cobra.Command, wallet *wallets.Keygener, containerGetter func() di.Container, version string) {
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

	// Descriptor commands (BTC only)
	btckeygen.AddDescriptorCommands(rootCmd, containerGetter)

	// MuSig2 command (BTC only)
	AddMuSig2Commands(rootCmd, containerGetter)

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
			case *btcwallet.BTCKeygen:
				btcapi.AddCommands(cmd, v.BTC)
			case *ethwallet.ETHKeygen:
				eth.AddCommands(cmd, v.ETH)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// This will run if no subcommand is given, e.g., `keygen api`
			return cmd.Help()
		},
	}
	rootCmd.AddCommand(apiCmd)
}
