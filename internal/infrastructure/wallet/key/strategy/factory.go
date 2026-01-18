package strategy

import (
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"

	portsWallet "github.com/hiromaily/go-crypto-wallet/internal/application/ports/wallet"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
)

// CreateCoinKeyStrategy creates the appropriate CoinKeyStrategy based on coin type
func CreateCoinKeyStrategy(
	coinTypeCode domainCoin.CoinTypeCode,
	conf *chaincfg.Params,
) (portsWallet.CoinKeyStrategy, error) {
	switch coinTypeCode {
	case domainCoin.BTC:
		return NewBTCKeyStrategy(conf), nil
	case domainCoin.BCH:
		return NewBCHKeyStrategy(conf), nil
	case domainCoin.ETH, domainCoin.ERC20:
		return NewETHKeyStrategy(), nil
	case domainCoin.XRP:
		return NewXRPKeyStrategy(), nil
	case domainCoin.LTC, domainCoin.HYT:
		return nil, fmt.Errorf("coin type %s is not implemented yet", coinTypeCode.String())
	default:
		return nil, fmt.Errorf("unsupported coin type: %s", coinTypeCode.String())
	}
}
