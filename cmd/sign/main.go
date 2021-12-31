package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/command"
	"github.com/hiromaily/go-crypto-wallet/pkg/command/sign"
	"github.com/hiromaily/go-crypto-wallet/pkg/config"
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
	authName   = "" // this value is supposed to be embedded when building
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
		walleter        wallets.Signer
	)
	flags := flag.NewFlagSet("main", flag.ContinueOnError)
	flags.StringVar(&confPath, "conf", "", "config file path")
	flags.StringVar(&coinTypeCode, "coin", "btc", "coin type code `btc`, `bch`")
	flags.StringVar(&btcWallet, "wallet", "", "specify wallet.dat in bitcoin core")
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

	// validate coinTypeCode
	if !coin.IsCoinTypeCode(coinTypeCode) {
		log.Fatal("coin args is invalid. `btc`, `bch` is allowed")
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

	// help
	if !isHelp && len(os.Args) > 1 {
		// config
		conf, err := config.NewWallet(confPath, walletType, coin.CoinTypeCode(coinTypeCode))
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

		// override conf.Bitcoin.Host
		if btcWallet != "" {
			conf.Bitcoin.Host = fmt.Sprintf("%s/wallet/%s", conf.Bitcoin.Host, btcWallet)
			log.Println("conf.Bitcoin.Host:", conf.Bitcoin.Host)
		}

		// create wallet
		regi := NewRegistry(conf, accountConf, walletType, authName)
		walleter = regi.NewSigner()
	}
	defer func() {
		walleter.Done()
	}()

	// sub command
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
