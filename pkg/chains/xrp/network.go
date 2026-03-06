package xrp

// PublicWSServer public ripple server
type PublicWSServer string

// public server
// https://xrpl.org/get-started-with-the-rippled-api.html
// downside to use public server is Admin method can not be used
// https://xrpl.org/admin-rippled-methods.html
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
	NetworkTypeXRPMainNet       NetworkTypeXRP = "mainnet"
	NetworkTypeXRPTestNet       NetworkTypeXRP = "testnet"
	NetworkTypeXRPDevNet        NetworkTypeXRP = "devnet"
	NetworkTypeXRPStandaloneNet NetworkTypeXRP = "standalone"
)

// String converter
func (n NetworkTypeXRP) String() string {
	return string(n)
}

// GetPublicWSServer returns public server url from network type.
// Returns "" for standalone mode — the caller must provide a local URL via config.
func GetPublicWSServer(networkType string) PublicWSServer {
	switch NetworkTypeXRP(networkType) {
	case NetworkTypeXRPMainNet:
		return PublicWSServerMainnet1
	case NetworkTypeXRPTestNet:
		return PublicWSServerTestnet
	case NetworkTypeXRPDevNet:
		return PublicWSServerDevnet
	case NetworkTypeXRPStandaloneNet:
		return "" // local rippled node; URL must be set via websocket_public_url in config
	default:
	}
	return ""
}
