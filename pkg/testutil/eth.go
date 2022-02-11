package testutil

import (
	"fmt"
	"os"

	"github.com/pkg/errors"

	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
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

	projPath := fmt.Sprintf("%s/src/github.com/hiromaily/go-crypto-wallet", os.Getenv("GOPATH"))
	confPath := fmt.Sprintf("%s/data/config/eth_watch.toml", projPath)
	conf, err := config.NewWallet(confPath, wallet.WalletTypeWatchOnly, coin.ETH)
	if err != nil {
		return nil, errors.Wrap(err, "fail to create config")
	}
	// TODO: if config should be overridden, here
	conf.CoinTypeCode = coin.ETH

	// logger
	logger := logger.NewZapLogger(&conf.Logger)
	// client
	client, err := ethgrp.NewRPCClient(&conf.Ethereum)
	if err != nil {
		return nil, errors.Wrap(err, "fail to create ethereum rpc client")
	}
	et, err = ethgrp.NewEthereum(client, &conf.Ethereum, logger, conf.CoinTypeCode)
	if err != nil {
		return nil, errors.Wrap(err, "fail to create eth instance")
	}
	return et, nil
}
