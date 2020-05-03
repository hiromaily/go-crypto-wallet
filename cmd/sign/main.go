package main

import (
	"flag"
	"fmt"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/coin"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/hiromaily/go-bitcoin/pkg/command"
	"github.com/hiromaily/go-bitcoin/pkg/command/sign"
	"github.com/hiromaily/go-bitcoin/pkg/config"
	"github.com/hiromaily/go-bitcoin/pkg/wallet"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets"
)

// sign wallet as cold wallet
//  - generate one key and seed for only auth accounts
//  - sing on unsigned transaction as second or more signature
//   (multisig addresses require signature)

//TODO: bitcoin functionalities
// - encrypt wallet itself by `encryptwallet` command
// - passphrase would be required when using secret key to sign unsigned transaction

var (
	walletType = wallet.WalletTypeSign
	appName    = walletType.String()
	appVersion = "2.3.0"
	authName   = "" // this value is supposed to be embedded when building
)

func main() {
	// command line
	var (
		confPath  string
		btcWallet string
		isHelp    bool
		isVersion bool
		walleter  wallets.Signer
	)
	flags := flag.NewFlagSet("main", flag.ContinueOnError)
	flags.StringVar(&confPath, "conf", os.Getenv("SIGN_WALLET_CONF"), "config file path")
	flags.StringVar(&btcWallet, "wallet", "", "specify wallet in bitcoin core")
	flags.BoolVar(&isVersion, "version", false, "show version")
	flags.BoolVar(&isHelp, "help", false, "show help")
	if err := flags.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}

	// version
	if isVersion {
		fmt.Printf("%s v%s for %s\n", appName, appVersion, authName)
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
		if os.Getenv("COIN_TYPE") != "" && coin.ValidateCoinTypeCode(os.Getenv("COIN_TYPE")){
			conf.CoinTypeCode = coin.CoinTypeCode(os.Getenv("COIN_TYPE"))
			log.Println("conf.CoinTypeCode:", conf.CoinTypeCode)
		}

		// create wallet
		regi := NewRegistry(conf, walletType, authName)
		walleter = regi.NewSigner()
	}
	defer func() {
		walleter.Done()
	}()

	//sub command
	args := flags.Args()
	cmds := sign.WalletSubCommands(walleter, appVersion)
	cl := command.CreateSubCommand(appName, appVersion, args, cmds)
	cl.HelpFunc = command.HelpFunc(cl.Name)

	flags.Usage = func() { fmt.Println(cl.HelpFunc(cl.Commands)) }

	code, err := cl.Run()
	if err != nil {
		log.Printf("fail to call Run() %s command: %v", appName, err)
	}
	os.Exit(code)
}
