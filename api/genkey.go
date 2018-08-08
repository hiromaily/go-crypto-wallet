package api

import (
	"errors"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
)

//TODO:WIP

// Network bitcoin系のnetwork parameters 情報
type Network struct {
	name        string
	symbol      string
	xpubkey     byte
	xprivatekey byte
}

var network = map[string]Network{
	"btc": {name: "bitcoin", symbol: "btc", xpubkey: 0x00, xprivatekey: 0x80},
	"ltc": {name: "litecoin", symbol: "ltc", xpubkey: 0x30, xprivatekey: 0xb0},
	"dgb": {name: "digibyte", symbol: "dgb", xpubkey: 0x1e, xprivatekey: 0x80},
	"rdd": {name: "reddcoin", symbol: "rdd", xpubkey: 0x3d, xprivatekey: 0xbd},
}

// GetNetworkParams network parametersを返す
//func (n Network) getNetworkParams(conf *chaincfg.Params) *chaincfg.Params {
//	networkParams := conf //Default is bitcoin
//	if n.name != "bitcoin" {
//		networkParams.PubKeyHashAddrID = n.xpubkey
//		networkParams.PrivateKeyID = n.xprivatekey
//	}
//
//	return networkParams
//}

// CreatePrivateKey private keyを作成する
func (n Network) createPrivateKey(conf *chaincfg.Params) (*btcutil.WIF, error) {
	secret, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return nil, err
	}
	return btcutil.NewWIF(secret, conf, true)
}

// publicなアドレスを返す
func (n Network) getAddressPubKey(wif *btcutil.WIF, conf *chaincfg.Params) (*btcutil.AddressPubKey, error) {
	return btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeCompressed(), conf)
}

//func (network Network) ImportPrivateKey(secretHex string) (*btcutil.WIF, error) {
//}

// ImportWIF stringをインポートし、WIFを生成する
// TODO:未
func (n Network) ImportWIF(wifStr string, conf *chaincfg.Params) (*btcutil.WIF, error) {
	wif, err := btcutil.DecodeWIF(wifStr)
	if err != nil {
		return nil, err
	}
	if !wif.IsForNet(conf) {
		return nil, errors.New("The WIF string is not valid for the `" + n.name + "` network")
	}
	return wif, nil
}

// TODO:こっちをServiceに移すべき
// GenerateKey 単一のWIF, 公開鍵のアドレス(string)を生成する
//  WIF(Wallet Import Format): 秘密鍵をより簡潔に表現したもの
func (b *Bitcoin) GenerateKey(symbol string) (*btcutil.WIF, string, error) {
	if symbol == "" {
		symbol = "btc"
	}

	// Create Private Key
	wif, err := network[symbol].createPrivateKey(b.GetChainConf())
	if err != nil {
		return nil, "", err
	}

	// Return Public Key address
	pubAddress, err := network[symbol].getAddressPubKey(wif, b.GetChainConf())
	if err != nil {
		return nil, "", err
	}

	//Debug
	//log.Printf("[WIF] %s - [Pub Address] %s\n", wif.String(), pubAddress.EncodeAddress())

	return wif, pubAddress.EncodeAddress(), nil
}
