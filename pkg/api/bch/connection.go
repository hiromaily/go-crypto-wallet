package bch

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/cpacia/bchutil"
	"github.com/hiromaily/go-bitcoin/pkg/enum"
)

// BitcoinCash includes Client to call Json-RPC
type BitcoinCash struct {
	client    *rpcclient.Client
	chainConf *chaincfg.Params
	//receipt           KeyInfo
	//payment           KeyInfo
	confirmationBlock int
	//feeRange          FeeAdjustmentRate
	version enum.BTCVersion //179900
}

// OverrideChainParamsByBCH chaincfgをBCH用に上書きする
func (b *BitcoinCash) OverrideChainParamsByBCH() {
	switch b.chainConf.Name {
	case chaincfg.TestNet3Params.Name:
		b.chainConf.Net = bchutil.TestnetMagic
	case chaincfg.RegressionNetParams.Name:
		b.chainConf.Net = bchutil.Regtestmagic
	default:
		//chaincfg.MainNetParams.Name
		b.chainConf.Net = bchutil.MainnetMagic
	}
}
