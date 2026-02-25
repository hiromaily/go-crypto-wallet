package testutil

import (
	"fmt"
	"os"

	"github.com/stretchr/testify/suite"

	apibtc "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/btc"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
	bitcoin "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/btc"
	cryptocurrency "github.com/hiromaily/go-crypto-wallet/pkg/chains"
	"github.com/hiromaily/go-crypto-wallet/pkg/config"
)

var bc apibtc.Bitcoiner

// GetBTC returns btc instance
// FIXME: hard coded config path
func GetBTC() (apibtc.Bitcoiner, error) {
	if bc != nil {
		return bc, nil
	}

	projPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-crypto-wallet"
	confPath := projPath + "/config/wallet/btc/watch.yaml"
	conf, err := config.NewWallet(confPath, wallet.WalletTypeWatchOnly, domainCoin.BTC)
	if err != nil {
		return nil, fmt.Errorf("fail to create config: %w", err)
	}
	// TODO: if config should be overridden, here
	conf.CoinTypeCode = domainCoin.BTC

	// client
	client, err := cryptocurrency.NewBitcoinRPCClient(&conf.Bitcoin)
	if err != nil {
		return nil, fmt.Errorf("fail to create bitcoin core client: %w", err)
	}
	// Bitcoin instance
	btcFull, err := bitcoin.NewBitcoin(client, &conf.Bitcoin, conf.CoinTypeCode)
	if err != nil {
		return nil, fmt.Errorf("fail to create btc instance: %w", err)
	}
	bc = btcFull
	return bc, nil
}

// BTCTestSuite is a test suite for BTC

type BTCTestSuite struct {
	suite.Suite
	BTC apibtc.Bitcoiner
}

func (bts *BTCTestSuite) SetupTest() {
	btc, err := GetBTC()
	bts.NoError(err)
	bts.BTC = btc
}

func (bts *BTCTestSuite) TearDownTest() {
	bts.BTC.Close()
}
