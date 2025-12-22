package testutil

import (
	"fmt"
	"os"

	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/btcgrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

var bc btcgrp.Bitcoiner

// GetBTC returns btc instance
// FIXME: hard coded config path
func GetBTC() (btcgrp.Bitcoiner, error) {
	if bc != nil {
		return bc, nil
	}

	projPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-crypto-wallet"
	confPath := projPath + "/data/config/btc_watch.toml"
	conf, err := config.NewWallet(confPath, wallet.WalletTypeWatchOnly, coin.BTC)
	if err != nil {
		return nil, fmt.Errorf("fail to create config: %w", err)
	}
	// TODO: if config should be overridden, here
	conf.CoinTypeCode = coin.BTC

	// logger
	log := logger.NewSlogFromConfig(conf.Logger.Env, conf.Logger.Level, conf.Logger.Service)
	// client
	client, err := btcgrp.NewRPCClient(&conf.Bitcoin)
	if err != nil {
		return nil, fmt.Errorf("fail to create bitcoin core client: %w", err)
	}
	bc, err = btcgrp.NewBitcoin(client, &conf.Bitcoin, log, conf.CoinTypeCode)
	if err != nil {
		return nil, fmt.Errorf("fail to create btc instance: %w", err)
	}
	return bc, nil
}
