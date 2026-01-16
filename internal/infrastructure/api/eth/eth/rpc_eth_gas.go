package eth

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// GasPrice returns the current price per gas in wei
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_gasprice
// - returns like `1000000000`
func (e *Ethereum) GasPrice(ctx context.Context) (*big.Int, error) {
	var gasPrice string
	err := e.rpcClient.CallContext(ctx, &gasPrice, "eth_gasPrice")
	if err != nil {
		return nil, fmt.Errorf("fail to call rpc.CallContext(eth_gasPrice): %w", err)
	}
	h, err := hexutil.DecodeBig(gasPrice)
	if err != nil {
		return nil, fmt.Errorf("fail to call hexutil.DecodeBig(): %w", err)
	}

	return h, nil
}

// EstimateGas generates and returns an estimate of how much gas is necessary to allow the transaction to complete
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_estimategas
//
//	is this value GasLimit??
func (e *Ethereum) EstimateGas(ctx context.Context, msg *ethereum.CallMsg) (*big.Int, error) {
	var estimated string
	err := e.rpcClient.CallContext(ctx, &estimated, "eth_estimateGas", toCallArg(msg))
	if err != nil {
		// Invalid params: Invalid bytes format. Expected a 0x-prefixed hex string with even length.
		return nil, fmt.Errorf("fail to call rpc.CallContext(eth_estimateGas): %w", err)
	}

	h, err := hexutil.DecodeBig(estimated)
	if err != nil {
		return nil, fmt.Errorf("fail to call hexutil.DecodeBig(): %w", err)
	}

	return h, nil
}
