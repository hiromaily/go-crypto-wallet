package api

import (
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/hiromaily/go-bitcoin/pkg/api/bch"
	"github.com/hiromaily/go-bitcoin/pkg/api/btc"
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/toml"
	"github.com/pkg/errors"
)

// Connection is to local bitcoin core RPC server using HTTP POST mode
//func Connection(conf *toml.BitcoinConf) (*Bitcoin, error) {
func Connection(conf *toml.BitcoinConf, coinType enum.CoinType) (Bitcoiner, error) {
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
		bit, err := btc.NewBitcoin(client, conf)
		if err != nil {
			return nil, errors.Errorf("btc.NewBitcoin() error: %s", err)
		}

		return bit, err
	} else if coinType == enum.BCH {
		//BCH
		bitc, err := bch.NewBitcoinCash(client, conf)
		if err != nil {
			return nil, errors.Errorf("bitc.NewBitcoinCash() error: %s", err)
		}

		return bitc, err
	}

	return nil, errors.New("coinType is out of range. It should be set by `btc`,`bch`")
}
