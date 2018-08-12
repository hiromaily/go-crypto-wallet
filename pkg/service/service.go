package service

import (
	"github.com/hiromaily/go-bitcoin/pkg/api"
	"github.com/jmoiron/sqlx"
)

// Wallet 基底オブジェクト
type Wallet struct {
	Btc *api.Bitcoin
	DB  *sqlx.DB
	//Db  *kvs.LevelDB
}
