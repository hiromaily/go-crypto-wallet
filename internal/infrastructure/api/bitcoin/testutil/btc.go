package testutil

import (
	"fmt"
	"os"

	"github.com/stretchr/testify/suite"

	portsBtc "github.com/hiromaily/go-crypto-wallet/internal/application/ports/btc"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/bitcoin"
	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	"github.com/hiromaily/go-crypto-wallet/pkg/cryptocurrency"
)

var bc portsBtc.Bitcoiner

// GetBTC returns btc instance
// FIXME: hard coded config path
func GetBTC() (portsBtc.Bitcoiner, error) {
	if bc != nil {
		return bc, nil
	}

	projPath := os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-crypto-wallet"
	confPath := projPath + "/data/config/btc_watch.toml"
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
	bc, err = bitcoin.NewBitcoin(client, &conf.Bitcoin, conf.CoinTypeCode)
	if err != nil {
		return nil, fmt.Errorf("fail to create btc instance: %w", err)
	}
	return bc, nil
}

// BTCTestSuite is a test suite for BTC

type BTCTestSuite struct {
	suite.Suite
	BTC portsBtc.Bitcoiner
}

func (bts *BTCTestSuite) SetupTest() {
	btc, err := GetBTC()
	bts.NoError(err)
	bts.BTC = btc
}

func (bts *BTCTestSuite) TearDownTest() {
	bts.BTC.Close()
}
