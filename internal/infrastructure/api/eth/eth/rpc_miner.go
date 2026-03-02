package eth

import (
	"context"
	"math/big"
)

// StartMining starts the CPU mining process with the given number of threads and generate a new DAG if need be
// https://github.com/ethereum/go-ethereum/wiki/Management-APIs#miner_start
func (e *Ethereum) StartMining(ctx context.Context) error {
	return e.pkgrpc.StartMining(ctx)
}

// StopMining stops the CPU mining operation
// https://github.com/ethereum/go-ethereum/wiki/Management-APIs#miner_stop
func (e *Ethereum) StopMining(ctx context.Context) error {
	return e.pkgrpc.StopMining(ctx)
}

// Mining returns true if client is actively mining new blocks
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_mining
func (e *Ethereum) Mining(ctx context.Context) (bool, error) {
	return e.pkgrpc.Mining(ctx)
}

// HashRate returns the number of hashes per second that the node is mining with
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_hashrate
func (e *Ethereum) HashRate(ctx context.Context) (*big.Int, error) {
	return e.pkgrpc.HashRate(ctx)
}
