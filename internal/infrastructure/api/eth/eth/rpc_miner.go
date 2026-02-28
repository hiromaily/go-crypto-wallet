package eth

import (
	"context"
	"math/big"

	ethrpc "github.com/hiromaily/go-crypto-wallet/pkg/chains/eth/rpc"
)

// StartMining starts the CPU mining process with the given number of threads and generate a new DAG if need be
// https://github.com/ethereum/go-ethereum/wiki/Management-APIs#miner_start
func (e *Ethereum) StartMining(ctx context.Context) error {
	return ethrpc.StartMining(ctx, e.rpcClient)
}

// StopMining stops the CPU mining operation
// https://github.com/ethereum/go-ethereum/wiki/Management-APIs#miner_stop
func (e *Ethereum) StopMining(ctx context.Context) error {
	return ethrpc.StopMining(ctx, e.rpcClient)
}

// Mining returns true if client is actively mining new blocks
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_mining
func (e *Ethereum) Mining(ctx context.Context) (bool, error) {
	return ethrpc.Mining(ctx, e.rpcClient)
}

// HashRate returns the number of hashes per second that the node is mining with
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_hashrate
func (e *Ethereum) HashRate(ctx context.Context) (*big.Int, error) {
	return ethrpc.HashRate(ctx, e.rpcClient)
}
