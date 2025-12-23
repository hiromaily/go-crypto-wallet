package testutil

import (
	"fmt"
	"os"

	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	"github.com/hiromaily/go-crypto-wallet/pkg/uuid"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/ethgrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

var et ethgrp.Ethereumer

// GetETH returns eth instance
// FIXME: hard coded
func GetETH() (ethgrp.Ethereumer, error) {
	if et != nil {
		return et, nil
	}

	projPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-crypto-wallet"
	confPath := projPath + "/data/config/eth_watch.toml"
	conf, err := config.NewWallet(confPath, wallet.WalletTypeWatchOnly, coin.ETH)
	if err != nil {
		return nil, fmt.Errorf("fail to create config: %w", err)
	}
	// TODO: if config should be overridden, here
	conf.CoinTypeCode = coin.ETH

	// logger
	log := logger.NewSlogFromConfig(conf.Logger.Env, conf.Logger.Level, conf.Logger.Service)
	// uuid handler
	uuidHandler := uuid.NewGoogleUUIDHandler()
	// client
	client, err := ethgrp.NewRPCClient(&conf.Ethereum)
	if err != nil {
		return nil, fmt.Errorf("fail to create ethereum rpc client: %w", err)
	}
	et, err = ethgrp.NewEthereum(client, &conf.Ethereum, log, conf.CoinTypeCode, uuidHandler)
	if err != nil {
		return nil, fmt.Errorf("fail to create eth instance: %w", err)
	}
	return et, nil
}
