// Package ethereum defines domain types for Ethereum blockchain operations.
//
// These types are used in application layer interfaces following Clean Architecture
// principles. Infrastructure implementations convert between these domain types
// and their infrastructure-specific representations.
package ethereum

import "math/big"

// UserAmount represents a user address and their balance amount.
type UserAmount struct {
	Address string
	Amount  uint64
}

// QuantityTag represents Ethereum block parameter tags.
// https://github.com/ethereum/wiki/wiki/JSON-RPC#the-default-block-parameter
type QuantityTag string

const (
	// QuantityTagLatest refers to the latest mined block
	QuantityTagLatest QuantityTag = "latest"
	// QuantityTagPending refers to the pending state/transactions
	QuantityTagPending QuantityTag = "pending"
)

// String returns the string representation of QuantityTag.
func (q QuantityTag) String() string {
	return string(q)
}

// BlockInfo represents information about an Ethereum block.
type BlockInfo struct {
	Number           *big.Int
	Hash             string
	ParentHash       string
	Nonce            *big.Int
	Sha3Uncles       string
	LogsBloom        string
	TransactionsRoot string
	StateRoot        string
	Miner            string
	Difficulty       *big.Int
	TotalDifficulty  *big.Int
	ExtraData        string
	Size             *big.Int
	GasLimit         *big.Int
	GasUsed          *big.Int
	Timestamp        *big.Int
	Transactions     []string
	Uncles           []string
	// BaseFeePerGas is the base fee per gas for EIP-1559 transactions (London hard fork and later).
	// This field will be nil for blocks before the London hard fork.
	BaseFeePerGas *big.Int
}

// ResponseGetTransaction represents the response from eth_getTransactionByHash.
type ResponseGetTransaction struct {
	BlockHash        string
	BlockNumber      int64
	From             string
	Gas              int64
	GasPrice         int64
	Hash             string
	Input            string
	Nonce            int64
	To               string
	TransactionIndex int64
	Value            int64
	V                int64
	R                string
	S                string
}

// ResponseGetTransactionReceipt represents the response from eth_getTransactionReceipt.
type ResponseGetTransactionReceipt struct {
	TransactionHash   string
	TransactionIndex  int64
	BlockHash         string
	BlockNumber       int64
	From              string
	To                string
	CumulativeGasUsed int64
	GasUsed           int64
	ContractAddress   string
	Logs              []string
	LogsBloom         string
	Status            int64
}

// ResponseSyncing represents the response from eth_syncing.
type ResponseSyncing struct {
	StartingBlock       int64
	HighestBlock        int64
	CurrentBlock        int64
	SyncedAccounts      int64
	SyncedAccountBytes  int64
	SyncedBytecodes     int64
	SyncedBytecodeBytes int64
	SyncedStorage       int64
	SyncedStorageBytes  int64
	HealingBytecode     int64
	HealedBytecodes     int64
	HealedBytecodeBytes int64
	HealingTrienodes    int64
	HealedTrienodes     int64
	HealedTrienodeBytes int64
}

// RawTx represents a raw Ethereum transaction before being submitted to the network.
type RawTx struct {
	UUID  string
	From  string
	To    string
	Value big.Int
	Nonce uint64
	TxHex string
	Hash  string
}
