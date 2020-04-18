package coin

import "github.com/btcsuite/btcd/chaincfg"

//----------------------------------------------------
// CoinType
//----------------------------------------------------

// CoinType creates a separate subtree for every cryptocoin
type CoinType uint32

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
	CoinTypeBitcoinCash          = 145 // Bitcoin Cash
)

//CoinTypeCode coin type
type CoinTypeCode string

// coin_type_code
const (
	BTC CoinTypeCode = "btc"
	BCH CoinTypeCode = "bch"
	LTC CoinTypeCode = "ltc"
	ETH CoinTypeCode = "eth"
	XRP CoinTypeCode = "xrp"
)

func (c CoinTypeCode) String() string {
	return string(c)
}

func (c CoinTypeCode) CoinType(conf *chaincfg.Params) CoinType {
	if conf.Name != NetworkTypeMainNet.String() {
		return CoinTypeTestnet
	}
	if coinType, ok := CoinTypeCodeValue[c]; ok {
		return coinType
	}
	// coinType could not found
	return CoinTypeTestnet
}

//CoinTypeCodeValue value
var CoinTypeCodeValue = map[CoinTypeCode]CoinType{
	BTC: CoinTypeBitcoin,
	BCH: CoinTypeBitcoinCash,
	LTC: CoinTypeLitecoin,
	ETH: CoinTypeEther,
	XRP: CoinTypeRipple,
}

// ValidateCoinTypeCode validate
func ValidateCoinTypeCode(val string) bool {
	if _, ok := CoinTypeCodeValue[CoinTypeCode(val)]; ok {
		return true
	}
	return false
}

//----------------------------------------------------
// BTCVersion
//----------------------------------------------------

//BTCVersion version
type BTCVersion int

// expected version
const (
	BTCVer17 BTCVersion = 170000
	BTCVer18 BTCVersion = 180000
	BTCVer19 BTCVersion = 190000
)

func (v BTCVersion) Int() int {
	return int(v)
}

const RequiredVersion = BTCVer19

//----------------------------------------------------
// NetworkType
//----------------------------------------------------

//NetworkType network type
type NetworkType string

// network type
const (
	NetworkTypeMainNet    NetworkType = "mainnet"
	NetworkTypeTestNet3   NetworkType = "testnet3"
	NetworkTypeRegTestNet NetworkType = "regtest"
)

func (n NetworkType) String() string {
	return string(n)
}
