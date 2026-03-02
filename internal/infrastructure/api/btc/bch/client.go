package bch

import (
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"

	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	apibtcimpl "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/btc/btc"
	pkgbch "github.com/hiromaily/go-crypto-wallet/pkg/chains/bch"
	"github.com/hiromaily/go-crypto-wallet/pkg/config"
)

// TODO: BitcoinCash specific func must be overridden by same func name to Bitcoin

// BitcoinCash embeds Bitcoin
type BitcoinCash struct {
	apibtcimpl.Bitcoin
}

// NewBitcoinCash bitcoin cash instance based on Bitcoin
func NewBitcoinCash(
	client *rpcclient.Client,
	conf *config.Bitcoin,
	coinTypeCode domainCoin.CoinTypeCode,
) (*BitcoinCash, error) {
	// bitcoin base
	bit, err := apibtcimpl.NewBitcoin(client, conf, coinTypeCode)
	if err != nil {
		return nil, fmt.Errorf("btc.NewBitcoin() error: %s", err)
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
		conf.Net = pkgbch.TestnetMagic
	case chaincfg.RegressionNetParams.Name:
		conf.Net = pkgbch.RegtestMagic
	default:
		// chaincfg.MainNetParams.Name
		conf.Net = pkgbch.MainnetMagic
	}
	b.SetChainConfNet(conf.Net)
}
