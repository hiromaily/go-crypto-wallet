package xrp

import (
	"context"
	"fmt"

	dtoxrp "github.com/hiromaily/go-crypto-wallet/internal/application/dto/xrp"
	xrpkg "github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp"
	xrprpcadmin "github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp/rpc/admin"
)

// Note: Admin commands are available only if you connect to rippled on a host and port that
// the rippled.cfg file identifies as admin

// https://xrpl.org/key-generation-methods.html

// Assign a Regular Key Pair
// https://xrpl.org/assign-a-regular-key-pair.html
// https://github.com/ripple/ripple-keypairs

// ValidationCreate calls validation_create method
func (w *WSClient) ValidationCreate(ctx context.Context, secret string) (*xrprpcadmin.ResponseValidationCreate, error) {
	if w.adminRPC == nil {
		return nil, XRPErrorDisabledAdminAPI
	}

	res, err := w.adminRPC.ValidationCreate(ctx, secret)
	if err != nil {
		return nil, fmt.Errorf("fail to call adminRPC.ValidationCreate: %w", err)
	}
	return res, nil
}

// WalletProposeWithKey calls wallet_propose method
func (w *WSClient) WalletProposeWithKey(
	ctx context.Context, seed string, keyType dtoxrp.XRPKeyType,
) (*xrprpcadmin.ResponseWalletPropose, error) {
	if w.adminRPC == nil {
		return nil, XRPErrorDisabledAdminAPI
	}

	res, err := w.adminRPC.WalletProposeWithKey(ctx, seed, xrpkg.KeyType(keyType))
	if err != nil {
		return nil, fmt.Errorf("fail to call adminRPC.WalletProposeWithKey: %w", err)
	}
	return res, nil
}

// WalletPropose calls wallet_propose method
// - result is same as long as using same passphrase
func (w *WSClient) WalletPropose(ctx context.Context, passphrase string) (*xrprpcadmin.ResponseWalletPropose, error) {
	if w.adminRPC == nil {
		return nil, XRPErrorDisabledAdminAPI
	}

	res, err := w.adminRPC.WalletPropose(ctx, passphrase)
	if err != nil {
		return nil, fmt.Errorf("fail to call adminRPC.WalletPropose: %w", err)
	}
	return res, nil
}
