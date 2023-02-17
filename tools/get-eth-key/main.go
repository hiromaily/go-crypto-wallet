package main

import (
	"crypto/ecdsa"
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/pkg/errors"
)

var (
	addr     = flag.String("addr", "", "address")
	keyDir   = flag.String("dir", "", "key directory")
	mnemonic = flag.String("mnemonic", "", "mnemonic")
	hdPath   = flag.String("hdpath", "", "hd path")
)

func main() {
	flag.Parse()

	fmt.Printf("addr: %s\n", *addr)
	fmt.Printf("keyDir: %s\n", *keyDir)
	fmt.Printf("mnemonic: %s\n", *mnemonic)
	fmt.Printf("hdPath: %s\n", *hdPath)

	// From address and keystore directory
	if *addr != "" && *keyDir != "" {
		fmt.Println("[Mode] From address and keystore directory")
		key, err := getPrivKey(*addr, "")
		if err != nil {
			panic(err)
		}
		fmt.Println(key.PrivateKey)
		return
	}
	// From mnemonic and hd path
	if *mnemonic != "" && *hdPath != "" {
		fmt.Println("From mnemonic and hd path")
		privKey, err := GetPrvKeyFromMnemonicAndHDWPath(*mnemonic, *hdPath)
		if err != nil {
			panic(err)
		}
		ethHexPrivKey := hexutil.Encode(crypto.FromECDSA(privKey))
		address, err := getAddress(privKey)
		if err != nil {
			panic(err)
		}
		fmt.Printf("privateKey: %s, address: %s\n", ethHexPrivKey, address)
		return
	}
	fmt.Println("[Mode] nothing to run")
}

func getAddress(privKey *ecdsa.PrivateKey) (string, error) {
	// pubkey, address
	ethPubkey := privKey.Public()
	pubkeyECDSA, ok := ethPubkey.(*ecdsa.PublicKey)
	if !ok {
		return "", errors.New("fail to call cast pubkey to ecsda pubkey")
	}
	// pubkey
	//ethHexPubKey := hexutil.Encode(crypto.FromECDSAPub(pubkeyECDSA))[4:]

	// address
	return crypto.PubkeyToAddress(*pubkeyECDSA).Hex(), nil
}

func getPrivKey(hexAddr, password string) (*keystore.Key, error) {
	keyJSON, err := os.ReadFile(fmt.Sprintf("%s/%s", *keyDir, hexAddr))
	if err != nil {
		return nil, err
	}
	if keyJSON == nil {
		// file is not found
		return nil, errors.New("private key file is not found")
	}

	key, err := keystore.DecryptKey(keyJSON, password)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to call keystore.DecryptKey()")
	}
	return key, nil
}
