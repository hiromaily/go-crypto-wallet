package admin

import "github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp/rpc"

// AdminRPC wraps a WSCaller to provide admin XRP node operations.
type AdminRPC struct {
	caller rpc.WSCaller
}

// NewAdminRPC creates a new AdminRPC with the given caller.
func NewAdminRPC(caller rpc.WSCaller) *AdminRPC {
	return &AdminRPC{caller: caller}
}

// Close closes the underlying WebSocket connection.
func (r *AdminRPC) Close() error {
	return r.caller.Close()
}
