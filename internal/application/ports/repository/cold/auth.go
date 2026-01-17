package cold

import (
	"context"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainAuth "github.com/hiromaily/go-crypto-wallet/internal/domain/auth"
)

// AuthFullPubkeyRepositorier is AuthFullPubkeyRepository interface
type AuthFullPubkeyRepositorier interface {
	GetOne(ctx context.Context, authType domainAccount.AuthType) (*domainAuth.AuthFullPubkey, error)
	GetOneByPurpose(authType domainAccount.AuthType, purpose domainAuth.Purpose) (*domainAuth.AuthFullPubkey, error)
	Insert(authType domainAccount.AuthType, fullPubKey string) error
	InsertBulk(items []*domainAuth.AuthFullPubkey) error
}

// AuthAccountKeyRepositorier is AuthAccountKeyRepository interface
type AuthAccountKeyRepositorier interface {
	GetOne(ctx context.Context, authType domainAccount.AuthType) (*domainAuth.AuthAccountKey, error)
	GetByAccount(
		authType domainAccount.AuthType,
		accountType domainAccount.AccountType,
	) (*domainAuth.AuthAccountKey, error)
	Insert(item *domainAuth.AuthAccountKey) error
	UpdateAddrStatus(addrStatus domainAddress.AddrStatus, strWIF string) (int64, error)
}
