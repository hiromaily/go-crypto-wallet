package watchrepo

import (
	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/action"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/tx"
)

// AddressRepositorier is AddressRepository interface
type AddressRepositorier interface {
	GetAll(accountType account.AccountType) ([]*models.Address, error)
	GetAllAddress(accountType account.AccountType) ([]string, error)
	GetOneUnAllocated(accountType account.AccountType) (*models.Address, error)
	InsertBulk(items []*models.Address) error
	UpdateIsAllocated(isAllocated bool, Address string) (int64, error)
}

// BTCTxRepositorier is BTCTxRepository interface
type BTCTxRepositorier interface {
	GetOne(id int64) (*models.BTCTX, error)
	GetCountByUnsignedHex(actionType action.ActionType, hex string) (int64, error)
	GetTxIDBySentHash(actionType action.ActionType, hash string) (int64, error)
	GetSentHashTx(actionType action.ActionType, txType tx.TxType) ([]string, error)
	InsertUnsignedTx(actionType action.ActionType, txItem *models.BTCTX) (int64, error)
	Update(txItem *models.BTCTX) (int64, error)
	UpdateAfterTxSent(txID int64, txType tx.TxType, signedHex, sentHashTx string) (int64, error)
	UpdateTxType(id int64, txType tx.TxType) (int64, error)
	UpdateTxTypeBySentHashTx(actionType action.ActionType, txType tx.TxType, sentHashTx string) (int64, error)
	DeleteAll() (int64, error)
}

// TxInputRepositorier is TxInputRepository interface
type TxInputRepositorier interface {
	GetOne(id int64) (*models.BTCTXInput, error)
	GetAllByTxID(id int64) ([]*models.BTCTXInput, error)
	Insert(txItem *models.BTCTXInput) error
	InsertBulk(txItems []*models.BTCTXInput) error
}

// TxOutputRepositorier is TxOutputRepository interface
type TxOutputRepositorier interface {
	GetOne(id int64) (*models.BTCTXOutput, error)
	GetAllByTxID(id int64) ([]*models.BTCTXOutput, error)
	Insert(txItem *models.BTCTXOutput) error
	InsertBulk(txItems []*models.BTCTXOutput) error
}

// TxRepositorier is TxRepository interface
type TxRepositorier interface {
	GetOne(id int64) (*models.TX, error)
	GetMaxID(actionType action.ActionType) (int64, error)
	InsertUnsignedTx(actionType action.ActionType) (int64, error)
	Update(txItem *models.TX) (int64, error)
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
	GetSentHashTx(txType tx.TxType) ([]string, error)
	Insert(txItem *models.EthDetailTX) error
	InsertBulk(txItems []*models.EthDetailTX) error
	UpdateAfterTxSent(uuid string, txType tx.TxType, signedHex, sentHashTx string) (int64, error)
	UpdateTxType(id int64, txType tx.TxType) (int64, error)
	UpdateTxTypeBySentHashTx(txType tx.TxType, sentHashTx string) (int64, error)
}

// XrpDetailTxRepositorier is XrpDetailTxRepository interface
type XrpDetailTxRepositorier interface {
	GetOne(id int64) (*models.XRPDetailTX, error)
	GetAllByTxID(id int64) ([]*models.XRPDetailTX, error)
	GetSentHashTx(txType tx.TxType) ([]string, error)
	Insert(txItem *models.XRPDetailTX) error
	InsertBulk(txItems []*models.XRPDetailTX) error
	UpdateAfterTxSent(
		uuid string, txType tx.TxType, signedTxID, signedTxBlob string, earlistLedgerVersion uint64,
	) (int64, error)
	UpdateTxType(id int64, txType tx.TxType) (int64, error)
	UpdateTxTypeBySentHashTx(txType tx.TxType, sentHashTx string) (int64, error)
}
