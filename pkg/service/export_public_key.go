package service

//Cold wallet

import (
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/key"
	"github.com/hiromaily/go-bitcoin/pkg/logger"
	"github.com/pkg/errors"
)

//ExportAccountKey AccountKeyテーブルをcsvとして出力する
//TODO:watch only walletにセットするアドレスは、clientの場合は、wallet_address, receipt/paymentの場合、`wallet_multisig_address`
func (w *Wallet) ExportAccountKey(accountType enum.AccountType, keyStatus enum.KeyStatus) (string, error) {
	//AccountType問わずexportは可能にしておくか？はじめ厳しくしておき、徐々に設定を緩める方向で

	//From coldwallet1
	//Client          -> key_status=1ならok, wallet_address          isMultisig=false
	//Receipt/Payment -> key_status=1ならok, full_public_key         isMultisig=false
	//Receipt/Payment -> key_status=3ならok, wallet_multisig_address isMultisig=true
	var updateKeyStatus enum.KeyStatus
	if accountType == enum.AccountTypeClient {
		updateKeyStatus = enum.KeyStatusAddressExported
	} else {
		if keyStatus == enum.KeyStatusImportprivkey {
			updateKeyStatus = enum.KeyStatusPubkeyExported
		} else if keyStatus == enum.KeyStatusMultiAddressImported {
			updateKeyStatus = enum.KeyStatusAddressExported
		}
	}
	if updateKeyStatus == "" {
		return "", errors.New("parameters are wrong to call ExportAccountKey()")
	}

	//TODO:From coldwallet2
	//history table
	//->これは別に定義しょう

	//DBから該当する全レコード
	accountKeyTable, err := w.DB.GetAllAccountKeyByKeyStatus(accountType, keyStatus)
	if err != nil {
		return "", errors.Errorf("DB.GetAllAccountKeyByKeyStatus() error: %s", err)
	}

	if len(accountKeyTable) == 0 {
		logger.Info("no record in table")
		return "", nil
	}

	//CSVに書き出す
	fileName, err := key.ExportAccountKeyTable(accountKeyTable, string(accountType),
		enum.KeyStatusValue[keyStatus])
	if err != nil {
		return "", errors.Errorf("key.ExportAccountKeyTable() error: %s", err)
	}
	logger.Infof("file name is %s", fileName)

	//DBの該当レコードをアップデート
	wifs := make([]string, len(accountKeyTable))
	for idx, record := range accountKeyTable {
		wifs[idx] = record.WalletImportFormat
	}
	_, err = w.DB.UpdateKeyStatusByWIFs(accountType, updateKeyStatus, wifs, nil, true)
	if err != nil {
		return "", errors.Errorf("DB.UpdateIsExprotedPubKey() error: %s", err)
	}

	return fileName, nil
}

//ExportAddedPubkeyHistory AddedPubkeyHistoryテーブルをcsvとして出力する
// coldwallet2から使用
func (w *Wallet) ExportAddedPubkeyHistory(accountType enum.AccountType) (string, error) {
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
	fileName, err := key.ExportAddedPubkeyHistoryTable(addedPubkeyHistoryTable, string(accountType),
		//enum.KeyStatusValue[enum.KeyStatusMultiAddressImported])
		enum.KeyStatusValue[enum.KeyStatusPubkeyExported])
	if err != nil {
		return "", errors.Errorf("key.ExportAccountKeyTable() error: %s", err)
	}
	logger.Infof("file name is %s", fileName)

	//DBの該当レコードをアップデート
	ids := make([]int64, len(addedPubkeyHistoryTable))
	for idx, record := range addedPubkeyHistoryTable {
		ids[idx] = record.ID
	}
	_, err = w.DB.UpdateIsExportedOnAddedPubkeyHistoryTable(accountType, ids, nil, true)
	if err != nil {
		return "", errors.Errorf("DB.UpdateIsExprotedPubKey() error: %s", err)
	}

	return fileName, nil
}
