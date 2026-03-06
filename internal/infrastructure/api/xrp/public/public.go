package public

import (
	"context"
	"fmt"

	xrprpcpublic "github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp/rpc/public"
)

type publicRPC struct {
	caller xrprpcpublic.XRPPublicer
}

func NewPublicRPC(caller xrprpcpublic.XRPPublicer) *publicRPC {
	return &publicRPC{
		caller: caller,
	}
}

// https://xrpl.org/account-methods.html
// error: https://xrpl.org/error-formatting.html#universal-errors

// AccountChannels calls account_channels method
func (r *publicRPC) AccountChannels(
	ctx context.Context, sender, receiver string,
) (*xrprpcpublic.ResponseAccountChannels, error) {
	res, err := r.caller.AccountChannels(ctx, sender, receiver)
	if err != nil {
		return nil, fmt.Errorf("fail to call AccountChannels: %w", err)
	}
	return res, nil
}

// AccountInfo calls account_info method
func (r *publicRPC) AccountInfo(ctx context.Context, address string) (*xrprpcpublic.ResponseAccountInfo, error) {
	res, err := r.caller.AccountInfo(ctx, address)
	if err != nil {
		return nil, fmt.Errorf("fail to call xrprpc.AccountInfo: %w", err)
	}
	return res, nil
}

// https://xrpl.org/server-info-methods.html

// ServerInfo calls server_info method
func (r *publicRPC) ServerInfo(ctx context.Context) (*xrprpcpublic.ResponseServerInfo, error) {
	res, err := r.caller.ServerInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("fail to call ServerInfo: %w", err)
	}
	return res, nil
}
