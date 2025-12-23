package watch

import (
	"github.com/spf13/cobra"

	"github.com/hiromaily/go-crypto-wallet/pkg/command/watch/api/btc"
	"github.com/hiromaily/go-crypto-wallet/pkg/command/watch/api/eth"
	"github.com/hiromaily/go-crypto-wallet/pkg/command/watch/api/xrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/command/watch/create"
	"github.com/hiromaily/go-crypto-wallet/pkg/command/watch/imports"
	"github.com/hiromaily/go-crypto-wallet/pkg/command/watch/monitor"
	"github.com/hiromaily/go-crypto-wallet/pkg/command/watch/send"
	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets/btcwallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets/ethwallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets/xrpwallet"
)

// AddCommands adds all watch subcommands to the root command
func AddCommands(rootCmd *cobra.Command, wallet *wallets.Watcher, version string, confPtr *config.WalletRoot) {
	// Import command
	importCmd := &cobra.Command{
		Use:   "import",
		Short: "import resources",
	}
	rootCmd.AddCommand(importCmd)
	imports.AddCommands(importCmd, wallet)

	// Create command
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "create resources",
	}
	rootCmd.AddCommand(createCmd)
	create.AddCommands(createCmd, wallet)

	// Send command
	rootCmd.AddCommand(send.AddCommand(wallet))

	// Monitor command
	monitorCmd := &cobra.Command{
		Use:   "monitor",
		Short: "monitor resources",
	}
	rootCmd.AddCommand(monitorCmd)
	monitor.AddCommands(monitorCmd, wallet)

	// API commands - wallet-type specific
	if *wallet == nil {
		return
	}

	switch v := (*wallet).(type) {
	case *btcwallet.BTCWatch:
		apiCmd := &cobra.Command{
			Use:   "api",
			Short: "Bitcoin API commands",
		}
		rootCmd.AddCommand(apiCmd)
		btc.AddCommands(apiCmd, v.BTC)
	case *ethwallet.ETHWatch:
		apiCmd := &cobra.Command{
			Use:   "api",
			Short: "Ethereum API commands",
		}
		rootCmd.AddCommand(apiCmd)
		eth.AddCommands(apiCmd, v.ETH)
	case *xrpwallet.XRPWatch:
		apiCmd := &cobra.Command{
			Use:   "api",
			Short: "Ripple API commands",
		}
		rootCmd.AddCommand(apiCmd)
		xrp.AddCommands(apiCmd, v.XRP, &confPtr.Ripple.API.TxData)
	}
}
