package xrpadmin

import (
	"context"
	"fmt"

	dtoxrp "github.com/hiromaily/go-crypto-wallet/internal/application/dto/xrp"
	xrpkg "github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp"
)

// WalletProposeWithKey calls wallet_propose with the given seed and key type.
// This method overrides the embedded AdminRPC method to accept dtoxrp.XRPKeyType
// (the application DTO type) instead of xrpkg.KeyType (the pkg type).
func (a *AdminXRP) WalletProposeWithKey(
	ctx context.Context, seed string, keyType dtoxrp.XRPKeyType,
) (*dtoxrp.WalletInfo, error) {
	res, err := a.AdminRPC.WalletProposeWithKey(ctx, seed, xrpkg.KeyType(keyType))
	if err != nil {
		return nil, fmt.Errorf("fail to call AdminRPC.WalletProposeWithKey: %w", err)
	}
	if res.Error != "" {
		return nil, fmt.Errorf("fail to call AdminRPC.WalletProposeWithKey: %s", res.Error)
	}
	return &dtoxrp.WalletInfo{
		AccountID:     res.Result.AccountID,
		KeyType:       res.Result.KeyType,
		MasterKey:     res.Result.MasterKey,
		MasterSeed:    res.Result.MasterSeed,
		MasterSeedHex: res.Result.MasterSeedHex,
		PublicKey:     res.Result.PublicKey,
		PublicKeyHex:  res.Result.PublicKeyHex,
	}, nil
}

// WalletPropose calls wallet_propose with a passphrase and converts the result to a DTO.
func (a *AdminXRP) WalletPropose(ctx context.Context, passphrase string) (*dtoxrp.WalletInfo, error) {
	res, err := a.AdminRPC.WalletPropose(ctx, passphrase)
	if err != nil {
		return nil, fmt.Errorf("fail to call AdminRPC.WalletPropose: %w", err)
	}
	if res.Error != "" {
		return nil, fmt.Errorf("fail to call AdminRPC.WalletPropose: %s", res.Error)
	}
	return &dtoxrp.WalletInfo{
		AccountID:     res.Result.AccountID,
		KeyType:       res.Result.KeyType,
		MasterKey:     res.Result.MasterKey,
		MasterSeed:    res.Result.MasterSeed,
		MasterSeedHex: res.Result.MasterSeedHex,
		PublicKey:     res.Result.PublicKey,
		PublicKeyHex:  res.Result.PublicKeyHex,
	}, nil
}

// ValidationCreate calls validation_create and converts the result to a DTO.
func (a *AdminXRP) ValidationCreate(ctx context.Context, secret string) (*dtoxrp.ValidationKey, error) {
	res, err := a.AdminRPC.ValidationCreate(ctx, secret)
	if err != nil {
		return nil, fmt.Errorf("fail to call AdminRPC.ValidationCreate: %w", err)
	}
	if res.Error != "" {
		return nil, fmt.Errorf("fail to call AdminRPC.ValidationCreate: %s", res.Error)
	}
	return &dtoxrp.ValidationKey{
		ValidationKey:       res.Result.ValidationKey,
		ValidationPublicKey: res.Result.ValidationPublicKey,
		ValidationSeed:      res.Result.ValidationSeed,
	}, nil
}

// LedgerAccept advances the ledger in standalone mode and converts the result to a DTO.
func (a *AdminXRP) LedgerAccept(ctx context.Context) (*dtoxrp.LedgerAccept, error) {
	res, err := a.AdminRPC.LedgerAccept(ctx)
	if err != nil {
		return nil, fmt.Errorf("fail to call AdminRPC.LedgerAccept: %w", err)
	}
	if res.Error != "" {
		return nil, fmt.Errorf("fail to call AdminRPC.LedgerAccept: %s", res.Error)
	}
	return &dtoxrp.LedgerAccept{
		LedgerCurrentIndex: res.Result.LedgerCurrentIndex,
	}, nil
}
