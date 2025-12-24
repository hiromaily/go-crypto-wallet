package ripple

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	domainCoin "github.com/hiromaily/go-crypto-wallet/pkg/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/ripple/xrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/network/websocket"
)

// NewWSClient try to connect Ripple Server by web socket
func NewWSClient(conf *config.Ripple) (*websocket.WS, *websocket.WS, error) {
	publicURL := conf.WebsocketPublicURL
	if publicURL == "" {
		if publicURL = xrp.GetPublicWSServer(conf.NetworkType).String(); publicURL == "" {
			return nil, nil, errors.New("websocket URL is not found")
		}
	}
	public, err := websocket.New(context.Background(), publicURL)
	if err != nil {
		return nil, nil, fmt.Errorf("fail to call websocket.New() for public API: %s: %w", publicURL, err)
	}

	// acceptable without adminClient
	adminURL := conf.WebsocketAdminURL
	if adminURL == "" {
		return public, nil, nil
	}
	admin, err := websocket.New(context.Background(), adminURL)
	if err != nil {
		return nil, nil, fmt.Errorf("fail to call websocket.New() for admin API: %s: %w", adminURL, err)
	}

	return public, admin, nil
}

// NewGRPCClient try to connect gRPC Server
func NewGRPCClient(conf *config.RippleAPI) (*grpc.ClientConn, error) {
	if conf.URL == "" {
		return nil, errors.New("url for grpc server is not defined in config")
	}
	var opts []grpc.DialOption
	if !conf.IsSecure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	conn, err := grpc.NewClient(conf.URL, opts...)
	if err != nil {
		return nil, fmt.Errorf("fail to call grpc.Dial: %s: %w", conf.URL, err)
	}
	return conn, nil
}

// NewRPCClient RPCClient, maybe not used
// func NewRPCClient(conf *config.Ripple) *jsonrpc.RPCClient {
//	if conf.JSONRpcURL == "" {
//		return nil
//	}
//	rpcClient := jsonrpc.NewClient(conf.JSONRpcURL)
//	return &rpcClient
//}

// NewRipple creates Ripple instance according to coinType
func NewRipple(
	wsPublic *websocket.WS, wsAdmin *websocket.WS, api *xrp.RippleAPI, conf *config.Ripple,
	coinTypeCode domainCoin.CoinTypeCode,
) (Rippler, error) {
	switch coinTypeCode {
	case domainCoin.XRP:
		ripple, err := xrp.NewRipple(context.Background(), wsPublic, wsAdmin, api, coinTypeCode, conf)
		if err != nil {
			return nil, fmt.Errorf("fail to call xrp.NewRipple(): %w", err)
		}
		return ripple, err
	case domainCoin.BTC, domainCoin.BCH, domainCoin.LTC, domainCoin.ETH, domainCoin.ERC20, domainCoin.HYT:
		return nil, fmt.Errorf("coinType %s is not defined", coinTypeCode.String())
	default:
		return nil, fmt.Errorf("coinType %s is not defined", coinTypeCode.String())
	}
}
