package eth

import "strings"

// NetworkTypeETH identifies the Ethereum network.
type NetworkTypeETH string

const (
	NetworkTypeETHMainNet NetworkTypeETH = "mainnet"
	NetworkTypeETHSepolia NetworkTypeETH = "sepolia"
	NetworkTypeETHHolesky NetworkTypeETH = "holesky"
	NetworkTypeETHAnvil   NetworkTypeETH = "anvil"
	NetworkTypeETHLocal   NetworkTypeETH = "local"
)

// String converter
func (n NetworkTypeETH) String() string {
	return string(n)
}

// ChainIDForNetwork returns the numeric EIP-155 chain ID for a given network type.
// Returns 0 for unknown network types.
func ChainIDForNetwork(network NetworkTypeETH) uint64 {
	switch network {
	case NetworkTypeETHMainNet:
		return 1
	case NetworkTypeETHSepolia:
		return 11155111
	case NetworkTypeETHHolesky:
		return 17000
	case NetworkTypeETHAnvil, NetworkTypeETHLocal:
		return 1337
	default:
		return 0
	}
}

// ClientVersion identifies the Ethereum client implementation.
type ClientVersion string

const (
	ClientVersionGeth  ClientVersion = "Geth"
	ClientVersionAnvil ClientVersion = "Anvil"
)

// String converter
func (c ClientVersion) String() string {
	return string(c)
}

// DetectClientType detects the Ethereum client type from a version string.
func DetectClientType(version string) ClientVersion {
	versionLower := strings.ToLower(version)
	if strings.Contains(versionLower, "anvil") {
		return ClientVersionAnvil
	}
	// Default to Geth (most common)
	return ClientVersionGeth
}

// GasLimit is the fixed gas limit for standard ETH transfers (EIP-21000).
const GasLimit uint64 = 21000

// ETHNodeType identifies the Ethereum node implementation.
type ETHNodeType string

const (
	// ETHNodeTypeAnvil represents a Foundry Anvil local development node.
	ETHNodeTypeAnvil ETHNodeType = "anvil"
	// ETHNodeTypeGeth represents a go-ethereum (Geth) node.
	ETHNodeTypeGeth ETHNodeType = "geth"
)

// String converter
func (n ETHNodeType) String() string {
	return string(n)
}
