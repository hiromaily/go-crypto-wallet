package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/hiromaily/go-bitcoin/pkg/command"
	"github.com/hiromaily/go-bitcoin/pkg/command/signature"
	"github.com/hiromaily/go-bitcoin/pkg/config"
	"github.com/hiromaily/go-bitcoin/pkg/wallets"
	"github.com/hiromaily/go-bitcoin/pkg/wallets/types"
)

// signature wallet as cold wallet
//  generate one key and seed for only authorization account
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
	walletType = types.WalletTypeSignature
	appName    = walletType.String()
	appVersion = "2.0.0"
)

func main() {
	// command line
	var (
		confPath  string
		isHelp    bool
		isVersion bool
		walleter  wallets.Signer
	)
	flags := flag.NewFlagSet("main", flag.ContinueOnError)
	flags.StringVar(&confPath, "conf", os.Getenv("SIGNATURE_WALLET_CONF"), "config file path")
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
		walleter = regi.NewSigner()
	}

	//sub command
	args := flags.Args()
	cmds := signature.WalletSubCommands(walleter, appVersion)
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
//	procedure.ShowColdWallet2()
//}

//TODO: after making sure command works, this code is deleted
//func checkImportFile() {
//	if opts.ImportFile == "" {
//		logger.Fatal("file path is required as option -i")
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
//func checkAccountOnlyAuth() {
//	if opts.Account != account.AccountTypeAuthorization.String() {
//		logger.Fatal("Account[authorization] should be set with -a option")
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
//		procedure.ShowColdWallet2()
//	}
//}
//
//// [coldwallet2]としてのKey関連機能群[k]
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
//		//[coldwallet2のみ]
//		//AuthorizationのKeyを作成する
//		//FIXME:とりあえず、1行しか追加できないようにしておくか？
//		logger.Info("Run: AuthorizationのKeyを作成する")
//
//		//checkAccountOnlyAuth()
//
//		//seed
//		bSeed, err := wallet.GenerateSeed()
//		if err != nil {
//			logger.Fatalf("%+v", err)
//		}
//
//		//generate
//		keys, err := wallet.GenerateAccountKey(account.AccountTypeAuthorization, ctype.BTC, bSeed, 1)
//		if err != nil {
//			logger.Fatalf("%+v", err)
//		}
//		grok.Value(keys)
//
//	case 20:
//		//[coldwallet2のみ]
//		//作成したAuthorizationのPrivateKeyをColdWalletにimportする
//		logger.Info("Run: 作成したAuthorizationのPrivateKeyをColdWalletにimportする")
//
//		//import private key to coldwallet
//		err := wallet.ImportPrivateKey(account.AccountTypeAuthorization)
//		if err != nil {
//			logger.Fatalf("%+v", err)
//		}
//		logger.Info("Done!")
//
//	case 30:
//		//[coldwallet2のみ]
//		//coldwallet1からexportしたAccountのpublicアドレスをcoldWallet2にimportする
//		logger.Info("Run: coldwallet1からexportしたAccountのpublicアドレスcoldWallet2にimportする")
//		checkImportFile()
//		checkAccountWithoutAuthAndClient()
//		logger.Infof("Run: Account[%s]", opts.Account)
//
//		//import public key to database
//		err := wallet.ImportPubKeyForColdWallet2(opts.ImportFile, account.AccountType(opts.Account))
//		if err != nil {
//			logger.Fatalf("%+v", err)
//		}
//		logger.Info("Done!")
//
//	case 40:
//		//[coldwallet2のみ]
//		//`addmultisigaddress`を実行する。パラメータは、accountのアドレス、authorizationのアドレス
//		logger.Info("Run: `addmultisigaddress`を実行する。パラメータは、accountのアドレス、authorizationのアドレス")
//
//		checkAccountWithoutAuthAndClient()
//		logger.Infof("Run: Account[%s]", opts.Account)
//
//		//execute addmultisigaddress
//		err := wallet.AddMultisigAddress(account.AccountType(opts.Account), ctype.AddressTypeP2shSegwit)
//		if err != nil {
//			logger.Fatalf("%+v", err)
//		}
//		logger.Info("Done!")
//
//	case 50:
//		//[coldwallet2のみ]
//		//作成したAccountのMultisigアドレスをcsvファイルとしてexportする
//		logger.Info("Run: 作成したAccountのMultisigアドレスをcsvファイルとしてexportする")
//
//		checkAccountWithoutAuthAndClient()
//		logger.Infof("Run: Account[%s]", opts.Account)
//
//		//export multisig address
//		fileName, err := wallet.ExportAddedPubkeyHistory(account.AccountType(opts.Account))
//		if err != nil {
//			logger.Fatalf("%+v", err)
//		}
//		logger.Infof("[fileName]: %s", fileName)
//
//	default:
//		logger.Warn("opts.Mode is out of range")
//		procedure.ShowColdWallet2()
//	}
//
//}
//
//// Debug 検証用[d]
//func debugForCheck(wallet *service.Wallet) {
//	//development(wallet)
//	procedure.ShowColdWallet2()
//}
