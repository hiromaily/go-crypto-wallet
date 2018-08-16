package file

import (
	"io/ioutil"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

var filePath = "./data/tx/"

// 出力されるファイルフォーマットについて
// [1_unsigned_timestamp]
// 最初の数字: 該当するtxReceiptID、つまりこの数値から、ファイル名のリレーションを追うことが可能
// unsigned, signed, sent: トランザクションのタイプ
// タイムスタンプ

// SetFilePath デフォルトの出入力に利用されるファイルパスをセットする
func SetFilePath(path string) {
	filePath = path
}

// WriteFileForUnsigned [Debug用] localにファイルを出力する(実運用では、未署名ファイルはGCSにUpload)
// 戻り値としてフィアル名を返す
func WriteFileForUnsigned(txReceiptID int64, hexTx string) string {
	filePrefix := strconv.FormatInt(txReceiptID, 10) + "_unsigned_"

	return writeFileOnLocal(hexTx, filePrefix)
}

// WriteFileForSigned localに署名済hexをファイルに出力する
func WriteFileForSigned(txReceiptID int64, hexTx string) string {
	filePrefix := strconv.FormatInt(txReceiptID, 10) + "_signed_"

	return writeFileOnLocal(hexTx, filePrefix)
}

func writeFileOnLocal(hexTx, filePrefix string) string {
	//localPath := "./data/tx/" //TODO:ここは外部ファイル(toml)に定義したほうがいい
	ts := strconv.FormatInt(time.Now().UnixNano(), 10)

	byteTx := []byte(hexTx)
	fileName := filePath + filePrefix + ts
	ioutil.WriteFile(fileName, byteTx, 0644)

	return fileName
}

// ReadFile ファイルを読み込み
func ReadFile(fileName string) (string, error) {
	//localPath := "./data/tx/"
	//ret, err := ioutil.ReadFile(localPath + fileName)
	ret, err := ioutil.ReadFile(filePath + fileName)
	if err != nil {
		return "", errors.Errorf("ioutil.ReadFile(%s): error: %v", fileName, err)
	}

	return string(ret), nil
}
