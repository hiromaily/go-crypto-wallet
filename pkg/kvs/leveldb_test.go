package kvs_test

import (
	"os"
	"testing"

	. "github.com/hiromaily/go-bitcoin/pkg/kvs"
)

var db *LevelDB

func setup() {
	// KVS
	var err error

	db, err = InitDB("../../data/kvs/db")
	if err != nil {
		panic(err)
	}
}

func teardown() {
	db.Close()
}

func TestMain(m *testing.M) {
	setup()

	code := m.Run()

	teardown()

	os.Exit(code)
}

func TestPutAndGet(t *testing.T) {
	//Put
	err := db.Put("unspent", "testkey1", []byte("data1234567890"))
	if err != nil {
		t.Fatal(err)
	}
	//Get
	val, err := db.Get("unspent", "testkey1")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("[Done] %s", string(val))

	//このトランザクションIDはDBに保存が必要 => ここで保存されたIDはconfirmationのチェックに使われる
	//err = w.Db.Put("unspent", tx.TxID+string(tx.Vout), nil)
	//if err != nil {
	//	//このタイミングでエラーがおきるのであれば、設計ミス
	//	log.Printf("[Error] Error by w.Db.Put(unspent). This error should not occurred.:, error:%v", err)
	//	continue
	//}
}
