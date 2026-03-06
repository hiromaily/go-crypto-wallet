package eth

import (
	"context"
	"math/big"

	dtoeth "github.com/hiromaily/go-crypto-wallet/internal/application/dto/eth"
	domainETH "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/eth"
	ethrpc "github.com/hiromaily/go-crypto-wallet/pkg/chains/eth/rpc"
)

// Syncing returns sync status or bool
//   - return false if not syncing (it means syncing is done)
//   - there seems 2 different responses
func (e *Ethereum) Syncing(ctx context.Context) (*dtoeth.SyncingStatus, bool, error) {
	result, isSyncing, err := e.pkgrpc.Syncing(ctx)
	if err != nil {
		return nil, false, err
	}
	if result == nil {
		return nil, isSyncing, nil
	}
	return &dtoeth.SyncingStatus{
		StartingBlock:       result.StartingBlock,
		HighestBlock:        result.HighestBlock,
		CurrentBlock:        result.CurrentBlock,
		SyncedAccounts:      result.SyncedAccounts,
		SyncedAccountBytes:  result.SyncedAccountBytes,
		SyncedBytecodes:     result.SyncedBytecodes,
		SyncedBytecodeBytes: result.SyncedBytecodeBytes,
		SyncedStorage:       result.SyncedStorage,
		SyncedStorageBytes:  result.SyncedStorageBytes,
		HealingBytecode:     result.HealingBytecode,
		HealedBytecodes:     result.HealedBytecodes,
		HealedBytecodeBytes: result.HealedBytecodeBytes,
		HealingTrienodes:    result.HealingTrienodes,
		HealedTrienodes:     result.HealedTrienodes,
		HealedTrienodeBytes: result.HealedTrienodeBytes,
	}, isSyncing, nil
}

// ProtocolVersion returns the current ethereum protocol version
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_protocolversion
func (e *Ethereum) ProtocolVersion(ctx context.Context) (uint64, error) {
	return e.pkgrpc.ProtocolVersion(ctx)
}

// Coinbase returns the client coinbase address
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_coinbase
func (e *Ethereum) Coinbase(ctx context.Context) (string, error) {
	return e.pkgrpc.Coinbase(ctx)
}

// Accounts returns a list of addresses owned by client
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_accounts
func (e *Ethereum) Accounts(ctx context.Context) ([]string, error) {
	return e.pkgrpc.Accounts(ctx)
}

// BlockNumber returns the number of most recent block
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_blocknumber
func (e *Ethereum) BlockNumber(ctx context.Context) (*big.Int, error) {
	return e.pkgrpc.BlockNumber(ctx)
}

// EnsureBlockNumber calls BlockNumber() several times
func (e *Ethereum) EnsureBlockNumber(ctx context.Context, loopCount int) (*big.Int, error) {
	return e.pkgrpc.EnsureBlockNumber(ctx, loopCount)
}

// GetBalance returns the balance of the account of given address
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getbalance
func (e *Ethereum) GetBalance(
	ctx context.Context, hexAddr string, quantityTag domainETH.QuantityTag,
) (*big.Int, error) {
	return e.pkgrpc.GetBalance(ctx, hexAddr, ethrpc.QuantityTag(quantityTag))
}

// GetTransactionCount returns the number of transactions sent from an address
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_gettransactioncount
func (e *Ethereum) GetTransactionCount(
	ctx context.Context, hexAddr string, quantityTag domainETH.QuantityTag,
) (*big.Int, error) {
	return e.pkgrpc.GetTransactionCount(ctx, hexAddr, ethrpc.QuantityTag(quantityTag))
}

// GetBlockTransactionCountByNumber returns the number of transactions in a block matching the given block number
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getblocktransactioncountbynumber
func (e *Ethereum) GetBlockTransactionCountByNumber(ctx context.Context, blockNumber uint64) (*big.Int, error) {
	return e.pkgrpc.GetBlockTransactionCountByNumber(ctx, blockNumber)
}

// GetUncleCountByBlockNumber returns the number of uncles in a block from a block matching the given block number
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getunclecountbyblocknumber
func (e *Ethereum) GetUncleCountByBlockNumber(ctx context.Context, blockNumber uint64) (*big.Int, error) {
	return e.pkgrpc.GetUncleCountByBlockNumber(ctx, blockNumber)
}

// GetBlockByNumber returns information about a block by block number
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getblockbynumber
func (e *Ethereum) GetBlockByNumber(ctx context.Context, blockNumber uint64) (*dtoeth.BlockInfo, error) {
	result, err := e.pkgrpc.GetBlockByNumber(ctx, blockNumber)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return &dtoeth.BlockInfo{
		Number:           result.Number,
		Hash:             result.Hash,
		ParentHash:       result.ParentHash,
		Nonce:            result.Nonce,
		Sha3Uncles:       result.Sha3Uncles,
		LogsBloom:        result.LogsBloom,
		TransactionsRoot: result.TransactionsRoot,
		StateRoot:        result.StateRoot,
		Miner:            result.Miner,
		Difficulty:       result.Difficulty,
		TotalDifficulty:  result.TotalDifficulty,
		ExtraData:        result.ExtraData,
		Size:             result.Size,
		GasLimit:         result.GasLimit,
		GasUsed:          result.GasUsed,
		Timestamp:        result.Timestamp,
		Transactions:     result.Transactions,
		Uncles:           result.Uncles,
		BaseFeePerGas:    result.BaseFeePerGas,
	}, nil
}
