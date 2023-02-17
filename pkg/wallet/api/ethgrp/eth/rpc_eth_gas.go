package eth

import (
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
)

// GasPrice returns the current price per gas in wei
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_gasprice
// - returns like `1000000000`
func (e *Ethereum) GasPrice() (*big.Int, error) {
	var gasPrice string
	err := e.rpcClient.CallContext(e.ctx, &gasPrice, "eth_gasPrice")
	if err != nil {
		return nil, errors.Wrap(err, "fail to call rpc.CallContext(eth_gasPrice)")
	}
	h, err := hexutil.DecodeBig(gasPrice)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call hexutil.DecodeBig()")
	}

	return h, nil
}

// EstimateGas generates and returns an estimate of how much gas is necessary to allow the transaction to complete
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_estimategas
//
//	is this value GasLimit??
func (e *Ethereum) EstimateGas(msg *ethereum.CallMsg) (*big.Int, error) {
	var estimated string
	err := e.rpcClient.CallContext(e.ctx, &estimated, "eth_estimateGas", toCallArg(msg))
	if err != nil {
		// Invalid params: Invalid bytes format. Expected a 0x-prefixed hex string with even length.
		return nil, errors.Wrap(err, "fail to call rpc.CallContext(eth_estimateGas)")
	}

	h, err := hexutil.DecodeBig(estimated)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call hexutil.DecodeBig()")
	}

	return h, nil
}
