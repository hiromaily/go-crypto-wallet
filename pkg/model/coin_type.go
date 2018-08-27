package model

import (
	"github.com/hiromaily/go-bitcoin/pkg/key"
	"time"
)

//CoinType coin_typeテーブル
type CoinType struct {
	ID          uint8      `db:"id"`
	Type        string     `db:"type"`
	Description string     `db:"description"`
	UpdatedAt   *time.Time `db:"updated_at"`
}

// GetCoinTypeByID 該当するIDのレコードを返す
func (m *DB) GetCoinTypeByID(coinType key.CoinType) (*CoinType, error) {
	sql := "SELECT * FROM coin_type WHERE id=?"

	ct := CoinType{}
	err := m.RDB.Get(&ct, sql, coinType)

	return &ct, err
}
