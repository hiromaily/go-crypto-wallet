package rdb

import (
	"fmt"

	"github.com/hiromaily/go-bitcoin/pkg/toml"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	_ "github.com/go-sql-driver/mysql"
)

//sqlx
// http://jmoiron.github.io/sqlx/#bindvars

// Connection MySQLサーバーに接続する
// TODO:リトライ機能も必要
func Connection(conf *toml.MySQLConf) (*sqlx.DB, error) {
	db, err := sqlx.Connect("mysql",
		fmt.Sprintf(
			"%s:%s@tcp(%s)/%s?parseTime=true",
			conf.User,
			conf.Pass,
			conf.Host,
			conf.DB))

	if err != nil {
		return nil, errors.Errorf("Connection(): error: %v", err)
	}
	return db, nil
}
