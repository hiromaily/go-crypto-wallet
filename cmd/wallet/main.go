package main

import (
	"log"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/hiromaily/go-bitcoin/btc/api"
	"github.com/hiromaily/go-bitcoin/btc/service"
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
	//接続情報
	Host string `short:"s" long:"server" default:"127.0.0.1:18332" description:"Host and Port of RPC Server"`
	User string `short:"u" long:"user" default:"xyz" description:"User of RPC Server"`
	Pass string `short:"p" long:"pass" default:"xyz" description:"Password of RPC Server"`
	//接続先: MainNet or TestNet
	IsMain bool `short:"m" long:"ismain" description:"Using MainNetParams as network permeters or Not"`
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
	// Connection
	//bit, err := bitcoin.Connection("127.0.0.1:18332", "xyz", "xyz", true, true)
	bit, err := api.Connection(opts.Host, opts.User, opts.Pass, true, true, opts.IsMain)
	if err != nil {
		log.Fatal(err)
	}
	defer bit.Close()

	// 処理をFunctionalityで切り替える
	//TODO:ここから呼び出すべきはService系のみに統一したい
	switch opts.Functionality {
	case 1:
		//TODO:cold wallet側の機能
		log.Print("Run: Keyの生成")
		//単一Keyの生成
		wif, pubAddress, err := bit.GenerateKey("btc")
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("[WIF] %s - [Pub Address] %s\n", wif.String(), pubAddress)
	case 2:
		//TODO:まだ検証中
		log.Print("Run: HDウォレット Keyの生成")
		bit.GenerateHDKey(opts.ParamSeed)
	case 3:
		log.Print("Run: 入金処理検知")
		//入金検知
		//TODO:処理中にして、再度対象としないようにしないといけない
		_, err := service.DetectReceivedCoin(bit)
		if err != nil {
			log.Fatal(err)
		}
	case 9:
		log.Print("Run: [Debug用]送金までの一連の流れを確認")
		//入金検知
		//TODO:処理中にして、再度対象としないようにしないといけない
		tx, err := service.DetectReceivedCoin(bit)
		if err != nil {
			log.Fatal(err)
		}
		//署名
		signedTx, err := bit.SignRawTransaction(tx)
		if err != nil {
			log.Fatal(err)
		}
		//送金
		hash, err := bit.SendRawTransaction(signedTx)
		if err != nil {
			log.Fatal(err)
		}
		//min relay fee not met
		log.Printf("[Hash] %v", hash)

	default:
		log.Print("Run: 検証コード")
		// for test
		callAPI(bit)
	}

}

// 検証用
func callAPI(bit *api.Bitcoin) {
	//txOut
	//txOut, err := bit.GetTxOutByTxID("d0f3b258dda46a5980a0a9e1e6f818eb421be572d12e4e641b7b77e699ecddca", 0)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Printf("TxOut: %v\n", txOut)
	//grok.Value(txOut)
}
