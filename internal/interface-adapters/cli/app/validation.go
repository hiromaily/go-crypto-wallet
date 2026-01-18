package app

import (
	"errors"

	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
)

// ValidateCoinType validates that the coin type code is supported.
func ValidateCoinType(coinTypeCode string) error {
	if coinTypeCode == "" {
		return errors.New("--coin flag is required. Valid values: btc, bch, eth, xrp, hyt")
	}
	if !domainCoin.IsCoinTypeCode(coinTypeCode) && !domainCoin.IsERC20Token(coinTypeCode) {
		return errors.New("coin args is invalid. `btc`, `bch`, `eth`, `xrp`, `hyt` is allowed")
	}
	return nil
}

// ValidateCoinTypeForSign validates coin type for sign wallet (BTC/BCH only).
func ValidateCoinTypeForSign(coinTypeCode string) error {
	if coinTypeCode == "" {
		return errors.New("--coin flag is required. Valid values for sign wallet: btc, bch")
	}
	code := domainCoin.CoinTypeCode(coinTypeCode)
	if code != domainCoin.BTC && code != domainCoin.BCH {
		return errors.New("coin args is invalid for sign wallet, only `btc` and `bch` are allowed")
	}
	return nil
}
