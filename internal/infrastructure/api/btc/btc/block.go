package btc

// GetBlockCount returns the number of blocks in the longest blockchain
func (b *Bitcoin) GetBlockCount() (int64, error) {
	return b.pkgrpc.GetBlockCount()
}
