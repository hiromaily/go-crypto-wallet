package bch

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/cpacia/bchutil"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/config"
	"github.com/hiromaily/go-bitcoin/pkg/wallets/api/btc"
	ctype "github.com/hiromaily/go-bitcoin/pkg/wallets/api/types"
)

//TODO: BitcoinCash特有の機能は同一func名でOverrideすること

// BitcoinCash embeds Bitcoin
type BitcoinCash struct {
	btc.Bitcoin
}

// NewBitcoinCash BitcoinCashオブジェクトを返す
func NewBitcoinCash(client *rpcclient.Client, conf *config.Bitcoin, logger *zap.Logger) (*BitcoinCash, error) {
	//bitcoin base
	bit, err := btc.NewBitcoin(client, conf, logger)
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

	bitc.SetCoinType(ctype.BCH)

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
