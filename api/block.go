package api

// GetBlockCount Blockのカウントを返す
func (b *Bitcoin) GetBlockCount() (int64, error) {
	return b.Client.GetBlockCount()
	//log.Printf("Block count: %d\n", blockCount)
}
