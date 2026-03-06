package admin

import (
	"context"

	"github.com/hiromaily/go-crypto-wallet/pkg/chains/xrp"
)

// XRPAdminer defines the interface for XRP admin node operations.
// These operations typically require admin access to the XRP node.
type XRPAdminer interface {
	// admin_keygen
	ValidationCreate(ctx context.Context, secret string) (*ResponseValidationCreate, error)
	WalletProposeWithKey(
		ctx context.Context,
		seed string,
		keyType xrp.KeyType,
	) (*ResponseWalletPropose, error)
	WalletPropose(ctx context.Context, passphrase string) (*ResponseWalletPropose, error)
}
