package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/command"
	wcmd "github.com/hiromaily/go-crypto-wallet/pkg/command/watch"
	"github.com/hiromaily/go-crypto-wallet/pkg/config"
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
	appVersion = "3.0.0"
)

func main() {
	// command line
	var (
		confPath        string
		accountConfPath string
		btcWallet       string
		coinTypeCode    string
		isHelp          bool
		isVersion       bool
		walleter        wallets.Watcher
	)
	flags := flag.NewFlagSet("main", flag.ContinueOnError)
	flags.StringVar(&confPath, "conf", "", "config file path")
	flags.StringVar(&coinTypeCode, "coin", "btc", "coin type code `btc`, `bch`, `eth`, `xrp`, `hyt`")
	flags.StringVar(&btcWallet, "wallet", "", "specify wallet.dat in bitcoin core")
	flags.BoolVar(&isVersion, "version", false, "show version")
	flags.BoolVar(&isHelp, "help", false, "show help")
	if err := flags.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}

	// version
	if isVersion {
		fmt.Printf("%s v%s\n", appName, appVersion)
		os.Exit(0)
	}

	// validate coinTypeCode
	if !coin.IsCoinTypeCode(coinTypeCode) && !coin.IsERC20Token(coinTypeCode) {
		log.Fatal("coin args is invalid. `btc`, `bch`, `eth`, `xrp`, `hyt` is allowed")
	}
	// for ERC-20 token
	var erc20Token string
	if coin.IsERC20Token(coinTypeCode) {
		erc20Token = coinTypeCode
		coinTypeCode = coin.ERC20.String()
	}

	// set config path if environment variable is existing
	if confPath == "" {
		switch coinTypeCode {
		case coin.BTC.String():
			confPath = os.Getenv("BTC_WATCH_WALLET_CONF")
		case coin.BCH.String():
			confPath = os.Getenv("BCH_WATCH_WALLET_CONF")
		case coin.ETH.String(), coin.ERC20.String():
			confPath = os.Getenv("ETH_WATCH_WALLET_CONF")
		case coin.XRP.String():
			confPath = os.Getenv("XRP_WATCH_WALLET_CONF")
		}
	}
	// account conf path for account settings
	if accountConfPath == "" {
		switch coinTypeCode {
		case coin.BTC.String():
			accountConfPath = os.Getenv("BTC_ACCOUNT_CONF")
		case coin.BCH.String():
			accountConfPath = os.Getenv("BCH_ACCOUNT_CONF")
		case coin.ETH.String(), coin.ERC20.String():
			accountConfPath = os.Getenv("ETH_ACCOUNT_CONF")
		case coin.XRP.String():
			accountConfPath = os.Getenv("XRP_ACCOUNT_CONF")
		}
	}

	// help
	var conf *config.WalletRoot
	if !isHelp && len(os.Args) > 1 {
		var err error

		// config
		conf, err = config.NewWallet(confPath, walletType, coin.CoinTypeCode(coinTypeCode))
		if err != nil {
			log.Fatal(err)
		}
		accountConf := &account.AccountRoot{}
		if accountConfPath != "" {
			accountConf, err = account.NewAccount(accountConfPath)
			if err != nil {
				log.Fatal(err)
			}
		}

		// override config
		conf.CoinTypeCode = coin.CoinTypeCode(coinTypeCode)
		if erc20Token != "" {
			conf.Ethereum.ERC20Token = coin.ERC20Token(erc20Token)
			if conf.ValidateERC20(conf.Ethereum.ERC20Token) != err {
				log.Fatal(err)
			}
		}

		// - conf.Bitcoin.Host
		if btcWallet != "" {
			conf.Bitcoin.Host = fmt.Sprintf("%s/wallet/%s", conf.Bitcoin.Host, btcWallet)
			log.Println("conf.Bitcoin.Host:", conf.Bitcoin.Host)
		}

		// create wallet
		reg := NewRegistry(conf, accountConf, walletType)
		walleter = reg.NewWalleter()
	}
	defer func() {
		walleter.Done()
	}()

	// sub command
	args := flags.Args()
	cmds := wcmd.WatchSubCommands(walleter, appVersion, conf)
	cl := command.CreateSubCommand(appName, appVersion, args, cmds)
	cl.HelpFunc = command.HelpFunc(cl.Name)

	flags.Usage = func() { fmt.Println(cl.HelpFunc(cl.Commands)) }

	code, err := cl.Run()
	if err != nil {
		log.Printf("fail to call Run() %s command: %v", appName, err)
	}
	os.Exit(code)
}
