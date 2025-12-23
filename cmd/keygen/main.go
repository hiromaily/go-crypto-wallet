package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/command/keygen"
	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	"github.com/hiromaily/go-crypto-wallet/pkg/di"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
)

// keygen wallet as cold wallet
//  - generate key and seed for accounts
//  - create multisig address with full pubkey of auth accounts
//  - sing on unsigned transaction as first signature
//   (signature would not be completed if address is multisig)

// TODO: bitcoin functionalities
// - encrypt wallet itself by `encryptwallet` command
// - passphrase would be required when using secret key to sign unsigned transaction

var (
	walletType = wallet.WalletTypeKeyGen
	appName    = walletType.String()
	appVersion = "5.0.0"

	// Global flags
	confPath        string
	accountConfPath string
	btcWallet       string
	coinTypeCode    string

	// Wallet instance
	walleter wallets.Keygener
)

func initializeWallet() error {
	// validate coinTypeCode
	if !coin.IsCoinTypeCode(coinTypeCode) {
		return errors.New("coin args is invalid. `btc`, `bch`, `eth`, `xrp`, `hyt` is allowed")
	}

	// set config path if environment variable is existing
	if confPath == "" {
		setConfigPathFromEnv()
	}

	// account conf path for account settings
	if accountConfPath == "" {
		setAccountConfPathFromEnv()
	}

	// base config
	conf, err := config.NewWallet(confPath, walletType, coin.CoinTypeCode(coinTypeCode))
	if err != nil {
		return err
	}

	// account config
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
	walleter = container.NewKeygener()

	return nil
}

func setConfigPathFromEnv() {
	switch {
	case coinTypeCode == coin.BTC.String():
		confPath = os.Getenv("BTC_KEYGEN_WALLET_CONF")
	case coinTypeCode == coin.BCH.String():
		confPath = os.Getenv("BCH_KEYGEN_WALLET_CONF")
	case coin.IsETHGroup(coin.CoinTypeCode(coinTypeCode)):
		confPath = os.Getenv("ETH_KEYGEN_WALLET_CONF")
	case coinTypeCode == coin.XRP.String():
		confPath = os.Getenv("XRP_KEYGEN_WALLET_CONF")
	}
}

func setAccountConfPathFromEnv() {
	switch coinTypeCode {
	case coin.BTC.String():
		accountConfPath = os.Getenv("BTC_ACCOUNT_CONF")
	case coin.BCH.String():
		accountConfPath = os.Getenv("BCH_ACCOUNT_CONF")
	}
}

func main() {
	rootCmd := &cobra.Command{
		Use:     appName,
		Short:   "Keygen wallet for key generation and first signature",
		Version: appVersion,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Skip initialization for help and version
			if cmd.Name() == "help" || cmd.Name() == "version" {
				return nil
			}
			return initializeWallet()
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			if walleter != nil {
				walleter.Done()
			}
		},
	}

	// Global flags
	rootCmd.PersistentFlags().StringVarP(&confPath, "conf", "c", "", "config file path")
	rootCmd.PersistentFlags().StringVar(&coinTypeCode, "coin", "btc",
		"coin type code `btc`, `bch`, `eth`, `xrp`, `hyt`")
	rootCmd.PersistentFlags().StringVarP(&btcWallet, "wallet", "w", "", "specify wallet.dat in bitcoin core")

	// Add subcommands
	keygen.AddCommands(rootCmd, &walleter, appVersion)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
