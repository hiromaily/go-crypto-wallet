package public

import "github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp/rpc"

type PublicRPC struct {
	//public *websocket.WS
	caller rpc.WSCaller
}
