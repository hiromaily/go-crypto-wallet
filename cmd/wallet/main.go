package main

import (
	"log"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/hiromaily/go-bitcoin/api"
	"github.com/hiromaily/go-bitcoin/service"
	"github.com/jessevdk/go-flags"
)

//TODO:ウォレットの定期バックアップ機能 + import機能
//TODO:coldウォレットへのデータ移行機能が必要なはず
//TODO:multisigの実装
//TODO:オフラインで可能機能と、不可能な機能の切り分けが必要

type Options struct {
	Host   string `short:"s" long:"server" default:"127.0.0.1:18332" description:"Host and Port of RPC Server"`
	User   string `short:"u" long:"user" default:"xyz" description:"User of RPC Server"`
	Pass   string `short:"p" long:"pass" default:"xyz" description:"Password of RPC Server"`
	IsMain bool   `short:"m" long:"ismain" description:"Using MainNetParams as network permeters or Not"`
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

	// for test
	//callAPI(bit)

	//入金検知
	//TODO:処理中にして、再度対象としないようにしないといけない
	err = service.DetectReceivedCoin(bit)
	if err != nil {
		log.Fatal(err)
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
