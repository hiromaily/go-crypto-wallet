package key

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
func CreateFilePath(strAccountType string, keyStatus uint8) string {

	// ./data/pubkey/client_1534744535097796209.csv
	ts := strconv.FormatInt(time.Now().UnixNano(), 10)

	return fmt.Sprintf("%s%s_%d_%s.csv", baseFilePath, strAccountType, keyStatus, ts)
}

// ImportPubKey pubkeyをファイルから読み込む
func ImportPubKey(fileName string) ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, errors.Errorf("os.Open(%s) error: %s", fileName, err)
	}
	defer file.Close()

	var pubKeys []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		pubKeys = append(pubKeys, scanner.Text())
	}

	return pubKeys, nil
}
