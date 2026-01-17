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
	GetOneMaxID(accountType domainAccount.AccountType) (*domainBitcoin.BtcAccountKey, error)
	GetAllAddrStatus(
		accountType domainAccount.AccountType, addrStatus domainAddress.AddrStatus,
	) ([]*domainBitcoin.BtcAccountKey, error)
	GetAllMultiAddr(accountType domainAccount.AccountType, addrs []string) ([]*domainBitcoin.BtcAccountKey, error)
	InsertBulk(items []*domainBitcoin.BtcAccountKey) error
	UpdateAddr(
		accountType domainAccount.AccountType, addr, keyAddress string,
	) (int64, error)
	UpdateAddrStatus(
		accountType domainAccount.AccountType, addrStatus domainAddress.AddrStatus, strWIFs []string,
	) (int64, error)
	UpdateMultisigAddr(accountType domainAccount.AccountType, item *domainBitcoin.BtcAccountKey) (int64, error)
	UpdateMultisigAddrs(accountType domainAccount.AccountType, items []*domainBitcoin.BtcAccountKey) (int64, error)
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
