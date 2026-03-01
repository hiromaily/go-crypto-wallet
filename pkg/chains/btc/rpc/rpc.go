package rpc

import "github.com/btcsuite/btcd/rpcclient"

// Client wraps *rpcclient.Client and exposes Bitcoin RPC operations as methods.
type Client struct {
	client *rpcclient.Client
}

// New creates a new Client wrapping the given rpcclient.Client.
func New(client *rpcclient.Client) *Client {
	return &Client{client: client}
}
