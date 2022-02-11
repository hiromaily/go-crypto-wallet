package btc

//----------------------------------------------------
// BTCVersion
//----------------------------------------------------

// BTCVersion version
type BTCVersion int

// expected version
const (
	BTCVer17 BTCVersion = 170000
	BTCVer18 BTCVersion = 180000
	BTCVer19 BTCVersion = 190000
	BTCVer20 BTCVersion = 200000
	BTCVer21 BTCVersion = 210000 // for BCH version
)

// Int converter
func (v BTCVersion) Int() int {
	return int(v)
}

// RequiredVersion is required version for this wallet apps
const RequiredVersion = BTCVer19

//----------------------------------------------------
// NetworkTypeBTC
//----------------------------------------------------

// NetworkTypeBTC network type
type NetworkTypeBTC string

// network type
const (
	NetworkTypeMainNet    NetworkTypeBTC = "mainnet"
	NetworkTypeTestNet3   NetworkTypeBTC = "testnet3"
	NetworkTypeRegTestNet NetworkTypeBTC = "regtest"
	NetworkTypeSigNet     NetworkTypeBTC = "signet"
)

// String converter
func (n NetworkTypeBTC) String() string {
	return string(n)
}
