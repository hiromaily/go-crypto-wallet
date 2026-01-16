package repository

import (
	"database/sql"

	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	domainXrp "github.com/hiromaily/go-crypto-wallet/internal/domain/xrp"
)

// XRPDetailTXRepositorier is XrpDetailTxRepository interface
type XRPDetailTXRepositorier interface {
	GetOne(id int64) (*domainXrp.XrpDetailTx, error)
	GetAllByTxID(id int64) ([]*domainXrp.XrpDetailTx, error)
	GetSentHashTx(txType domainTx.TxType) ([]string, error)
	Insert(txItem *domainXrp.XrpDetailTx) error
	InsertBulk(txItems []*domainXrp.XrpDetailTx) error
	UpdateAfterTxSent(
		uuid string, txType domainTx.TxType, signedTxID, signedTxBlob string, earlistLedgerVersion uint64,
	) (int64, error)
	UpdateTxType(id int64, txType domainTx.TxType) (int64, error)
	UpdateTxTypeBySentHashTx(txType domainTx.TxType, sentHashTx string) (int64, error)
	WithTx(tx *sql.Tx) XRPDetailTXRepositorier
}
