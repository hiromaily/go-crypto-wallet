package eth

import (
	"context"
)

// ClientVersion returns client version
// https://github.com/ethereum/wiki/wiki/JSON-RPC#web3_clientversion
func (e *Ethereum) ClientVersion(ctx context.Context) (string, error) {
	return e.pkgrpc.ClientVersion(ctx)
}

// SHA3 returns Keccak-256 (not the standardized SHA3-256) of the given data
// https://github.com/ethereum/wiki/wiki/JSON-RPC#web3_sha3
func (e *Ethereum) SHA3(ctx context.Context, data string) (string, error) {
	return e.pkgrpc.SHA3(ctx, data)
}
