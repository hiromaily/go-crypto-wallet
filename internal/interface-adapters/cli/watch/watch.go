package watch

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/hiromaily/go-crypto-wallet/internal/di"
	"github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/cli/watch/api/btc"
	"github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/cli/watch/api/eth"
	"github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/cli/watch/api/xrp"
	"github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/cli/watch/create"
	"github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/cli/watch/imports"
	"github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/cli/watch/monitor"
	"github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/cli/watch/send"
	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets/btcwallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets/ethwallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets/xrpwallet"
)

// AddCommands adds all watch subcommands to the root command
func AddCommands(
	rootCmd *cobra.Command,
	wallet *wallets.Watcher,
	container di.Container,
	version string,
	confPtr *config.WalletRoot,
) {
	// Import command
	importCmd := &cobra.Command{
		Use:   "import",
		Short: "import resources",
	}
	rootCmd.AddCommand(importCmd)
	imports.AddCommands(importCmd, wallet, container)

	// Create command
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "create resources",
	}
	rootCmd.AddCommand(createCmd)
	create.AddCommands(createCmd, wallet, container)

	// Send command
	sendCmd := send.AddCommand(wallet, container)
	rootCmd.AddCommand(sendCmd)

	// Monitor command
	monitorCmd := &cobra.Command{
		Use:   "monitor",
		Short: "monitor resources",
	}
	rootCmd.AddCommand(monitorCmd)
	monitor.AddCommands(monitorCmd, wallet, container)

	// API command - wallet-type specific, dynamically configured
	apiCmd := &cobra.Command{
		Use:   "api",
		Short: "API commands for the selected coin",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if *wallet == nil {
				return errors.New("wallet not initialized, check --coin flag")
			}
			if confPtr == nil {
				return errors.New("config not initialized")
			}
			// Clear existing subcommands to handle multiple runs in tests
			cmd.ResetCommands()
			switch v := (*wallet).(type) {
			case *btcwallet.BTCWatch:
				btc.AddCommands(cmd, v.BTC)
			case *ethwallet.ETHWatch:
				eth.AddCommands(cmd, v.ETH)
			case *xrpwallet.XRPWatch:
				xrp.AddCommands(cmd, v.XRP, &confPtr.Ripple.API.TxData)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// This will run if no subcommand is given, e.g., `watch api`
			return cmd.Help()
		},
	}
	rootCmd.AddCommand(apiCmd)
}
