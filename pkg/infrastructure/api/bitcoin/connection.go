package bitcoin

import (
	"fmt"

	"github.com/btcsuite/btcd/rpcclient"

	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	domainCoin "github.com/hiromaily/go-crypto-wallet/pkg/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/bitcoin/bch"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/bitcoin/btc"
)

// NewRPCClient try to connect bitcoin core RPCserver to create client instance
// using HTTP POST mode
func NewRPCClient(conf *config.Bitcoin) (*rpcclient.Client, error) {
	connCfg := &rpcclient.ConnConfig{
		Host:         conf.Host,
		User:         conf.User,
		Pass:         conf.Pass,
		HTTPPostMode: conf.PostMode,   // Bitcoin core only supports HTTP POST mode
		DisableTLS:   conf.DisableTLS, // Bitcoin core does not provide TLS by default
	}

	// Notice the notification parameter is nil since notifications are
	// not supported in HTTP POST mode.
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		return nil, fmt.Errorf("rpcclient.New() error: %s", err)
	}
	return client, err
}

// NewBitcoin creates bitcoin/bitcoin cash instance according to coinType
func NewBitcoin(
	client *rpcclient.Client, conf *config.Bitcoin, coinTypeCode domainCoin.CoinTypeCode,
) (Bitcoiner, error) {
	switch coinTypeCode {
	case domainCoin.BTC:
		bit, err := btc.NewBitcoin(client, conf, coinTypeCode)
		if err != nil {
			return nil, fmt.Errorf("fail to call btc.NewBitcoin(): %w", err)
		}

		return bit, err
	case domainCoin.BCH:
		// BCH
		bitc, err := bch.NewBitcoinCash(client, coinTypeCode, conf)
		if err != nil {
			return nil, fmt.Errorf("fail to call bch.NewBitcoinCash(): %w", err)
		}

		return bitc, err
	case domainCoin.LTC, domainCoin.ETH, domainCoin.XRP, domainCoin.ERC20, domainCoin.HYT:
		return nil, fmt.Errorf("coinType %s is not defined", coinTypeCode.String())
	default:
		return nil, fmt.Errorf("coinType %s is not defined", coinTypeCode.String())
	}
}
