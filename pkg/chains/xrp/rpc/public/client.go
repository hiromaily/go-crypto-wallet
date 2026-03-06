package public

import "github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp/rpc"

// PublicRPC wraps a WSCaller to provide public XRP node operations.
type PublicRPC struct {
	caller rpc.WSCaller
}

// NewPublicRPC creates a new PublicRPC with the given caller.
func NewPublicRPC(caller rpc.WSCaller) *PublicRPC {
	return &PublicRPC{caller: caller}
}
