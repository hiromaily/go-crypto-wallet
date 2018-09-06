package service

import (
	"bufio"
	"os"
	"strings"

	"github.com/hiromaily/go-bitcoin/pkg/key"
	"github.com/hiromaily/go-bitcoin/pkg/model"
	"github.com/pkg/errors"
)

// ExportAddedPubkeyHistoryTable AddedPubkeyHistoryテーブルをcsvとして出力する
func (w *Wallet) ExportAddedPubkeyHistoryTable(addedPubkeyHistoryTable []model.AddedPubkeyHistoryTable, strAccountType string, keyStatus uint8) (string, error) {
	//fileName
	fileName := key.CreateFilePath(strAccountType, keyStatus)

	file, err := os.Create(fileName)
	if err != nil {
		return "", errors.Errorf("os.Create(%s) error: %s", fileName, err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	//type AddedPubkeyHistoryTable struct {
	//	ID                    int64      `db:"id"`
	//	FullPublicKey         string     `db:"full_public_key"`
	//	AuthAddress1          string     `db:"auth_address1"`
	//	AuthAddress2          string     `db:"auth_address2"`
	//	WalletMultisigAddress string     `db:"wallet_multisig_address"`
	//	RedeemScript          string     `db:"redeem_script"`
	//	IsExported            bool       `db:"is_exported"`
	//	UpdatedAt             *time.Time `db:"updated_at"`
	//}
	for _, record := range addedPubkeyHistoryTable {
		//csvファイル
		tmpData := []string{
			record.FullPublicKey,
			record.AuthAddress1,
			record.AuthAddress2,
			record.WalletMultisigAddress,
			record.RedeemScript,
		}
		_, err = writer.WriteString(strings.Join(tmpData[:], ",") + "\n")
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
