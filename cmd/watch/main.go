package main

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"

	"github.com/hiromaily/go-crypto-wallet/internal/di"
	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/cli/app"
	wcmd "github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/cli/watch"
	wallets "github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/config"
)

// watch as watch only wallet
//   this wallet works online, so bitcoin network is required to call APIs
//   create unsigned transaction
//   send signed transaction
// Debug: incremental number: 3

var (
	walletType = domainWallet.WalletTypeWatchOnly
	appName    = walletType.String()
	// appVersion is set via ldflags at build time: -X main.appVersion=<version>
	appVersion = "dev"

	// CLI options
	opts app.Options

	// Wallet and config instances (initialized in PersistentPreRunE)
	walleter  wallets.Watcher
	conf      *config.WalletRoot
	container di.Container
)

func main() {
	rootCmd := &cobra.Command{
		Use:     appName,
		Short:   "Watch-only wallet for creating and sending transactions",
		Version: appVersion,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Skip initialization for help command
			if cmd.Name() == "help" {
				return nil
			}

			// Initialize the application
			application, err := app.NewApp(walletType, opts)
			if err != nil {
				return err
			}

			// Store config pointer (needed for API commands)
			conf = application.Config

			// Store container and create wallet
			container = application.Container
			walleter = container.NewWalleter()

			return nil
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			if walleter != nil {
				walleter.Done()
			}
		},
	}

	// Add global flags
	app.AddGlobalFlags(rootCmd, &opts)

	// Add subcommands
	wcmd.AddCommands(rootCmd, &walleter, func() di.Container { return container }, appVersion, conf)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
