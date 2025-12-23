package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	wcmd "github.com/hiromaily/go-crypto-wallet/pkg/command/watch"
	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	"github.com/hiromaily/go-crypto-wallet/pkg/di"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
)

// watch as watch only wallet
//  this wallet works online, so bitcoin network is required to call APIs
//  create unsigned transaction
//  send signed transaction

// TODO: bitcoin functionalities
// - back up wallet data periodically and import functionality
// - generated key must be encrypted
// - transfer for monitoring
// TODO:
// - logger interface: stdout(ui), log format, open tracing
// - btc command for mock is required

var (
	walletType = wallet.WalletTypeWatchOnly
	appName    = walletType.String()
	appVersion = "5.0.0"

	// Global flags
	confPath        string
	accountConfPath string
	btcWallet       string
	coinTypeCode    string

	// Wallet and config instances
	walleter wallets.Watcher
	conf     *config.WalletRoot
)

func initializeWallet() error {
	// validate coinTypeCode
	if !coin.IsCoinTypeCode(coinTypeCode) && !coin.IsERC20Token(coinTypeCode) {
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

	var err error

	// config
	conf, err = config.NewWallet(confPath, walletType, coin.CoinTypeCode(coinTypeCode))
	if err != nil {
		return fmt.Errorf("failed to load wallet config: %w", err)
	}

	accountConf := &account.AccountRoot{}
	if accountConfPath != "" {
		accountConf, err = account.NewAccount(accountConfPath)
		if err != nil {
			return fmt.Errorf("failed to load account config: %w", err)
		}
	}

	// override config
	conf.CoinTypeCode = coin.CoinTypeCode(coinTypeCode)
	if coin.IsERC20Token(coinTypeCode) {
		if err := conf.ValidateERC20(coin.ERC20Token(coinTypeCode)); err != nil {
			return fmt.Errorf("failed to validate ERC20 token: %w", err)
		}
		conf.Ethereum.ERC20Token = coin.ERC20Token(coinTypeCode)
	}

	// - conf.Bitcoin.Host
	if btcWallet != "" {
		conf.Bitcoin.Host = fmt.Sprintf("%s/wallet/%s", conf.Bitcoin.Host, btcWallet)
		log.Println("conf.Bitcoin.Host:", conf.Bitcoin.Host)
	}

	// create wallet
	container := di.NewContainer(conf, accountConf, walletType)
	walleter = container.NewWalleter()

	return nil
}

func setConfigPathFromEnv() {
	switch {
	case coinTypeCode == coin.BTC.String():
		confPath = os.Getenv("BTC_WATCH_WALLET_CONF")
	case coinTypeCode == coin.BCH.String():
		confPath = os.Getenv("BCH_WATCH_WALLET_CONF")
	case coin.IsETHGroup(coin.CoinTypeCode(coinTypeCode)):
		confPath = os.Getenv("ETH_WATCH_WALLET_CONF")
	case coinTypeCode == coin.XRP.String():
		confPath = os.Getenv("XRP_WATCH_WALLET_CONF")
	}
}

func setAccountConfPathFromEnv() {
	switch {
	case coinTypeCode == coin.BTC.String():
		accountConfPath = os.Getenv("BTC_ACCOUNT_CONF")
	case coinTypeCode == coin.BCH.String():
		accountConfPath = os.Getenv("BCH_ACCOUNT_CONF")
	case coin.IsETHGroup(coin.CoinTypeCode(coinTypeCode)):
		accountConfPath = os.Getenv("ETH_ACCOUNT_CONF")
	case coinTypeCode == coin.XRP.String():
		accountConfPath = os.Getenv("XRP_ACCOUNT_CONF")
	}
}

func main() {
	rootCmd := &cobra.Command{
		Use:     appName,
		Short:   "Watch-only wallet for creating and sending transactions",
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
	wcmd.AddCommands(rootCmd, &walleter, appVersion, conf)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
