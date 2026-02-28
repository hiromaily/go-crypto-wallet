package xrp

import (
	"context"
	"fmt"

	dtoxrp "github.com/hiromaily/go-crypto-wallet/internal/application/dto/xrp"
	xrprpc "github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp/rpc"
)

// https://xrpl.org/account-methods.html
// error: https://xrpl.org/error-formatting.html#universal-errors

// AccountChannels calls account_channels method
func (r *XRP) AccountChannels(
	ctx context.Context, sender, receiver string,
) (*dtoxrp.ResponseAccountChannels, error) {
	res, err := xrprpc.AccountChannels(ctx, r.wsPublic, sender, receiver)
	if err != nil {
		return nil, fmt.Errorf("fail to call xrprpc.AccountChannels: %w", err)
	}
	return ToDTOResponseAccountChannelsFromPkg(res), nil
}

// AccountInfo calls account_info method
func (r *XRP) AccountInfo(ctx context.Context, address string) (*dtoxrp.ResponseAccountInfo, error) {
	res, err := xrprpc.AccountInfo(ctx, r.wsPublic, address)
	if err != nil {
		return nil, fmt.Errorf("fail to call xrprpc.AccountInfo: %w", err)
	}
	return ToDTOResponseAccountInfoFromPkg(res), nil
}
