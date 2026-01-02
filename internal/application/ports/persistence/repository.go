// Package persistence defines interfaces for data persistence operations.
//
// This package follows the Dependency Inversion Principle of Clean Architecture
// by defining interfaces in the application layer that are implemented by the
// infrastructure layer.
package persistence

import (
	"context"
	"database/sql"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/mysql/sqlcgen"
)

// Repository interfaces for cold wallet (keygen and sign wallets)

// SeedRepositorier is SeedRepository interface
type SeedRepositorier interface {
	GetOne(ctx context.Context) (*sqlcgen.Seed, error)
	Insert(ctx context.Context, strSeed string) error
}

// BTCAccountKeyRepositorier is BtcAccountKeyRepository interface for BTC/BCH
type BTCAccountKeyRepositorier interface {
	GetMaxIndex(accountType domainAccount.AccountType) (int64, error)
	GetOneMaxID(accountType domainAccount.AccountType) (*sqlcgen.BtcAccountKey, error)
	GetAllAddrStatus(
		accountType domainAccount.AccountType, addrStatus domainAddress.AddrStatus,
	) ([]*sqlcgen.BtcAccountKey, error)
	GetAllMultiAddr(accountType domainAccount.AccountType, addrs []string) ([]*sqlcgen.BtcAccountKey, error)
	InsertBulk(items []*sqlcgen.BtcAccountKey) error
	UpdateAddr(
		accountType domainAccount.AccountType, addr, keyAddress string,
	) (int64, error)
	UpdateAddrStatus(
		accountType domainAccount.AccountType, addrStatus domainAddress.AddrStatus, strWIFs []string,
	) (int64, error)
	UpdateMultisigAddr(accountType domainAccount.AccountType, item *sqlcgen.BtcAccountKey) (int64, error)
	UpdateMultisigAddrs(accountType domainAccount.AccountType, items []*sqlcgen.BtcAccountKey) (int64, error)
}

// ETHAccountKeyRepositorier is EthAccountKeyRepository interface for ETH
type ETHAccountKeyRepositorier interface {
	GetMaxIndex(accountType domainAccount.AccountType) (int64, error)
	GetOneMaxID(accountType domainAccount.AccountType) (*sqlcgen.EthAccountKey, error)
	GetAllAddrStatus(
		accountType domainAccount.AccountType, addrStatus domainAddress.AddrStatus,
	) ([]*sqlcgen.EthAccountKey, error)
	GetByAddress(address string) (*sqlcgen.EthAccountKey, error)
	InsertBulk(items []*sqlcgen.EthAccountKey) error
	UpdateAddrStatus(
		accountType domainAccount.AccountType, addrStatus domainAddress.AddrStatus, privateKeys []string,
	) (int64, error)
}

// XRPAccountKeyRepositorier is XRPAccountKeyRepository interface
type XRPAccountKeyRepositorier interface {
	GetAllAddrStatus(
		ctx context.Context, accountType domainAccount.AccountType, addrStatus domainAddress.AddrStatus,
	) ([]*sqlcgen.XrpAccountKey, error)
	GetSecret(ctx context.Context, accountType domainAccount.AccountType, addr string) (string, error)
	InsertBulk(ctx context.Context, items []*sqlcgen.XrpAccountKey) error
	UpdateAddrStatus(
		ctx context.Context,
		accountType domainAccount.AccountType,
		addrStatus domainAddress.AddrStatus,
		strWIFs []string,
	) (int64, error)
}

// AuthFullPubkeyRepositorier is AuthFullPubkeyRepository interface
type AuthFullPubkeyRepositorier interface {
	GetOne(authType domainAccount.AuthType) (*sqlcgen.AuthFullpubkey, error)
	Insert(authType domainAccount.AuthType, fullPubKey string) error
	InsertBulk(items []*sqlcgen.AuthFullpubkey) error
}

// AuthAccountKeyRepositorier is AuthAccountKeyRepository interface
type AuthAccountKeyRepositorier interface {
	GetOne(authType domainAccount.AuthType) (*sqlcgen.AuthAccountKey, error)
	Insert(item *sqlcgen.AuthAccountKey) error
	UpdateAddrStatus(addrStatus domainAddress.AddrStatus, strWIF string) (int64, error)
}

// Repository interfaces for watch wallet

// AddressRepositorier is AddressRepository interface
type AddressRepositorier interface {
	GetAll(accountType domainAccount.AccountType) ([]*sqlcgen.Address, error)
	GetAllAddress(accountType domainAccount.AccountType) ([]string, error)
	GetOneUnAllocated(accountType domainAccount.AccountType) (*sqlcgen.Address, error)
	InsertBulk(ctx context.Context, items []*sqlcgen.Address) error
	UpdateIsAllocated(isAllocated bool, Address string) (int64, error)
}

// BTCTxRepositorier is BTCTxRepository interface
type BTCTxRepositorier interface {
	GetOne(id int64) (*sqlcgen.BtcTx, error)
	GetCountByUnsignedHex(actionType domainTx.ActionType, hex string) (int64, error)
	GetTxIDBySentHash(actionType domainTx.ActionType, hash string) (int64, error)
	GetSentHashTx(actionType domainTx.ActionType, txType domainTx.TxType) ([]string, error)
	InsertUnsignedTx(actionType domainTx.ActionType, txItem *sqlcgen.BtcTx) (int64, error)
	Update(txItem *sqlcgen.BtcTx) (int64, error)
	UpdateAfterTxSent(txID int64, txType domainTx.TxType, signedHex, sentHashTx string) (int64, error)
	UpdateTxType(id int64, txType domainTx.TxType) (int64, error)
	UpdateTxTypeBySentHashTx(actionType domainTx.ActionType, txType domainTx.TxType, sentHashTx string) (int64, error)
	DeleteAll() (int64, error)
}

// TxInputRepositorier is TxInputRepository interface
type TxInputRepositorier interface {
	GetOne(id int64) (*sqlcgen.BtcTxInput, error)
	GetAllByTxID(id int64) ([]*sqlcgen.BtcTxInput, error)
	Insert(txItem *sqlcgen.BtcTxInput) error
	InsertBulk(txItems []*sqlcgen.BtcTxInput) error
}

// TxOutputRepositorier is TxOutputRepository interface
type TxOutputRepositorier interface {
	GetOne(id int64) (*sqlcgen.BtcTxOutput, error)
	GetAllByTxID(id int64) ([]*sqlcgen.BtcTxOutput, error)
	Insert(txItem *sqlcgen.BtcTxOutput) error
	InsertBulk(txItems []*sqlcgen.BtcTxOutput) error
}

// TxRepositorier is TxRepository interface
type TxRepositorier interface {
	GetOne(id int64) (*sqlcgen.Tx, error)
	GetMaxID(actionType domainTx.ActionType) (int64, error)
	InsertUnsignedTx(actionType domainTx.ActionType) (int64, error)
	Update(txItem *sqlcgen.Tx) (int64, error)
	DeleteAll() (int64, error)
	WithTx(tx *sql.Tx) TxRepositorier
}

// PaymentRequestRepositorier is PaymentRequestRepository interface
type PaymentRequestRepositorier interface {
	GetAll() ([]*sqlcgen.PaymentRequest, error)
	GetAllByPaymentID(paymentID int64) ([]*sqlcgen.PaymentRequest, error)
	InsertBulk(items []*sqlcgen.PaymentRequest) error
	UpdatePaymentID(paymentID int64, ids []int64) (int64, error)
	UpdateIsDone(paymentID int64) (int64, error)
	DeleteAll() (int64, error)
	WithTx(tx *sql.Tx) PaymentRequestRepositorier
}

// ETHDetailTXRepositorier is ETHDetailTXRepository interface
type ETHDetailTXRepositorier interface {
	GetOne(id int64) (*sqlcgen.EthDetailTx, error)
	GetAllByTxID(id int64) ([]*sqlcgen.EthDetailTx, error)
	GetSentHashTx(txType domainTx.TxType) ([]string, error)
	Insert(txItem *sqlcgen.EthDetailTx) error
	InsertBulk(txItems []*sqlcgen.EthDetailTx) error
	UpdateAfterTxSent(uuid string, txType domainTx.TxType, signedHex, sentHashTx string) (int64, error)
	UpdateTxType(id int64, txType domainTx.TxType) (int64, error)
	UpdateTxTypeBySentHashTx(txType domainTx.TxType, sentHashTx string) (int64, error)
	WithTx(tx *sql.Tx) ETHDetailTXRepositorier
}

// XRPDetailTXRepositorier is XrpDetailTxRepository interface
type XRPDetailTXRepositorier interface {
	GetOne(id int64) (*sqlcgen.XrpDetailTx, error)
	GetAllByTxID(id int64) ([]*sqlcgen.XrpDetailTx, error)
	GetSentHashTx(txType domainTx.TxType) ([]string, error)
	Insert(txItem *sqlcgen.XrpDetailTx) error
	InsertBulk(txItems []*sqlcgen.XrpDetailTx) error
	UpdateAfterTxSent(
		uuid string, txType domainTx.TxType, signedTxID, signedTxBlob string, earlistLedgerVersion uint64,
	) (int64, error)
	UpdateTxType(id int64, txType domainTx.TxType) (int64, error)
	UpdateTxTypeBySentHashTx(txType domainTx.TxType, sentHashTx string) (int64, error)
	WithTx(tx *sql.Tx) XRPDetailTXRepositorier
}
