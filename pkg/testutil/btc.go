package testutil

import (
	"fmt"
	"os"

	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	domainCoin "github.com/hiromaily/go-crypto-wallet/pkg/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/pkg/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/bitcoin"
)

var bc bitcoin.Bitcoiner

// GetBTC returns btc instance
// FIXME: hard coded config path
func GetBTC() (bitcoin.Bitcoiner, error) {
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
	client, err := bitcoin.NewRPCClient(&conf.Bitcoin)
	if err != nil {
		return nil, fmt.Errorf("fail to create bitcoin core client: %w", err)
	}
	bc, err = bitcoin.NewBitcoin(client, &conf.Bitcoin, conf.CoinTypeCode)
	if err != nil {
		return nil, fmt.Errorf("fail to create btc instance: %w", err)
	}
	return bc, nil
}
