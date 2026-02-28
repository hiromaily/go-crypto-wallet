package btc

import (
	"fmt"

	btcrpc "github.com/hiromaily/go-crypto-wallet/pkg/chains/btc/rpc"
)

// GetBlockCount gets block count
//
//	e.g. 1383526
func (b *Bitcoin) GetBlockCount() (int64, error) {
	blockCnt, err := btcrpc.GetBlockCount(b.Client)
	if err != nil {
		return 0, fmt.Errorf("fail to call btcrpc.GetBlockCount(): %w", err)
	}

	return blockCnt, nil
}
