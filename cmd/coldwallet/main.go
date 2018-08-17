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

//TODO:encryptwalletコマンドによって、walletを暗号化した場合、秘密鍵を使用するタイミング(未署名トランザクションに署名する)
// でパスフレーズの入力が必要になり

// Options コマンドラインオプション
type Options struct {
	//Configパス
	ConfPath string `short:"c" long:"conf" default:"./data/toml/config.toml" description:"Path for configuration toml file"`
	//実行される機能
	Functionality uint8 `short:"f" long:"function" description:"Functionality: 1: generate key, 2: detect received coin, other: debug"`
	//HDウォレット用Key生成のためのseed情報
	ParamSeed string `short:"d" long:"seed" default:"" description:"backup seed"`
	//HDウォレット用Key生成のためのseed情報
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
	wallet := service.Wallet{BTC: bit, DB: nil}

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
		wif, pubAddress, err := key.GenerateKey(wallet.BTC.GetChainConf())
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("[WIF] %s - [Pub Address] %s\n", wif.String(), pubAddress)
	case 2:
		//TODO: HDウォレットによるKeyの作成 (まだ検証中)
		log.Print("Run: HDウォレット Keyの生成")
		key.GenerateHDKey(opts.ParamSeed, wallet.BTC.GetChainConf())
	case 3:
		//TODO:Multisigの作成

		//事前準備
		//getnewaddress taro 2N7ZwUXpo841GZDpxLGFqrhr1xwMzTba7ZP
		//getnewaddress boss1 2NAm558FWpiaJQLz838vbzBPpqmKxyeyxsu
		//TODO:ここで、AddMultisigAddressを使うのにパラメータとしてaccout名も渡さないといけない。。これをどうすべきか。。。
		//TODO: => おそらくBlankでもいい

		//TODO: Multisigアドレス作成 (まだ検証中)
		resAddr, err := wallet.BTC.CreateMultiSig(2, []string{"2N7ZwUXpo841GZDpxLGFqrhr1xwMzTba7ZP", "2NAm558FWpiaJQLz838vbzBPpqmKxyeyxsu"}, "multi01")
		if err != nil {
			log.Fatalf("%+v", err)
		}
		log.Printf("multisig address: %s, redeemScript: %s", resAddr.Address, resAddr.RedeemScript)
		//multisig address: 2N4Rm1aLPxCcg1H1V96bBzH69vMAipADLCQ, redeemScript: 522103d69e07dbf6da065e6fae1ef5761d029b9ff9143e75d579ffc439d47484044bed2103748797877523b8b36add26c9e0fb6a023f05083dd4056aedc658d2932df1eb6052ae

		//TODO:ここで生成されたアドレスに送金してみる。
		// https://testnet.manu.backend.hamburg/faucet
		//  Sent! TX ID: e278ce9750da9b89972001c3c221aa178e8ed4c187d5bef2513023e5a4bdcb9d
		// https://live.blockcypher.com/btc-testnet/tx/e278ce9750da9b89972001c3c221aa178e8ed4c187d5bef2513023e5a4bdcb9d/
		// 現時点で、hokan以外ではlistunspentで取得できないはず
		// これで、DetectReceivedCoin()を実行し、hexを取得
		// 02000000019dcbbda4e5233051f2bed587c1d48e8e17aa21c2c3012097899bda5097ce78e20100000000ffffffff01042bbf070000000017a9148191d41a7415a6a1f6ee14337e039f50b949e80e8700000000

		// service.MultiSigByHex(hex)を実行してみる。
	case 4:
		//TODO:HEXから署名を行う()
		log.Print("Run: HEXから署名を行う")
		//hex := "02000000019dcbbda4e5233051f2bed587c1d48e8e17aa21c2c3012097899bda5097ce78e20100000000ffffffff01042bbf070000000017a9148191d41a7415a6a1f6ee14337e039f50b949e80e8700000000"
		//hex := "02000000032e0183cd8e082c185030b8eed4bf19bace65936960fe79736dc21f3b0586b7640100000000ffffffff8afd01d2ecdfeb1657ae7a0ecee9e89b86feb58ed10803cdf6c95d25354161ff0100000000ffffffffc6f7645941324cfe9e36194a6443e0f50fe9117c4964ad942f39833da60363ba0000000000ffffffff01f0be8e0d0000000017a9148191d41a7415a6a1f6ee14337e039f50b949e80e8700000000"
		hex := "02000000021ed288be4c4d7923a0d044bb500a15b2eb0f2b3c5503293f251f7c94939a3f9f0000000000ffffffff557624120cdf3f4d092f35e5cd6b75418b76c3e3fd4c398357374e93cfe5c4200000000000ffffffff05c03b47030000000017a91419e70491572c55fb08ce90b0c6bf5cfe45a5420e87809698000000000017a9146b8902fc7a6a0bccea9dbd80a4c092c314227f618734e133070000000017a9148191d41a7415a6a1f6ee14337e039f50b949e80e87005a62020000000017a9149c877d6f21d5800ca60a7660ee56745f239b222b87002d31010000000017a914f575a0d1ddcfb98a11628826f1632453d718ff618700000000"
		//hexTx, err := wallet.MultiSigByHex(hex) //これはもう呼び出さない
		hexTx, isSigned, generatedFileName, err := wallet.SignatureByHex(hex, 10)
		if err != nil {
			log.Fatalf("%+v", err)
		}
		log.Printf("hex: %s\n, 署名完了: %t\n, fileName: %s", hexTx, isSigned, generatedFileName)
		//TODO:isSigned: 送信までした署名はfalseになる??
	case 5:
		//TODO: importしたファイルからhex値を取得し、署名を行う
		// ./coldwallet -f 5 -i 5_unsigned_1534466246366489473
		log.Print("Run: Importしたファイルからhex値を取得し、署名を行う")
		if opts.ImportFile == "" {
			log.Fatal("file path is required as argument file when running")
		}

		hexTx, isSigned, generatedFileName, err := wallet.SignatureFromFile(opts.ImportFile)
		if err != nil {
			log.Fatalf("%+v", err)
		}
		log.Printf("hex: %s\n, 署名完了: %t\n, fileName: %s", hexTx, isSigned, generatedFileName)

	default:
		log.Print("Run: 検証コード")
		// for test
		callAPI(wallet)
	}

}

func callAPI(wallet *service.Wallet) {

}
