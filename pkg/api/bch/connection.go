package bch

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/cpacia/bchutil"
	"github.com/hiromaily/go-bitcoin/pkg/api/btc"
	"github.com/hiromaily/go-bitcoin/pkg/toml"
)

//TODO:bitcoin cashもBitcoinオブジェクトをそのまま使用するかも
//その場合、confの上書きを行う

// BitcoinCash includes Client to call Json-RPC
type BitcoinCash struct {
	BTC btc.Bitcoin
}

// Connection is to local bitcoin core RPC server using HTTP POST mode
//func Connection(host, user, pass string, postMode, tls, isMain bool) (*Bitcoin, error) {
func Connection(conf *toml.BitcoinConf) (*BitcoinCash, error) {
	return nil, nil
}

// OverrideChainParamsByBCH chaincfgをBCH用に上書きする
func (b *BitcoinCash) OverrideChainParamsByBCH() {
	conf := b.BTC.GetChainConf()

	switch conf.Name {
	case chaincfg.TestNet3Params.Name:
		conf.Net = bchutil.TestnetMagic
	case chaincfg.RegressionNetParams.Name:
		conf.Net = bchutil.Regtestmagic
	default:
		//chaincfg.MainNetParams.Name
		conf.Net = bchutil.MainnetMagic
	}
	b.BTC.SetChainConfNet(conf.Net)
}
