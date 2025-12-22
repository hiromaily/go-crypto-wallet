package btc

import (
	"fmt"
)

// GetBlockCount gets block count
//
//	e.g. 1383526
func (b *Bitcoin) GetBlockCount() (int64, error) {
	blockCnt, err := b.Client.GetBlockCount()
	if err != nil {
		return 0, fmt.Errorf("fail to call client.GetBlockCount(): %w", err)
	}

	return blockCnt, nil
}
