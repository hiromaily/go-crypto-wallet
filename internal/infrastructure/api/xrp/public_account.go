package xrp

import (
	"context"
	"fmt"

	xrprpcpublic "github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp/rpc/public"
)

// https://xrpl.org/account-methods.html
// error: https://xrpl.org/error-formatting.html#universal-errors

// AccountChannels calls account_channels method
func (w *WSClient) AccountChannels(
	ctx context.Context, sender, receiver string,
) (*xrprpcpublic.ResponseAccountChannels, error) {
	res, err := w.publicRPC.AccountChannels(ctx, sender, receiver)
	if err != nil {
		return nil, fmt.Errorf("fail to call publicRPC.AccountChannels: %w", err)
	}
	return res, nil
}

// AccountInfo calls account_info method
func (w *WSClient) AccountInfo(ctx context.Context, address string) (*xrprpcpublic.ResponseAccountInfo, error) {
	res, err := w.publicRPC.AccountInfo(ctx, address)
	if err != nil {
		return nil, fmt.Errorf("fail to call publicRPC.AccountInfo: %w", err)
	}
	return res, nil
}
