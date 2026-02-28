package eth

import (
	"context"
	"math/big"

	ethrpc "github.com/hiromaily/go-crypto-wallet/pkg/chains/eth/rpc"
)

// NetVersion returns the current network id
// "1": Ethereum Mainnet
// "5": Goerli Testnet
func (e *Ethereum) NetVersion(ctx context.Context) (uint16, error) {
	return ethrpc.NetVersion(ctx, e.rpcClient)
}

// NetListening returns true if client is actively listening for network connections
// https://github.com/ethereum/wiki/wiki/JSON-RPC#net_listening
func (e *Ethereum) NetListening(ctx context.Context) (bool, error) {
	return ethrpc.NetListening(ctx, e.rpcClient)
}

// NetPeerCount returns number of peers currently connected to the client
// https://github.com/ethereum/wiki/wiki/JSON-RPC#net_peercount
func (e *Ethereum) NetPeerCount(ctx context.Context) (*big.Int, error) {
	return ethrpc.NetPeerCount(ctx, e.rpcClient)
}
