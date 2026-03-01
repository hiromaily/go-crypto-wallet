package btc

import (
	"fmt"

	btcrpc "github.com/hiromaily/go-crypto-wallet/pkg/chains/btc/rpc"
)

// Logging calls RPC `logging`
func (b *Bitcoin) Logging() (*btcrpc.LoggingResult, error) {
	result, err := b.RPC.Logging()
	if err != nil {
		return nil, fmt.Errorf("fail to call btcrpc.Logging(): %w", err)
	}

	return result, nil
}
