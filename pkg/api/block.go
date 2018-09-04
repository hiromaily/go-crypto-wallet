package api

import (
	"github.com/pkg/errors"
)

// GetBlockCount Blockのカウントを返す
//  e.g. 1383526
func (b *Bitcoin) GetBlockCount() (int64, error) {
	blockCnt, err := b.client.GetBlockCount()
	if err != nil {
		return 0, errors.Errorf("client.GetBlockCount(): error: %s", err)
	}

	return blockCnt, nil
}
