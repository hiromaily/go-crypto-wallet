package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/hiromaily/go-bitcoin/pkg/command"
	wcmd "github.com/hiromaily/go-bitcoin/pkg/command/wallet"
	"github.com/hiromaily/go-bitcoin/pkg/config"
	"github.com/hiromaily/go-bitcoin/pkg/wallet"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/coin"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets"
)

// watch as watch only wallet
//  this wallet works online, so bitcoin network is required to call APIs
//  create unsigned transaction
//  send signed transaction

//TODO: bitcoin functionalities
// - back up wallet data periodically and import functionality
// - generated key must be encrypted
// - transfer for monitoring
//TODO:
// - logger interface: stdout(ui), log format, open tracing
// - btc command for mock is required

var (
	walletType = wallet.WalletTypeWatchOnly
	appName    = walletType.String()
	appVersion = "2.3.0"
)

func main() {
	// command line
	var (
		confPath  string
		btcWallet string
		isHelp    bool
		isVersion bool
		walleter  wallets.Watcher
	)
	flags := flag.NewFlagSet("main", flag.ContinueOnError)
	flags.StringVar(&confPath, "conf", os.Getenv("WATCH_WALLET_CONF"), "config file path")
	flags.StringVar(&btcWallet, "wallet", "", "specify wallet in bitcoin core")
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

	// help
	if !isHelp && len(os.Args) > 1 {
		// config
		conf, err := config.New(confPath, walletType)
		if err != nil {
			log.Fatal(err)
		}
		// override conf.Bitcoin.Host
		if btcWallet != "" {
			conf.Bitcoin.Host = fmt.Sprintf("%s/wallet/%s", conf.Bitcoin.Host, btcWallet)
			log.Println("conf.Bitcoin.Host:", conf.Bitcoin.Host)
		}
		// override conf.CoinTypeCode
		if os.Getenv("COIN_TYPE") != "" && coin.ValidateCoinTypeCode(os.Getenv("COIN_TYPE")) {
			conf.CoinTypeCode = coin.CoinTypeCode(os.Getenv("COIN_TYPE"))
			log.Println("conf.CoinTypeCode:", conf.CoinTypeCode)
		}

		// create wallet
		regi := NewRegistry(conf, walletType)
		walleter = regi.NewWalleter()
	}
	defer func() {
		walleter.Done()
	}()

	//sub command
	args := flags.Args()
	cmds := wcmd.WalletSubCommands(walleter, appVersion)
	cl := command.CreateSubCommand(appName, appVersion, args, cmds)
	cl.HelpFunc = command.HelpFunc(cl.Name)

	flags.Usage = func() { fmt.Println(cl.HelpFunc(cl.Commands)) }

	code, err := cl.Run()
	if err != nil {
		log.Printf("fail to call Run() %s command: %v", appName, err)
	}
	os.Exit(code)
}
