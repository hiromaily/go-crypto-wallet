package coldwallet

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/address"
	"github.com/hiromaily/go-bitcoin/pkg/model/rdb/coldrepo"
	"github.com/hiromaily/go-bitcoin/pkg/wallets/types"
)

//ExportAccountKey export key information in account_key_table as csv file
// required address in watch only wallet
//  - acount client: `wallet_address`
//  - acount others: `wallet_multisig_address`
// this func is expected to be used by only keygen
func (w *ColdWallet) ExportAccountKey(accountType account.AccountType, keyStatus address.AddressStatus) (string, error) {
	//TODO:remove it
	if w.wtype != types.WalletTypeKeyGen {
		return "", errors.New("it's available on Coldwallet1")
	}

	//Note: condition of data in database at keygen wallet
	// - account: client, key_status==1, isMultisig==false then export address for `wallet_address`
	// - acccunt: others, key_status==1, isMultisig==false then export address for `full_public_key`
	// - account: others, key_status==3, isMultisig==true then export address for `wallet_multisig_address`

	// exptected key status for update
	updateAddressStatus := getAddressStatus(keyStatus, accountType)
	if updateAddressStatus == "" {
		return "", errors.New("it can't export file anymore")
	}

	// get account key
	accountKeyTable, err := w.storager.GetAllAccountKeyByAddressStatus(accountType, keyStatus)
	if err != nil {
		return "", errors.Wrap(err, "fail to call storager.GetAllAccountKeyByAddressStatus()")
	}
	if len(accountKeyTable) == 0 {
		w.logger.Info("no records in account_key table")
		return "", nil
	}

	//export csv file
	fileName, err := w.exportAccountKeyTable(accountKeyTable, accountType,
		address.AddressStatusValue[keyStatus])
	if err != nil {
		return "", errors.Wrap(err, "fail to call w.exportAccountKeyTable()")
	}

	//update table
	wifs := make([]string, len(accountKeyTable))
	for idx, record := range accountKeyTable {
		wifs[idx] = record.WalletImportFormat
	}
	_, err = w.storager.UpdateAddressStatusByWIFs(accountType, updateAddressStatus, wifs, nil, true)
	if err != nil {
		return "", errors.Wrap(err, "fail to call storager.UpdateAddressStatusByWIFs()")
	}

	w.logger.Debug(
		"address type of account",
		zap.String("accountType", accountType.String()),
		zap.Bool("isMultisig", account.AccountTypeMultisig[accountType]))

	return fileName, nil
}

func getAddressStatus(currentKey address.AddressStatus, accountType account.AccountType) address.AddressStatus {
	//TODO: Though file is already exported, allow to export again?? Yes
	// if you wanna export file again, update keystatus in database manually
	//if keystatus.AddressStatusValue[currentKey] >= keystatus.AddressStatusValue[keystatus.AddressStatusAddressExported]{
	//	return ""
	//}
	if !account.AccountTypeMultisig[accountType] {
		// not multisig account
		//TODO: current key status should be checked as well
		return address.AddressStatusAddressExported //4
	} else {
		// multisig account
		if currentKey == address.AddressStatusPrivKeyImported { //1
			return address.AddressStatusPubkeyExported //2
		} else if currentKey == address.AddressStatusMultiAddressImported { //3
			return address.AddressStatusAddressExported //4
		}
	}
	return ""
}

// exportAccountKeyTable export account_key_table as csv file
func (w *ColdWallet) exportAccountKeyTable(accountKeyTable []coldrepo.AccountKeyTable, accountType account.AccountType, keyStatusVal uint8) (string, error) {
	//create fileName
	//fileName := key.CreateFilePath(accountType.String(), keyStatusVal)
	fileName := w.addrFileStorager.CreateFilePath(accountType, keyStatusVal)

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
