package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/command/sign"
	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	"github.com/hiromaily/go-crypto-wallet/pkg/di"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
)

// sign wallet as cold wallet
//  - generate one key and seed for only auth accounts
//  - sing on unsigned transaction as second or more signature
//   (multisig addresses require signature)

// TODO: bitcoin functionalities
// - encrypt wallet itself by `encryptwallet` command
// - passphrase would be required when using secret key to sign unsigned transaction

var (
	walletType = wallet.WalletTypeSign
	appName    = walletType.String()
	appVersion = "5.0.0"
	// used as account name like client, deposit, payment
	// this value is supposed to be embedded when building
	authName = ""

	// Global flags
	confPath        string
	accountConfPath string
	btcWallet       string
	coinTypeCode    string

	// Wallet instance
	walleter wallets.Signer
)

func main() {
	rootCmd := &cobra.Command{
		Use:     appName,
		Short:   "Sign wallet for additional signatures on multisig transactions",
		Version: appVersion,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Skip initialization for help and version
			if cmd.Name() == "help" || cmd.Name() == "version" {
				return nil
			}

			// validate coinTypeCode
			if !coin.IsCoinTypeCode(coinTypeCode) {
				return errors.New("coin args is invalid. `btc`, `bch` is allowed")
			}

			// set config path if environment variable is existing
			if confPath == "" {
				switch coinTypeCode {
				case coin.BTC.String():
					confPath = os.Getenv("BTC_SIGN_WALLET_CONF")
				case coin.BCH.String():
					confPath = os.Getenv("BCH_SIGN_WALLET_CONF")
				}
			}

			// account conf path for multisig
			if accountConfPath == "" {
				switch coinTypeCode {
				case coin.BTC.String():
					accountConfPath = os.Getenv("BTC_ACCOUNT_CONF")
				case coin.BCH.String():
					accountConfPath = os.Getenv("BCH_ACCOUNT_CONF")
				}
			}

			// config
			conf, err := config.NewWallet(confPath, walletType, coin.CoinTypeCode(coinTypeCode))
			if err != nil {
				return err
			}

			accountConf := &account.AccountRoot{}
			if accountConfPath != "" {
				accountConf, err = account.NewAccount(accountConfPath)
				if err != nil {
					return err
				}
			}

			// override config
			conf.CoinTypeCode = coin.CoinTypeCode(coinTypeCode)

			// override conf.Bitcoin.Host
			if btcWallet != "" {
				conf.Bitcoin.Host = fmt.Sprintf("%s/wallet/%s", conf.Bitcoin.Host, btcWallet)
				log.Println("conf.Bitcoin.Host:", conf.Bitcoin.Host)
			}

			// create wallet
			container := di.NewContainer(conf, accountConf, walletType)
			walleter = container.NewSigner(authName)

			return nil
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			if walleter != nil {
				walleter.Done()
			}
		},
	}

	// Global flags
	rootCmd.PersistentFlags().StringVarP(&confPath, "conf", "c", "", "config file path")
	rootCmd.PersistentFlags().StringVar(&coinTypeCode, "coin", "btc", "coin type code `btc`, `bch`")
	rootCmd.PersistentFlags().StringVarP(&btcWallet, "wallet", "w", "", "specify wallet.dat in bitcoin core")

	// Add subcommands
	sign.AddCommands(rootCmd, &walleter, appVersion)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
