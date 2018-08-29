package key

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hiromaily/go-bitcoin/pkg/model"
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

// ImportPubKey pubkeyをファイルから読み込む
func ImportPubKey(fileName string) ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, errors.Errorf("os.Open(%s) error: %v", fileName, err)
	}
	defer file.Close()

	var pubKeys []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		pubKeys = append(pubKeys, scanner.Text())
	}

	return pubKeys, nil
}

// ExportAccountKeyTable AccountKeyTableをファイルとして出力する
func ExportAccountKeyTable(accountKeyTable []model.AccountKeyTable, strAccountType string) (string, error) {
	//fileName
	fileName := CreateFilePath(strAccountType)

	file, err := os.Create(fileName)
	//file, _ := os.OpenFile(*fileName, os.O_WRONLY | os.O_APPEND, 0644)
	if err != nil {
		return "", errors.Errorf("os.Create(%s) error: %v", fileName, err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	//type AccountKeyTable struct {
	//	ID                    int64      `db:"id"`
	//	WalletAddress         string     `db:"wallet_address"`
	//	WalletMultisigAddress string     `db:"wallet_multisig_address"`
	//	RedeemScript          string     `db:"redeem_script"`
	//	WalletImportFormat    string     `db:"wallet_import_format"`
	//	Account               string     `db:"account"`
	//	KeyType               uint8      `db:"key_type"`
	//	Idx                   uint32     `db:"idx"`
	//	IsImprotedPrivKey     bool       `db:"is_imported_priv_key"`
	//	IsExprotedPubKey      bool       `db:"is_exported_pub_key"`
	//	UpdatedAt             *time.Time `db:"updated_at"`
	//}

	for _, record := range accountKeyTable {
		//csvファイル
		tmpData := []string{
			record.WalletAddress,
			record.Account,
			strconv.Itoa(int(record.KeyType)),
			strconv.Itoa(int(record.Idx)),
		}
		_, err = writer.WriteString(strings.Join(tmpData[:], ",") + "\n")
		//_, err = writer.WriteString(record. + "\n")
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
