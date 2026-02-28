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
func NetVersion(ctx context.Context, caller RPCCaller) (uint16, error) {
	var resNetVersion string
	err := caller.CallContext(ctx, &resNetVersion, "net_version")
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
func NetListening(ctx context.Context, caller RPCCaller) (bool, error) {
	var isConnected bool
	err := caller.CallContext(ctx, &isConnected, "net_listening")
	if err != nil {
		return false, fmt.Errorf("fail to call rpc.CallContext(net_listening): %w", err)
	}

	return isConnected, nil
}

// NetPeerCount returns the number of peers currently connected to the client.
// https://github.com/ethereum/wiki/wiki/JSON-RPC#net_peercount
func NetPeerCount(ctx context.Context, caller RPCCaller) (*big.Int, error) {
	var resPeerNumber string
	err := caller.CallContext(ctx, &resPeerNumber, "net_peerCount")
	if err != nil {
		return nil, fmt.Errorf("fail to call client.CallContext(net_peerCount): %w", err)
	}
	return hexutil.DecodeBig(resPeerNumber)
}
