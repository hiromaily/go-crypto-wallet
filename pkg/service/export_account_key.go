package service

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/hiromaily/go-bitcoin/pkg/key"
	"github.com/hiromaily/go-bitcoin/pkg/model"
	"github.com/pkg/errors"
)

// ExportAccountKeyTable AccountKeyTableをファイルとして出力する
func (w *Wallet) ExportAccountKeyTable(accountKeyTable []model.AccountKeyTable, strAccountType string, keyStatus uint8) (string, error) {
	//fileName
	fileName := key.CreateFilePath(strAccountType, keyStatus)

	file, err := os.Create(fileName)
	//file, _ := os.OpenFile(*fileName, os.O_WRONLY | os.O_APPEND, 0644)
	if err != nil {
		return "", errors.Errorf("os.Create(%s) error: %s", fileName, err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	//type AccountKeyTable struct {
	//	ID                    int64  `db:"id"`
	//	WalletAddress         string `db:"wallet_address"`
	//	P2shSegwitAddress     string `db:"p2sh_segwit_address"`
	//	FullPublicKey         string `db:"full_public_key"`
	//	WalletMultisigAddress string `db:"wallet_multisig_address"`
	//	RedeemScript          string `db:"redeem_script"`
	//	WalletImportFormat    string `db:"wallet_import_format"`
	//	Account               string `db:"account"`
	//	KeyType               uint8  `db:"key_type"`
	//	Idx                   uint32 `db:"idx"`
	//	KeyStatus             uint8  `db:"key_status"`
	//	UpdatedAt *time.Time `db:"updated_at"`
	//}
	for _, record := range accountKeyTable {
		//csvファイル
		tmpData := []string{
			record.WalletAddress,
			record.P2shSegwitAddress,
			record.FullPublicKey,
			record.WalletMultisigAddress,
			record.Account,
			strconv.Itoa(int(record.KeyType)),
			strconv.Itoa(int(record.Idx)),
		}
		_, err = writer.WriteString(strings.Join(tmpData[:], ",") + "\n")
		//_, err = writer.WriteString(record. + "\n")
		if err != nil {
			return "", errors.Errorf("writer.WriteString(%s) error: %s", fileName, err)
		}
	}
	err = writer.Flush()
	if err != nil {
		return "", errors.Errorf("writer.Flush(%s) error: %s", fileName, err)
	}

	return fileName, nil
}

// ExportPubKey pubkeyをファイルとして出力する
//func ExportPubKey(pubKeys []string, strAccountType string) (string, error) {
//	//fileName
//	fileName := CreateFilePath(strAccountType)
//
//	file, err := os.Create(fileName)
//	//file, _ := os.OpenFile(*fileName, os.O_WRONLY | os.O_APPEND, 0644)
//	if err != nil {
//		return "", errors.Errorf("os.Create(%s) error: %s", fileName, err)
//	}
//	defer file.Close()
//
//	writer := bufio.NewWriter(file)
//	for _, key := range pubKeys {
//		//csvファイルとしてだが、このケースでは、pubkeyの1カラムのみ
//		_, err = writer.WriteString(key + "\n")
//		if err != nil {
//			return "", errors.Errorf("writer.WriteString(%s) error: %s", fileName, err)
//		}
//	}
//	err = writer.Flush()
//	if err != nil {
//		return "", errors.Errorf("writer.Flush(%s) error: %s", fileName, err)
//	}
//
//	return fileName, nil
//}
