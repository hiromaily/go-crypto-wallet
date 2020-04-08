package wallet

import (
	"bufio"
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/logger"
	"github.com/hiromaily/go-bitcoin/pkg/model"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/key"
)

//ExportAddedPubkeyHistory AddedPubkeyHistoryテーブルをcsvとして出力する
// coldwallet2から使用
func (w *Wallet) ExportAddedPubkeyHistory(accountType account.AccountType) (string, error) {
	//TODO:remove it
	if w.Type != WalletTypeSignature {
		return "", errors.New("it's available on Coldwallet2")
	}

	//DBから該当する全レコード
	//is_exported=falseで且つ、multisig_addressが生成済のレコードが対象
	addedPubkeyHistoryTable, err := w.DB.GetAddedPubkeyHistoryTableByNotExported(accountType)
	if err != nil {
		return "", errors.Errorf("DB.GetAddedPubkeyHistoryTableByNotExported() error: %s", err)
	}

	if len(addedPubkeyHistoryTable) == 0 {
		logger.Info("no record in table")
		return "", nil
	}

	//CSVに書き出す
	//TODO:何がわかりやすいか, このために新たなステータスを追加したほうがいいか
	fileName, err := w.exportAddedPubkeyHistoryTable(addedPubkeyHistoryTable, string(accountType),
		enum.KeyStatusValue[enum.KeyStatusPubkeyExported])
	if err != nil {
		return "", errors.Errorf("key.ExportAddedPubkeyHistoryTable() error: %s", err)
	}
	logger.Infof("file name is %s", fileName)

	//DBの該当レコードをアップデート
	ids := make([]int64, len(addedPubkeyHistoryTable))
	for idx, record := range addedPubkeyHistoryTable {
		ids[idx] = record.ID
	}
	_, err = w.DB.UpdateIsExportedOnAddedPubkeyHistoryTable(accountType, ids, nil, true)
	if err != nil {
		return "", errors.Errorf("DB.UpdateIsExportedOnAddedPubkeyHistoryTable() error: %s", err)
	}

	return fileName, nil
}

// ExportAddedPubkeyHistoryTable AddedPubkeyHistoryテーブルをcsvとして出力する
func (w *Wallet) exportAddedPubkeyHistoryTable(addedPubkeyHistoryTable []model.AddedPubkeyHistoryTable, strAccountType string, keyStatus uint8) (string, error) {
	//fileName
	fileName := key.CreateFilePath(strAccountType, keyStatus)

	file, err := os.Create(fileName)
	if err != nil {
		return "", errors.Errorf("os.Create(%s) error: %s", fileName, err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

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
