package kvs

import (
	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
)

// 検証用に追加しただけで使わないかと

//GoDoc
//https://godoc.org/github.com/syndtr/goleveldb/leveldb
//https://www.sambaiz.net/article/45/

// LevelDB connection for LevelDB
type LevelDB struct {
	conn *leveldb.DB
}

// ここに定義したものを常にprefix(仮想テーブル名)として、keyにaddして利用する。
var tableList = map[string]string{
	"unspent": "unspent",
}

// InitDB 接続処理
func InitDB(path string) (*LevelDB, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &LevelDB{conn: db}, nil
}

// Conn connectionオブジェクトを返す
func (d *LevelDB) Conn() *leveldb.DB {
	return d.conn
}

// Close 接続
func (d *LevelDB) Close() {
	d.conn.Close()
}

// SetKey keyを作成する
func (d *LevelDB) SetKey(prefix, key string) ([]byte, error) {
	if _, ok := tableList[prefix]; ok {
		return []byte(prefix + key), nil
	}
	return nil, errors.Errorf("table in leveldb is not found: %s", prefix)
}

// Put Putのwrapper
func (d *LevelDB) Put(prefix, key string, val []byte) error {
	bKey, err := d.SetKey(prefix, key)
	if err != nil {
		return err
	}
	err = d.conn.Put(bKey, val, nil)
	if err != nil {
		return err
	}

	return nil
}

// Get Getのwrapper
func (d *LevelDB) Get(prefix, key string) ([]byte, error) {
	bKey, err := d.SetKey(prefix, key)
	if err != nil {
		return nil, err
	}
	val, err := d.conn.Get(bKey, nil)
	if err != nil {
		return nil, err
	}

	return val, err
}
