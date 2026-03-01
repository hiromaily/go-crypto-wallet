package btc

import (
	btcrpc "github.com/hiromaily/go-crypto-wallet/pkg/chains/btc/rpc"
)

// Logging gets and sets the logging configuration for Bitcoin Core
func (b *Bitcoin) Logging() (*btcrpc.LoggingResult, error) {
	return b.pkgrpc.Logging()
}
