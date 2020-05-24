package xrpgrp

import (
	"context"

	"github.com/pkg/errors"
	"github.com/rubblelabs/ripple/websockets"
	"go.uber.org/zap"

	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/xrpgrp/xrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
	"github.com/hiromaily/go-crypto-wallet/pkg/ws"
)

// NewWSClient try to connect Ripple Server by web socket
func NewWSClient(conf *config.Ripple) (*ws.WS, error) {
	url := conf.WebsocketURL
	if url == "" {
		if url = xrp.GetPublicWSServer(conf.NetworkType).String(); url == "" {
			return nil, errors.New("websocket URL is not found")
		}
	}
	return ws.New(context.Background(), url), nil
}

// NewWSRemote try to connect Ripple Server by web socket
func NewWSRemote(conf *config.Ripple) (*websockets.Remote, error) {
	url := conf.WebsocketURL
	if url == "" {
		if url = xrp.GetPublicWSServer(conf.NetworkType).String(); url == "" {
			return nil, errors.New("websocket URL is not found")
		}
	}
	return websockets.NewRemote(url)
}

// NewRPCClient RPCClient, maybe not used
//func NewRPCClient(conf *config.Ripple) *jsonrpc.RPCClient {
//	if conf.JSONRpcURL == "" {
//		return nil
//	}
//	rpcClient := jsonrpc.NewClient(conf.JSONRpcURL)
//	return &rpcClient
//}

// NewRipple creates Ripple instance according to coinType
func NewRipple(wsClient *ws.WS, wsRemote *websockets.Remote, conf *config.Ripple, logger *zap.Logger, coinTypeCode coin.CoinTypeCode) (Rippler, error) {
	switch coinTypeCode {
	case coin.XRP:
		eth, err := xrp.NewRipple(context.Background(), wsClient, wsRemote, coinTypeCode, conf, logger)
		if err != nil {
			return nil, errors.Wrap(err, "fail to call xrp.NewRipple()")
		}
		return eth, err
	}
	return nil, errors.Errorf("coinType %s is not defined", coinTypeCode.String())
}
