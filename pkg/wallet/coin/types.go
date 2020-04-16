package coin

//----------------------------------------------------
// CoinType
//----------------------------------------------------

//CoinType Bitcoin種別(CayenneWalletで取引するcoinの種別)
type CoinType string

// coin_type
const (
	BTC CoinType = "btc"
	BCH CoinType = "bch"
	//ETH CoinType = "eth"
)

//CoinTypeValue coin_typeの値
var CoinTypeValue = map[CoinType]uint8{
	BTC: 1,
	BCH: 2,
	//ETH: 3,
}

func (c CoinType) String() string {
	return string(c)
}

// ValidateBitcoinType BitcoinTypeのバリデーションを行う
func ValidateBitcoinType(val string) bool {
	if _, ok := CoinTypeValue[CoinType(val)]; ok {
		return true
	}
	return false
}

//----------------------------------------------------
// BTCVersion
//----------------------------------------------------

//BTCVersion 実行環境
type BTCVersion int

// environment
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

//NetworkType ネットワーク種別
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
