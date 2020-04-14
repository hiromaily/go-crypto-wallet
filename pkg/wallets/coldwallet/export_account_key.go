package coldwallet

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/key"
	"github.com/hiromaily/go-bitcoin/pkg/model/rdb/coldrepo"
	"github.com/hiromaily/go-bitcoin/pkg/wallets/types"
)

//ExportAccountKey export key information in account_key_table as csv file
// required address in watch only wallet
//  - acount client: `wallet_address`
//  - acount others: `wallet_multisig_address`
// this func is expected to be used by only keygen
func (w *ColdWallet) ExportAccountKey(accountType account.AccountType, keyStatus key.KeyStatus) (string, error) {
	//TODO:remove it
	if w.wtype != types.WalletTypeKeyGen {
		return "", errors.New("it's available on Coldwallet1")
	}

	//Note: condition of data in database at keygen wallet
	// - account: client, key_status==1, isMultisig==false then export address for `wallet_address`
	// - acccunt: others, key_status==1, isMultisig==false then export address for `full_public_key`
	// - account: others, key_status==3, isMultisig==true then export address for `wallet_multisig_address`

	// exptected key status for update
	updateKeyStatus := getKeyStatus(keyStatus, accountType)
	if updateKeyStatus == "" {
		return "", errors.New("it can't export file anymore")
	}

	// get account key
	accountKeyTable, err := w.storager.GetAllAccountKeyByKeyStatus(accountType, keyStatus)
	if err != nil {
		return "", errors.Wrap(err, "fail to call storager.GetAllAccountKeyByKeyStatus()")
	}
	if len(accountKeyTable) == 0 {
		w.logger.Info("no records in account_key table")
		return "", nil
	}

	//export csv file
	fileName, err := exportAccountKeyTable(accountKeyTable, accountType,
		key.KeyStatusValue[keyStatus])
	if err != nil {
		return "", errors.Wrap(err, "fail to call w.exportAccountKeyTable()")
	}

	//update table
	wifs := make([]string, len(accountKeyTable))
	for idx, record := range accountKeyTable {
		wifs[idx] = record.WalletImportFormat
	}
	_, err = w.storager.UpdateKeyStatusByWIFs(accountType, updateKeyStatus, wifs, nil, true)
	if err != nil {
		return "", errors.Wrap(err, "fail to call storager.UpdateKeyStatusByWIFs()")
	}

	w.logger.Debug(
		"address type of account",
		zap.String("accountType", accountType.String()),
		zap.Bool("isMultisig", account.AccountTypeMultisig[accountType]))

	return fileName, nil
}

func getKeyStatus(currentKey key.KeyStatus, accountType account.AccountType) key.KeyStatus {
	//TODO: Though file is already exported, allow to export again?? Yes
	// if you wanna export file again, update keystatus in database manually
	//if keystatus.KeyStatusValue[currentKey] >= keystatus.KeyStatusValue[keystatus.KeyStatusAddressExported]{
	//	return ""
	//}
	if !account.AccountTypeMultisig[accountType] {
		// not multisig account
		//TODO: current key status should be checked as well
		return key.KeyStatusAddressExported //4
	} else {
		// multisig account
		if currentKey == key.KeyStatusImportprivkey { //1
			return key.KeyStatusPubkeyExported //2
		} else if currentKey == key.KeyStatusMultiAddressImported { //3
			return key.KeyStatusAddressExported //4
		}
	}
	return ""
}

// exportAccountKeyTable export account_key_table as csv file
func exportAccountKeyTable(accountKeyTable []coldrepo.AccountKeyTable, accountType account.AccountType, keyStatusVal uint8) (string, error) {
	//create fileName
	fileName := key.CreateFilePath(accountType.String(), keyStatusVal)

	file, err := os.Create(fileName)
	if err != nil {
		return "", errors.Wrapf(err, "fail to call os.Create(%s)", fileName)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for _, record := range accountKeyTable {
		//each line of csv data
		tmpData := []string{
			record.WalletAddress,
			record.P2shSegwitAddress,
			record.FullPublicKey,
			record.WalletMultisigAddress,
			record.Account,
			strconv.Itoa(int(record.Idx)),
		}
		_, err = writer.WriteString(strings.Join(tmpData[:], ",") + "\n")
		if err != nil {
			return "", errors.Wrapf(err, "fail to call writer.WriteString(%s)", fileName)
		}
	}
	err = writer.Flush()
	if err != nil {
		return "", errors.Wrapf(err, "fail to call writer.Flush(%s)", fileName)
	}

	return fileName, nil
}
