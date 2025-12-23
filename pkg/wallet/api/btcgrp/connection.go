package btcgrp

import (
	"fmt"

	"github.com/btcsuite/btcd/rpcclient"

	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/btcgrp/bch"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/btcgrp/btc"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
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
	client *rpcclient.Client, conf *config.Bitcoin, coinTypeCode coin.CoinTypeCode,
) (Bitcoiner, error) {
	switch coinTypeCode {
	case coin.BTC:
		bit, err := btc.NewBitcoin(client, conf, coinTypeCode)
		if err != nil {
			return nil, fmt.Errorf("fail to call btc.NewBitcoin(): %w", err)
		}

		return bit, err
	case coin.BCH:
		// BCH
		bitc, err := bch.NewBitcoinCash(client, coinTypeCode, conf)
		if err != nil {
			return nil, fmt.Errorf("fail to call bch.NewBitcoinCash(): %w", err)
		}

		return bitc, err
	case coin.LTC, coin.ETH, coin.XRP, coin.ERC20, coin.HYC:
		return nil, fmt.Errorf("coinType %s is not defined", coinTypeCode.String())
	default:
		return nil, fmt.Errorf("coinType %s is not defined", coinTypeCode.String())
	}
}
