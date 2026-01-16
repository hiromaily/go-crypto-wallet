package bitcoin

import (
	"fmt"

	"github.com/btcsuite/btcd/rpcclient"

	portsBitcoin "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/bitcoin"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/bitcoin/bch"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/bitcoin/btc"
	"github.com/hiromaily/go-crypto-wallet/pkg/config"
)

// NewBitcoin creates bitcoin/bitcoin cash instance according to coinType
func NewBitcoin(
	client *rpcclient.Client, conf *config.Bitcoin, coinTypeCode domainCoin.CoinTypeCode,
) (portsBitcoin.Bitcoiner, error) {
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
