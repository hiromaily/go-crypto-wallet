package repository

import (
	"database/sql"

	domainTx "github.com/hiromaily/go-crypto-wallet/internal/domain/transaction"
)

// TxRepositorier is TxRepository interface
type TxRepositorier interface {
	GetOne(id int64) (*domainTx.Transaction, error)
	GetMaxID(actionType domainTx.ActionType) (int64, error)
	InsertUnsignedTx(actionType domainTx.ActionType) (int64, error)
	Update(txItem *domainTx.Transaction) (int64, error)
	DeleteAll() (int64, error)
	WithTx(tx *sql.Tx) TxRepositorier
}
