package rdb_test

//import (
//	_ "github.com/go-sql-driver/mysql"
//)

//TODO: test as integration test
//var (
//	db       *DB
//	confPath = flag.String("conf", "../../data/config/btc/wallet.toml", "Path for configuration toml file")
//)
//
//func setup() {
//	// RDB
//	if db != nil {
//		return
//	}
//	flag.Parse()
//
//	conf, err := config.New(*confPath)
//	if err != nil {
//		panic(err)
//	}
//
//	rds, err := rdb.NewMySQL(&conf.MySQL)
//	if err != nil {
//		panic(err)
//	}
//
//	db = NewDB(rds)
//}
//
//func teardown() {
//	db.RDB.Close()
//}
//
//func TestMain(m *testing.M) {
//	setup()
//
//	code := m.Run()
//
//	teardown()
//
//	os.Exit(code)
//}
