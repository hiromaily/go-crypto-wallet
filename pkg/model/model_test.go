package model_test

import (
	"flag"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	. "github.com/hiromaily/go-bitcoin/pkg/model"
	"github.com/hiromaily/go-bitcoin/pkg/rdb"
	"github.com/hiromaily/go-bitcoin/pkg/toml"
)

var (
	db       *DB
	confPath = flag.String("conf", "../../data/toml/local_watch_only.toml", "Path for configuration toml file")
)

func setup() {
	// RDB
	if db != nil {
		return
	}
	flag.Parse()

	conf, err := toml.New(*confPath)
	if err != nil {
		panic(err)
	}

	rds, err := rdb.Connection(&conf.MySQL)
	if err != nil {
		panic(err)
	}

	db = NewDB(rds)
}

func teardown() {
	db.RDB.Close()
}

func TestMain(m *testing.M) {
	setup()

	code := m.Run()

	teardown()

	os.Exit(code)
}
