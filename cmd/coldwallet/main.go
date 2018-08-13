package main

import (
	"log"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/hiromaily/go-bitcoin/pkg/api"
	"github.com/hiromaily/go-bitcoin/pkg/key"
	"github.com/hiromaily/go-bitcoin/pkg/service"
	"github.com/hiromaily/go-bitcoin/pkg/toml"
	"github.com/jessevdk/go-flags"
)

// HDウォレットとしてseed作成、keyを指定した数だけ生成し、出力する
// これは、ネットワーク環境下のwallet側から、

// Options コマンドラインオプション
type Options struct {
	//Configパス
	ConfPath string `short:"c" long:"conf" default:"./data/toml/config.toml" description:"Path for configuration toml file"`
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

	// Connection to Bitcoin core
	//bit, err := api.Connection(conf.Bitcoin.Host, conf.Bitcoin.User, conf.Bitcoin.Pass, true, true, conf.Bitcoin.IsMain)
	bit, err := api.Connection(&conf.Bitcoin)
	if err != nil {
		log.Fatal(err)
	}
	defer bit.Close()

	//Wallet Object
	wallet := service.Wallet{Btc: bit, DB: nil}

	//switch
	switchFunction(&wallet)
}

func switchFunction(wallet *service.Wallet) {
	// 処理をFunctionalityで切り替える
	//TODO:ここから呼び出すべきはService系のみに統一したい
	switch opts.Functionality {
	case 1:
		//TODO: 通常のKeyの生成
		log.Print("Run: Keyの生成")
		//単一Keyの生成
		wif, pubAddress, err := key.GenerateKey(wallet.Btc.GetChainConf())
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("[WIF] %s - [Pub Address] %s\n", wif.String(), pubAddress)
	case 2:
		//TODO: HDウォレットによるKeyの作成 (まだ検証中)
		log.Print("Run: HDウォレット Keyの生成")
		key.GenerateHDKey(opts.ParamSeed, wallet.Btc.GetChainConf())
	case 3:
		//事前準備
		//getnewaddress taro 2N7ZwUXpo841GZDpxLGFqrhr1xwMzTba7ZP
		//getnewaddress boss1 2NAm558FWpiaJQLz838vbzBPpqmKxyeyxsu
		//TODO:ここで、AddMultisigAddressを使うのにパラメータとしてaccout名も渡さないといけない。。これをどうすべきか。。。
		//TODO: => おそらくBlankでもいい

		//TODO: Multisigアドレス作成 (まだ検証中)
		resAddr, err := wallet.Btc.CreateMultiSig(2, []string{"2N7ZwUXpo841GZDpxLGFqrhr1xwMzTba7ZP", "2NAm558FWpiaJQLz838vbzBPpqmKxyeyxsu"}, "multi01")
		if err != nil {
			log.Fatalf("%+v", err)
		}
		log.Printf("multisig address: %s, redeemScript: %s", resAddr.Address, resAddr.RedeemScript)
		//multisig address: 2N4Rm1aLPxCcg1H1V96bBzH69vMAipADLCQ, redeemScript: 522103d69e07dbf6da065e6fae1ef5761d029b9ff9143e75d579ffc439d47484044bed2103748797877523b8b36add26c9e0fb6a023f05083dd4056aedc658d2932df1eb6052ae

		//TODO:ここで生成されたアドレスに送金してみる。
		// https://testnet.manu.backend.hamburg/faucet
		//  Sent! TX ID: e278ce9750da9b89972001c3c221aa178e8ed4c187d5bef2513023e5a4bdcb9d
		// 現時点で、hokan以外ではlistunspentで取得できない。
		// これで、DetectReceivedCoin()を実行し、hexを取得

		// service.MultiSigByHex(hex)を実行してみる。TODO:WIP
	case 4:
		//TODO:ImportしたHEXから署名を行う()
		log.Print("Run: ImportしたHEXから署名を行う")

	default:
		log.Print("Run: 検証コード")
		// for test
		callAPI(wallet)
	}

}

func callAPI(wallet *service.Wallet) {

}
