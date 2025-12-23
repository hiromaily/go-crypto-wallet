package eth

import (
	"context"
	"errors"
	"fmt"
	"math/big"
)

// StartMining starts the CPU mining process with the given number of threads and generate a new DAG if need be
// https://github.com/ethereum/go-ethereum/wiki/Management-APIs#miner_start
func (e *Ethereum) StartMining(ctx context.Context) error {
	var r []byte
	// TODO: Result needs to be verified
	err := e.rpcClient.CallContext(ctx, &r, "miner_start")
	if err != nil {
		return fmt.Errorf("fail to call rpc.CallContext(miner_start): %w", err)
	}
	return err
}

// StopMining stops the CPU mining operation
// https://github.com/ethereum/go-ethereum/wiki/Management-APIs#miner_stop
func (e *Ethereum) StopMining(ctx context.Context) error {
	err := e.rpcClient.CallContext(ctx, nil, "miner_stop")
	if err != nil {
		return errors.New("fail to call rpc.CallContext(miner_start)")
	}
	return err
}

// Mining returns true if client is actively mining new blocks
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_mining
func (e *Ethereum) Mining(ctx context.Context) (bool, error) {
	var bRet bool
	err := e.rpcClient.CallContext(ctx, &bRet, "eth_mining")
	if err != nil {
		return false, errors.New("fail to call rpc.CallContext(eth_mining)")
	}
	return bRet, err
}

// HashRate returns the number of hashes per second that the node is mining with
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_hashrate
func (e *Ethereum) HashRate(ctx context.Context) (*big.Int, error) {
	var hashCount string
	err := e.rpcClient.CallContext(ctx, &hashCount, "eth_hashrate")
	if err != nil {
		return nil, fmt.Errorf("fail to call rpc.CallContext(eth_hashrate): %w", err)
	}
	h, err := e.DecodeBig(hashCount)
	if err != nil {
		return nil, fmt.Errorf("fail to call hexutil.DecodeBig(): %w", err)
	}

	return h, nil
}
