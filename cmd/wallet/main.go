package main

import (
	"github.com/hiromaily/go-bitcoin/pkg/testdata"
	"log"

	"github.com/bookerzzz/grok"
	"github.com/btcsuite/btcd/chaincfg"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/logger"
	"github.com/hiromaily/go-bitcoin/pkg/procedure"
	"github.com/hiromaily/go-bitcoin/pkg/service"
	"github.com/jessevdk/go-flags"
)

//こちらはHotwallet、ただし、Watch Only Walletとしての機能を実装していく。
//ネットワークへの接続はGCP上のBitcoin Core
//Watch Only Walletとしてのセットアップが必要
// - Cold Wallet側から生成したPublic Key をMultisigアドレス変換後、`importaddress xxxxx`でimportする

//TODO:coldwallet側(非ネットワーク環境)側の機能と明確に分ける
//TODO:オフラインで可能機能と、不可能な機能の切り分けが必要
//TODO:ウォレットの定期バックアップ機能 + import機能
//TODO:coldウォレットへのデータ移行機能が必要なはず
//TODO:multisigの実装
//TODO:生成したkeyの暗号化処理のpkgが必要になるはず
//TODO:入金時にMultisigでの送金は不要な気がする

// Options コマンドラインオプション
type Options struct {
	//Configパス
	ConfPath string `short:"c" long:"conf" default:"./data/toml/watch_config.toml" description:"Path for configuration toml file"`

	//Keyモード
	Key bool `short:"k" long:"key" description:"for adding key"`
	//入金モード
	Receipt bool `short:"r" long:"receipt" description:"for receipt"`
	//出金モード
	Payment bool `short:"p" long:"payment" description:"for payment"`
	//送金モード
	Send bool `short:"s" long:"sending" description:"for sending transaction"`
	//ステータスチェックモード
	Monitor bool `short:"n" long:"monitor" description:"for monitoring transaction"`
	//bitcoin Commandモード
	Cmd bool `short:"b" long:"bitcoin-command" description:"for bitcoin command"`
	//Debugモード
	Debug bool `short:"d" long:"debug" description:"for only development use"`

	//実行される機能
	Mode uint8 `short:"m" long:"mode" description:"Mode i.e.Functionality"`

	//txファイルパス
	ImportFile string `short:"i" long:"import" default:"" description:"import file path for hex"`
	//調整fee
	Fee float64 `short:"f" long:"fee" default:"" description:"adjustment fee"`
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
	//initialSettings()
	wallet, err := service.InitialSettings(opts.ConfPath)
	if err != nil {
		logger.Fatal(err)
	}
	defer wallet.Done()

	if opts.Key {
		//キー関連機能
		keyFunctionalities(wallet)
	} else if opts.Receipt {
		//入金関連機能
		receiptFunctionalities(wallet)
	} else if opts.Payment {
		//出金関連機能
		paymentFunctionalities(wallet)
	} else if opts.Send {
		//署名送信関連機能
		sendingFunctionalities(wallet)
	} else if opts.Monitor {
		//transaction監視関連機能
		monitoringFunctionalities(wallet)
	} else if opts.Cmd {
		//BTCコマンド実行
		btcCommand(wallet)
	} else if opts.Debug {
		//debug用 機能確認
		debugForCheck(wallet)
	} else {
		logger.Warn("either sign:-s, key:-k, debug:-d should be set as main function")
		procedure.ShowWallet()
	}
}

//キー関連機能
func keyFunctionalities(wallet *service.Wallet) {
	switch opts.Mode {
	case 1:
		//TODO:imporot後、getaddressesbyaccount "" で内容を確認??
		logger.Info("Run: coldwalletで生成した[client]アドレスをwalletにimportする")
		if opts.ImportFile == "" {
			logger.Fatal("file path is required as argument file when running")
		}
		err := wallet.ImportPublicKeyForWatchWallet(opts.ImportFile, enum.AccountTypeClient)
		if err != nil {
			logger.Fatalf("%+v", err)
		}
	case 2:
		//TODO:imporot後、getaddressesbyaccount "" で内容を確認??
		logger.Info("Run: coldwalletで生成した[receipt]アドレスをwalletにimportする")
		if opts.ImportFile == "" {
			logger.Fatal("file path is required as argument file when running")
		}
		err := wallet.ImportPublicKeyForWatchWallet(opts.ImportFile, enum.AccountTypeReceipt)
		if err != nil {
			logger.Fatalf("%+v", err)
		}
	case 3:
		//TODO:imporot後、getaddressesbyaccount "" で内容を確認??
		logger.Info("Run: coldwalletで生成した[payment]アドレスをwalletにimportする")
		if opts.ImportFile == "" {
			logger.Fatal("file path is required as argument file when running")
		}
		err := wallet.ImportPublicKeyForWatchWallet(opts.ImportFile, enum.AccountTypePayment)
		if err != nil {
			logger.Fatalf("%+v", err)
		}

	default:
		logger.Warn("opts.Mode is out of range")
		procedure.ShowWallet()
	}

	//clientのaddress, receipt,paymentのmultisigアドレスをimportする
	//DBにinsert後、bitcoin commandで登録する
	//importmulti or importaddress
}

//入金関連機能
func receiptFunctionalities(wallet *service.Wallet) {
	switch opts.Mode {
	case 1:
		logger.Info("Run: 入金処理検知 + 未署名トランザクション作成")
		//実際には署名処理は手動なので、ユーザーの任意のタイミングで走らせたほうがいい。
		//入金検知 + 未署名トランザクション作成
		hex, fileName, err := wallet.DetectReceivedCoin(opts.Fee)
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		if hex == "" {
			logger.Info("No utxo")
			return
		}
		logger.Infof("[hex]: %s\n[fileName]: %s", hex, fileName)
	case 2:
		logger.Info("Run: 入金処理検知 (確認のみ)")
		//TODO:WIP

	case 10:
		logger.Info("Run: [Debug用]入金から送金までの一連の流れを確認")

		//入金検知 + 未署名トランザクション作成
		logger.Info("[1]Run: 入金検知")
		hex, fileName, err := wallet.DetectReceivedCoin(opts.Fee)
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		if hex == "" {
			logger.Info("No utxo")
			return
		}
		logger.Infof("[hex]: %s\n[fileName]: %s", hex, fileName)

		//署名(本来はColdWalletの機能)
		logger.Info("\n[2]Run: 署名")
		hexTx, isSigned, generatedFileName, err := wallet.SignatureFromFile(fileName)
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		logger.Infof("[hex]: %s\n[署名完了]: %t\n[fileName]: %s", hexTx, isSigned, generatedFileName)

		//送信
		logger.Info("\n[3]Run: 送信")
		txID, err := wallet.SendFromFile(generatedFileName)
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		logger.Infof("[Done]送信までDONE!! txID: %s", txID)

		//一連の署名から送信までの流れをチェック
		//[WIF] cUW7ZSF9WX7FUTeHkuw5L9Rj26V5Kz8yCkYjZamyvATTwsu7KUCi - [Pub Address] muVSWToBoNWusjLCbxcQNBWTmPjioRLpaA
		//hash, tx, err := wallet.BTC.SequentialTransaction(hex)
		//if err != nil {
		//	logger.Fatalf("%+v", err)
		//}
		////tx.MsgTx()
		//logger.Debugf("送信までDONE!! %s, %v", hash.String(), tx)

	default:
		logger.Warn("opts.Mode is out of range")
		procedure.ShowWallet()
	}

}

//出金関連機能
func paymentFunctionalities(wallet *service.Wallet) {
	switch opts.Mode {
	case 1:
		logger.Info("Run:出金のための未署名トランザクション作成")
		hex, fileName, err := wallet.CreateUnsignedTransactionForPayment(opts.Fee)
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		if hex == "" {
			logger.Info("No utxo")
			return
		}
		logger.Infof("[hex]: %s, \n[fileName]: %s", hex, fileName)
	case 10:
		logger.Info("Run: [Debug用]出金から送金までの一連の流れを確認")

		//出金準備
		logger.Info("[1]Run:出金のための未署名トランザクション作成")
		hex, fileName, err := wallet.CreateUnsignedTransactionForPayment(opts.Fee)
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		if hex == "" {
			logger.Info("No utxo")
			return
		}
		logger.Infof("[hex]: %s, \n[fileName]: %s", hex, fileName)

		//署名(本来はColdWalletの機能)
		logger.Info("\n[2]Run: 署名")
		hexTx, isSigned, generatedFileName, err := wallet.SignatureFromFile(fileName)
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		logger.Infof("[hex]: %s\n[署名完了]: %t\n[fileName]: %s", hexTx, isSigned, generatedFileName)

		//送信
		logger.Info("\n[3]Run: 送信")
		txID, err := wallet.SendFromFile(generatedFileName)
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		logger.Infof("[Done]送信までDONE!! txID: %s", txID)

	default:
		logger.Warn("opts.Mode is out of range")
		procedure.ShowWallet()
	}
}

//署名の送信 関連機能
func sendingFunctionalities(wallet *service.Wallet) {
	switch opts.Mode {
	case 1:
		logger.Info("Run: ファイルから署名済みtxを送信する")
		// 1.GPSにupload(web管理画面から行う??)
		// 2.Uploadされたtransactionファイルから、送信する？
		if opts.ImportFile == "" {
			logger.Fatal("file path is required as argument file when running")
		}
		// フルパスを指定する
		txID, err := wallet.SendFromFile(opts.ImportFile)
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		logger.Infof("[Done]送信までDONE!! txID: %s", txID)
	default:
		logger.Warn("opts.Mode is out of range")
		procedure.ShowWallet()
	}
}

//transactionの監視 関連機能
func monitoringFunctionalities(wallet *service.Wallet) {
	switch opts.Mode {
	case 1:
		logger.Info("Run: 送信済ステータスのトランザクションを監視する")
		err := wallet.UpdateStatus()
		if err != nil {
			logger.Fatalf("%+v", err)
		}
	default:
		logger.Warn("opts.Mode is out of range")
		procedure.ShowWallet()
	}

}

// bitcoin RPC command
func btcCommand(wallet *service.Wallet) {
	switch opts.Mode {
	case 1:
		//入金検知処理後、lockされたunspenttransactionの解除を行う
		logger.Info("Run: lockされたトランザクションの解除")
		err := wallet.BTC.UnlockAllUnspentTransaction()
		if err != nil {
			logger.Fatalf("%+v", err)
		}
	case 2:
		//手数料算出
		logger.Info("Run: 手数料算出 estimatesmartfee")
		feePerKb, err := wallet.BTC.EstimateSmartFee()
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		logger.Infof("Estimatesmartfee: %f\n", feePerKb)
	case 3:
		//ロギング
		logger.Info("Run: ロギング logging")
		logData, err := wallet.BTC.Logging()
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		grok.Value(logData)
	case 4:
		//getnetworkinfoの呼び出し
		logger.Info("Run: INFO getnetworkinfo")
		infoData, err := wallet.BTC.GetNetworkInfo()
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		grok.Value(infoData)
		logger.Infof("Relayfee: %f", infoData.Relayfee)
	case 5:
		//ValidateAddress
		logger.Info("Run: AddressのValidationチェック")
		_, err := wallet.BTC.ValidateAddress("2NFXSXxw8Fa6P6CSovkdjXE6UF4hupcTHtr")
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		_, err = wallet.BTC.ValidateAddress("4VHGkbQTGg2vN5P6yHZw3UJhmsBh9igsSos")
		if err == nil {
			logger.Fatal("something is wrong")
		}
	default:
		logger.Warn("opts.Mode is out of range")
		procedure.ShowWallet()
	}
}

// 検証用
func debugForCheck(wallet *service.Wallet) {
	switch opts.Mode {
	case 1:
		//[Debug用]payment_requestテーブルを作成する
		logger.Info("Run: payment_requestテーブルを作成する")
		err := testdata.CreateInitialTestData(wallet.DB, wallet.BTC)
		if err != nil {
			logger.Fatal(err)
		}
	case 2:
		//[Debug用]payment_requestテーブルの情報を初期化する
		logger.Info("Run: payment_requestテーブルの情報を初期化する")
		_, err := wallet.DB.ResetAnyFlagOnPaymentRequestForTestOnly(nil, true)
		if err != nil {
			log.Fatalf("%+v", err)
		}
	default:
		logger.Warn("opts.Mode is out of range")
		procedure.ShowWallet()
	}

}
