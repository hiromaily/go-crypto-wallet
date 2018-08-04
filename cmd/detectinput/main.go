package main

import (
	"log"

	"github.com/hiromaily/go-bitcoin"
)

func main() {
	// Connection
	bit, err := bitcoin.Connection("127.0.0.1:18332", "xyz", "xyz", true, true)
	if err != nil{
		log.Fatal(err)
	}
	defer bit.Close()

	// Get the current block count.
	blockCount, err := bit.Client.GetBlockCount()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Block count: %d", blockCount)
}
