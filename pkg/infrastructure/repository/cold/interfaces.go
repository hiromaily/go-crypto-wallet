package cold

import (
	"github.com/hiromaily/go-crypto-wallet/pkg/address"
	domainAccount "github.com/hiromaily/go-crypto-wallet/pkg/domain/account"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
)

// SeedRepositorier is SeedRepository interface
type SeedRepositorier interface {
	GetOne() (*models.Seed, error)
	Insert(strSeed string) error
}

// AccountKeyRepositorier is AccountKeyRepository interface
type AccountKeyRepositorier interface {
	GetMaxIndex(accountType domainAccount.AccountType) (int64, error)
	GetOneMaxID(accountType domainAccount.AccountType) (*models.AccountKey, error)
	GetAllAddrStatus(accountType domainAccount.AccountType, addrStatus address.AddrStatus) ([]*models.AccountKey, error)
	GetAllMultiAddr(accountType domainAccount.AccountType, addrs []string) ([]*models.AccountKey, error)
	InsertBulk(items []*models.AccountKey) error
	UpdateAddr(
		accountType domainAccount.AccountType, addr, keyAddress string,
	) (int64, error)
	UpdateAddrStatus(
		accountType domainAccount.AccountType, addrStatus address.AddrStatus, strWIFs []string,
	) (int64, error)
	UpdateMultisigAddr(accountType domainAccount.AccountType, item *models.AccountKey) (int64, error)
	UpdateMultisigAddrs(accountType domainAccount.AccountType, items []*models.AccountKey) (int64, error)
}

// XRPAccountKeyRepositorier is XRPAccountKeyRepository interface
type XRPAccountKeyRepositorier interface {
	GetAllAddrStatus(
		accountType domainAccount.AccountType, addrStatus address.AddrStatus,
	) ([]*models.XRPAccountKey, error)
	GetSecret(accountType domainAccount.AccountType, addr string) (string, error)
	InsertBulk(items []*models.XRPAccountKey) error
	UpdateAddrStatus(
		accountType domainAccount.AccountType, addrStatus address.AddrStatus, strWIFs []string,
	) (int64, error)
}

// AuthFullPubkeyRepositorier is AuthFullPubkeyRepository interface
type AuthFullPubkeyRepositorier interface {
	GetOne(authType domainAccount.AuthType) (*models.AuthFullpubkey, error)
	Insert(authType domainAccount.AuthType, fullPubKey string) error
	InsertBulk(items []*models.AuthFullpubkey) error
}

// AuthAccountKeyRepositorier is AuthAccountKeyRepository interface
type AuthAccountKeyRepositorier interface {
	GetOne(authType domainAccount.AuthType) (*models.AuthAccountKey, error)
	Insert(item *models.AuthAccountKey) error
	UpdateAddrStatus(addrStatus address.AddrStatus, strWIF string) (int64, error)
}

// GetRedeemScriptByAddress returns redeem script by address
func GetRedeemScriptByAddress(accountKeys []*models.AccountKey, addr string) string {
	for _, val := range accountKeys {
		if val.MultisigAddress == addr {
			return val.RedeemScript
		}
	}
	return ""
}
