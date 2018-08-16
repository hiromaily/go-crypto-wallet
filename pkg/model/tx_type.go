package model

import (
	"time"
)

type TxType struct {
	ID        int        `db:"id"`
	Type      string     `db:"type"`
	UpdatedAt *time.Time `db:"updated_at"`
}

// GetTxType GetTxTypeテーブル全体を返す
func (m *DB) GetTxType() ([]TxType, error) {
	txTypes := []TxType{}
	err := m.RDB.Select(&txTypes, "SELECT * FROM tx_type")

	return txTypes, err
}
