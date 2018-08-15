package file

import (
	"io/ioutil"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// WriteFileForUnsigned [Debug用] localにファイルを出力する(実運用では、未署名ファイルはGCSにUpload)
func WriteFileForUnsigned(hexTx string) {
	filePrefix := "unsigned_"

	writeFileOnLocal(hexTx, filePrefix)
}

// WriteFileForSigned localに署名済hexをファイルに出力する
func WriteFileForSigned(hexTx string) {
	filePrefix := "signed_"

	writeFileOnLocal(hexTx, filePrefix)
}

func writeFileOnLocal(hexTx, filePrefix string) {
	localPath := "./data/tx/" //TODO:ここは外部ファイル(toml)に定義したほうがいい
	ts := strconv.FormatInt(time.Now().UnixNano(), 10)

	byteTx := []byte(hexTx)
	ioutil.WriteFile(localPath+filePrefix+ts, byteTx, 0644)
}

// ReadFile ファイルを読み込み
func ReadFile(filePath string) (string, error) {
	//localPath := "./data/tx/"
	//ret, err := ioutil.ReadFile(localPath + fileName)
	ret, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", errors.Errorf("ioutil.ReadFile(%s): error: %v", filePath, err)
	}

	return string(ret), nil
}
