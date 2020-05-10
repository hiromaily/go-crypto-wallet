package eth

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
	"github.com/pkg/errors"
)

// GasPrice returns the current price per gas in wei
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_gasprice
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

//eth_estimateGas
//https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_estimategas