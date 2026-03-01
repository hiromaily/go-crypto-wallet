package btc

import (
	"fmt"

	btcrpc "github.com/hiromaily/go-crypto-wallet/pkg/chains/btc/rpc"
)

// GetNetworkInfo call RPC `getnetworkinfo`
func (b *Bitcoin) GetNetworkInfo() (*btcrpc.GetNetworkInfoResult, error) {
	result, err := btcrpc.GetNetworkInfo(b.Client)
	if err != nil {
		return nil, fmt.Errorf("fail to call btcrpc.GetNetworkInfo(): %w", err)
	}

	return result, nil
}

// GetBlockchainInfo call RPC `getblockchaininfo`
func (b *Bitcoin) GetBlockchainInfo() (*btcrpc.GetBlockchainInfoResult, error) {
	result, err := btcrpc.GetBlockchainInfo(b.Client)
	if err != nil {
		return nil, fmt.Errorf("fail to call btcrpc.GetBlockchainInfo(): %w", err)
	}

	return result, nil
}
