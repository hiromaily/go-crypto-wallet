package bch

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/btcgrp/btc"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

// TODO: BitcoinCash specific func must be overridden by same func name to Bitcoin

// refer to original [github.com/cpacia/bchutil](https://github.com/cpacia/bchutil/blob/master/protocol.go)

const (
	// MainNet represents the main bitcoin network.
	MainnetMagic wire.BitcoinNet = 0xe8f3e1e3

	// Testnet represents the test network (version 3).
	TestnetMagic wire.BitcoinNet = 0xf4f3e5f4

	// Regtest represents the regression test network.
	Regtestmagic wire.BitcoinNet = 0xfabfb5da
)

// BitcoinCash embeds Bitcoin
type BitcoinCash struct {
	btc.Bitcoin
}

// NewBitcoinCash bitcoin cash instance based on Bitcoin
func NewBitcoinCash(
	client *rpcclient.Client,
	coinTypeCode coin.CoinTypeCode,
	conf *config.Bitcoin,
	logger *zap.Logger,
) (*BitcoinCash, error) {
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
		conf.Net = TestnetMagic
	case chaincfg.RegressionNetParams.Name:
		conf.Net = Regtestmagic
	default:
		// chaincfg.MainNetParams.Name
		conf.Net = MainnetMagic
	}
	b.SetChainConfNet(conf.Net)
}
