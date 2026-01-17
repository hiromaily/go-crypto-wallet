package main

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"

	"github.com/hiromaily/go-crypto-wallet/internal/di"
	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/cli/app"
	"github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/cli/keygen"
	wallets "github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/wallet"
)

// keygen wallet as cold wallet
//   - generate key and seed for accounts
//   - create multisig address with full pubkey of auth accounts
//   - sign on unsigned transaction as first signature
//     (signature would not be completed if address is multisig)

var (
	walletType = domainWallet.WalletTypeKeyGen
	appName    = walletType.String()
	// appVersion is set via ldflags at build time: -X main.appVersion=<version>
	appVersion = "dev"

	// CLI options
	opts app.Options

	// Wallet instance (initialized in PersistentPreRunE)
	walleter  wallets.Keygener
	container di.Container
)

func main() {
	rootCmd := &cobra.Command{
		Use:     appName,
		Short:   "Keygen wallet for key generation and first signature",
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

			// Store container and create wallet
			container = application.Container
			walleter = container.NewKeygener()

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
	keygen.AddCommands(rootCmd, &walleter, func() di.Container { return container }, appVersion)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
