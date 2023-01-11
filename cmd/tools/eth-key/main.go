package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
)

func main() {
	// args
	var (
		dir      = flag.String("keydir", "./tmp", "keystore directory, `tmp` dir would be preferable")
		keyfile  = flag.String("keyfile", "", "keystore file full path in keystore")
		password = flag.String("password", "", "password")
	)
	flag.Parse()

	// validate args
	if *dir == "" {
		log.Fatal("args `keydir` must not be empty")
	}
	if *keyfile == "" {
		log.Fatal("args `keyfile` must not be empty")
	}
	if *password == "" {
		log.Fatal("args `password` must not be empty")
	}

	// create key store
	ks := keystore.NewKeyStore(*dir, keystore.StandardScryptN, keystore.StandardScryptP)

	// read files
	jsonBytes, err := os.ReadFile(*keyfile)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("run ks.Import()")
	// FIXME: cannot unmarshal object into Go struct field CryptoJSON.crypto.kdf of type string
	account, err := ks.Import(jsonBytes, *password, *password)
	if err != nil {
		log.Fatal(err)
	}

	// print address
	fmt.Println(account.Address.Hex())
}
