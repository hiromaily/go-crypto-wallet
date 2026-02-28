package bitcoin

import (
	"fmt"

	"github.com/btcsuite/btcd/rpcclient"

	apibtc "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/btc"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	apibchimpl "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/btc/bch"
	apibtcimpl "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/btc/btc"
	"github.com/hiromaily/go-crypto-wallet/pkg/config"
)

// NewBitcoin creates bitcoin/bitcoin cash instance according to coinType
func NewBitcoin(
	client *rpcclient.Client, conf *config.Bitcoin, coinTypeCode domainCoin.CoinTypeCode,
) (apibtc.Bitcoiner, error) {
	switch coinTypeCode {
	case domainCoin.BTC:
		bit, err := apibtcimpl.NewBitcoin(client, conf, coinTypeCode)
		if err != nil {
			return nil, fmt.Errorf("fail to call btc.NewBitcoin(): %w", err)
		}

		return bit, err
	case domainCoin.BCH:
		// BCH
		bitc, err := apibchimpl.NewBitcoinCash(client, conf, coinTypeCode)
		if err != nil {
			return nil, fmt.Errorf("fail to call bch.NewBitcoinCash(): %w", err)
		}

		return bitc, err
	case domainCoin.LTC, domainCoin.ETH, domainCoin.XRP, domainCoin.HYT:
		return nil, fmt.Errorf("coinType %s is not defined", coinTypeCode.String())
	default:
		return nil, fmt.Errorf("coinType %s is not defined", coinTypeCode.String())
	}
}
