package xrp

import (
	"github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp/rpc/admin"
	"github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp/rpc/public"
	"github.com/hiromaily/go-crypto-wallet/pkg/websocket"
)

// WSClient encapsulates the WebSocket connections to the XRP node.
// It implements all operations that communicate directly with the node over WebSocket.
type WSClient struct {
	wsPublic  *websocket.WS
	wsAdmin   *websocket.WS
	publicRPC *public.PublicRPC
	adminRPC  *admin.AdminRPC
}

// newWSClient creates a new WSClient.
func newWSClient(wsPublic, wsAdmin *websocket.WS) *WSClient {
	wsc := &WSClient{
		wsPublic: wsPublic,
		wsAdmin:  wsAdmin,
	}
	if wsPublic != nil {
		wsc.publicRPC = public.NewPublicRPC(wsPublic)
	}
	if wsAdmin != nil {
		wsc.adminRPC = admin.NewAdminRPC(wsAdmin)
	}
	return wsc
}

// Close disconnects both WebSocket connections.
func (w *WSClient) Close() error {
	if w.wsPublic != nil {
		_ = w.wsPublic.Close()
	}
	if w.wsAdmin != nil {
		_ = w.wsAdmin.Close()
	}
	return nil
}
