package service

import (
	"github.com/bookerzzz/grok"
	"github.com/hiromaily/go-bitcoin/pkg/api"
	"github.com/hiromaily/go-bitcoin/pkg/file"
	"github.com/hiromaily/go-bitcoin/pkg/model"
	"github.com/hiromaily/go-bitcoin/pkg/rdb"
	"github.com/hiromaily/go-bitcoin/pkg/toml"
	"github.com/pkg/errors"
)

// Wallet 基底オブジェクト
type Wallet struct {
	BTC *api.Bitcoin
	DB  *model.DB
	//DB  *sqlx.DB
	//Db  *kvs.LevelDB
}

//InitialSettings 実行前に必要なすべての設定をこちらで行う
//TODO:hotwalletとColdwalletで設定が異なるので要調整
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
	rds, err := rdb.Connection(&conf.MySQL)
	if err != nil {
		return nil, errors.Errorf("rds.Connection() error: %v", err)
	}
	//defer rds.Close()

	// File
	if conf.File.ReceiptPath != "" && conf.File.PaymentPath != "" {
		file.SetFilePath(conf.File.ReceiptPath, conf.File.PaymentPath)
	}

	// Connection to Bitcoin core
	//bit, err := api.Connection(conf.Bitcoin.Host, conf.Bitcoin.User, conf.Bitcoin.Pass, true, true, conf.Bitcoin.IsMain)
	bit, err := api.Connection(&conf.Bitcoin)
	if err != nil {
		return nil, errors.Errorf("api.Connection error: %v", err)
	}
	//defer bit.Close()

	//Wallet Object
	wallet := Wallet{BTC: bit, DB: model.NewDB(rds)}
	return &wallet, nil
}

// Done 終了時に必要な処理
func (w *Wallet) Done() {
	w.DB.RDB.Close()
	w.BTC.Close()
}
