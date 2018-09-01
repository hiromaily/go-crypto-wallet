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
	//pubKeys, err := w.DB.GetPubkeyNotExportedPubKey(accountType, isMultisig)
	accountKeyTable, err := w.DB.GetAllByKeyStatus(accountType, keyStatus)
	if err != nil {
		return "", errors.Errorf("key.GetAllByKeyStatus() error: %s", err)
	}

	if len(accountKeyTable) == 0 {
		logger.Info("no record in table")
		return "", nil
	}

	//CSVに書き出す
	//fileName, err := key.ExportPubKey(pubKeys, string(accountType))
	fileName, err := key.ExportAccountKeyTable(accountKeyTable, string(accountType),
		enum.KeyStatusValue[keyStatus])
	if err != nil {
		return "", errors.Errorf("key.ExportPubKey() error: %s", err)
	}
	logger.Infof("file name is %s", fileName)

	//DBの該当レコードをアップデート
	//TODO
	wifs := make([]string, len(accountKeyTable))
	for idx, record := range accountKeyTable {
		wifs[idx] = record.WalletImportFormat
	}

	logger.Debug(updateKeyStatus)

	//_, err = w.DB.UpdateIsExprotedPubKey(accountType, pubKeys, isMultisig, nil, true)
	_, err = w.DB.UpdateKeyStatusByWIFs(accountType, updateKeyStatus, wifs, nil, true)
	if err != nil {
		return "", errors.Errorf("DB.UpdateIsExprotedPubKey() error: %s", err)
	}

	return fileName, nil
}