package rdb

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/hiromaily/go-bitcoin/pkg/config"
)

//sqlx
// http://jmoiron.github.io/sqlx/#bindvars

//sqlboiler
// https://github.com/volatiletech/sqlboiler

// NewMySQL connect to MySQL server
// TODO:
//  - retry functionality and retry count should be configured in config file
//  - change sqlx.DB to basic one because it would be replaced sqlboiler
func NewMySQL(conf *config.MySQL) (*sql.DB, error) {
	db, err := sql.Open("mysql",
		fmt.Sprintf(
			"%s:%s@tcp(%s)/%s?parseTime=true&charset=utf8mb4",
			conf.User,
			conf.Pass,
			conf.Host,
			conf.DB))

	if err != nil {
		return nil, errors.Errorf("Connection(): error: %v", err)
	}
	return db, nil
}

// NewMySQLX connect to MySQL server
// TODO:
//  - retry functionality and retry count should be configured in config file
//  - change sqlx.DB to basic one because it would be replaced sqlboiler
func NewMySQLX(conf *config.MySQL) (*sqlx.DB, error) {
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
