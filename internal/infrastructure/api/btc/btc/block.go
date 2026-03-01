package btc

import (
	"fmt"
)

// GetBlockCount gets block count
//
//	e.g. 1383526
func (b *Bitcoin) GetBlockCount() (int64, error) {
	blockCnt, err := b.RPC.GetBlockCount()
	if err != nil {
		return 0, fmt.Errorf("fail to call btcrpc.GetBlockCount(): %w", err)
	}

	return blockCnt, nil
}
