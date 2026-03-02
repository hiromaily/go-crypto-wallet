package rpc

import (
	"context"
	"fmt"
)

// ClientVersion returns the client version string.
// https://github.com/ethereum/wiki/wiki/JSON-RPC#web3_clientversion
func (c *rpcClient) ClientVersion(ctx context.Context) (string, error) {
	var resClientVersion string
	err := c.caller.CallContext(ctx, &resClientVersion, "web3_clientVersion")
	if err != nil {
		return "", fmt.Errorf("fail to call client.CallContext(web3_clientVersion): %w", err)
	}
	return resClientVersion, nil
}

// SHA3 returns the Keccak-256 hash of the given data.
// https://github.com/ethereum/wiki/wiki/JSON-RPC#web3_sha3
func (c *rpcClient) SHA3(ctx context.Context, data string) (string, error) {
	var res string
	err := c.caller.CallContext(ctx, &res, "web3_sha3", data)
	if err != nil {
		return "", fmt.Errorf("fail to call client.CallContext(web3_sha3): %w", err)
	}
	return res, nil
}
