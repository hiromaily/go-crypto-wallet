package xrp

// PublicWSServer public ripple server
type PublicWSServer string

// public server
// https://xrpl.org/get-started-with-the-rippled-api.html
const (
	PublicWSServerMainnet1 PublicWSServer = "wss://s1.ripple.com:51234"
	PublicWSServerMainnet2 PublicWSServer = "wss://s2.ripple.com:51234"
	PublicWSServerTestnet  PublicWSServer = "wss://s.altnet.rippletest.net:51233"
	PublicWSServerDevnet   PublicWSServer = "wss://s.devnet.rippletest.net:51233"
)

// String converter
func (p PublicWSServer) String() string {
	return string(p)
}

// NetworkTypeXRP network type
type NetworkTypeXRP string

// network type
const (
	NetworkTypeXRPMainNet NetworkTypeXRP = "mainnet"
	NetworkTypeXRPTestNet NetworkTypeXRP = "testnet"
	NetworkTypeXRPDevNet  NetworkTypeXRP = "devnet"
)

// String converter
func (n NetworkTypeXRP) String() string {
	return string(n)
}

// GetPublicWSServer returns public server url from network type
func GetPublicWSServer(networkType string) PublicWSServer {
	switch NetworkTypeXRP(networkType) {
	case NetworkTypeXRPMainNet:
		return PublicWSServerMainnet1
	case NetworkTypeXRPTestNet:
		return PublicWSServerTestnet
	case NetworkTypeXRPDevNet:
		return PublicWSServerDevnet
	default:
	}
	return ""
}

// XRPKeyType key type
type XRPKeyType string

// network type
const (
	XRPKeyTypeSECP256 XRPKeyType = "secp256k1"
	XRPKeyTypeED25519 XRPKeyType = "ed25519"
)

// String converter
func (k XRPKeyType) String() string {
	return string(k)
}
