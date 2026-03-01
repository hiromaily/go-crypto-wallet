package xrp

import (
	"context"
	"fmt"

	dtoxrp "github.com/hiromaily/go-crypto-wallet/internal/application/dto/xrp"
	xrprpc "github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp/rpc"
)

// Note: Admin commands are available only if you connect to rippled on a host and port that
// the rippled.cfg file identifies as admin

// https://xrpl.org/key-generation-methods.html

// Assign a Regular Key Pair
// https://xrpl.org/assign-a-regular-key-pair.html
// https://github.com/ripple/ripple-keypairs

// ValidationCreate calls validation_create method
func (r *XRP) ValidationCreate(ctx context.Context, secret string) (*xrprpc.ResponseValidationCreate, error) {
	if r.wsAdmin == nil {
		return nil, XRPErrorDisabledAdminAPI
	}

	res, err := xrprpc.ValidationCreate(ctx, r.wsAdmin, secret)
	if err != nil {
		return nil, fmt.Errorf("fail to call xrprpc.ValidationCreate: %w", err)
	}
	return res, nil
}

// WalletProposeWithKey calls wallet_propose method
func (r *XRP) WalletProposeWithKey(
	ctx context.Context, seed string, keyType dtoxrp.XRPKeyType,
) (*xrprpc.ResponseWalletPropose, error) {
	if r.wsAdmin == nil {
		return nil, XRPErrorDisabledAdminAPI
	}

	res, err := xrprpc.WalletProposeWithKey(ctx, r.wsAdmin, seed, xrprpc.KeyType(keyType))
	if err != nil {
		return nil, fmt.Errorf("fail to call xrprpc.WalletProposeWithKey: %w", err)
	}
	return res, nil
}

// WalletPropose calls wallet_propose method
// - result is same as long as using same passphrase
func (r *XRP) WalletPropose(ctx context.Context, passphrase string) (*xrprpc.ResponseWalletPropose, error) {
	if r.wsAdmin == nil {
		return nil, XRPErrorDisabledAdminAPI
	}

	res, err := xrprpc.WalletPropose(ctx, r.wsAdmin, passphrase)
	if err != nil {
		return nil, fmt.Errorf("fail to call xrprpc.WalletPropose: %w", err)
	}
	return res, nil
}
