package xrpadmin

import (
	"context"
	"fmt"

	dtoxrp "github.com/hiromaily/go-crypto-wallet/internal/application/dto/xrp"
	xrpkg "github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp"
	xrprpcadmin "github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp/rpc/admin"
)

// WalletProposeWithKey calls wallet_propose with the given seed and key type.
// This method overrides the embedded AdminRPC method to accept dtoxrp.XRPKeyType
// (the application DTO type) instead of xrpkg.KeyType (the pkg type).
func (a *AdminXRP) WalletProposeWithKey(
	ctx context.Context, seed string, keyType dtoxrp.XRPKeyType,
) (*xrprpcadmin.ResponseWalletPropose, error) {
	res, err := a.AdminRPC.WalletProposeWithKey(ctx, seed, xrpkg.KeyType(keyType))
	if err != nil {
		return nil, fmt.Errorf("fail to call AdminRPC.WalletProposeWithKey: %w", err)
	}
	return res, nil
}
