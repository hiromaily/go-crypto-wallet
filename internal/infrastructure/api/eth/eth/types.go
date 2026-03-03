package eth

import "strings"

//----------------------------------------------------
// QuantityTag
//----------------------------------------------------

// QuantityTag quantity tag
type QuantityTag string

// quantity-tag
// https://github.com/ethereum/wiki/wiki/JSON-RPC#the-default-block-parameter
const (
	QuantityTagLatest  QuantityTag = "latest"  // for the latest mined block
	QuantityTagPending QuantityTag = "pending" // for the pending state/transactions
	// QuantityTagEarliest QuantityTag = "earliest" // for the earliest/genesis block
)

// String converter
func (q QuantityTag) String() string {
	return string(q)
}

//----------------------------------------------------
// NetworkTypeETH
//----------------------------------------------------

// NetworkTypeETH network type
type NetworkTypeETH string

// network type
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

//----------------------------------------------------
// ClientVersion
//----------------------------------------------------

// ClientVersion returns client version
type ClientVersion string

// client-version
const (
	ClientVersionGeth  ClientVersion = "Geth"
	ClientVersionAnvil ClientVersion = "Anvil"
)

// String converter
func (c ClientVersion) String() string {
	return string(c)
}

// DetectClientType detects the Ethereum client type from version string
func DetectClientType(version string) ClientVersion {
	versionLower := strings.ToLower(version)
	if strings.Contains(versionLower, "anvil") {
		return ClientVersionAnvil
	}
	// Default to Geth (most common)
	return ClientVersionGeth
}

// GasLimit fixed GasLimit
const GasLimit uint64 = 21000

//----------------------------------------------------
// EthNodeType
//----------------------------------------------------

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
