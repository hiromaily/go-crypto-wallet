package csv

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

var (
	baseFilePath = "./data/pubkey/"
)

// SetFilePath デフォルトの出入力に利用されるファイルパスをセットする
func SetFilePath(basePath string) {
	baseFilePath = basePath
}

// CreateFilePath ファイルパスを作成する
func CreateFilePath(strAccountType string) string {

	// ./data/pubkey/client_1534744535097796209.csv
	ts := strconv.FormatInt(time.Now().UnixNano(), 10)

	return fmt.Sprintf("%s%s_%s.csv", baseFilePath, strAccountType, ts)
}

// ExportPubKey pubkeyをファイルとして出力する
func ExportPubKey(pubKeys []string, strAccountType string) (string, error) {
	//fileName
	fileName := CreateFilePath(strAccountType)

	file, err := os.Create(fileName)
	//file, _ := os.OpenFile(*fileName, os.O_WRONLY | os.O_APPEND, 0644)
	if err != nil {
		return "", errors.Errorf("os.Create(%s) error: %v", fileName, err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, key := range pubKeys {
		//csvファイルとしてだが、このケースでは、pubkeyの1カラムのみ
		_, err = writer.WriteString(key + "\n")
		if err != nil {
			return "", errors.Errorf("writer.WriteString(%s) error: %v", fileName, err)
		}
	}
	err = writer.Flush()
	if err != nil {
		return "", errors.Errorf("writer.Flush(%s) error: %v", fileName, err)
	}

	return fileName, nil
}
