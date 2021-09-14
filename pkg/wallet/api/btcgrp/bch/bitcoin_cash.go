package bch

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/cpacia/bchutil"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/btcgrp/btc"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

// TODO: BitcoinCash specific func must be overridden by same func name to Bitcoin

// BitcoinCash embeds Bitcoin
type BitcoinCash struct {
	btc.Bitcoin
}

// NewBitcoinCash bitcoin cash instance based on Bitcoin
func NewBitcoinCash(
	client *rpcclient.Client,
	coinTypeCode coin.CoinTypeCode,
	conf *config.Bitcoin,
	logger *zap.Logger) (*BitcoinCash, error) {
	// bitcoin base
	bit, err := btc.NewBitcoin(client, coinTypeCode, conf, logger)
	if err != nil {
		return nil, errors.Errorf("btc.NewBitcoin() error: %s", err)
	}

	bitc := BitcoinCash{Bitcoin: *bit}
	bitc.initChainParams()

	return &bitc, nil
}

// initChainParams overrides chain parms as for bitcoin cash
func (b *BitcoinCash) initChainParams() {
	conf := b.GetChainConf()

	switch conf.Name {
	case chaincfg.TestNet3Params.Name:
		conf.Net = bchutil.TestnetMagic
	case chaincfg.RegressionNetParams.Name:
		conf.Net = bchutil.Regtestmagic
	default:
		// chaincfg.MainNetParams.Name
		conf.Net = bchutil.MainnetMagic
	}
	b.SetChainConfNet(conf.Net)
}
