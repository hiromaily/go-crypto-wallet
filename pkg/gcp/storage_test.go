package gcp_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hiromaily/go-bitcoin/pkg/enum"
	. "github.com/hiromaily/go-bitcoin/pkg/gcp"
	"github.com/hiromaily/go-bitcoin/pkg/txfile"
)

//[gcs]
//storage_key_path = "./data/api_keys/cayenne-dev-strage.json"
//receipt_bucket_name = "cayenne-dev-exchanges-yasui-bucket"
//payment_bucket_name = "cayenne-dev-exchanges-yasui-bucket"

func isGcpDir() bool {
	dir, _ := os.Getwd()
	if s := strings.Split(dir, "/"); s[len(s)-1] == "gcp" {
		return true
	}
	return false
}

func initialStorage(t *testing.T) *ExtClient {
	//TODO: tomlから読み込むように修正する。そのうち
	bucketName := "cayenne-dev-exchanges-yasui-bucket"
	key := "./data/api_keys/cayenne-dev-strage.json"

	//PWDで実行環境に応じてパスを切り替える
	if isGcpDir() {
		key = "../../" + key
	}

	//初期化処理
	storage := NewStorage(bucketName, key)
	cli, err := storage.NewClient(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	return cli
}

func TestSaveAndRead(t *testing.T) {
	txReceiptID := int64(999)
	hex := "storage_test"
	path := txfile.CreateFilePath(enum.ActionTypeReceipt, enum.TxTypeUnsigned, txReceiptID, false)

	//初期化処理
	cli := initialStorage(t)

	//書き込み
	generatedFileName, err := cli.Write(path, []byte(hex))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("generatedFileName: %s", generatedFileName)

	//読み込み
	outputPath := fmt.Sprintf("./data/gcs/%s", generatedFileName)
	if isGcpDir() {
		outputPath = "../../" + outputPath
	}

	err = cli.ReadAndSave(generatedFileName, outputPath, 0666)
	if err != nil {
		t.Fatal(err)
	}

	//Close
	err = cli.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestRead(t *testing.T) {
	//初期化処理
	cli := initialStorage(t)

	fileName := "payment_5_unsigned_1535011799610609638"

	//読み込み
	outputPath := fmt.Sprintf("./data/gcs/%s", fileName)
	if isGcpDir() {
		outputPath = "../../" + outputPath
	}

	err := cli.ReadAndSave(fileName, outputPath, 0666)
	if err != nil {
		t.Fatal(err)
	}

	//Close
	err = cli.Close()
	if err != nil {
		t.Fatal(err)
	}
}
