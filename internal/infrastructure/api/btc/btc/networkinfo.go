package btc

import (
	"fmt"

	dtobtc "github.com/hiromaily/go-crypto-wallet/internal/application/dto/btc"
	btcrpc "github.com/hiromaily/go-crypto-wallet/pkg/chains/btc/rpc"
)

// GetNetworkInfo call RPC `getnetworkinfo`
func (b *Bitcoin) GetNetworkInfo() (*dtobtc.NetworkInfo, error) {
	result, err := btcrpc.GetNetworkInfo(b.Client)
	if err != nil {
		return nil, fmt.Errorf("fail to call btcrpc.GetNetworkInfo(): %w", err)
	}

	return ToNetworkInfoFromPkg(result), nil
}

// GetBlockchainInfo call RPC `getblockchaininfo`
func (b *Bitcoin) GetBlockchainInfo() (*dtobtc.BlockchainInfo, error) {
	result, err := btcrpc.GetBlockchainInfo(b.Client)
	if err != nil {
		return nil, fmt.Errorf("fail to call btcrpc.GetBlockchainInfo(): %w", err)
	}

	return ToBlockchainInfoFromPkg(result), nil
}
