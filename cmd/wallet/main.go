package main

import (
	"log"
	"os"

	"github.com/bookerzzz/grok"
	"github.com/btcsuite/btcd/chaincfg"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jessevdk/go-flags"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/logger"
	"github.com/hiromaily/go-bitcoin/pkg/procedure"
	"github.com/hiromaily/go-bitcoin/pkg/testdata"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/service"
)

//Watch Only Wallet
//Bitcoin Core would be on GCP network for production
//setup is required as watch only wallet

//TODO: bitcoin functionalities
// - back up wallet data periodically and import functionality
// - generated key must be encrypted
// - transfer with amount
//TODO:
// - logger interface: stdout(ui), log format
// - repository interface (from mysql, mock, redis and so on)
// - service interface
// - procedure pkg move to help of command
// - change Japanese to English

// Options is command line options
type Options struct {
	//Configパス
	ConfPath string `short:"c" long:"conf" default:"" description:"Path for configuration toml file"`

	//Keyモード
	Key bool `short:"k" long:"key" description:"for adding key"`
	//入金モード
	Receipt bool `short:"r" long:"receipt" description:"for receipt"`
	//出金モード
	Payment bool `short:"p" long:"payment" description:"for payment"`
	//内部アカウント送金モード
	Transfer bool `short:"t" long:"transfer" description:"for transfer"`
	//送金モード
	Send bool `short:"s" long:"sending" description:"for sending transaction"`
	//ステータスチェックモード
	Monitor bool `short:"n" long:"monitor" description:"for monitoring transaction"`
	//bitcoin Commandモード
	Cmd bool `short:"b" long:"bitcoin-command" description:"for bitcoin command"`
	//Debugモード
	Debug bool `short:"d" long:"debug" description:"for only development use"`

	//実行されるサブ機能
	Mode uint8 `short:"m" long:"mode" description:"Mode i.e.Functionality"`

	//txファイルパス
	ImportFile string `short:"i" long:"import" default:"" description:"import file path for hex"`
	//調整fee
	Fee float64 `short:"f" long:"fee" default:"" description:"adjustment fee"`
	//import時にscanするかどうか
	IsRescan bool `short:"x" long:"rescan" description:"scan blocks when importing key"`
	//アカウント (アカウント間送金の場合は、from account)
	Account string `short:"a" long:"account" description:"account like client, receipt, payment"`
	//アカウント (to account)
	Account2 string `short:"z" long:"account_to" description:"account like client, receipt, payment"`
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
	confPath := os.Getenv("WATCH_WALLET_CONF")
	if opts.ConfPath != "" {
		confPath = opts.ConfPath
	}

	//initialSettings()
	wallet, err := service.InitialSettings(confPath)
	if err != nil {
		// ここでエラーが出た場合、まだloggerの初期化が終わってない
		//logger.Fatal(err)
		log.Fatal(err)
	}
	wallet.Type = enum.WalletTypeWatchOnly
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
	} else if opts.Transfer {
		//内部アカウント転送関連機能
		transferFunctionalities(wallet)
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

func checkImportFile() {
	if opts.ImportFile == "" {
		logger.Fatal("file path is required as option -i")
	}
}

func checkAccountWithoutAuth() {
	if opts.Account == "" || !account.ValidateAccountType(opts.Account) ||
		opts.Account == account.AccountTypeAuthorization.String() {
		logger.Fatal("Account[client, receipt, payment, quoine, fee, stored] should be set with -a option")
	}
}

func checkAccountWithoutAuthAndClient(acnt string) {
	if acnt == "" || !account.ValidateAccountType(acnt) ||
		acnt == account.AccountTypeAuthorization.String() || acnt == account.AccountTypeClient.String() {
		logger.Fatal("Account[receipt, payment, quoine, fee, stored] should be set with -a option")
	}
}

//キー関連機能[k]
func keyFunctionalities(wallet *service.Wallet) {
	switch opts.Mode {
	case 1:
		logger.Info("Run: coldwalletで生成したAccountのアドレスをwalletにimportする")
		checkImportFile()
		checkAccountWithoutAuth()
		logger.Infof("Run: Account[%s]", opts.Account)

		//import public key
		err := wallet.ImportPublicKeyForWatchWallet(opts.ImportFile, account.AccountType(opts.Account), opts.IsRescan)
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		logger.Info("Done!")

	default:
		logger.Warn("opts.Mode is out of range")
		procedure.ShowWallet()
	}
}

//入金関連機能[r]
func receiptFunctionalities(wallet *service.Wallet) {
	switch opts.Mode {
	case 1:
		logger.Info("Run: 入金処理検知 + 未署名トランザクション作成")
		//実際には署名処理は手動なので、ユーザーの任意のタイミングで実行する
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

	default:
		logger.Warn("opts.Mode is out of range")
		procedure.ShowWallet()
	}

}

//出金関連機能[p]
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

//内部アカウント転送関連機能[t]
func transferFunctionalities(wallet *service.Wallet) {
	switch opts.Mode {
	case 1:
		logger.Info("Run:内部アカウント転送のための未署名トランザクション作成")
		checkAccountWithoutAuthAndClient(opts.Account)
		checkAccountWithoutAuthAndClient(opts.Account2)

		hex, fileName, err := wallet.SendToAccount(account.AccountType(opts.Account), account.AccountType(opts.Account2), 0)
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		if hex == "" {
			logger.Info("No utxo")
			return
		}
		logger.Infof("[hex]: %s, \n[fileName]: %s", hex, fileName)
	}
}

//署名の送信 関連機能[s]
func sendingFunctionalities(wallet *service.Wallet) {
	switch opts.Mode {
	case 1:
		logger.Info("Run: ファイルから署名済みtxを送信する")
		if opts.ImportFile == "" {
			logger.Fatal("file path is required as argument file when running")
		}
		// 送信: フルパスを指定する
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

//transactionの監視 関連機能[n]
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

// bitcoin RPC command[b]
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
	case 6:
		//入金検知のみ: listunspentのみ
		logger.Info("Run: listunspentのみ")
		unspentList, err := wallet.BTC.Client().ListUnspentMin(wallet.BTC.ConfirmationBlock()) //6
		if err != nil {
			logger.Fatalf("%+v", err)
		}
		logger.Debug("List Unspent")
		grok.Value(unspentList) //Debug
	default:
		logger.Warn("opts.Mode is out of range")
		procedure.ShowWallet()
	}
}

// 検証用[d]
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
			logger.Fatalf("%+v", err)
		}
	case 3:
		//[Debug用]payment_requestテーブルの情報を初期化する
		logger.Info("Run: I/Fが変わってエラーが出るようになったのでテスト")
		logger.Debugf("account: %s, confirmation block: %d", account.AccountTypePayment.String(), wallet.BTC.ConfirmationBlock())
		//FIXME:wallet.BTC.GetBalanceByAccountAndMinConf()の呼び出しをやめて、GetReceivedByAccountAndMinConf()をcallするように変更する
		//balance, err := wallet.BTC.GetBalanceByAccountAndMinConf(string(enum.AccountTypePayment), wallet.BTC.ConfirmationBlock())
		balance, err := wallet.BTC.GetReceivedByAccountAndMinConf(account.AccountTypePayment.String(), wallet.BTC.ConfirmationBlock())
		if err != nil {
			log.Fatalf("%+v", err)
		}
		logger.Infof("balance: %v", balance)
	case 4:
		//bash戻り値検証用
		logger.Info("aaabbb")
		logger.Info("[fileName]: ./data/pubkey/stored_1_1537438241658920000.csv")
	default:
		logger.Warn("opts.Mode is out of range")
		procedure.ShowWallet()
	}

}
