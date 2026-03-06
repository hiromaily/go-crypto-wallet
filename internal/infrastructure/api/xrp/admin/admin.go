package admin

import (
	"context"
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp"
	xrprpcadmin "github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp/rpc/admin"
)

// Note: Admin commands are available only if you connect to rippled on a host and port that
// the rippled.cfg file identifies as admin

// https://xrpl.org/key-generation-methods.html

// Assign a Regular Key Pair
// https://xrpl.org/assign-a-regular-key-pair.html
// https://github.com/ripple/ripple-keypairs

type adminRPC struct {
	caller xrprpcadmin.XRPAdminer
}

func NewAdminRPC(caller xrprpcadmin.XRPAdminer) *adminRPC {
	return &adminRPC{
		caller: caller,
	}
}

// ValidationCreate calls validation_create method
func (r *adminRPC) ValidationCreate(ctx context.Context, secret string) (*xrprpcadmin.ResponseValidationCreate, error) {
	res, err := r.caller.ValidationCreate(ctx, secret)
	if err != nil {
		return nil, fmt.Errorf("fail to call ValidationCreate: %w", err)
	}
	return res, nil
}

// WalletProposeWithKey calls wallet_propose method
func (r *adminRPC) WalletProposeWithKey(
	ctx context.Context, seed string, keyType dtoxrp.XRPKeyType,
) (*xrprpcadmin.ResponseWalletPropose, error) {
	res, err := r.caller.WalletProposeWithKey(ctx, seed, xrp.KeyType(keyType))
	if err != nil {
		return nil, fmt.Errorf("fail to call WalletProposeWithKey: %w", err)
	}
	return res, nil
}

// WalletPropose calls wallet_propose method
// - result is same as long as using same passphrase
func (r *adminRPC) WalletPropose(ctx context.Context, passphrase string) (*xrprpcadmin.ResponseWalletPropose, error) {
	res, err := r.caller.WalletPropose(ctx, passphrase)
	if err != nil {
		return nil, fmt.Errorf("fail to call WalletPropose: %w", err)
	}
	return res, nil
}
