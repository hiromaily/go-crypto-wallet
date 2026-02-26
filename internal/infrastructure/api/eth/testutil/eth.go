package testutil

import (
	"fmt"
	"os"

	"github.com/stretchr/testify/suite"

	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
	ethereum "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/eth"
	apiethimpl "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/eth/eth"
	cryptocurrency "github.com/hiromaily/go-crypto-wallet/pkg/chains"
	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	"github.com/hiromaily/go-crypto-wallet/pkg/uuid"
)

var et *apiethimpl.Ethereum

// GetETH returns eth instance
// FIXME: hard coded
func GetETH() (*apiethimpl.Ethereum, error) {
	if et != nil {
		return et, nil
	}

	projPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-crypto-wallet"
	confPath := projPath + "/config/wallet/eth/watch.yaml"
	conf, err := config.NewWallet(confPath, wallet.WalletTypeWatchOnly, domainCoin.ETH)
	if err != nil {
		return nil, fmt.Errorf("fail to create config: %w", err)
	}
	// TODO: if config should be overridden, here
	conf.CoinTypeCode = domainCoin.ETH

	// uuid handler
	uuidHandler := uuid.NewGoogleUUIDHandler()
	// client
	client, err := cryptocurrency.NewEthereumRPCClient(&conf.Ethereum)
	if err != nil {
		return nil, fmt.Errorf("fail to create ethereum rpc client: %w", err)
	}
	et, err = ethereum.NewEthereum(client, &conf.Ethereum, conf.CoinTypeCode, uuidHandler)
	if err != nil {
		return nil, fmt.Errorf("fail to create eth instance: %w", err)
	}
	return et, nil
}

// ETHTestSuite is a test suite for ETH
type ETHTestSuite struct {
	suite.Suite
	ETH *apiethimpl.Ethereum
}

func (ets *ETHTestSuite) SetupTest() {
	eth, err := GetETH()
	ets.NoError(err)
	ets.ETH = eth
}

func (ets *ETHTestSuite) TearDownTest() {
	ets.ETH.Close()
}
