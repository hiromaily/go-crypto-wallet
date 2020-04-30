package signature

import (
	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets"
)

// Signature signature wallet object
//  it is almost same to Wallet object, difference is storager interface
type Signature struct {
	wallets.Coldwalleter
	authAccount account.AuthType
}

// NewSignature returns Signature instance
func NewSignature(
	coldWallter wallets.Coldwalleter,
	authAccount account.AuthType) *Signature {

	return &Signature{
		Coldwalleter: coldWallter,
		authAccount:  authAccount,
	}
}

// GetAuthType gets auth_type
func (w *Signature) GetAuthType() account.AuthType {
	return w.authAccount
}
