package xrp

import "github.com/hiromaily/go-crypto-wallet/pkg/websocket"

// WSClient encapsulates the WebSocket connections to the XRP node.
// It implements all operations that communicate directly with the node over WebSocket,
// keeping them separate from the gRPC-based operations in XRPAPI.
type WSClient struct {
	public *websocket.WS
	admin  *websocket.WS
}

// newWSClient creates a new WSClient.
func newWSClient(public, admin *websocket.WS) *WSClient {
	return &WSClient{
		public: public,
		admin:  admin,
	}
}

// Close disconnects both WebSocket connections.
func (w *WSClient) Close() error {
	if w.public != nil {
		_ = w.public.Close()
	}
	if w.admin != nil {
		_ = w.admin.Close()
	}
	return nil
}
