package main

import (
	"log"

	"github.com/hiromaily/go-bitcoin/api"
	"github.com/jessevdk/go-flags"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

type Options struct {
	Host string `short:"s" long:"host" default:"127.0.0.1:18332" description:"Host and Port of RPC Server"`
	User string `short:"u" long:"user" default:"xyz" description:"User of RPC Server"`
	Pass string `short:"p" long:"pass" default:"xyz" description:"Password of RPC Server"`
}

var (
	opts Options
)

func init() {
	if _, err := flags.Parse(&opts); err != nil {
		panic(err)
	}
}

func main() {
	// Connection
	//bit, err := bitcoin.Connection("127.0.0.1:18332", "xyz", "xyz", true, true)
	bit, err := api.Connection(opts.Host, opts.User, opts.Pass, true, true)
	if err != nil{
		log.Fatal(err)
	}
	defer bit.Close()

	//
	callAPI(bit)
}

// For example
func callAPI(bit *api.Bitcoin) {
	//Block
	blockCount, err := bit.Client.GetBlockCount()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Block count: %d", blockCount)

	// Getnewaddress
	addr, err := bit.Client.GetNewAddress("ben")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("created address: %v", addr)


	// Unspent
	//[]btcjson.ListUnspentResult, error
	list, err := bit.Client.ListUnspent()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("List Unspent: %v", list)

	// Listaccounts
	// map[string]btcutil.Amount, error
	listAcnt, err := bit.Client.ListAccounts()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("List Accounts: %v", listAcnt)

	// Getbalance
	// btcutil.Amount, error
	amount, err := bit.Client.GetBalance("hiroki")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Hiroki's Accounts: %v", amount)

	// Gettransaction
	hash, err := chainhash.NewHashFromStr("5fe20dace7be113a73e5324194e20d24ae39307dd749b623fd7fe3f65115cadb")
	if err != nil {
		log.Fatal(err)
	}
	tran, err := bit.Client.GetTransaction(hash)
	log.Printf("Transactions: %v", tran)

	// Gettxout
	// txHash *chainhash.Hash, index uint32, mempool bool
	txOut, err := bit.Client.GetTxOut(hash, 0, false)
	log.Printf("TxOut: %v", txOut)

	// CreateRawTransaction
	//func (c *Client) CreateRawTransaction(inputs []btcjson.TransactionInput,
	//	amounts map[btcutil.Address]btcutil.Amount, lockTime *int64) (*wire.MsgTx, error) {
	//createrawtransaction "[{\"txid\":\"myid\",\"vout\":0}]" "{\"address\":0.01}"
	//btcjson.TransactionInput{Txid:"5fe20dace7be113a73e5324194e20d24ae39307dd749b623fd7fe3f65115cadb", }
	//bit.Client.CreateRawTransaction(
	//	[]btcjson.TransactionInput{{Txid: "123", Vout: 1}},
	//	map[btcutil.Address]btcutil.Amount{"456": .0123}, nil)

}