package kvs

import (
	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
)

//GoDoc
//https://godoc.org/github.com/syndtr/goleveldb/leveldb
//https://www.sambaiz.net/article/45/

// LevelDB connection for LevelDB
type LevelDB struct {
	conn *leveldb.DB
}

var tableList = map[string]string{
	"unspent": "unspent",
}

func InitDB(path string) (*LevelDB, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &LevelDB{conn: db}, nil
}

func (d *LevelDB) Conn() *leveldb.DB {
	return d.conn
}

func (d *LevelDB) Close() {
	d.conn.Close()
}

func (d *LevelDB) SetKey(prefix, key string) ([]byte, error) {
	if _, ok := tableList[prefix]; ok {
		return []byte(prefix + key), nil
	}
	return nil, errors.Errorf("table in leveldb is not found: %s", prefix)
}

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
