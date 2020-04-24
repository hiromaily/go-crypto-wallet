package repository

import (
	models "github.com/hiromaily/go-bitcoin/pkg/models/rdb"
	"github.com/hiromaily/go-bitcoin/pkg/tx"
)

type TxRepository interface {
	GetOne(id uint64) (*models.TX, error)
	GetCount(hex string) (int64, error)
	GetTxID(hash string) (uint64, error)
	GetSentHashTx(txType tx.TxType) ([]string, error)
	InsertUnsignedTx(txItem *models.TX) error
	UpdateTx(txItem *models.TX) (int64, error)
}
