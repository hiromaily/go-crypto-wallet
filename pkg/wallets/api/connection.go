package api

import (
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/config"
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/wallets/api/bch"
	"github.com/hiromaily/go-bitcoin/pkg/wallets/api/btc"
)

// Connection is to local bitcoin core RPC server using HTTP POST mode
//func Connection(conf *toml.BitcoinConf) (*Bitcoin, error) {
func NewBitcoin(conf *config.Bitcoin, logger *zap.Logger, coinType enum.CoinType) (Bitcoiner, error) {
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

	//BTC
	if coinType == enum.BTC {
		// New
		bit, err := btc.NewBitcoin(client, conf, logger)
		if err != nil {
			return nil, errors.Errorf("btc.NewBitcoin() error: %s", err)
		}

		return bit, err
	} else if coinType == enum.BCH {
		//BCH
		bitc, err := bch.NewBitcoinCash(client, conf, logger)
		if err != nil {
			return nil, errors.Errorf("bitc.NewBitcoinCash() error: %s", err)
		}

		return bitc, err
	}

	return nil, errors.New("coinType is out of range. It should be set by `btc`,`bch`")
}
