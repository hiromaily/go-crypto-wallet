package rdb

import (
	"github.com/jmoiron/sqlx"
)

// DB データベースオブジェクト
// 複数のDBがある場合、こちらを拡張していく
type DB struct {
	RDB *sqlx.DB
}

// NewDB DBオブジェクトを返す
func NewDB(d *sqlx.DB) *DB {
	db := DB{RDB: d}
	return &db
}
