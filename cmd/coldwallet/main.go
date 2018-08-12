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
	wallet := service.Wallet{Btc: bit, Db: nil}

	//switch
	switchFunction(&wallet)
}

func switchFunction(wallet *service.Wallet) {
	// 処理をFunctionalityで切り替える
	//TODO:ここから呼び出すべきはService系のみに統一したい
	switch opts.Functionality {
	case 1:
		//TODO:cold wallet側の機能
		log.Print("Run: Keyの生成")
		//単一Keyの生成
		wif, pubAddress, err := key.GenerateKey(wallet.Btc.GetChainConf())
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("[WIF] %s - [Pub Address] %s\n", wif.String(), pubAddress)
	case 2:
		//TODO:まだ検証中
		log.Print("Run: HDウォレット Keyの生成")
		key.GenerateHDKey(opts.ParamSeed, wallet.Btc.GetChainConf())
	case 3:
		//TODO:ImportしたHEXから署名を行う
		log.Print("Run: ImportしたHEXから署名を行う")

	default:
		log.Print("Run: 検証コード")
		// for test
		callAPI(wallet)
	}

}

func callAPI(wallet *service.Wallet) {

}
