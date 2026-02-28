package btc

import (
	"fmt"

	dtobtc "github.com/hiromaily/go-crypto-wallet/internal/application/dto/btc"
	btcrpc "github.com/hiromaily/go-crypto-wallet/pkg/chains/btc/rpc"
)

// Logging calls RPC `logging`
func (b *Bitcoin) Logging() (*dtobtc.LoggingResult, error) {
	result, err := btcrpc.Logging(b.Client)
	if err != nil {
		return nil, fmt.Errorf("fail to call btcrpc.Logging(): %w", err)
	}

	return ToLoggingResultFromPkg(result), nil
}
