package app

import (
	"errors"

	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
)

// ValidateCoinType validates that the coin type code is supported.
func ValidateCoinType(coinTypeCode string) error {
	if !domainCoin.IsCoinTypeCode(coinTypeCode) && !domainCoin.IsERC20Token(coinTypeCode) {
		return errors.New("coin args is invalid. `btc`, `bch`, `eth`, `xrp`, `hyt` is allowed")
	}
	return nil
}

// ValidateCoinTypeForSign validates coin type for sign wallet (BTC/BCH only).
func ValidateCoinTypeForSign(coinTypeCode string) error {
	if !domainCoin.IsCoinTypeCode(coinTypeCode) {
		return errors.New("coin args is invalid. `btc`, `bch` is allowed")
	}
	// Sign wallet only supports BTC and BCH
	code := domainCoin.CoinTypeCode(coinTypeCode)
	if code != domainCoin.BTC && code != domainCoin.BCH {
		return errors.New("sign wallet only supports btc and bch")
	}
	return nil
}
