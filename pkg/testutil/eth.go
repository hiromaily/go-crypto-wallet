package testutil

import (
	"fmt"
	"os"

	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	domainCoin "github.com/hiromaily/go-crypto-wallet/pkg/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/pkg/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/ethereum"
	"github.com/hiromaily/go-crypto-wallet/pkg/uuid"
)

var et ethereum.Ethereumer

// GetETH returns eth instance
// FIXME: hard coded
func GetETH() (ethereum.Ethereumer, error) {
	if et != nil {
		return et, nil
	}

	projPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-crypto-wallet"
	confPath := projPath + "/data/config/eth_watch.toml"
	conf, err := config.NewWallet(confPath, wallet.WalletTypeWatchOnly, domainCoin.ETH)
	if err != nil {
		return nil, fmt.Errorf("fail to create config: %w", err)
	}
	// TODO: if config should be overridden, here
	conf.CoinTypeCode = domainCoin.ETH

	// uuid handler
	uuidHandler := uuid.NewGoogleUUIDHandler()
	// client
	client, err := ethereum.NewRPCClient(&conf.Ethereum)
	if err != nil {
		return nil, fmt.Errorf("fail to create ethereum rpc client: %w", err)
	}
	et, err = ethereum.NewEthereum(client, &conf.Ethereum, conf.CoinTypeCode, uuidHandler)
	if err != nil {
		return nil, fmt.Errorf("fail to create eth instance: %w", err)
	}
	return et, nil
}
