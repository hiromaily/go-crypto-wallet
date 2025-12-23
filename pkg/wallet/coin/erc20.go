package coin

import (
	domainCoin "github.com/hiromaily/go-crypto-wallet/pkg/domain/coin"
)

// Deprecated: Use github.com/hiromaily/go-crypto-wallet/pkg/domain/coin instead.
// This package provides backward compatibility aliases.

// ERC20Token erc20 token
//
// Deprecated: Use domain/coin.ERC20Token
type ERC20Token = domainCoin.ERC20Token

// erc20_token
//
// Deprecated: Use constants from domain/coin package
const (
	HYT = domainCoin.HYT
	BAT = domainCoin.BAT
)

// ERC20Map map of ERC20 tokens
//
// Deprecated: Use domain/coin.ERC20Map
var ERC20Map = domainCoin.ERC20Map

// IsERC20Token validate
//
// Deprecated: Use domain/coin.IsERC20Token
func IsERC20Token(val string) bool {
	return domainCoin.IsERC20Token(val)
}
