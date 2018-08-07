package main

import (
	"log"

	"github.com/bookerzzz/grok"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"
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

	err = service.DetectReceivedCoin(bit)
	if err != nil {
		log.Fatal(err)
	}

}

// For example
func callAPI(bit *api.Bitcoin) {
	//Block
	blockCount, err := bit.Client.GetBlockCount()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Block count: %d\n", blockCount)

	// Getnewaddress
	//addr, err := bit.Client.GetNewAddress("ben")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Printf("created address: %v\n", addr)

	// Unspent => このレスポンスを元に処理を進める
	//[]btcjson.ListUnspentResult, error
	list, err := bit.Client.ListUnspent()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("List Unspent: %v\n", list)
	grok.Value(list)

	// Listaccounts => これは単純にアカウントの資産一覧が表示されるだけ
	//listAcnt, err := bit.Client.ListAccounts()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Printf("List Accounts: %v\n", listAcnt)
	//grok.Value(listAcnt)

	// Getbalance -> 個人レベルでしか取得できない
	// btcutil.Amount, error
	//amount, err := bit.Client.GetBalance("hiroki")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Printf("Hiroki's Accounts: %v\n", amount)
	//grok.Value(amount)

	// Gettransaction => トランザクション詳細情報を取得するのだが、ここから必要となる情報がどこで必要となるか？？
	hash, err := chainhash.NewHashFromStr("5fe20dace7be113a73e5324194e20d24ae39307dd749b623fd7fe3f65115cadb")
	if err != nil {
		log.Fatal(err)
	}
	tran, err := bit.Client.GetTransaction(hash)
	log.Printf("Transactions: %v\n", tran)
	grok.Value(tran)
	//TODO:Txid, Vout should be retrieved from result.

	// Gettxout
	// txHash *chainhash.Hash, index uint32, mempool bool
	//txOut, err := bit.Client.GetTxOut(hash, 0, false)
	//log.Printf("TxOut: %v\n", txOut)
	//grok.Value(txOut)

	// CreateRawTransaction
	//Getaddressesbyaccount "account"
	//TODO:集約用アドレス情報を、どこか内部的に保持しておく必要がある。
	//TODO:ドキュメント見る限り、ここはオンラインでいいはず
	//required:送金先のアドレス
	//required:送金金額
	//TODO:これはサンプルロジックであって、実際はこういうやり方はしない
	addrs, err := bit.Client.GetAddressesByAccount("hokan")
	if err != nil || len(addrs) == 0 {
		log.Fatal(err)
	}
	//sendAddr, err := btcutil.DecodeAddress(addrs[0].String(), &chaincfg.TestNet3Params) //for test(TODO:切り替えが必要)
	sendAddr, err := btcutil.DecodeAddress(addrs[0].String(), bit.GetChainConf())
	//1Satoshi＝0.00000001BTC
	//TODO:変換ロジック bitcoin to satoshiがあると楽
	msgTx, err := bit.Client.CreateRawTransaction(
		[]btcjson.TransactionInput{{Txid: "5fe20dace7be113a73e5324194e20d24ae39307dd749b623fd7fe3f65115cadb", Vout: 1}},
		map[btcutil.Address]btcutil.Amount{sendAddr: 80000000}, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("CreateRawTransaction: %v\n", msgTx)
	grok.Value(msgTx)

	//Signrawtransaction
	//TODO: It should be implemented on Cold Strage
	//この処理がHotwallet内で動くということは、重要な情報がwallet内に含まれてしまっているということでは？
	signed, isSigned, err := bit.Client.SignRawTransaction(msgTx)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Signrawtransaction: %v\n", signed)
	log.Printf("Signrawtransaction isSigned: %v\n", isSigned)

	//Sendrawtransaction

	//TODO:トランザクションのkbに応じて、手数料を算出

}
