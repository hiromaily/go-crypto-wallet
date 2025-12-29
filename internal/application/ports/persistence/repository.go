// Package persistence defines interfaces for data persistence operations.
//
// This package follows the Dependency Inversion Principle of Clean Architecture
// by defining interfaces in the application layer that are implemented by the
// infrastructure layer.
package persistence

import (
	"context"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	models "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/sqlc"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/address"
)

// Repository interfaces for cold wallet (keygen and sign wallets)

// SeedRepositorier is SeedRepository interface
type SeedRepositorier interface {
	GetOne(ctx context.Context) (*sqlc.Seed, error)
	Insert(ctx context.Context, strSeed string) error
}

// AccountKeyRepositorier is AccountKeyRepository interface
type AccountKeyRepositorier interface {
	GetMaxIndex(accountType domainAccount.AccountType) (int64, error)
	GetOneMaxID(accountType domainAccount.AccountType) (*sqlc.AccountKey, error)
	GetAllAddrStatus(accountType domainAccount.AccountType, addrStatus address.AddrStatus) ([]*sqlc.AccountKey, error)
	GetAllMultiAddr(accountType domainAccount.AccountType, addrs []string) ([]*sqlc.AccountKey, error)
	InsertBulk(items []*sqlc.AccountKey) error
	UpdateAddr(
		accountType domainAccount.AccountType, addr, keyAddress string,
	) (int64, error)
	UpdateAddrStatus(
		accountType domainAccount.AccountType, addrStatus address.AddrStatus, strWIFs []string,
	) (int64, error)
	UpdateMultisigAddr(accountType domainAccount.AccountType, item *sqlc.AccountKey) (int64, error)
	UpdateMultisigAddrs(accountType domainAccount.AccountType, items []*sqlc.AccountKey) (int64, error)
}

// XRPAccountKeyRepositorier is XRPAccountKeyRepository interface
type XRPAccountKeyRepositorier interface {
	GetAllAddrStatus(
		ctx context.Context, accountType domainAccount.AccountType, addrStatus address.AddrStatus,
	) ([]*sqlc.XrpAccountKey, error)
	GetSecret(ctx context.Context, accountType domainAccount.AccountType, addr string) (string, error)
	InsertBulk(ctx context.Context, items []*sqlc.XrpAccountKey) error
	UpdateAddrStatus(
		ctx context.Context, accountType domainAccount.AccountType, addrStatus address.AddrStatus, strWIFs []string,
	) (int64, error)
}

// AuthFullPubkeyRepositorier is AuthFullPubkeyRepository interface
type AuthFullPubkeyRepositorier interface {
	GetOne(authType domainAccount.AuthType) (*sqlc.AuthFullpubkey, error)
	Insert(authType domainAccount.AuthType, fullPubKey string) error
	InsertBulk(items []*sqlc.AuthFullpubkey) error
}

// AuthAccountKeyRepositorier is AuthAccountKeyRepository interface
type AuthAccountKeyRepositorier interface {
	GetOne(authType domainAccount.AuthType) (*sqlc.AuthAccountKey, error)
	Insert(item *sqlc.AuthAccountKey) error
	UpdateAddrStatus(addrStatus address.AddrStatus, strWIF string) (int64, error)
}

// Repository interfaces for watch wallet

// AddressRepositorier is AddressRepository interface
type AddressRepositorier interface {
	GetAll(accountType domainAccount.AccountType) ([]*sqlc.Address, error)
	GetAllAddress(accountType domainAccount.AccountType) ([]string, error)
	GetOneUnAllocated(accountType domainAccount.AccountType) (*sqlc.Address, error)
	InsertBulk(ctx context.Context, items []*sqlc.Address) error
	UpdateIsAllocated(isAllocated bool, Address string) (int64, error)
}

// BTCTxRepositorier is BTCTxRepository interface
type BTCTxRepositorier interface {
	GetOne(id int64) (*sqlc.BtcTx, error)
	GetCountByUnsignedHex(actionType domainTx.ActionType, hex string) (int64, error)
	GetTxIDBySentHash(actionType domainTx.ActionType, hash string) (int64, error)
	GetSentHashTx(actionType domainTx.ActionType, txType domainTx.TxType) ([]string, error)
	InsertUnsignedTx(actionType domainTx.ActionType, txItem *sqlc.BtcTx) (int64, error)
	Update(txItem *sqlc.BtcTx) (int64, error)
	UpdateAfterTxSent(txID int64, txType domainTx.TxType, signedHex, sentHashTx string) (int64, error)
	UpdateTxType(id int64, txType domainTx.TxType) (int64, error)
	UpdateTxTypeBySentHashTx(actionType domainTx.ActionType, txType domainTx.TxType, sentHashTx string) (int64, error)
	DeleteAll() (int64, error)
}

// TxInputRepositorier is TxInputRepository interface
type TxInputRepositorier interface {
	GetOne(id int64) (*sqlc.BtcTxInput, error)
	GetAllByTxID(id int64) ([]*sqlc.BtcTxInput, error)
	Insert(txItem *sqlc.BtcTxInput) error
	InsertBulk(txItems []*sqlc.BtcTxInput) error
}

// TxOutputRepositorier is TxOutputRepository interface
type TxOutputRepositorier interface {
	GetOne(id int64) (*sqlc.BtcTxOutput, error)
	GetAllByTxID(id int64) ([]*sqlc.BtcTxOutput, error)
	Insert(txItem *sqlc.BtcTxOutput) error
	InsertBulk(txItems []*sqlc.BtcTxOutput) error
}

// TxRepositorier is TxRepository interface
type TxRepositorier interface {
	GetOne(id int64) (*sqlc.Tx, error)
	GetMaxID(actionType domainTx.ActionType) (int64, error)
	InsertUnsignedTx(actionType domainTx.ActionType) (int64, error)
	Update(txItem *sqlc.Tx) (int64, error)
	DeleteAll() (int64, error)
}

// PaymentRequestRepositorier is PaymentRequestRepository interface
type PaymentRequestRepositorier interface {
	GetAll() ([]*models.PaymentRequest, error)
	GetAllByPaymentID(paymentID int64) ([]*models.PaymentRequest, error)
	InsertBulk(items []*models.PaymentRequest) error
	UpdatePaymentID(paymentID int64, ids []int64) (int64, error)
	UpdateIsDone(paymentID int64) (int64, error)
	DeleteAll() (int64, error)
}

// EthDetailTxRepositorier is EthDetailTxRepository interface
type EthDetailTxRepositorier interface {
	GetOne(id int64) (*models.EthDetailTX, error)
	GetAllByTxID(id int64) ([]*models.EthDetailTX, error)
	GetSentHashTx(txType domainTx.TxType) ([]string, error)
	Insert(txItem *models.EthDetailTX) error
	InsertBulk(txItems []*models.EthDetailTX) error
	UpdateAfterTxSent(uuid string, txType domainTx.TxType, signedHex, sentHashTx string) (int64, error)
	UpdateTxType(id int64, txType domainTx.TxType) (int64, error)
	UpdateTxTypeBySentHashTx(txType domainTx.TxType, sentHashTx string) (int64, error)
}

// XrpDetailTxRepositorier is XrpDetailTxRepository interface
type XrpDetailTxRepositorier interface {
	GetOne(id int64) (*models.XRPDetailTX, error)
	GetAllByTxID(id int64) ([]*models.XRPDetailTX, error)
	GetSentHashTx(txType domainTx.TxType) ([]string, error)
	Insert(txItem *models.XRPDetailTX) error
	InsertBulk(txItems []*models.XRPDetailTX) error
	UpdateAfterTxSent(
		uuid string, txType domainTx.TxType, signedTxID, signedTxBlob string, earlistLedgerVersion uint64,
	) (int64, error)
	UpdateTxType(id int64, txType domainTx.TxType) (int64, error)
	UpdateTxTypeBySentHashTx(txType domainTx.TxType, sentHashTx string) (int64, error)
}
