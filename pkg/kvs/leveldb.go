package kvs

import (
	"github.com/syndtr/goleveldb/leveldb"
)

// LevelDB connection for LevelDB
type LevelDB struct {
	conn *leveldb.DB
}

func InitDB(path string) (*LevelDB, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &LevelDB{conn: db}, nil
}

func (d *LevelDB) Close() {
	d.conn.Close()
}
