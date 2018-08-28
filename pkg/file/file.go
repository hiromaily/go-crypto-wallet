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

// WatchOnlyWalletにおいては、開発時にしか使われないはず？
// でもGCSからDLした先の保存ディレクトリは必要になるな

var (
	baseFilePath = "./data/tx/"
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
func CreateFilePath(actionType enum.ActionType, txType enum.TxType, txID int64, withPath bool) string {

	// ./data/tx/receipt/receipt_8_unsigned_1534744535097796209
	if withPath {
		baseDir := fmt.Sprintf("%s%s/", baseFilePath, string(actionType))
		return fmt.Sprintf("%s%s_%d_%s_", baseDir, string(actionType), txID, txType)
	}
	return fmt.Sprintf("%s_%d_%s_", string(actionType), txID, txType)
}

// ParseFile ファイル名を解析する
func ParseFile(filePath string, txType enum.TxType) (int64, enum.ActionType, string, error) {
	//フルパスが渡されることも想定
	tmp := strings.Split(filePath, "/")
	fileName := tmp[len(tmp)-1]

	//receipt_5_unsigned_1534466246366489473
	//length
	s := strings.Split(fileName, "_")
	if len(s) != 4 {
		return 0, "", "", errors.Errorf("error: invalid file: %s", fileName)
	}

	//Action
	if !enum.ValidateActionType(s[0]) {
		return 0, "", "", errors.Errorf("error: invalid file: %s", fileName)
	}

	//receiptID
	txReceiptID, err := strconv.ParseInt(s[1], 10, 64)
	if err != nil {
		return 0, "", "", errors.Errorf("error: invalid file: %s", fileName)
	}

	//txType
	if s[2] != string(txType) {
		return 0, "", "", errors.Errorf("error: invalid file: %s", fileName)
	}

	return txReceiptID, enum.ActionType(s[0]), s[2], nil
}

// WriteFile ファイルに書き込む
// TODO:ioをパラメータに持つか
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
