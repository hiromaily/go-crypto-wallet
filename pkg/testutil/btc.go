package testutil

import (
	"fmt"
	"log"
	"os"

	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/btcgrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

var bc btcgrp.Bitcoiner

// GetBTC returns btc instance
//FIXME: hard coded
func GetBTC() btcgrp.Bitcoiner {
	if bc != nil {
		return bc
	}

	projPath := fmt.Sprintf("%s/src/github.com/hiromaily/go-crypto-wallet", os.Getenv("GOPATH"))
	confPath := fmt.Sprintf("%s/data/config/btc_watch.toml", projPath)
	conf, err := config.New(confPath, wallet.WalletTypeWatchOnly, coin.BTC)
	if err != nil {
		log.Fatalf("fail to create config: %v", err)
	}
	//TODO: if config should be overridden, here

	// logger
	logger := logger.NewZapLogger(&conf.Logger)
	// client
	client, err := btcgrp.NewRPCClient(&conf.Bitcoin)
	if err != nil {
		log.Fatalf("fail to create bitcoin core client: %v", err)
	}
	bc, err = btcgrp.NewBitcoin(client, &conf.Bitcoin, logger, conf.CoinTypeCode)
	if err != nil {
		log.Fatalf("fail to create btc instance: %v", err)
	}
	return bc
}
