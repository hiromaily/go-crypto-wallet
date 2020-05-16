package testutil

import (
	"fmt"
	"log"
	"os"

	"github.com/hiromaily/go-bitcoin/pkg/config"
	"github.com/hiromaily/go-bitcoin/pkg/logger"
	"github.com/hiromaily/go-bitcoin/pkg/wallet"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api/ethgrp"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/coin"
)

var et ethgrp.Ethereumer

// GetETH returns eth instance
//FIXME: hard coded
func GetETH() ethgrp.Ethereumer {
	if et != nil {
		return et
	}

	projPath := fmt.Sprintf("%s/src/github.com/hiromaily/go-bitcoin", os.Getenv("GOPATH"))
	confPath := fmt.Sprintf("%s/data/config/eth_watch.toml", projPath)
	conf, err := config.New(confPath, wallet.WalletTypeWatchOnly, coin.ETH)
	if err != nil {
		log.Fatalf("fail to create config: %v", err)
	}
	//TODO: if config should be overridden, here

	// logger
	logger := logger.NewZapLogger(&conf.Logger)
	// client
	client, err := ethgrp.NewRPCClient(&conf.Ethereum)
	if err != nil {
		log.Fatalf("fail to create ethereum rpc client: %v", err)
	}
	et, err = ethgrp.NewEthereum(client, &conf.Ethereum, logger, conf.CoinTypeCode)
	if err != nil {
		log.Fatalf("fail to create eth instance: %v", err)
	}
	return et
}
