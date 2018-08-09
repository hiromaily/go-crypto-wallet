package service

import (
	"github.com/hiromaily/go-bitcoin/pkg/api"
	"github.com/hiromaily/go-bitcoin/pkg/kvs"
)

type Wallet struct {
	Btc *api.Bitcoin
	Db  *kvs.LevelDB
}

//HokanAddress 保管用アドレスだが、これをどこに保持すべきか TODO
const HokanAddress = "2N54KrNdyuAkqvvadqSencgpr9XJZnwFYKW"
