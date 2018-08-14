package service

import (
	"github.com/bookerzzz/grok"
	"github.com/hiromaily/go-bitcoin/pkg/api"
	"github.com/hiromaily/go-bitcoin/pkg/rds"
	"github.com/hiromaily/go-bitcoin/pkg/toml"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Wallet 基底オブジェクト
type Wallet struct {
	Btc *api.Bitcoin
	DB  *sqlx.DB
	//Db  *kvs.LevelDB
}

//InitialSettings 実行前に必要なすべての設定をこちらで行う
func InitialSettings(confPath string) (*Wallet, error) {
	// Config
	conf, err := toml.New(confPath)
	if err != nil {
		return nil, errors.Errorf("toml.New() error: %v", err)
	}
	grok.Value(conf)

	// KVS
	//db, err := kvs.InitDB(conf.LevelDB.Path)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer db.Close()

	// MySQL
	rds, err := rds.Connection(&conf.MySQL)
	if err != nil {
		return nil, errors.Errorf("rds.Connection() error: %v", err)
	}
	//defer rds.Close()

	// Connection to Bitcoin core
	//bit, err := api.Connection(conf.Bitcoin.Host, conf.Bitcoin.User, conf.Bitcoin.Pass, true, true, conf.Bitcoin.IsMain)
	bit, err := api.Connection(&conf.Bitcoin)
	if err != nil {
		return nil, errors.Errorf("api.Connection error: %v", err)
	}
	//defer bit.Close()

	//Wallet Object
	wallet := Wallet{Btc: bit, DB: rds}
	return &wallet, nil
}

// Done 終了時に必要な処理
func (w *Wallet) Done() {
	w.DB.Close()
	w.Btc.Close()
}
