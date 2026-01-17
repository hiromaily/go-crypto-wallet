package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"

	"github.com/hiromaily/go-crypto-wallet/internal/di"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/cli/app"
	"github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/cli/sign"
	wallets "github.com/hiromaily/go-crypto-wallet/internal/interface-adapters/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/config"
)

// sign wallet as cold wallet
//   - generate one key and seed for only auth accounts
//   - sign on unsigned transaction as second or more signature
//     (multisig addresses require signature)

var (
	walletType = domainWallet.WalletTypeSign
	appName    = walletType.String()
	appVersion = "5.0.0"
	// used as account name like client, deposit, payment
	// this value is supposed to be embedded when building
	authName = ""

	// CLI options
	opts app.Options

	// Wallet instance (initialized in PersistentPreRunE)
	walleter  wallets.Signer
	container di.Container
)

func initializeSignWallet() error {
	// Validate coin type for sign wallet (BTC/BCH only)
	if err := app.ValidateCoinTypeForSign(opts.CoinTypeCode); err != nil {
		return err
	}

	// Validate config path is provided
	if opts.ConfigPath == "" {
		return errors.New("--config flag is required")
	}

	// Load wallet config
	conf, err := config.NewWallet(opts.ConfigPath, walletType, domainCoin.CoinTypeCode(opts.CoinTypeCode))
	if err != nil {
		return fmt.Errorf("failed to load wallet config: %w", err)
	}

	// Load account config (optional)
	accountConf := &config.AccountRoot{}
	if opts.AccountConfigPath != "" {
		accountConf, err = config.NewAccount(opts.AccountConfigPath)
		if err != nil {
			return fmt.Errorf("failed to load account config: %w", err)
		}
	}

	// Override config with CLI values
	conf.CoinTypeCode = domainCoin.CoinTypeCode(opts.CoinTypeCode)

	// Override Bitcoin host for specific wallet
	if opts.BTCWallet != "" {
		conf.Bitcoin.Host = fmt.Sprintf("%s/wallet/%s", conf.Bitcoin.Host, opts.BTCWallet)
		log.Println("conf.Bitcoin.Host:", conf.Bitcoin.Host)
	}

	// Create DI container and wallet
	container = di.NewContainer(conf, accountConf, walletType)
	walleter = container.NewSigner(authName)

	return nil
}

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
			return initializeSignWallet()
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			if walleter != nil {
				walleter.Done()
			}
		},
	}

	// Add global flags with sign-specific coin help
	app.AddGlobalFlagsWithCoinOptions(rootCmd, &opts, "coin type code: btc, bch")

	// Add subcommands
	sign.AddCommands(rootCmd, &walleter, func() di.Container { return container }, appVersion)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
