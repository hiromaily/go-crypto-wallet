package xrp

import (
	"context"
	"fmt"

	xrprpcpublic "github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp/rpc/public"
)

// https://xrpl.org/server-info-methods.html

// ServerInfo calls server_info method
func (w *WSClient) ServerInfo(ctx context.Context) (*xrprpcpublic.ResponseServerInfo, error) {
	res, err := w.publicRPC.ServerInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("fail to call publicRPC.ServerInfo: %w", err)
	}
	return res, nil
}
