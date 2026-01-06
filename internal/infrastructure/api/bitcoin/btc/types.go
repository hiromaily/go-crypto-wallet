package btc

//----------------------------------------------------
// BTCVersion
//----------------------------------------------------

// BTCVersion version
type BTCVersion int

// expected version
const (
	BTCVer28 BTCVersion = 280000
	BTCVer29 BTCVersion = 290000
	BTCVer30 BTCVersion = 300000
)

// Int converter
func (v BTCVersion) Int() int {
	return int(v)
}

// RequiredVersion is required version for this wallet apps
const RequiredVersion = BTCVer28

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
