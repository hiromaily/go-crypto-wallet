package model

import (
	"github.com/hiromaily/go-bitcoin/pkg/key"
	"time"
)

//AccountType account_typeテーブル
type AccountType struct {
	ID          uint8      `db:"id"`
	Type        string     `db:"type"`
	Description string     `db:"description"`
	UpdatedAt   *time.Time `db:"updated_at"`
}

// GetAccountTypeByID 該当するIDのレコードを返す
func (m *DB) GetAccountTypeByID(accountType key.AccountType) (*AccountType, error) {
	sql := "SELECT * FROM account_type WHERE id=?"

	at := AccountType{}
	err := m.RDB.Get(&at, sql, accountType)

	return &at, err
}
