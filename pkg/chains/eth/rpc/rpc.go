package rpc

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// RPCCaller abstracts the raw JSON-RPC transport and connection lifecycle.
// *ethrpc.Client (go-ethereum/rpc) satisfies this interface without modification.
type RPCCaller interface {
	CallContext(ctx context.Context, result any, method string, args ...any) error
	Close()
}

// ETHCaller abstracts higher-level typed Ethereum client operations.
// *ethclient.Client (go-ethereum/ethclient) satisfies this interface without modification.
// Used for methods that require binary-encoded transaction submission or receipt lookup.
type ETHCaller interface {
	SendTransaction(ctx context.Context, tx *types.Transaction) error
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
}

// rpcClient wraps RPCCaller and exposes Ethereum JSON-RPC operations as methods.
type rpcClient struct {
	caller RPCCaller
}

// NewRPCClient creates a new client wrapping the given RPCCaller.
func NewRPCClient(caller RPCCaller) *rpcClient {
	return &rpcClient{caller: caller}
}

// Close closes the underlying RPC connection.
func (c *rpcClient) Close() {
	c.caller.Close()
}
