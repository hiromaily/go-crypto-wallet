package btc

import (
	btcrpc "github.com/hiromaily/go-crypto-wallet/pkg/chains/btc/rpc"
)

// GetNetworkInfo returns information about the node's connection to the network
func (b *Bitcoin) GetNetworkInfo() (*btcrpc.GetNetworkInfoResult, error) {
	return b.pkgrpc.GetNetworkInfo()
}

// GetBlockchainInfo returns an object containing various state info regarding blockchain processing
func (b *Bitcoin) GetBlockchainInfo() (*btcrpc.GetBlockchainInfoResult, error) {
	return b.pkgrpc.GetBlockchainInfo()
}
