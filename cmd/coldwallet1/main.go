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
	ConfPath string `short:"c" long:"conf" default:"./data/toml/local_cold1.toml" description:"Path for configuration toml file"`

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
		procedure.ShowColdWallet1()
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
		procedure.ShowColdWallet1()
	}
}

// [coldwallet1]としてのKey関連機能群[k]
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

	case 10:
		//[coldwallet1のみ]
		//ClientのKeyを作成する
		logger.Info("Run: ClientのKeyを作成する")
		bSeed, err := wallet.GenerateSeed()
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		keys, err := wallet.GenerateAccountKey(enum.AccountTypeClient, enum.BTC, bSeed, 10)
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		grok.Value(keys)
	case 11:
		//[coldwallet1のみ]
		//ReceiptのKeyを作成する
		logger.Info("Run: ReceiptのKeyを作成する")
		bSeed, err := wallet.GenerateSeed()
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		keys, err := wallet.GenerateAccountKey(enum.AccountTypeReceipt, enum.BTC, bSeed, 5)
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		grok.Value(keys)
	case 12:
		//[coldwallet1のみ]
		//PaymentのKeyを作成する
		logger.Info("Run: PaymentのKeyを作成する")
		bSeed, err := wallet.GenerateSeed()
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		keys, err := wallet.GenerateAccountKey(enum.AccountTypePayment, enum.BTC, bSeed, 5)
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		grok.Value(keys)

	case 20:
		//[coldwallet1のみ]
		//作成したClientのPrivateKeyをColdWalletにimportする
		logger.Info("Run: 作成したClientのPrivateKeyをColdWalletにimportする")
		err := wallet.ImportPrivateKey(enum.AccountTypeClient)
		if err != nil {
			logger.Fatalf("%+v", err)
		}
	case 21:
		//[coldwallet1のみ]
		//作成したReceiptのPrivateKeyをColdWalletにimportする
		logger.Info("Run: 作成したReceiptのPrivateKeyをColdWalletにimportする")
		err := wallet.ImportPrivateKey(enum.AccountTypeReceipt)
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		// getaddressesbyaccount "receipt" で確認
	case 22:
		//[coldwallet1のみ]
		//作成したPaymentのPrivateKeyをColdWalletにimportする
		logger.Info("Run: 作成したPaymentのPrivateKeyをColdWalletにimportする")
		err := wallet.ImportPrivateKey(enum.AccountTypePayment)
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		//getaddressesbyaccount "payment" で内容を確認

	case 30:
		//[coldwallet1のみ]
		//作成したClientのPublicKeyをcsvファイルとしてexportする (watch only wallet用)
		logger.Info("Run: 作成したClientのPublicアドレスをcsvファイルとしてexportする")
		fileName, err := wallet.ExportAccountKey(enum.AccountTypeClient, enum.KeyStatusImportprivkey)
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		logger.Infof("fileName: %s", fileName)
	case 31:
		//[coldwallet1のみ]
		//作成したReceiptのPublicKeyをcsvファイルとしてexportする (coldwallet2用)
		logger.Info("Run: 作成したReceiptのPublicアドレスをcsvファイルとしてexportする")
		fileName, err := wallet.ExportAccountKey(enum.AccountTypeReceipt, enum.KeyStatusImportprivkey)
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		logger.Infof("fileName: %s", fileName)
	case 32:
		//[coldwallet1のみ]
		//作成したPaymentのPublicKeyをcsvファイルとしてexportする　(coldwallet2用)
		logger.Info("Run: 作成したPaymentのPublicアドレスをcsvファイルとしてexportする")
		fileName, err := wallet.ExportAccountKey(enum.AccountTypePayment, enum.KeyStatusImportprivkey)
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		logger.Infof("fileName: %s", fileName)

	case 40:
		//[coldwallet1のみ]
		//TODO:coldwallet2からexportしたReceiptのmultisigアドレスをcoldWallet1にimportする
		logger.Info("Run: coldwallet2からexportしたReceiptのmultisigアドレスをcoldWallet1にimportする")
		if opts.ImportFile == "" {
			logger.Fatal("file path is required as argument file when running")
		}
		err := wallet.ImportMultisigAddrForColdWallet1(opts.ImportFile, enum.AccountTypeReceipt)
		if err != nil {
			logger.Fatalf("%+v", err)
		}

	case 41:
		//[coldwallet1のみ]
		//TODO:coldwallet2からexportしたPaymentのmultisigアドレスをcoldWallet1にimportする
		logger.Info("Run: coldwallet2からexportしたPaymentのmultisigアドレスをcoldWallet1にimportする")
		if opts.ImportFile == "" {
			logger.Fatal("file path is required as argument file when running")
		}
		err := wallet.ImportMultisigAddrForColdWallet1(opts.ImportFile, enum.AccountTypePayment)
		if err != nil {
			logger.Fatalf("%+v", err)
		}
	case 50:
		//[coldwallet1のみ]
		//multisigimport後、ReceiptのMultisigをcsvファイルとしてexportする (DBに出力済を登録する必要がある)
		//=>TODO:しかし、coldwallet2側から出力されたファイルがそのまま使えるような？？しかし、情報の管理のために、一度coldwallet1にimportが必要
		logger.Info("Run: 作成したReceiptのMultisigアドレスをcsvファイルとしてexportする")
		fileName, err := wallet.ExportAccountKey(enum.AccountTypeReceipt, enum.KeyStatusMultiAddressImported)
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		logger.Infof("fileName: %s", fileName)
	case 51:
		//[coldwallet1のみ]
		//multisigimport後、PaymentのMultisigをcsvファイルとしてexportする (DBに出力済を登録する必要がある)
		//=>TODO:しかし、coldwallet2側から出力されたファイルがそのまま使えるような？？しかし、情報の管理のために、一度coldwallet1にimportが必要
		logger.Info("Run: 作成したPaymentのMultisigアドレスをcsvファイルとしてexportする")
		fileName, err := wallet.ExportAccountKey(enum.AccountTypePayment, enum.KeyStatusMultiAddressImported)
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		logger.Infof("fileName: %s", fileName)
	default:
		logger.Warn("opts.Mode is out of range")
		procedure.ShowColdWallet1()
	}

}

// Debug 検証用[d]
func debugForCheck(wallet *service.Wallet) {
	switch opts.Mode {
	case 1:
		//通常のKeyの生成(実運用では使わない)
		logger.Info("Run: Keyの生成")
		//単一Keyの生成
		wif, pubAddress, err := key.GenerateKey(wallet.BTC.GetChainConf())
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		logger.Infof("[WIF] %s - [Pub Address] %s\n", wif.String(), pubAddress)
	case 10:
		//TODO:Multisigの作成
		logger.Info("Run: Multisigの作成")

		//事前準備
		//getnewaddress taro 2N7ZwUXpo841GZDpxLGFqrhr1xwMzTba7ZP
		//getnewaddress boss1 2NAm558FWpiaJQLz838vbzBPpqmKxyeyxsu
		//TODO:ここで、AddMultisigAddressを使うのにパラメータとしてaccout名も渡さないといけない。。これをどうすべきか。。。
		//TODO: => おそらくBlankでもいい

		//TODO: Multisigアドレス作成
		resAddr, err := wallet.BTC.CreateMultiSig(2, []string{"2N7ZwUXpo841GZDpxLGFqrhr1xwMzTba7ZP", "2NAm558FWpiaJQLz838vbzBPpqmKxyeyxsu"}, "multi01", enum.AddressTypeP2shSegwit)
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		logger.Infof("multisig address: %s, redeemScript: %s", resAddr.Address, resAddr.RedeemScript)

	case 20:
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
		logger.Warn("opts.Mode is out of range")
		//procedure.Show()
		//開発用
		development(wallet)

	}
}

//coldwallet1 -d
func development(wallet *service.Wallet) {
	//------------------------------------------------------------------
	//coldwallet1の機能
	//------------------------------------------------------------------
	//seed
	bSeed, err := wallet.GenerateSeed()
	if err != nil {
		logger.Fatalf("%+v", err)
	}

	//1.generate key
	_, err = wallet.GenerateAccountKey(enum.AccountTypeClient, enum.BTC, bSeed, 10)
	if err != nil {
		logger.Fatalf("%+v", err)
	}
	_, err = wallet.GenerateAccountKey(enum.AccountTypeReceipt, enum.BTC, bSeed, 5)
	if err != nil {
		logger.Fatalf("%+v", err)
	}
	_, err = wallet.GenerateAccountKey(enum.AccountTypePayment, enum.BTC, bSeed, 5)
	if err != nil {
		logger.Fatalf("%+v", err)
	}

	//2.importprivkey
	err = wallet.ImportPrivateKey(enum.AccountTypeClient)
	if err != nil {
		logger.Fatalf("%+v", err)
	}
	err = wallet.ImportPrivateKey(enum.AccountTypeReceipt)
	if err != nil {
		logger.Fatalf("%+v", err)
	}
	err = wallet.ImportPrivateKey(enum.AccountTypePayment)
	if err != nil {
		logger.Fatalf("%+v", err)
	}

	//3.export full public key
	fileName, err := wallet.ExportAccountKey(enum.AccountTypeClient, enum.KeyStatusImportprivkey)
	if err != nil {
		logger.Fatalf("%+v", err)
	}
	fileName1, err := wallet.ExportAccountKey(enum.AccountTypeReceipt, enum.KeyStatusImportprivkey)
	if err != nil {
		logger.Fatalf("%+v", err)
	}
	fileName2, err := wallet.ExportAccountKey(enum.AccountTypePayment, enum.KeyStatusImportprivkey)
	if err != nil {
		logger.Fatalf("%+v", err)
	}

	//------------------------------------------------------------------
	//coldwallet2の機能だが、同一DBで検証
	//------------------------------------------------------------------
	//1.generate key
	_, err = wallet.GenerateAccountKey(enum.AccountTypeAuthorization, enum.BTC, bSeed, 1)
	if err != nil {
		logger.Fatalf("%+v", err)
	}

	//2.importprivkey
	err = wallet.ImportPrivateKey(enum.AccountTypeAuthorization)
	if err != nil {
		logger.Fatalf("%+v", err)
	}

	//3.import publickey to history table
	err = wallet.ImportPublicKeyForColdWallet2(fileName1, enum.AccountTypeReceipt)
	if err != nil {
		logger.Fatalf("%+v", err)
	}
	err = wallet.ImportPublicKeyForColdWallet2(fileName2, enum.AccountTypePayment)
	if err != nil {
		logger.Fatalf("%+v", err)
	}

	//4.addmultisig
	err = wallet.AddMultisigAddressByAuthorization(enum.AccountTypeReceipt, enum.AddressTypeP2shSegwit)
	if err != nil {
		logger.Fatalf("%+v", err)
	}
	err = wallet.AddMultisigAddressByAuthorization(enum.AccountTypePayment, enum.AddressTypeP2shSegwit)
	if err != nil {
		logger.Fatalf("%+v", err)
	}

	//5.multisigアドレスとキーとなるアドレスをwatch only walletのためにexport
	fileName3, err := wallet.ExportAddedPubkeyHistory(enum.AccountTypeReceipt)
	if err != nil {
		logger.Fatalf("%+v", err)
	}
	fileName4, err := wallet.ExportAddedPubkeyHistory(enum.AccountTypePayment)
	if err != nil {
		logger.Fatalf("%+v", err)
	}

	//------------------------------------------------------------------
	//coldwallet1の機能
	//------------------------------------------------------------------
	//6.coldwallet2からexportしたReceiptのmultisigアドレスをcoldWallet1にimportする
	err = wallet.ImportMultisigAddrForColdWallet1(fileName3, enum.AccountTypeReceipt)
	if err != nil {
		logger.Fatalf("%+v", err)
	}
	err = wallet.ImportMultisigAddrForColdWallet1(fileName4, enum.AccountTypePayment)
	if err != nil {
		logger.Fatalf("%+v", err)
	}

	//multisigimport後、ReceiptのMultisigをcsvファイルとしてexportする
	fileName6, err := wallet.ExportAccountKey(enum.AccountTypeReceipt, enum.KeyStatusMultiAddressImported)
	if err != nil {
		logger.Fatalf("%+v", err)
	}
	logger.Infof("fileName: %s", fileName6)

	fileName7, err := wallet.ExportAccountKey(enum.AccountTypePayment, enum.KeyStatusMultiAddressImported)
	if err != nil {
		logger.Fatalf("%+v", err)
	}
	logger.Infof("fileName: %s", fileName7)

	//------------------------------------------------------------------
	//walletの機能
	//------------------------------------------------------------------
	//coldwalletで生成したアドレスをwalletにimportする
	//TODO:一台のPCで実行すると、service/import_public_key.go:50でエラーが出る
	//error: -4: The wallet already contains the private key for this address or script
	//err = wallet.ImportPublicKeyForWatchWallet(fileName, enum.AccountTypeClient)
	//if err != nil {
	//	logger.Fatalf("%+v", err)
	//}
	//err = wallet.ImportPublicKeyForWatchWallet(fileName6, enum.AccountTypeReceipt)
	//if err != nil {
	//	logger.Fatalf("%+v", err)
	//}
	//err = wallet.ImportPublicKeyForWatchWallet(fileName7, enum.AccountTypePayment)
	//if err != nil {
	//	logger.Fatalf("%+v", err)
	//}

	logger.Infof("next you should run\nwallet -k -m 1 -i %s,\nwallet -k -m 2 -i %s,\nwallet -k -m 3 -i %s", fileName, fileName6, fileName7)
	//wallet -k -m 1 -i ./data/pubkey/xxx.csv
	//wallet -k -m 2 -i ./data/pubkey/xxx.csv
	//wallet -k -m 3 -i ./data/pubkey/xxx.csv
}
