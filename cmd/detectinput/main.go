package main

import (
	"log"

	"github.com/hiromaily/go-bitcoin/api"
	"github.com/jessevdk/go-flags"
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
	blockCount, err := bit.Client.GetBlockCount()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Block count: %d", blockCount)

}