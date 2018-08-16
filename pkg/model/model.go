package model

import (
	"github.com/jmoiron/sqlx"
)

// Model データベースオブジェクト
// 複数のDBがある場合、こちらを拡張していく
type DB struct {
	RDB *sqlx.DB
}

// New DBオブジェクトを返す
func New(d *sqlx.DB) *DB {
	db := DB{RDB: d}
	return &db
}
