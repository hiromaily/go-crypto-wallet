package cold

import (
	"context"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainBitcoin "github.com/hiromaily/go-crypto-wallet/internal/domain/bitcoin"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainEth "github.com/hiromaily/go-crypto-wallet/internal/domain/ethereum"
	domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
	domainXrp "github.com/hiromaily/go-crypto-wallet/internal/domain/xrp"
)

// BTCAccountKeyRepositorier is BtcAccountKeyRepository interface for BTC/BCH
type BTCAccountKeyRepositorier interface {
	GetMaxIndex(ctx context.Context, accountType domainAccount.AccountType) (int64, error)
	GetOneMaxID(accountType domainAccount.AccountType) (*domainBitcoin.BTCAccountKey, error)
	GetAllAddrStatus(
		accountType domainAccount.AccountType, addrStatus domainAddress.AddrStatus,
	) ([]*domainBitcoin.BTCAccountKey, error)
	GetAllMultiAddr(accountType domainAccount.AccountType, addrs []string) ([]*domainBitcoin.BTCAccountKey, error)
	InsertBulk(items []*domainBitcoin.BTCAccountKey) error
	UpdateAddr(
		accountType domainAccount.AccountType, addr, keyAddress string,
	) (int64, error)
	UpdateAddrStatus(
		accountType domainAccount.AccountType, addrStatus domainAddress.AddrStatus, strWIFs []string,
	) (int64, error)
	UpdateMultisigAddr(accountType domainAccount.AccountType, item *domainBitcoin.BTCAccountKey) (int64, error)
	UpdateMultisigAddrs(accountType domainAccount.AccountType, items []*domainBitcoin.BTCAccountKey) (int64, error)
}

// ETHAccountKeyRepositorier is EthAccountKeyRepository interface for ETH
type ETHAccountKeyRepositorier interface {
	GetMaxIndex(ctx context.Context, accountType domainAccount.AccountType) (int64, error)
	GetOneMaxID(accountType domainAccount.AccountType) (*domainEth.ETHAccountKey, error)
	GetAllAddrStatus(
		accountType domainAccount.AccountType, addrStatus domainAddress.AddrStatus,
	) ([]*domainEth.ETHAccountKey, error)
	GetByAddress(address string) (*domainEth.ETHAccountKey, error)
	InsertBulk(items []*domainEth.ETHAccountKey) error
	UpdateAddrStatus(
		accountType domainAccount.AccountType, addrStatus domainAddress.AddrStatus, privateKeys []string,
	) (int64, error)
}

// XRPAccountKeyRepositorier is XRPAccountKeyRepository interface
type XRPAccountKeyRepositorier interface {
	GetAllAddrStatus(
		ctx context.Context, accountType domainAccount.AccountType, addrStatus domainAddress.AddrStatus,
	) ([]*domainXrp.XRPAccountKey, error)
	GetSecret(ctx context.Context, accountType domainAccount.AccountType, addr string) (string, error)
	InsertBulk(ctx context.Context, items []*domainXrp.XRPAccountKey) error
	UpdateAddrStatus(
		ctx context.Context,
		accountType domainAccount.AccountType,
		addrStatus domainAddress.AddrStatus,
		strWIFs []string,
	) (int64, error)
}

// XRPRegularKeyRepositorier is the repository interface for XRP regular key management.
// Regular keys allow accounts to sign transactions without using the master key.
// Reference: https://xrpl.org/docs/concepts/accounts/cryptographic-keys#regular-key-pair
type XRPRegularKeyRepositorier interface {
	// GetByAccountID returns the active regular key for an account
	GetByAccountID(ctx context.Context, accountID string) (*domainXrp.XRPRegularKey, error)
	// GetAllByAccountID returns all regular keys (active and inactive) for an account
	GetAllByAccountID(ctx context.Context, accountID string) ([]*domainXrp.XRPRegularKey, error)
	// GetActiveKeys returns all currently active regular keys
	GetActiveKeys(ctx context.Context) ([]*domainXrp.XRPRegularKey, error)
	// GetByAddress returns a regular key by its address
	GetByAddress(ctx context.Context, regularKeyAddress string) (*domainXrp.XRPRegularKey, error)
	// Insert creates a new regular key record
	Insert(ctx context.Context, key *domainXrp.XRPRegularKey) (int64, error)
	// UpdateStatus updates the active status and rotated_at timestamp
	UpdateStatus(ctx context.Context, id int64, isActive bool) error
	// DeactivateByAccountID deactivates all regular keys for an account
	DeactivateByAccountID(ctx context.Context, accountID string) error
	// UpdateTxHash updates the SetRegularKey transaction hash
	UpdateTxHash(ctx context.Context, id int64, txHash string) error
}

// HDWalletRepo is an interface for HD wallet key storage operations.
// It abstracts over key storage for different account types (e.g., regular accounts
// and authorization accounts), allowing the same use case code to work with either.
type HDWalletRepo interface {
	GetMaxIndex(ctx context.Context, accountType domainAccount.AccountType) (int64, error)
	Insert(
		keys []domainKey.WalletKey,
		accountXpriv string,
		idx int64,
		coinTypeCode domainCoin.CoinTypeCode,
		accountType domainAccount.AccountType,
		keyType domainKey.KeyType,
	) error
}
