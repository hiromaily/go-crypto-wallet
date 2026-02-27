package watch

import (
	domainXRP "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/xrp"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	dbtx "github.com/hiromaily/go-crypto-wallet/pkg/db/tx"
)

// XRPDetailTXRepositorier is XRPDetailTxRepository interface
type XRPDetailTXRepositorier interface {
	GetOne(id int64) (*domainXRP.XRPDetailTx, error)
	GetAllByTxID(id int64) ([]*domainXRP.XRPDetailTx, error)
	GetSentHashTx(txType domainTx.TxType) ([]string, error)
	Insert(txItem *domainXRP.XRPDetailTx) error
	InsertBulk(txItems []*domainXRP.XRPDetailTx) error
	UpdateAfterTxSent(
		uuid string, txType domainTx.TxType, signedTxID, signedTxBlob string, earlistLedgerVersion uint64,
	) (int64, error)
	UpdateTxType(id int64, txType domainTx.TxType) (int64, error)
	UpdateTxTypeBySentHashTx(txType domainTx.TxType, sentHashTx string) (int64, error)
	WithTransaction(tx dbtx.Transaction) (XRPDetailTXRepositorier, error)
}
