package eth

import (
	"math/big"

	"github.com/pkg/errors"
)

// StartMining starts the CPU mining process with the given number of threads and generate a new DAG if need be
// https://github.com/ethereum/go-ethereum/wiki/Management-APIs#miner_start
func (e *Ethereum) StartMining() error {
	var r []byte
	// TODO: Result needs to be verified
	err := e.rpcClient.CallContext(e.ctx, &r, "miner_start")
	if err != nil {
		return errors.Wrap(err, "fail to call rpc.CallContext(miner_start)")
	}
	return err
}

// StopMining stops the CPU mining operation
// https://github.com/ethereum/go-ethereum/wiki/Management-APIs#miner_stop
func (e *Ethereum) StopMining() error {
	err := e.rpcClient.CallContext(e.ctx, nil, "miner_stop")
	if err != nil {
		return errors.New("fail to call rpc.CallContext(miner_start)")
	}
	return err
}

// HashRate returns the number of hashes per second that the node is mining with
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_hashrate
func (e *Ethereum) HashRate() (*big.Int, error) {
	var hashCount string
	err := e.rpcClient.CallContext(e.ctx, &hashCount, "eth_hashrate")
	if err != nil {
		return nil, errors.Wrap(err, "fail to call rpc.CallContext(eth_hashrate)")
	}
	h, err := e.DecodeBig(hashCount)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call hexutil.DecodeBig()")
	}

	return h, nil
}
