package xrp

import (
	"context"
	"fmt"

	xrprpc "github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp/rpc"
)

// https://xrpl.org/account-methods.html
// error: https://xrpl.org/error-formatting.html#universal-errors

// AccountChannels calls account_channels method
func (w *WSClient) AccountChannels(
	ctx context.Context, sender, receiver string,
) (*xrprpc.ResponseAccountChannels, error) {
	res, err := xrprpc.AccountChannels(ctx, w.public, sender, receiver)
	if err != nil {
		return nil, fmt.Errorf("fail to call xrprpc.AccountChannels: %w", err)
	}
	return res, nil
}

// AccountInfo calls account_info method
func (w *WSClient) AccountInfo(ctx context.Context, address string) (*xrprpc.ResponseAccountInfo, error) {
	res, err := xrprpc.AccountInfo(ctx, w.public, address)
	if err != nil {
		return nil, fmt.Errorf("fail to call xrprpc.AccountInfo: %w", err)
	}
	return res, nil
}
