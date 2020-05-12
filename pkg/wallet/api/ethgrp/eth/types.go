package eth

//----------------------------------------------------
// QuantityTag
//----------------------------------------------------

// QuantityTag quantity tag
type QuantityTag string

// quantity-tag
// https://github.com/ethereum/wiki/wiki/JSON-RPC#the-default-block-parameter
const (
	QuantityTagLatest   QuantityTag = "latest"   // for the earliest/genesis block
	QuantityTagEarliest QuantityTag = "earliest" // for the latest mined block
	QuantityTagPending  QuantityTag = "pending"  // for the pending state/transactions
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
	NetworkTypeETHMainNet    NetworkTypeETH = "mainnet"
	NetworkTypeETHGoerli     NetworkTypeETH = "goerli"
	NetworkTypeETHRegRinkeby NetworkTypeETH = "rinkeby"
	NetworkTypeETHRopsten    NetworkTypeETH = "ropsten"
)

// String converter
func (n NetworkTypeETH) String() string {
	return string(n)
}

//----------------------------------------------------
// ChainID
//----------------------------------------------------

// ChainID type of network ID // not net-version
// https://github.com/ethereum/EIPs/blob/master/EIPS/eip-155.md
type ChainID string

// chain-id
const (
	ChainIDMainNet       ChainID = "Ethereum mainnet"
	ChainIDMorden        ChainID = "Morden Expanse mainnet"
	ChainIDRopsten       ChainID = "Ropsten"
	ChainIDRinkeby       ChainID = "Rinkeby"
	ChainIDGoerli        ChainID = "Goerli"
	ChainIDKovan         ChainID = "Kovan"
	ChainIDPrivateChains ChainID = "Geth private chains"
)

// ChainIDMap chainID mapping
var ChainIDMap = map[uint16]ChainID{
	1:    ChainIDMainNet,
	2:    ChainIDMorden,
	3:    ChainIDRopsten,
	4:    ChainIDRinkeby,
	5:    ChainIDGoerli,
	42:   ChainIDKovan,
	1337: ChainIDPrivateChains,
}

// String converter
func (c ChainID) String() string {
	return string(c)
}

//----------------------------------------------------
// ClientVersion
//----------------------------------------------------

// ClientVersion Client(Node)のバージョンを返す
type ClientVersion string

// client-version
const (
	ClientVersionGeth   ClientVersion = "Geth"
	ClientVersionParity ClientVersion = "Parity-Ethereum"
)

// String converter
func (c ClientVersion) String() string {
	return string(c)
}

// GasLimit fixed GasLimit
const GasLimit uint64 = 21000
