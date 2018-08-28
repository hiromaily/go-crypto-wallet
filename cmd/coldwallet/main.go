package main

import (
	"github.com/bookerzzz/grok"
	"github.com/btcsuite/btcd/chaincfg"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/key"
	"github.com/hiromaily/go-bitcoin/pkg/logger"
	"github.com/hiromaily/go-bitcoin/pkg/service"
	"github.com/jessevdk/go-flags"
)

// HDウォレットとしてseed作成、keyを指定した数だけ生成し、出力する
// これは、ネットワーク環境下のwallet側から、

//TODO:encryptwalletコマンドによって、walletを暗号化した場合、秘密鍵を使用するタイミング(未署名トランザクションに署名する)
// でパスフレーズの入力が必要になり

// Options コマンドラインオプション
type Options struct {
	//Configパス
	ConfPath string `short:"c" long:"conf" default:"./data/toml/cold1_config.toml" description:"Path for configuration toml file"`
	//実行される機能
	Mode uint8 `short:"m" long:"mode" description:"Mode i.e.Functionality"`
	//HDウォレット用Key生成のためのseed情報
	//ParamSeed string `short:"d" long:"seed" default:"" description:"backup seed"`
	//txファイルパス
	ImportFile string `short:"i" long:"import" default:"" description:"import file path for hex"`
	//Debugモード
	Debug bool `short:"d" long:"debug" description:"for only development use"`
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

	if opts.Debug {
		//debug用 機能確認
		debugForCheck(wallet)
	} else {
		//switch mode
		switchFunction(wallet)
	}
}

// 実運用上利用するもののみ、こちらに定義する
func switchFunction(wallet *service.Wallet) {
	// 処理をFunctionalityで切り替える
	//TODO:ここから呼び出すべきはService系のみに統一したい
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
		logger.Info("該当Mode無し")
	}
}

// 検証用
func debugForCheck(wallet *service.Wallet) {
	switch opts.Mode {
	case 1:
		//通常のKeyの生成
		logger.Info("Run: Keyの生成")
		//単一Keyの生成
		wif, pubAddress, err := key.GenerateKey(wallet.BTC.GetChainConf())
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		logger.Infof("[WIF] %s - [Pub Address] %s\n", wif.String(), pubAddress)
	case 2:
		//HDウォレットによるSeedの作成
		logger.Info("Run: HDウォレット Seedの生成")
		bSeed, err := wallet.GenerateSeed()
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		logger.Infof("seed: %s", key.SeedToString(bSeed))
	case 3:
		//ClientのKeyを作成する
		logger.Info("Run: ClientのKeyを作成する")
		bSeed, err := wallet.GenerateSeed()
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		keys, err := wallet.GenerateAccountKey(enum.AccountTypeClient, bSeed, 10)
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		grok.Value(keys)
	case 4:
		//ReceiptのKeyを作成する
		logger.Info("Run: ReceiptのKeyを作成する")
		bSeed, err := wallet.GenerateSeed()
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		keys, err := wallet.GenerateAccountKey(enum.AccountTypeReceipt, bSeed, 10)
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		grok.Value(keys)
	case 5:
		//PaymentのKeyを作成する
		logger.Info("Run: PaymentのKeyを作成する")
		bSeed, err := wallet.GenerateSeed()
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		keys, err := wallet.GenerateAccountKey(enum.AccountTypePayment, bSeed, 5)
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		grok.Value(keys)
	case 6:
		//TODO:これはcoldwallet2(承認用)の機能
		//AuthorizationのKeyを作成する
		logger.Info("Run: AuthorizationのKeyを作成する")
		bSeed, err := wallet.GenerateSeed()
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		keys, err := wallet.GenerateAccountKey(enum.AccountTypeAuthorization, bSeed, 2)
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		grok.Value(keys)
	case 10:
		//作成したPrivateKeyをWalletにimportする
		err := wallet.ImportPrivateKey(enum.AccountTypeClient)
		if err != nil {
			logger.Fatalf("%+v", err)
		}
	case 11:
		//作成したPublicKeyをcsvファイルとしてexportする
		err := wallet.ExportPublicKey(enum.AccountTypeClient)
		if err != nil {
			logger.Fatalf("%+v", err)
		}
	case 20:
		//TODO:Multisigの作成
		logger.Info("Run: Multisigの作成")

		//事前準備
		//getnewaddress taro 2N7ZwUXpo841GZDpxLGFqrhr1xwMzTba7ZP
		//getnewaddress boss1 2NAm558FWpiaJQLz838vbzBPpqmKxyeyxsu
		//TODO:ここで、AddMultisigAddressを使うのにパラメータとしてaccout名も渡さないといけない。。これをどうすべきか。。。
		//TODO: => おそらくBlankでもいい

		//TODO: Multisigアドレス作成 (まだ検証中)
		resAddr, err := wallet.BTC.CreateMultiSig(2, []string{"2N7ZwUXpo841GZDpxLGFqrhr1xwMzTba7ZP", "2NAm558FWpiaJQLz838vbzBPpqmKxyeyxsu"}, "multi01")
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		logger.Infof("multisig address: %s, redeemScript: %s", resAddr.Address, resAddr.RedeemScript)

	case 30:
		//[Debug用]HEXから署名を行う
		logger.Info("Run: HEXから署名を行う")
		hex := "02000000021ed288be4c4d7923a0d044bb500a15b2eb0f2b3c5503293f251f7c94939a3f9f0000000000ffffffff557624120cdf3f4d092f35e5cd6b75418b76c3e3fd4c398357374e93cfe5c4200000000000ffffffff05c03b47030000000017a91419e70491572c55fb08ce90b0c6bf5cfe45a5420e87809698000000000017a9146b8902fc7a6a0bccea9dbd80a4c092c314227f618734e133070000000017a9148191d41a7415a6a1f6ee14337e039f50b949e80e87005a62020000000017a9149c877d6f21d5800ca60a7660ee56745f239b222b87002d31010000000017a914f575a0d1ddcfb98a11628826f1632453d718ff618700000000"
		hexTx, isSigned, generatedFileName, err := wallet.SignatureByHex(enum.ActionTypeReceipt, hex, 10)
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		logger.Infof("hex: %s\n, 署名完了: %t\n, fileName: %s", hexTx, isSigned, generatedFileName)
		//TODO:isSigned: 送信までした署名はfalseになる??
	default:
		logger.Info("該当Mode無し")
	}
}
