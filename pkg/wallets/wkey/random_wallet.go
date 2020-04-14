package wkey

import (
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
)

//ランダムウォレット
//実運用では利用しない

// CreatePrivateKey private keyを作成する
func createPrivateKey(conf *chaincfg.Params) (*btcutil.WIF, error) {
	secret, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return nil, err
	}
	return btcutil.NewWIF(secret, conf, true)
}

// publicなアドレスを返す
func getAddressPubKey(wif *btcutil.WIF, conf *chaincfg.Params) (*btcutil.AddressPubKey, error) {
	return btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeCompressed(), conf)
}

//func ImportPrivateKey(secretHex string) (*btcutil.WIF, error) {
//}

// ImportWIF stringをインポートし、WIFを生成する
func ImportWIF(wifStr string, conf *chaincfg.Params) (*btcutil.WIF, error) {
	//  WIF(Wallet Import Format): 秘密鍵をより簡潔に表現したもの
	wif, err := btcutil.DecodeWIF(wifStr)
	if err != nil {
		return nil, err
	}
	if !wif.IsForNet(conf) {
		return nil, errors.New("The WIF string is not valid")
	}
	return wif, nil
}

// GenerateKey 単一のWIF, 公開鍵のアドレス(string)を生成する
func GenerateKey(conf *chaincfg.Params) (*btcutil.WIF, string, error) {
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

	//Debug
	//log.Printf("[WIF] %s - [Pub Address] %s\n", wif.String(), pubAddress.EncodeAddress())

	return wif, pubAddress.EncodeAddress(), nil
}
