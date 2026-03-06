package rpc

import "context"

// WSCaller abstracts the WebSocket call transport.
// *websocket.WS satisfies this interface without modification.
type WSCaller interface {
	Call(ctx context.Context, req, res any) error
	Close() error
}
