package main

import (
	"log"

	"github.com/bookerzzz/grok"
	"github.com/btcsuite/btcd/chaincfg"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hiromaily/go-bitcoin/pkg/service"
	"github.com/jessevdk/go-flags"
)

//こちらはHotwallet、ただし、Watch Only Walletとしての機能を実装していく。
//ネットワークへの接続はGCP上のBitcoin Core
//Watch Only Walletとしてのセットアップが必要
// - Cold Wallet側から生成したPublic Key をMultisigアドレス変換後、`importaddress xxxxx`でimportする
//   これがかなり時間がかかる。。。実運用ではどうすべきか。rescanしなくても最初はOKかと

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
	ConfPath string `short:"c" long:"conf" default:"./data/toml/config.toml" description:"Path for configuration toml file"`
	//実行される機能
	Functionality uint8 `short:"f" long:"function" description:"Functionality: 1: generate key, 2: detect received coin, other: debug"`
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
		log.Fatal(err)
	}
	defer wallet.Done()

	//switch
	switchFunction(wallet)
}

func switchFunction(wallet *service.Wallet) {
	// 処理をFunctionalityで切り替える
	//TODO:ここから呼び出すべきはService系のみに統一したい
	switch opts.Functionality {
	case 1:
		//[Debug用]入金検知処理後、lock解除を行う
		log.Print("Run: lockされたトランザクションの解除")
		wallet.Btc.UnlockAllUnspentTransaction()
	case 2:
		//[Debug用]手数料算出
		log.Print("Run: 手数料算出 estimatesmartfee")
		feePerKb, err := wallet.Btc.EstimateSmartFee()
		if err != nil {
			log.Fatalf("%+v", err)
		}
		log.Printf("Estimatesmartfee: %f\n", feePerKb)
	case 3:
		//[Debug用]手数料算出
		log.Print("Run: ロギング logging")
		logData, err := wallet.Btc.Logging()
		if err != nil {
			log.Fatalf("%+v", err)
		}
		//Debug
		grok.Value(logData)
	case 4:
		//[Debug用]ValidateAddress
		log.Print("Run: AddressのValidationチェック")
		err := wallet.Btc.ValidateAddress("2NFXSXxw8Fa6P6CSovkdjXE6UF4hupcTHtr")
		if err != nil {
			log.Fatalf("%+v", err)
		}
		err = wallet.Btc.ValidateAddress("4VHGkbQTGg2vN5P6yHZw3UJhmsBh9igsSos")
		if err == nil {
			log.Fatal("something is wrong")
		}

		log.Printf("error is %v", err)
		log.Print("Done!")
	case 11:
		log.Print("Run: 入金処理検知")
		//実際には署名処理は手動なので、ユーザーの任意のタイミングで走らせたほうがいい。

		//入金検知 + 未署名トランザクション作成
		//TODO:この中でLoopする必要はない。実行するtaskrunner側で実行間隔を調整する。
		hex, err := wallet.DetectReceivedCoin()
		if err != nil {
			log.Fatalf("%+v", err)
		}
		if hex == "" {
			log.Printf("No utxo")
			return
		}
		log.Printf("hex: %s", hex)
	case 12:
		log.Print("Run: 署名済みtxを送信する")

		hex := "020000000001019dcbbda4e5233051f2bed587c1d48e8e17aa21c2c3012097899bda5097ce78e201000000232200208e1343e11e4def66d7102d9b0f36f019188118df5a5f30dacdd1008928b12f5fffffffff01042bbf070000000017a9148191d41a7415a6a1f6ee14337e039f50b949e80e870400483045022100f4975a5ea23e5799b1df65d699f85236b9d00bcda8da333731ffa508285d3c59022037285857821ee68cbe5f74239299170686b108ce44e724a9a280a3ef9291746901483045022100f94ce83946b4698b8dfbb7cb75eece12932c5097017e70e60d924aeae1ec829a02206e7b2437e9747a9c28a3a3d7291ea16db1d2f0a60482cdb8eca91c28c01aba790147522103d69e07dbf6da065e6fae1ef5761d029b9ff9143e75d579ffc439d47484044bed2103748797877523b8b36add26c9e0fb6a023f05083dd4056aedc658d2932df1eb6052ae00000000"
		hash, err := wallet.Btc.SendTransactionByHex(hex)
		if err != nil {
			log.Fatalf("%+v", err)
		}
		log.Printf("[Debug] 送信までDONE!! %s", hash.String())
	case 13:
		//TODO:実際にはUSBで署名済の16進数で書かれたtransactionの文字列を処理するのだが、
		// 1.GPSにupload(web管理画面から行う??)
		// 2.Uploadされたtransactionファイルから、送信する？
		// 3. unsignedトランザクションファイル名と、signedトランザクションファイル名のリレーションをDBに保存したほうがいいか？

		//[Debug]とりあえず、localのファイルから読みこんで送信してみる
		txID, err := wallet.SendFromFile("./data/tx/signed_1534303789374409841")
		if err != nil {
			log.Fatalf("%+v", err)
		}
		log.Printf("[Debug] 送信までDONE!! %s", txID)
	case 14:
		log.Print("Run:出金トランザクション作成")
		wallet.CreateUnsignedTransactionForPayment()

	case 20:
		log.Print("Run: [Debug用]送金までの一連の流れを確認")

		//Debug中のみ
		//wallet.Btc.UnlockAllUnspentTransaction()

		//入金検知 + 未署名トランザクション作成
		hex, err := wallet.DetectReceivedCoin()
		if err != nil {
			log.Fatal(err)
		}
		if hex == "" {
			log.Printf("No utxo")
			return
		}
		log.Printf("hex: %s", hex)

		//一連の署名から送信までの流れをチェック
		//[WIF] cUW7ZSF9WX7FUTeHkuw5L9Rj26V5Kz8yCkYjZamyvATTwsu7KUCi - [Pub Address] muVSWToBoNWusjLCbxcQNBWTmPjioRLpaA
		hash, tx, err := wallet.Btc.SequentialTransaction(hex)
		if err != nil {
			log.Fatalf("%+v", err)
		}
		//tx.MsgTx()
		log.Printf("[Debug] 送信までDONE!! %s, %v", hash.String(), tx)

	default:
		log.Print("Run: 検証コード")
		// for test
		callAPI(wallet)
	}

}

// 検証用
func callAPI(wallet *service.Wallet) {
	//txOut
	//txOut, err := bit.GetTxOutByTxID("d0f3b258dda46a5980a0a9e1e6f818eb421be572d12e4e641b7b77e699ecddca", 0)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Printf("TxOut: %v\n", txOut)
	//grok.Value(txOut)
}

//func checkLevelDB(wallet *service.Wallet) {
//	//Put
//	err := wallet.DB.Put("unspent", "testkey1", []byte("data1234567890"))
//	if err != nil {
//		log.Println(err)
//	}
//	//Get
//	val, err := wallet.DB.Get("unspent", "testkey1")
//	if err != nil {
//		log.Println(err)
//	}
//	log.Printf("[Done] %s", string(val))
//}
