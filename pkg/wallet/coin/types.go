package coin

import "github.com/btcsuite/btcd/chaincfg"

// CoinType creates a separate subtree for every cryptocoin
type CoinType uint32

// Uint32 is converter
func (c CoinType) Uint32() uint32 {
	return uint32(c)
}

// coin_type
// https://github.com/satoshilabs/slips/blob/master/slip-0044.md
const (
	CoinTypeBitcoin     CoinType = 0   // Bitcoin
	CoinTypeTestnet     CoinType = 1   // Testnet (all coins)
	CoinTypeLitecoin    CoinType = 2   // Litecoin
	CoinTypeEther       CoinType = 60  // Ether
	CoinTypeRipple      CoinType = 144 // Ripple
	CoinTypeBitcoinCash CoinType = 145 // Bitcoin Cash
	// ERC20
	CoinTypeERC20    CoinType = 9000 // TODO: temporary
	CoinTypeERC20HYT CoinType = 9001 // TODO: temporary
)

// CoinTypeCode coin type code
type CoinTypeCode string

// coin_type_code
const (
	BTC   CoinTypeCode = "btc"
	BCH   CoinTypeCode = "bch"
	LTC   CoinTypeCode = "ltc"
	ETH   CoinTypeCode = "eth"
	XRP   CoinTypeCode = "xrp"
	ERC20 CoinTypeCode = "erc20"
	HYC   CoinTypeCode = "hyc"
)

// String converter
func (c CoinTypeCode) String() string {
	return string(c)
}

// CoinType returns CoinType
func (c CoinTypeCode) CoinType(conf *chaincfg.Params) CoinType {
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
var CoinTypeCodeValue = map[CoinTypeCode]CoinType{
	BTC:   CoinTypeBitcoin,
	BCH:   CoinTypeBitcoinCash,
	LTC:   CoinTypeLitecoin,
	ETH:   CoinTypeEther,
	XRP:   CoinTypeRipple,
	ERC20: CoinTypeERC20,
	HYC:   CoinTypeERC20HYT,
}

// IsCoinTypeCode validate
func IsCoinTypeCode(val string) bool {
	if _, ok := CoinTypeCodeValue[CoinTypeCode(val)]; ok {
		return true
	}
	return false
}
