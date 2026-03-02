package eth

import (
	"context"
	"math/big"
)

// NetVersion returns the current network id
// "1": Ethereum Mainnet
// "5": Goerli Testnet
func (e *Ethereum) NetVersion(ctx context.Context) (uint16, error) {
	return e.pkgrpc.NetVersion(ctx)
}

// NetListening returns true if client is actively listening for network connections
// https://github.com/ethereum/wiki/wiki/JSON-RPC#net_listening
func (e *Ethereum) NetListening(ctx context.Context) (bool, error) {
	return e.pkgrpc.NetListening(ctx)
}

// NetPeerCount returns number of peers currently connected to the client
// https://github.com/ethereum/wiki/wiki/JSON-RPC#net_peercount
func (e *Ethereum) NetPeerCount(ctx context.Context) (*big.Int, error) {
	return e.pkgrpc.NetPeerCount(ctx)
}
