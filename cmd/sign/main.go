package main

import (
	"flag"
	"fmt"
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
//  generate one key and seed for only authorization account
//  target account: client, deposit, payment

// procedure
//  1. create seed
//  2. create key
//  3. run `importprivkey`
//  4. export pubkey from DB
//  5. sing on unsigned transaction
//   sign for unsigned transaction (multisig addresses are required to sign by multiple devices)

//TODO: bitcoin functionalities
// - encrypt wallet itself by `encryptwallet` command
// - passphrase would be required when using secret key to sign unsigned transaction
// - multisig with bigger number e.g. 3:5
var (
	walletType = wallet.WalletTypeSignature
	appName    = walletType.String()
	appVersion = "2.2.0"
	authName   = "" // this account is supposed to be embedded when building
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
		}
		log.Println("conf.Bitcoin.Host:", conf.Bitcoin.Host)

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
