package model

import (
	"github.com/jmoiron/sqlx"
)

// Model データベースオブジェクト
type DB struct {
	DB *sqlx.DB
}

// New DBオブジェクトを返す
func New(d *sqlx.DB) *DB {
	db := DB{DB: d}
	return &db
}
