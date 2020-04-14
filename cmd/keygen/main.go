package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/hiromaily/go-bitcoin/pkg/command"
	"github.com/hiromaily/go-bitcoin/pkg/command/keygen"
	"github.com/hiromaily/go-bitcoin/pkg/config"
	"github.com/hiromaily/go-bitcoin/pkg/wallets"
	"github.com/hiromaily/go-bitcoin/pkg/wallets/types"
)

// keygen wallet as cold wallet
//  generate key and seed for account
//  target account: client, receipt, payment

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

var (
	walletType = types.WalletTypeKeyGen
	appName    = walletType.String()
	appVersion = "2.0.0"
)

func main() {
	// command line
	var (
		confPath  string
		isHelp    bool
		isVersion bool
		walleter  wallets.Keygener
	)
	flags := flag.NewFlagSet("main", flag.ContinueOnError)
	flags.StringVar(&confPath, "conf", os.Getenv("KEYGEN_WALLET_CONF"), "config file path")
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
		// create wallet
		regi := NewRegistry(conf, walletType)
		walleter = regi.NewKeygener()
	}

	//sub command
	args := flags.Args()
	cmds := keygen.WalletSubCommands(walleter, appVersion)
	cl := command.CreateSubCommand(appName, appVersion, args, cmds)
	cl.HelpFunc = command.HelpFunc(cl.Name)

	flags.Usage = func() { fmt.Println(cl.HelpFunc(cl.Commands)) }

	code, err := cl.Run()
	if err != nil {
		log.Printf("fail to call Run() %s command: %v", appName, err)
	}
	os.Exit(code)
}

//TODO: after making sure command works, this code is deleted
//type Options struct {
//	//Configパス
//	ConfPath string `short:"c" long:"conf" default:"" description:"Path for configuration toml file"`
//
//	//署名モード
//	Sign bool `short:"s" long:"sign" description:"for signature"`
//	//Keyモード
//	Key bool `short:"k" long:"key" description:"for key related use (generate/import/export)"`
//	//Debugモード
//	Debug bool `short:"d" long:"debug" description:"for only development use"`
//
//	//実行されるサブ機能
//	Mode uint8 `short:"m" long:"mode" description:"Mode: detailed functionalities"`
//
//	//txファイルパス
//	ImportFile string `short:"i" long:"import" default:"" description:"import file path for hex"`
//	//key生成時に発行する数
//	KeyNumber uint32 `short:"n" long:"keynumber" description:"key number for generation"`
//	//アカウント
//	Account string `short:"a" long:"account" description:"account like client, receipt, payment"`
//}

//TODO: after making sure command works, this code is deleted
//if opts.Sign {
//	//sign関連機能
//	signFunctionalities(wallet)
//} else if opts.Key {
//	//key関連機能
//	keyFunctionalities(wallet)
//} else if opts.Debug {
//	//debug用 機能確認
//	debugForCheck(wallet)
//} else {
//	//logger.Warn("either sign:-s, key:-k, debug:-d should be set as main function")
//	procedure.ShowColdWallet1()
//}

//TODO: after making sure command works, this code is deleted
//func checkImportFile() {
//	if opts.ImportFile == "" {
//		logger.Fatal("file path is required as option -i")
//	}
//}
//
//func checkAccountWithoutAuth() {
//	if opts.Account == "" || !account.ValidateAccountType(opts.Account) ||
//		opts.Account == account.AccountTypeAuthorization.String() {
//		logger.Fatal("Account[client, receipt, payment, quoine, fee, stored] should be set with -a option")
//	}
//}
//
//func checkAccountWithoutAuthAndClient() {
//	if opts.Account == "" || !account.ValidateAccountType(opts.Account) ||
//		opts.Account == account.AccountTypeAuthorization.String() || opts.Account == account.AccountTypeClient.String() {
//		logger.Fatal("Account[receipt, payment, quoine, fee, stored] should be set with -a option")
//	}
//}
//
//// [coldwallet1]としての署名機能群 入金時の署名/出金時の署名[s]
//// TODO:出金時の署名は、coldwallet1/coldwallet2でそれぞれで署名が必要
//func signFunctionalities(wallet *service.Wallet) {
//	// 処理をModeで切り替える
//	switch opts.Mode {
//	case 1:
//		// importしたファイルからhex値を取得し、署名を行う(ReceiptかPaymentかはfileNameから判別))
//		logger.Info("Run: Importしたファイルからhex値を取得し、署名を行う(Receipt)")
//		checkImportFile()
//
//		//出金/入金の判別はファイル名から行う
//		hexTx, isSigned, generatedFileName, err := wallet.SignatureFromFile(opts.ImportFile)
//		if err != nil {
//			logger.Fatalf("%+v", err)
//		}
//		logger.Infof("[hex]: %s\n[署名完了]: %t\n[fileName]: %s", hexTx, isSigned, generatedFileName)
//	default:
//		logger.Warn("opts.Mode is out of range")
//		procedure.ShowColdWallet1()
//	}
//}
//
//// [coldwallet1]としてのKey関連機能群[k]
//func keyFunctionalities(wallet *service.Wallet) {
//	switch opts.Mode {
//	case 1:
//		//[coldwallet共通]
//		//HDウォレットによるSeedの作成
//		logger.Info("Run: HDウォレット Seedの生成")
//		bSeed, err := wallet.GenerateSeed()
//		if err != nil {
//			logger.Fatalf("%+v", err)
//		}
//		logger.Infof("seed: %s", key.SeedToString(bSeed))
//
//	case 10:
//		//[coldwallet1のみ]
//		//Keyを作成する(client, receipt, payment, quoine, fee, stored)
//		logger.Info("Run: Keyを作成する")
//		if opts.KeyNumber == 0 {
//			logger.Fatal("key number should be set with -n option")
//		}
//
//		checkAccountWithoutAuth()
//		logger.Infof("Run: Account[%s]", opts.Account)
//
//		//seed
//		bSeed, err := wallet.GenerateSeed()
//		if err != nil {
//			logger.Fatalf("%+v", err)
//		}
//
//		//generate
//		keys, err := wallet.GenerateAccountKey(account.AccountType(opts.Account), ctype.BTC, bSeed, opts.KeyNumber)
//		if err != nil {
//			logger.Fatalf("%+v", err)
//		}
//		grok.Value(keys)
//
//	case 20:
//		//[coldwallet1のみ]
//		//作成したAccountのPrivateKeyをColdWalletにimportする(client, receipt, payment, quoine, fee, stored)
//		logger.Info("Run: 作成したAccountのPrivateKeyをColdWallet1にimportする")
//
//		checkAccountWithoutAuth()
//		logger.Infof("Run: Account[%s]", opts.Account)
//
//		//import private key to coldwallet
//		err := wallet.ImportPrivateKey(account.AccountType(opts.Account))
//		if err != nil {
//			logger.Fatalf("%+v", err)
//		}
//		logger.Info("Done!")
//
//	case 30:
//		//[coldwallet1のみ]
//		//作成したAccountのPublicKeyをcsvファイルとしてexportする (watch only walletで利用するcsvファイル)
//		logger.Info("Run: 作成したAccountのPublicアドレスをcsvファイルとしてexportする")
//
//		checkAccountWithoutAuth()
//		logger.Infof("Run: Account[%s]", opts.Account)
//
//		//export public key as csv
//		fileName, err := wallet.ExportAccountKey(account.AccountType(opts.Account), keystatus.AddressStatusImportprivkey)
//		if err != nil {
//			logger.Fatalf("%+v", err)
//		}
//		logger.Infof("[fileName]: %s", fileName)
//
//	case 40:
//		//[coldwallet1のみ]
//		//TODO:coldwallet2からexportしたAccountのmultisigアドレスをcoldWallet1にimportする
//		logger.Info("Run: coldwallet2からexportしたAccountのmultisigアドレスをcoldWallet1にimportする")
//		checkImportFile()
//		checkAccountWithoutAuthAndClient()
//		logger.Infof("Run: Account[%s]", opts.Account)
//
//		//import multisig address from csv to database
//		err := wallet.ImportMultisigAddrForColdWallet1(opts.ImportFile, account.AccountType(opts.Account))
//		if err != nil {
//			logger.Fatalf("%+v", err)
//		}
//		logger.Info("Done!")
//
//	case 50:
//		//[coldwallet1のみ]
//		//Multisig Addressをcsvファイルとしてexportする (DBに出力済であるフラグを登録する必要がある)
//		logger.Info("Run: 作成したAccountのMultisigアドレスをcsvファイルとしてexportする")
//
//		checkAccountWithoutAuthAndClient()
//		logger.Infof("Run: Account[%s]", opts.Account)
//
//		//export account key
//		fileName, err := wallet.ExportAccountKey(account.AccountType(opts.Account), keystatus.AddressStatusMultiAddressImported)
//		if err != nil {
//			logger.Fatalf("%+v", err)
//		}
//		logger.Infof("[fileName]: %s", fileName)
//
//	default:
//		logger.Warn("opts.Mode is out of range")
//		procedure.ShowColdWallet1()
//	}
//
//}
//
//// Debug 検証用[d]
//func debugForCheck(wallet *service.Wallet) {
//	switch opts.Mode {
//	case 1:
//		//通常のKeyの生成(実運用では使わない)
//		logger.Info("Run: Keyの生成")
//		//単一Keyの生成
//		wif, pubAddress, err := key.GenerateKey(wallet.BTC.GetChainConf())
//		if err != nil {
//			logger.Fatalf("%+v", err)
//		}
//		logger.Infof("[WIF] %s - [Pub Address] %s\n", wif.String(), pubAddress)
//	case 10:
//		//Multisigの作成
//		logger.Info("Run: Multisigの作成")
//
//		//Multisigアドレス作成
//		resAddr, err := wallet.BTC.AddMultisigAddress(2, []string{"2N7ZwUXpo841GZDpxLGFqrhr1xwMzTba7ZP", "2NAm558FWpiaJQLz838vbzBPpqmKxyeyxsu"}, "multi01", ctype.AddressTypeP2shSegwit)
//		if err != nil {
//			logger.Fatalf("%+v", err)
//		}
//		logger.Infof("multisig address: %s, redeemScript: %s", resAddr.Address, resAddr.RedeemScript)
//	default:
//		logger.Warn("opts.Mode is out of range")
//		procedure.ShowColdWallet1()
//		//開発用
//		//development(wallet)
//	}
//}
