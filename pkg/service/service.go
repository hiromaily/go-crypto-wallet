package service

import (
	"github.com/hiromaily/go-bitcoin/pkg/api"
	"github.com/hiromaily/go-bitcoin/pkg/kvs"
)

// Wallet 基底オブジェクト
type Wallet struct {
	Btc *api.Bitcoin
	Db  *kvs.LevelDB
}
