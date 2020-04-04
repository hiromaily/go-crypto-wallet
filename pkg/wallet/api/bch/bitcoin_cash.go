package bch

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/cpacia/bchutil"
	"github.com/pkg/errors"

	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/toml"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api/btc"
)

//TODO: BitcoinCash特有の機能は同一func名でOverrideすること

// BitcoinCash embeds Bitcoin
type BitcoinCash struct {
	btc.Bitcoin
}

// NewBitcoinCash BitcoinCashオブジェクトを返す
func NewBitcoinCash(client *rpcclient.Client, conf *toml.BitcoinConf) (*BitcoinCash, error) {
	//bitcoin base
	bit, err := btc.NewBitcoin(client, conf)
	if err != nil {
		return nil, errors.Errorf("btc.NewBitcoin() error: %s", err)
	}

	bitc := BitcoinCash{Bitcoin: *bit}
	//if conf.IsMain {
	//	bitc.SetChainConf(&chaincfg.MainNetParams)
	//} else {
	//	bitc.SetChainConf(&chaincfg.TestNet3Params)
	//}
	bitc.initChainParams()

	//Bitcoinのバージョンを入れておく
	//netInfo, err := bitc.GetNetworkInfo()
	//if err != nil {
	//	return nil, errors.Errorf("bit.GetNetworkInfo() error: %s", err)
	//}
	//bitc.SetVersion(netInfo.Version)
	//logger.Infof("bitcoin server version: %d", netInfo.Version)

	bitc.SetCoinType(enum.BCH)

	return &bitc, nil
}

// initChainParams bitcoin cash用に書き換える
func (b *BitcoinCash) initChainParams() {
	conf := b.GetChainConf()

	switch conf.Name {
	case chaincfg.TestNet3Params.Name:
		conf.Net = bchutil.TestnetMagic
	case chaincfg.RegressionNetParams.Name:
		conf.Net = bchutil.Regtestmagic
	default:
		//chaincfg.MainNetParams.Name
		conf.Net = bchutil.MainnetMagic
	}
	b.SetChainConfNet(conf.Net)
}
