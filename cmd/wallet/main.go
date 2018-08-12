package main

import (
	"log"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/hiromaily/go-bitcoin/pkg/api"
	"github.com/hiromaily/go-bitcoin/pkg/key"
	"github.com/hiromaily/go-bitcoin/pkg/kvs"
	"github.com/hiromaily/go-bitcoin/pkg/service"
	"github.com/hiromaily/go-bitcoin/pkg/toml"
	"github.com/jessevdk/go-flags"
)

//TODO:coldwallet側(非ネットワーク環境)側の機能と明確に分ける
//TODO:オフラインで可能機能と、不可能な機能の切り分けが必要
//TODO:ウォレットの定期バックアップ機能 + import機能
//TODO:coldウォレットへのデータ移行機能が必要なはず
//TODO:multisigの実装
//TODO:生成したkeyの暗号化周りのpkgが必要になるはず

// Options コマンドラインオプション
type Options struct {
	//Configパス
	ConfPath string `short:"c" long:"conf" default:"./data/toml/config.toml" description:"Path for configuration toml file"`
	//接続情報
	Host string `short:"s" long:"server" default:"127.0.0.1:18332" description:"Host and Port of RPC Server"`
	User string `short:"u" long:"user" default:"xyz" description:"User of RPC Server"`
	Pass string `short:"p" long:"pass" default:"xyz" description:"Password of RPC Server"`
	//接続先: MainNet or TestNet
	IsMain bool `short:"m" long:"ismain" description:"Using MainNetParams as network permeters or Not"`
	//KVSのstoreされるファイルパス
	DBPath string `short:"k" long:"kvspath" default:"./data/kvs/db" description:"Path for stored data by KVS"`
	//実行される機能
	Functionality uint8 `short:"f" long:"function" description:"Functionality: 1: generate key, 2: detect received coin, other: debug"`
	//HDウォレット用Key生成のためのseed情報
	ParamSeed string `short:"d" long:"seed" default:"" description:"backup seed"`
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
	conf, err := toml.New(opts.ConfPath)
	if err != nil {
		log.Fatal(err)
	}

	// KVS
	db, err := kvs.InitDB(opts.DBPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Connection to Bitcoin core
	bit, err := api.Connection(opts.Host, opts.User, opts.Pass, true, true, opts.IsMain)
	if err != nil {
		log.Fatal(err)
	}
	defer bit.Close()

	bit.SetConfiguration(conf)

	//Wallet Object
	wallet := service.Wallet{Btc: bit, Db: db}

	//switch
	switchFunction(&wallet)
}

func switchFunction(wallet *service.Wallet) {
	// 処理をFunctionalityで切り替える
	//TODO:ここから呼び出すべきはService系のみに統一したい
	switch opts.Functionality {
	case 1:
		//入金検知処理後、lock解除を行う
		log.Print("Run: lockされたトランザクションの解除")
		wallet.Btc.UnlockAllUnspentTransaction()
	case 2:
		//TODO:未実装
		log.Print("Run: 手数料算出 estimatesmartfee")
		feePerKb, err := wallet.Btc.EstimateSmartFee()
		if err != nil {
			log.Fatalf("%+v", err)
		}
		log.Printf("Estimatesmartfee: %f\n", feePerKb)

	case 11:
		log.Print("Run: 入金処理検知")

		//Debug中のみ
		wallet.Btc.UnlockAllUnspentTransaction()

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
		log.Print("Run: [Debug用]送金までの一連の流れを確認")

		//Debug中のみ
		wallet.Btc.UnlockAllUnspentTransaction()

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

	case 21:
		//TODO:cold wallet側の機能
		log.Print("Run: Keyの生成")
		//単一Keyの生成
		wif, pubAddress, err := key.GenerateKey(wallet.Btc.GetChainConf())
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("[WIF] %s - [Pub Address] %s\n", wif.String(), pubAddress)
	case 22:
		//TODO:まだ検証中
		log.Print("Run: HDウォレット Keyの生成")
		key.GenerateHDKey(opts.ParamSeed, wallet.Btc.GetChainConf())

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

func checkLevelDB(wallet *service.Wallet) {
	//Put
	err := wallet.Db.Put("unspent", "testkey1", []byte("data1234567890"))
	if err != nil {
		log.Println(err)
	}
	//Get
	val, err := wallet.Db.Get("unspent", "testkey1")
	if err != nil {
		log.Println(err)
	}
	log.Printf("[Done] %s", string(val))
}