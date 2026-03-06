package admin

import "github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp/rpc"

type AdminRPC struct {
	//admin *websocket.WS
	caller rpc.WSCaller
}
