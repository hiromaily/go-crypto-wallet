package eth

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"

	ethrpc "github.com/hiromaily/go-crypto-wallet/pkg/chains/eth/rpc"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// GasPrice returns the current price per gas in wei
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_gasprice
func (e *Ethereum) GasPrice(ctx context.Context) (*big.Int, error) {
	return ethrpc.GasPrice(ctx, e.rpcClient)
}

// SuggestGasTipCap returns a suggested gas tip cap (max priority fee per gas) for EIP-1559 transactions.
// If the RPC call fails, it falls back to the configured MaxPriorityFeePerGas value (default: 2 Gwei).
func (e *Ethereum) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	tip, err := e.ethClient.SuggestGasTipCap(ctx)
	if err != nil {
		const (
			defaultPriorityFeeGwei = 2
			gweiInWei              = 1_000_000_000
		)
		logger.Warn("SuggestGasTipCap RPC failed, using config fallback", "error", err)
		priorityFeeGwei := e.conf.MaxPriorityFeePerGas
		if priorityFeeGwei == 0 {
			priorityFeeGwei = defaultPriorityFeeGwei
		}
		return new(big.Int).SetUint64(priorityFeeGwei * gweiInWei), nil
	}
	return tip, nil
}

// EstimateGas generates and returns an estimate of how much gas is necessary to allow the transaction to complete
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_estimategas
func (e *Ethereum) EstimateGas(ctx context.Context, msg *ethereum.CallMsg) (*big.Int, error) {
	return ethrpc.EstimateGas(ctx, e.rpcClient, msg)
}
