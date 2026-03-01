package rpc

import "github.com/btcsuite/btcd/rpcclient"

// rpcClient wraps *rpcclient.Client and exposes Bitcoin RPC operations as methods.
type rpcClient struct {
	client *rpcclient.Client
}

// NewRPCClient creates a new client wrapping the given rpcclient.Client.
func NewRPCClient(client *rpcclient.Client) *rpcClient {
	return &rpcClient{client: client}
}
