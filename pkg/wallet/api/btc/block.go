package btc

import (
	"github.com/pkg/errors"
)

// GetBlockCount gets block count
//  e.g. 1383526
func (b *Bitcoin) GetBlockCount() (int64, error) {
	blockCnt, err := b.client.GetBlockCount()
	if err != nil {
		return 0, errors.Wrap(err, "fail to call client.GetBlockCount()")
	}

	return blockCnt, nil
}
