package repository

import (
	"database/sql"

	domainEth "github.com/hiromaily/go-crypto-wallet/internal/domain/ethereum"
	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
)

// ETHDetailTXRepositorier is ETHDetailTXRepository interface
type ETHDetailTXRepositorier interface {
	GetOne(id int64) (*domainEth.EthDetailTx, error)
	GetAllByTxID(id int64) ([]*domainEth.EthDetailTx, error)
	GetSentHashTx(txType domainTx.TxType) ([]string, error)
	Insert(txItem *domainEth.EthDetailTx) error
	InsertBulk(txItems []*domainEth.EthDetailTx) error
	UpdateAfterTxSent(uuid string, txType domainTx.TxType, signedHex, sentHashTx string) (int64, error)
	UpdateTxType(id int64, txType domainTx.TxType) (int64, error)
	UpdateTxTypeBySentHashTx(txType domainTx.TxType, sentHashTx string) (int64, error)
	WithTx(tx *sql.Tx) ETHDetailTXRepositorier
}
