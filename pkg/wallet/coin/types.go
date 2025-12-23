package coin

import (
	"github.com/btcsuite/btcd/chaincfg"

	domainCoin "github.com/hiromaily/go-crypto-wallet/pkg/domain/coin"
)

// Deprecated: Use github.com/hiromaily/go-crypto-wallet/pkg/domain/coin instead.
// This package provides backward compatibility aliases.

// CoinType creates a separate subtree for every cryptocoin
//
// Deprecated: Use domain/coin.CoinType
type CoinType = domainCoin.CoinType

// coin_type
//
// Deprecated: Use constants from domain/coin package
const (
	CoinTypeBitcoin     = domainCoin.CoinTypeBitcoin
	CoinTypeTestnet     = domainCoin.CoinTypeTestnet
	CoinTypeLitecoin    = domainCoin.CoinTypeLitecoin
	CoinTypeEther       = domainCoin.CoinTypeEther
	CoinTypeRipple      = domainCoin.CoinTypeRipple
	CoinTypeBitcoinCash = domainCoin.CoinTypeBitcoinCash
	CoinTypeERC20       = domainCoin.CoinTypeERC20
	CoinTypeERC20HYT    = domainCoin.CoinTypeERC20HYT
)

// CoinTypeCode coin type code
//
// Deprecated: Use domain/coin.CoinTypeCode
type CoinTypeCode = domainCoin.CoinTypeCode

// coin_type_code
//
// Deprecated: Use constants from domain/coin package
const (
	BTC   = domainCoin.BTC
	BCH   = domainCoin.BCH
	LTC   = domainCoin.LTC
	ETH   = domainCoin.ETH
	XRP   = domainCoin.XRP
	ERC20 = domainCoin.ERC20
	HYT   = domainCoin.HYT
	HYC   = domainCoin.HYT // Deprecated: Use HYT instead
)

// GetCoinType returns CoinType based on network configuration
// This function has infrastructure dependency (chaincfg) and remains in this package
func GetCoinType(c CoinTypeCode, conf *chaincfg.Params) CoinType {
	if conf.Name != "mainnet" {
		return CoinTypeTestnet
	}
	if coinType, ok := CoinTypeCodeValue[c]; ok {
		return coinType
	}
	// coinType could not found
	return CoinTypeTestnet
}

// CoinTypeCodeValue value
//
// Deprecated: Use domain/coin.CoinTypeCodeValue
var CoinTypeCodeValue = domainCoin.CoinTypeCodeValue

// IsCoinTypeCode validate
//
// Deprecated: Use domain/coin.IsCoinTypeCode
func IsCoinTypeCode(val string) bool {
	return domainCoin.IsCoinTypeCode(val)
}

// IsBTCGroup validates bitcoin group
//
// Deprecated: Use domain/coin.IsBTCGroup
func IsBTCGroup(val CoinTypeCode) bool {
	return domainCoin.IsBTCGroup(val)
}

// IsETHGroup validates ethereum, ERC20 group
//
// Deprecated: Use domain/coin.IsETHGroup
func IsETHGroup(val CoinTypeCode) bool {
	return domainCoin.IsETHGroup(val)
}
