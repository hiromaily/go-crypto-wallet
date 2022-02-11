package testutil

import (
	"fmt"
	"os"

	"github.com/pkg/errors"

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

	projPath := fmt.Sprintf("%s/src/github.com/hiromaily/go-crypto-wallet", os.Getenv("GOPATH"))
	confPath := fmt.Sprintf("%s/data/config/btc_watch.toml", projPath)
	conf, err := config.NewWallet(confPath, wallet.WalletTypeWatchOnly, coin.BTC)
	if err != nil {
		return nil, errors.Wrap(err, "fail to create config")
	}
	// TODO: if config should be overridden, here
	conf.CoinTypeCode = coin.BTC

	// logger
	logger := logger.NewZapLogger(&conf.Logger)
	// client
	client, err := btcgrp.NewRPCClient(&conf.Bitcoin)
	if err != nil {
		return nil, errors.Wrap(err, "fail to create bitcoin core client")
	}
	bc, err = btcgrp.NewBitcoin(client, &conf.Bitcoin, logger, conf.CoinTypeCode)
	if err != nil {
		return nil, errors.Wrap(err, "fail to create btc instance")
	}
	return bc, nil
}
