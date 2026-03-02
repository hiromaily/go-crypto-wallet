package rpc

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// GasPrice returns the current price per gas in wei.
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_gasprice
func (c *rpcClient) GasPrice(ctx context.Context) (*big.Int, error) {
	var gasPrice string
	err := c.caller.CallContext(ctx, &gasPrice, "eth_gasPrice")
	if err != nil {
		return nil, fmt.Errorf("fail to call rpc.CallContext(eth_gasPrice): %w", err)
	}
	h, err := hexutil.DecodeBig(gasPrice)
	if err != nil {
		return nil, fmt.Errorf("fail to call hexutil.DecodeBig(): %w", err)
	}

	return h, nil
}

// EstimateGas returns an estimate of the gas needed to execute the given transaction.
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_estimategas
func (c *rpcClient) EstimateGas(ctx context.Context, msg *ethereum.CallMsg) (*big.Int, error) {
	var estimated string
	err := c.caller.CallContext(ctx, &estimated, "eth_estimateGas", toCallArg(msg))
	if err != nil {
		return nil, fmt.Errorf("fail to call rpc.CallContext(eth_estimateGas): %w", err)
	}

	h, err := hexutil.DecodeBig(estimated)
	if err != nil {
		return nil, fmt.Errorf("fail to call hexutil.DecodeBig(): %w", err)
	}

	return h, nil
}
