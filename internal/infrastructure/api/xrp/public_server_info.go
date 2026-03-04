package xrp

import (
	"context"
	"fmt"

	xrprpc "github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp/rpc"
)

// https://xrpl.org/server-info-methods.html

// ServerInfo calls server_info method
func (w *WSClient) ServerInfo(ctx context.Context) (*xrprpc.ResponseServerInfo, error) {
	res, err := xrprpc.ServerInfo(ctx, w.public)
	if err != nil {
		return nil, fmt.Errorf("fail to call xrprpc.ServerInfo: %w", err)
	}
	return res, nil
}
