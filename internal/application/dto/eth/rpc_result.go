package eth

import "math/big"

// SyncingStatus holds the Ethereum node sync progress.
// Replaces ethrpc.ResponseSyncing in port interface signatures.
type SyncingStatus struct {
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

// BlockInfo holds decoded Ethereum block fields.
// Replaces ethrpc.BlockInfo in port interface signatures.
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
	BaseFeePerGas    *big.Int // EIP-1559 (London hard fork)
}

// TransactionInfo holds the response of eth_getTransactionByHash.
// Replaces ethrpc.ResponseGetTransaction in port interface signatures.
type TransactionInfo struct {
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

// RawTransactionReceipt holds the wire-format response of eth_getTransactionReceipt.
// Replaces ethrpc.ResponseGetTransactionReceipt in port interface signatures.
type RawTransactionReceipt struct {
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
