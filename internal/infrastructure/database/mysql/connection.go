package mysql

import (
	"database/sql"
	"fmt"

	"github.com/hiromaily/go-crypto-wallet/pkg/config"
)

// NewMySQL connect to MySQL server
// TODO: retry functionality and retry count should be configured in config file
func NewMySQL(conf *config.MySQL) (*sql.DB, error) {
	db, err := sql.Open("mysql",
		fmt.Sprintf(
			"%s:%s@tcp(%s)/%s?parseTime=true&charset=utf8mb4",
			conf.User,
			conf.Pass,
			conf.Host,
			conf.DB))
	if err != nil {
		return nil, fmt.Errorf("Connection(): error: %v", err)
	}
	return db, nil
}
