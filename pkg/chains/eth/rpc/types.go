package rpc

import "math/big"

// QuantityTag is the default block parameter for Ethereum RPC calls.
// https://github.com/ethereum/wiki/wiki/JSON-RPC#the-default-block-parameter
type QuantityTag string

const (
	// QuantityTagLatest refers to the latest mined block.
	QuantityTagLatest QuantityTag = "latest"
	// QuantityTagPending refers to the pending state/transactions.
	QuantityTagPending QuantityTag = "pending"
)

// String returns the string representation of the QuantityTag.
func (q QuantityTag) String() string {
	return string(q)
}

// ResponseSyncing is the wire-format response of eth_syncing.
// When the node is not syncing, the RPC returns false (handled by the bool return value).
type ResponseSyncing struct {
	StartingBlock       int64 `json:"startingBlock" mapstructure:"startingBlock"`
	HighestBlock        int64 `json:"highestBlock" mapstructure:"highestBlock"`
	CurrentBlock        int64 `json:"currentBlock" mapstructure:"currentBlock"`
	SyncedAccounts      int64 `json:"syncedAccounts" mapstructure:"syncedAccounts"`
	SyncedAccountBytes  int64 `json:"syncedAccountBytes" mapstructure:"syncedAccountBytes"`
	SyncedBytecodes     int64 `json:"syncedBytecodes" mapstructure:"syncedBytecodes"`
	SyncedBytecodeBytes int64 `json:"syncedBytecodeBytes" mapstructure:"syncedBytecodeBytes"`
	SyncedStorage       int64 `json:"syncedStorage" mapstructure:"syncedStorage"`
	SyncedStorageBytes  int64 `json:"syncedStorageBytes" mapstructure:"syncedStorageBytes"`

	HealingBytecode     int64 `json:"healingBytecode" mapstructure:"healingBytecode"`
	HealedBytecodes     int64 `json:"healedBytecodes" mapstructure:"healedBytecodes"`
	HealedBytecodeBytes int64 `json:"healedBytecodeBytes" mapstructure:"healedBytecodeBytes"`

	HealingTrienodes    int64 `json:"healingTrienodes" mapstructure:"healingTrienodes"`
	HealedTrienodes     int64 `json:"healedTrienodes" mapstructure:"healedTrienodes"`
	HealedTrienodeBytes int64 `json:"healedTrienodeBytes" mapstructure:"healedTrienodeBytes"`
}

// BlockRawInfo holds the raw hex-encoded block fields as returned by eth_getBlockByNumber.
type BlockRawInfo struct {
	Number           string   `json:"number"`
	Hash             string   `json:"hash"`
	ParentHash       string   `json:"parentHash"`
	Nonce            string   `json:"nonce"`
	Sha3Uncles       string   `json:"sha3Uncles"`
	LogsBloom        string   `json:"logsBloom"`
	TransactionsRoot string   `json:"transactionsRoot"`
	StateRoot        string   `json:"stateRoot"`
	Miner            string   `json:"miner"`
	Difficulty       string   `json:"difficulty"`
	TotalDifficulty  string   `json:"totalDifficulty"`
	ExtraData        string   `json:"extraData"`
	Size             string   `json:"size"`
	GasLimit         string   `json:"gasLimit"`
	GasUsed          string   `json:"gasUsed"`
	Timestamp        string   `json:"timestamp"`
	Transactions     []string `json:"transactions"`
	Uncles           []string `json:"uncles"`
	BaseFeePerGas    string   `json:"baseFeePerGas,omitempty"` // EIP-1559 (London hard fork)
}

// BlockInfo holds decoded block fields with numeric types.
type BlockInfo struct {
	Number           *big.Int `json:"number"`
	Hash             string   `json:"hash"`
	ParentHash       string   `json:"parentHash"`
	Nonce            *big.Int `json:"nonce"`
	Sha3Uncles       string   `json:"sha3Uncles"`
	LogsBloom        string   `json:"logsBloom"`
	TransactionsRoot string   `json:"transactionsRoot"`
	StateRoot        string   `json:"stateRoot"`
	Miner            string   `json:"miner"`
	Difficulty       *big.Int `json:"difficulty"`
	TotalDifficulty  *big.Int `json:"totalDifficulty"`
	ExtraData        string   `json:"extraData"`
	Size             *big.Int `json:"size"`
	GasLimit         *big.Int `json:"gasLimit"`
	GasUsed          *big.Int `json:"gasUsed"`
	Timestamp        *big.Int `json:"timestamp"`
	Transactions     []string `json:"transactions"`
	Uncles           []string `json:"uncles"`
	BaseFeePerGas    *big.Int `json:"baseFeePerGas,omitempty"` // EIP-1559 (London hard fork)
}

// ResponseGetTransaction is the wire-format response of eth_getTransactionByHash.
type ResponseGetTransaction struct {
	BlockHash        string `json:"blockHash"`
	BlockNumber      int64  `json:"blockNumber"`
	From             string `json:"from"`
	Gas              int64  `json:"gas"`
	GasPrice         int64  `json:"gasPrice"`
	Hash             string `json:"hash"`
	Input            string `json:"input"`
	Nonce            int64  `json:"nonce"`
	To               string `json:"to"`
	TransactionIndex int64  `json:"transactionIndex"`
	Value            int64  `json:"value"`
	V                int64  `json:"v"`
	R                string `json:"r"`
	S                string `json:"s"`
}

// ResponseGetTransactionReceipt is the wire-format response of eth_getTransactionReceipt.
type ResponseGetTransactionReceipt struct {
	TransactionHash   string   `json:"transactionHash"`
	TransactionIndex  int64    `json:"transactionIndex"`
	BlockHash         string   `json:"blockHash"`
	BlockNumber       int64    `json:"blockNumber"`
	From              string   `json:"from"`
	To                string   `json:"to"`
	CumulativeGasUsed int64    `json:"cumulativeGasUsed"`
	GasUsed           int64    `json:"gasUsed"`
	ContractAddress   string   `json:"contractAddress"`
	Logs              []string `json:"logs"`
	LogsBloom         string   `json:"logsBloom"`
	Status            int64    `json:"status"`
}

// TransactionReceipt is the decoded transaction receipt with unsigned integer fields.
// Returned by GetTxReceipt; matches the domain TransactionReceipt structure.
type TransactionReceipt struct {
	TransactionHash   string
	TransactionIndex  uint64
	BlockHash         string
	BlockNumber       uint64
	From              string
	To                string
	CumulativeGasUsed uint64
	GasUsed           uint64
	ContractAddress   string
	Status            uint64
}
