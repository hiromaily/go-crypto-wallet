package rpc

import "github.com/btcsuite/btcd/rpcclient"

// rpcClient wraps *rpcclient.Client and exposes Bitcoin RPC operations as methods.
type rpcClient struct {
	btcdClient *rpcclient.Client
}

// NewRPCClient creates a new client wrapping the given rpcclient.Client.
func NewRPCClient(btcdClient *rpcclient.Client) *rpcClient {
	return &rpcClient{btcdClient: btcdClient}
}
