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

// これらの定義はどこにすべきか
const (
	//HokanAddress 保管用アドレス
	HokanAddress         = "2N54KrNdyuAkqvvadqSencgpr9XJZnwFYKW"
	ConfirmationBlockNum = 6
)
