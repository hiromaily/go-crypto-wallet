package key

import (
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
)

// random wallet
// - this is not used, just for test

// GenerateWIF generate one WIF, pubKey
//  WIF(Wallet Import Format) something simple format for privKey
func GenerateWIF(conf *chaincfg.Params) (*btcutil.WIF, string, error) {
	// Create Private Key
	wif, err := createPrivateKey(conf)
	if err != nil {
		return nil, "", err
	}

	// Return Public Key address
	pubAddress, err := getAddressPubKey(wif, conf)
	if err != nil {
		return nil, "", err
	}

	// Debug
	// log.Printf("[WIF] %s - [Pub Address] %s\n", wif.String(), pubAddress.EncodeAddress())

	return wif, pubAddress.EncodeAddress(), nil
}

func createPrivateKey(conf *chaincfg.Params) (*btcutil.WIF, error) {
	secret, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return nil, err
	}
	return btcutil.NewWIF(secret, conf, true)
}

func getAddressPubKey(wif *btcutil.WIF, conf *chaincfg.Params) (*btcutil.AddressPubKey, error) {
	return btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeCompressed(), conf)
}

// ImportWIF generate WIF from string
func ImportWIF(wifStr string, conf *chaincfg.Params) (*btcutil.WIF, error) {
	wif, err := btcutil.DecodeWIF(wifStr)
	if err != nil {
		return nil, err
	}
	if !wif.IsForNet(conf) {
		return nil, errors.New("WIF string is not valid")
	}
	return wif, nil
}

//func ImportPrivateKey(secretHex string) (*btcutil.WIF, error) {
//}
