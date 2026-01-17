package main

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"

	"github.com/hiromaily/go-crypto-wallet/internal/di"
	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/cli/app"
	"github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/cli/sign"
	wallets "github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/wallet"
)

// sign wallet as cold wallet
//   - generate one key and seed for only auth accounts
//   - sign on unsigned transaction as second or more signature
//     (multisig addresses require signature)

var (
	walletType = domainWallet.WalletTypeSign
	appName    = walletType.String()
	// appVersion is set via ldflags at build time: -X main.appVersion=<version>
	appVersion = "dev"
	// authName is set via ldflags at build time: -X main.authName=<name>
	// used as account name like client, deposit, payment
	authName = ""

	// CLI options
	opts app.Options

	// Wallet instance (initialized in PersistentPreRunE)
	walleter  wallets.Signer
	container di.Container
)

func main() {
	rootCmd := &cobra.Command{
		Use:     appName,
		Short:   "Sign wallet for additional signatures on multisig transactions",
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
			walleter, err = container.NewSigner(authName)
			if err != nil {
				return err
			}

			return nil
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			if walleter != nil {
				walleter.Done()
			}
		},
	}

	// Add global flags with sign-specific coin help
	app.AddGlobalFlags(rootCmd, &opts, "coin type code: btc, bch")

	// Add subcommands
	sign.AddCommands(rootCmd, &walleter, func() di.Container { return container }, appVersion)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
