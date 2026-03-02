package rpc

import (
	"context"
	"fmt"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

// NetVersion returns the current network ID.
// "1" = Ethereum Mainnet, "5" = Goerli Testnet, etc.
// https://github.com/ethereum/wiki/wiki/JSON-RPC#net_version
func (c *rpcClient) NetVersion(ctx context.Context) (uint16, error) {
	var resNetVersion string
	err := c.caller.CallContext(ctx, &resNetVersion, "net_version")
	if err != nil {
		return 0, fmt.Errorf("fail to call client.CallContext(net_version): %w", err)
	}
	u, err := strconv.ParseUint(resNetVersion, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("fail to call strconv.ParseUint(%s): %w", resNetVersion, err)
	}

	return uint16(u), nil
}

// NetListening returns true if the client is actively listening for network connections.
// https://github.com/ethereum/wiki/wiki/JSON-RPC#net_listening
func (c *rpcClient) NetListening(ctx context.Context) (bool, error) {
	var isConnected bool
	err := c.caller.CallContext(ctx, &isConnected, "net_listening")
	if err != nil {
		return false, fmt.Errorf("fail to call rpc.CallContext(net_listening): %w", err)
	}

	return isConnected, nil
}

// NetPeerCount returns the number of peers currently connected to the client.
// https://github.com/ethereum/wiki/wiki/JSON-RPC#net_peercount
func (c *rpcClient) NetPeerCount(ctx context.Context) (*big.Int, error) {
	var resPeerNumber string
	err := c.caller.CallContext(ctx, &resPeerNumber, "net_peerCount")
	if err != nil {
		return nil, fmt.Errorf("fail to call client.CallContext(net_peerCount): %w", err)
	}
	return hexutil.DecodeBig(resPeerNumber)
}
