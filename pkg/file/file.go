package file

import (
	"io/ioutil"
	"strconv"
	"time"

	"fmt"
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/pkg/errors"
	"strings"
)

// これは開発時にしか使われないはず

var (
	baseFilePath = "./data/tx/%s/"
)

// 出力されるファイルフォーマットについて
// [receipt_1_unsigned_timestamp]
// - receipt or payment: 入金/出金
// - int: 該当するtxReceiptID、つまりこの数値から、ファイル名のリレーションを追うことが可能
// - unsigned, signed, sent: トランザクションのタイプ
// - タイムスタンプ

// SetFilePath デフォルトの出入力に利用されるファイルパスをセットする
func SetFilePath(basePath string) {
	baseFilePath = basePath
}

// CreateFilePath 書き込み用として、ファイルパスを生成する(読み込みは渡されたパスをそのまま利用するのみ)
// TODO:Actionも名前として考慮すること
func CreateFilePath(actionType enum.ActionType, txType enum.TxType, txID int64) string {
	basePath := fmt.Sprintf(baseFilePath, actionType)

	return fmt.Sprintf("%s%d_%s_", basePath, txID, txType)
}

// ParseFile ファイル名を解析する
func ParseFile(filePath string, txType enum.TxType) (int64, string, error) {
	//フルパスが渡されることも想定
	tmp := strings.Split(filePath, "/")
	fileName := tmp[len(tmp)-1]

	//receipt_5_unsigned_1534466246366489473
	s := strings.Split(fileName, "_")
	if len(s) != 4 {
		return 0, "", errors.Errorf("error: invalid file: %s", fileName)
	}
	txReceiptID, err := strconv.ParseInt(s[1], 10, 64)
	if err != nil {
		return 0, "", errors.Errorf("error: invalid file: %s", fileName)
	}
	if s[2] != string(txType) {
		return 0, "", errors.Errorf("error: invalid file: %s", fileName)
	}

	return txReceiptID, s[2], nil
}

// WriteFileForUnsigned [Debug用] localにファイルを出力する(実運用では、未署名ファイルはGCSにUpload)
// 戻り値としてファイル名を返す
//func WriteFileForUnsigned(txReceiptID int64, path, hexTx string) string {
//	filePrefix := strconv.FormatInt(txReceiptID, 10) + "_unsigned_"
//
//	return writeFileOnLocal(hexTx, path, filePrefix)
//}

// WriteFileForSigned localに署名済hexをファイルに出力する
//func WriteFileForSigned(txReceiptID int64, path, hexTx string) string {
//	filePrefix := strconv.FormatInt(txReceiptID, 10) + "_signed_"
//
//	return writeFileOnLocal(hexTx, path, filePrefix)
//}

//func WriteFileOnLocal(hexTx, path, filePrefix string) string {
//	ts := strconv.FormatInt(time.Now().UnixNano(), 10)
//
//	byteTx := []byte(hexTx)
//	fileName := filePath + path + filePrefix + ts
//	ioutil.WriteFile(fileName, byteTx, 0644)
//
//	return fileName
//}

// WriteFile ファイルに書き込む
func WriteFile(path, hexTx string) (string, error) {
	ts := strconv.FormatInt(time.Now().UnixNano(), 10)

	byteTx := []byte(hexTx)
	fileName := path + ts
	err := ioutil.WriteFile(fileName, byteTx, 0644)
	if err != nil {
		return "", errors.Errorf("ioutil.WriteFile(%s) error:%v", fileName, err)
	}

	return fileName, nil
}

// ReadFile ファイルを読み込み
func ReadFile(path string) (string, error) {
	ret, err := ioutil.ReadFile(path)
	if err != nil {
		return "", errors.Errorf("ioutil.ReadFile(%s): error: %v", path, err)
	}

	return string(ret), nil
}
