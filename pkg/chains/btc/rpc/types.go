package rpc

//----------------------------------------------------
// NodeVersion
//----------------------------------------------------

// NodeVersion represents the Bitcoin Core node version number
type NodeVersion int

// Known node versions
const (
	NodeVer28 NodeVersion = 280000
	NodeVer29 NodeVersion = 290000
	NodeVer30 NodeVersion = 300000
)

// Int converts NodeVersion to int
func (v NodeVersion) Int() int {
	return int(v)
}

// RequiredNodeVersion is the minimum required Bitcoin Core version for this wallet
const RequiredNodeVersion = NodeVer28

//----------------------------------------------------
// NetworkType
//----------------------------------------------------

// NetworkType represents the Bitcoin network type
type NetworkType string

// Network type values
const (
	NetworkTypeMainNet  NetworkType = "mainnet"
	NetworkTypeTestNet3 NetworkType = "testnet3"
	NetworkTypeRegTest  NetworkType = "regtest"
	NetworkTypeSigNet   NetworkType = "signet"
)

// String returns the string representation of NetworkType
func (n NetworkType) String() string {
	return string(n)
}
