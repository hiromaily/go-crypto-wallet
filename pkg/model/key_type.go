package model

import (
	"time"

	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/logger"
)

//KeyType key_typeテーブル
type KeyType struct {
	ID          uint8      `db:"id"`
	Purpose     uint8      `db:"purpose"`
	CoinType    uint8      `db:"coin_type"`
	AccountType uint8      `db:"account_type"`
	ChangeType  uint8      `db:"change_type"`
	Description string     `db:"description"`
	UpdatedAt   *time.Time `db:"updated_at"`
}

// GetKeyTypeByCoinAndAccountType 該当するIDのレコードを返す
func (m *DB) GetKeyTypeByCoinAndAccountType(coinType enum.CoinType, accountType enum.AccountType) (*KeyType, error) {
	sql := "SELECT * FROM key_type WHERE coin_type=? AND account_type=? LIMIT 1"
	logger.Debugf("sql: %s", sql)

	kt := KeyType{}
	err := m.RDB.Get(&kt, sql, coinType, accountType)

	return &kt, err
}
