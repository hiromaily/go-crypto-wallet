package coldwallet

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/keystatus"
	"github.com/hiromaily/go-bitcoin/pkg/model/rdb/coldrepo"
	"github.com/hiromaily/go-bitcoin/pkg/wallets/key"
	"github.com/hiromaily/go-bitcoin/pkg/wallets/types"
)

//ExportAccountKey AccountKeyテーブルをcsvとして出力する
//TODO:watch only walletにセットするアドレスは、clientの場合は、wallet_address, receipt/paymentの場合、`wallet_multisig_address`
func (w *ColdWallet) ExportAccountKey(accountType account.AccountType, keyStatus keystatus.KeyStatus) (string, error) {
	//TODO:remove it
	if w.wtype != types.WalletTypeKeyGen {
		return "", errors.New("it's available on Coldwallet1")
	}

	//From coldwallet1
	//Client          -> key_status=1ならok, wallet_address          isMultisig=false
	//Receipt/Payment -> key_status=1ならok, full_public_key         isMultisig=false
	//Receipt/Payment -> key_status=3ならok, wallet_multisig_address isMultisig=true

	//TODO:Multisig対応かどうかのジャッジ
	var updateKeyStatus keystatus.KeyStatus
	if !account.AccountTypeMultisig[accountType] {
		updateKeyStatus = keystatus.KeyStatusAddressExported //4
	} else {
		if keyStatus == keystatus.KeyStatusImportprivkey { //1
			updateKeyStatus = keystatus.KeyStatusPubkeyExported //2
		} else if keyStatus == keystatus.KeyStatusMultiAddressImported { //3
			updateKeyStatus = keystatus.KeyStatusAddressExported //4
		}
	}
	if updateKeyStatus == "" {
		return "", errors.New("parameters are wrong to call ExportAccountKey()")
	}

	//DBから該当する全レコード
	accountKeyTable, err := w.storager.GetAllAccountKeyByKeyStatus(accountType, keyStatus)
	if err != nil {
		return "", errors.Errorf("DB.GetAllAccountKeyByKeyStatus() error: %s", err)
	}

	if len(accountKeyTable) == 0 {
		w.logger.Info("no record in table")
		return "", nil
	}

	//CSVに書き出す
	fileName, err := w.exportAccountKeyTable(accountKeyTable, string(accountType),
		keystatus.KeyStatusValue[keyStatus])
	if err != nil {
		return "", errors.Errorf("key.exportAccountKeyTable() error: %s", err)
	}
	w.logger.Info("call exportAccountKeyTable()", zap.String("fileName", fileName))

	//DBの該当レコードをアップデート
	wifs := make([]string, len(accountKeyTable))
	for idx, record := range accountKeyTable {
		wifs[idx] = record.WalletImportFormat
	}
	_, err = w.storager.UpdateKeyStatusByWIFs(accountType, updateKeyStatus, wifs, nil, true)
	if err != nil {
		return "", errors.Errorf("DB.UpdateKeyStatusByWIFs() error: %s", err)
	}

	//Multisig対応かどうかのジャッジ
	w.logger.Info(
		"Is this account[%s] for multisig: %t",
		zap.String("accountType", accountType.String()),
		zap.Bool("isMultisig", account.AccountTypeMultisig[accountType]))

	return fileName, nil
}

// exportAccountKeyTable AccountKeyTableをファイルとして出力する
func (w *ColdWallet) exportAccountKeyTable(accountKeyTable []coldrepo.AccountKeyTable, strAccountType string, keyStatus uint8) (string, error) {
	//fileName
	fileName := key.CreateFilePath(strAccountType, keyStatus)

	file, err := os.Create(fileName)
	//file, _ := os.OpenFile(*fileName, os.O_WRONLY | os.O_APPEND, 0644)
	if err != nil {
		return "", errors.Errorf("os.Create(%s) error: %s", fileName, err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for _, record := range accountKeyTable {
		//csvファイル
		tmpData := []string{
			record.WalletAddress,
			record.P2shSegwitAddress,
			record.FullPublicKey,
			record.WalletMultisigAddress,
			record.Account,
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
