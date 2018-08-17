package file

import (
	"io/ioutil"
	"strconv"
	"strings"
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
	ts := strconv.FormatInt(time.Now().UnixNano(), 10)

	byteTx := []byte(hexTx)
	fileName := filePath + filePrefix + ts
	ioutil.WriteFile(fileName, byteTx, 0644)

	return fileName
}

// ReadFile ファイルを読み込み
func ReadFile(fileName string) (string, error) {
	ret, err := ioutil.ReadFile(filePath + fileName)
	if err != nil {
		return "", errors.Errorf("ioutil.ReadFile(%s): error: %v", fileName, err)
	}

	return string(ret), nil
}

// ParseFile ファイル名を解析する
func ParseFile(fileName, txType string) (int64, string, error) {
	//5_unsigned_1534466246366489473
	s := strings.Split(fileName, "_")
	if len(s) != 3 {
		return 0, "", errors.Errorf("error: invalid file: %s", fileName)
	}
	txReceiptID, err := strconv.ParseInt(s[0], 10, 64)
	if err != nil {
		return 0, "", errors.Errorf("error: invalid file: %s", fileName)
	}
	if s[1] != txType {
		return 0, "", errors.Errorf("error: invalid file: %s", fileName)
	}

	return txReceiptID, s[1], nil
}
