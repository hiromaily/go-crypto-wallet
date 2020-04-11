package api

import (
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/hiromaily/go-bitcoin/pkg/wallets/api/btc"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/config"
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/wallets/api/bch"
)

//NewRPCClient try to connect bitcoin core RPCserver to create client instance
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
		return nil, errors.Errorf("rpcclient.New() error: %s", err)
	}
	return client, err
}

// NewBitcoin creates bitcoin/bitcoin cash instance according to coinType
func NewBitcoin(client *rpcclient.Client, conf *config.Bitcoin, logger *zap.Logger, coinType enum.CoinType) (Bitcoiner, error) {
	switch coinType {
	case enum.BTC:
		bit, err := btc.NewBitcoin(client, conf, logger)
		if err != nil {
			return nil, errors.Errorf("btc.NewBitcoin() error: %s", err)
		}

		return bit, err
	case enum.BCH:
		//BCH
		bitc, err := bch.NewBitcoinCash(client, conf, logger)
		if err != nil {
			return nil, errors.Errorf("bitc.NewBitcoinCash() error: %s", err)
		}

		return bitc, err
	}
	return nil, errors.New("coinType is out of range. It should be set by `btc`,`bch`")
}
