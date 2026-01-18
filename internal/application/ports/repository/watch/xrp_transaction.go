package watch

import (
	"github.com/hiromaily/go-crypto-wallet/internal/application/ports/persistence"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
	domainXrp "github.com/hiromaily/go-crypto-wallet/internal/domain/xrp"
)

// XRPDetailTXRepositorier is XRPDetailTxRepository interface
type XRPDetailTXRepositorier interface {
	GetOne(id int64) (*domainXrp.XRPDetailTx, error)
	GetAllByTxID(id int64) ([]*domainXrp.XRPDetailTx, error)
	GetSentHashTx(txType domainTx.TxType) ([]string, error)
	Insert(txItem *domainXrp.XRPDetailTx) error
	InsertBulk(txItems []*domainXrp.XRPDetailTx) error
	UpdateAfterTxSent(
		uuid string, txType domainTx.TxType, signedTxID, signedTxBlob string, earlistLedgerVersion uint64,
	) (int64, error)
	UpdateTxType(id int64, txType domainTx.TxType) (int64, error)
	UpdateTxTypeBySentHashTx(txType domainTx.TxType, sentHashTx string) (int64, error)
	WithTransaction(tx persistence.Transaction) XRPDetailTXRepositorier
}
