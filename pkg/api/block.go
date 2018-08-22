package api

import (
	"github.com/pkg/errors"
)

// GetBlockCount Blockのカウントを返す
//  e.g. 1383526
func (b *Bitcoin) GetBlockCount() (int64, error) {
	blockCnt, err := b.client.GetBlockCount()
	if err != nil {
		return 0, errors.Errorf("GetBlockCount(): error: %v", err)
	}

	return blockCnt, nil
}
