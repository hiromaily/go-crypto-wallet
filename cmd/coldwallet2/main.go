package main

import (
	"github.com/bookerzzz/grok"
	"github.com/btcsuite/btcd/chaincfg"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/key"
	"github.com/hiromaily/go-bitcoin/pkg/logger"
	"github.com/hiromaily/go-bitcoin/pkg/procedure"
	"github.com/hiromaily/go-bitcoin/pkg/service"
	"github.com/jessevdk/go-flags"
)

// HDウォレットとしてseed作成、keyを指定した数だけ生成し、出力する
// 対象アカウント: client, receipt, payment
// 1. create seed
// 2. create key
// 3. run `importprivkey`
// 4. export pubkey from DB

// 5. 未署名トランザクションへの署名

//TODO:encryptwalletコマンドによって、walletを暗号化した場合、秘密鍵を使用するタイミング(未署名トランザクションに署名する)
// でパスフレーズの入力が必要になる

// Options コマンドラインオプション
type Options struct {
	//Configパス
	ConfPath string `short:"c" long:"conf" default:"./data/toml/local_cold2.toml" description:"Path for configuration toml file"`

	//署名モード
	Sign bool `short:"s" long:"sign" description:"for signature"`
	//Keyモード
	Key bool `short:"k" long:"key" description:"for key related use (generate/import/export)"`
	//Debugモード
	Debug bool `short:"d" long:"debug" description:"for only development use"`

	//実行される詳細機能
	Mode uint8 `short:"m" long:"mode" description:"Mode: detailed functionalities"`

	//txファイルパス
	ImportFile string `short:"i" long:"import" default:"" description:"import file path for hex"`
}

var (
	opts      Options
	chainConf *chaincfg.Params
)

func init() {
	if _, err := flags.Parse(&opts); err != nil {
		panic(err)
	}
}

func main() {
	// Config
	wallet, err := service.InitialSettings(opts.ConfPath)
	if err != nil {
		logger.Fatal(err)
	}
	defer wallet.Done()

	if opts.Sign {
		//sign関連機能
		signFunctionalities(wallet)
	} else if opts.Key {
		//key関連機能
		keyFunctionalities(wallet)
	} else if opts.Debug {
		//debug用 機能確認
		debugForCheck(wallet)
	} else {
		//logger.Warn("either sign:-s, key:-k, debug:-d should be set as main function")
		procedure.ShowColdWallet2()
	}

}

// [coldwallet1]としての署名機能群 入金時の署名/出金時の署名[s]
// TODO:出金時の署名は、coldwallet1/coldwallet2でそれぞれで署名が必要
func signFunctionalities(wallet *service.Wallet) {
	// 処理をModeで切り替える
	switch opts.Mode {
	case 1:
		// importしたファイルからhex値を取得し、署名を行う(ReceiptかPaymentかはfileNameから判別))
		logger.Info("Run: Importしたファイルからhex値を取得し、署名を行う(Receipt)")
		if opts.ImportFile == "" {
			logger.Fatal("file path is required as argument file when running")
		}

		//出金/入金の判別はファイル名から行う
		hexTx, isSigned, generatedFileName, err := wallet.SignatureFromFile(opts.ImportFile)
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		logger.Infof("[hex]: %s\n[署名完了]: %t\n[fileName]: %s", hexTx, isSigned, generatedFileName)
	default:
		logger.Warn("opts.Mode is out of range")
		procedure.ShowColdWallet2()
	}
}

// [coldwallet2]としてのKey関連機能群[k]
func keyFunctionalities(wallet *service.Wallet) {
	switch opts.Mode {
	case 1:
		//[coldwallet共通]
		//HDウォレットによるSeedの作成
		logger.Info("Run: HDウォレット Seedの生成")
		bSeed, err := wallet.GenerateSeed()
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		logger.Infof("seed: %s", key.SeedToString(bSeed))

	case 13:
		//[coldwallet2のみ]
		//AuthorizationのKeyを作成する
		//FIXME:とりあえず、1行しか追加できないようにしておくか？
		logger.Info("Run: AuthorizationのKeyを作成する")
		bSeed, err := wallet.GenerateSeed()
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		keys, err := wallet.GenerateAccountKey(enum.AccountTypeAuthorization, enum.BTC, bSeed, 1)
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		grok.Value(keys)

	case 23:
		//[coldwallet2のみ]
		//作成したAuthorizationのPrivateKeyをColdWalletにimportする
		logger.Info("Run: 作成したAuthorizationのPrivateKeyをColdWalletにimportする")
		err := wallet.ImportPrivateKey(enum.AccountTypeAuthorization)
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		//getaddressesbyaccount "authorization" で内容を確認

	case 33:
		//[coldwallet2のみ]
		//coldwallet1からexportしたReceiptのpublicアドレスをcoldWallet2にimportする
		logger.Info("Run: coldwallet1からexportしたReceiptのpublicアドレスcoldWallet2にimportする")
		if opts.ImportFile == "" {
			logger.Fatal("file path is required as argument file when running")
		}
		err := wallet.ImportPublicKeyForColdWallet2(opts.ImportFile, enum.AccountTypeReceipt)
		if err != nil {
			logger.Fatalf("%+v", err)
		}

	case 34:
		//[coldwallet2のみ]
		//coldwallet1からexportしたPaymentのpublicアドレスをcoldWallet2にimportする
		logger.Info("Run: coldwallet1からexportしたPaymentのpublicアドレスcoldWallet2にimportする")
		if opts.ImportFile == "" {
			logger.Fatal("file path is required as argument file when running")
		}
		err := wallet.ImportPublicKeyForColdWallet2(opts.ImportFile, enum.AccountTypePayment)
		if err != nil {
			logger.Fatalf("%+v", err)
		}

	case 50:
		//[coldwallet2のみ]
		//`addmultisigaddress`を実行する。パラメータは、receiptのアドレス、authorizationのアドレス
		logger.Info("Run: `addmultisigaddress`を実行する。パラメータは、receiptのアドレス、authorizationのアドレス")
		err := wallet.AddMultisigAddressByAuthorization(enum.AccountTypeReceipt, enum.AddressTypeP2shSegwit)
		if err != nil {
			logger.Fatalf("%+v", err)
		}

	case 51:
		//[coldwallet2のみ]
		//`addmultisigaddress`を実行する。パラメータは、paymentのアドレス、authorizationのアドレス
		logger.Info("Run: `addmultisigaddress`を実行する。パラメータは、receiptのアドレス、authorizationのアドレス")
		err := wallet.AddMultisigAddressByAuthorization(enum.AccountTypePayment, enum.AddressTypeP2shSegwit)
		if err != nil {
			logger.Fatalf("%+v", err)
		}

	case 60:
		//[coldwallet2のみ]
		//TODO:作成したReceiptのMultisigアドレスをcsvファイルとしてexportする
		logger.Info("Run: 作成したReceiptのMultisigアドレスをcsvファイルとしてexportする")
		fileName, err := wallet.ExportAddedPubkeyHistory(enum.AccountTypeReceipt)
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		logger.Infof("fileName: %s", fileName)

	case 61:
		//[coldwallet2のみ]
		//TODO:作成したPaymentのMultisigアドレスをcsvファイルとしてexportする
		logger.Info("Run: 作成したPaymentのMultisigアドレスをcsvファイルとしてexportする")

		fileName, err := wallet.ExportAddedPubkeyHistory(enum.AccountTypePayment)
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		logger.Infof("fileName: %s", fileName)

	default:
		logger.Warn("opts.Mode is out of range")
		procedure.ShowColdWallet2()
	}

}

// Debug 検証用[d]
func debugForCheck(wallet *service.Wallet) {
	//development(wallet)
	procedure.ShowColdWallet2()
}
