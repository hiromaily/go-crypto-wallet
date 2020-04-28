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
	models "github.com/hiromaily/go-bitcoin/pkg/models/rdb"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/types"
)

//ExportAccountKey export key information in account_key_table as csv file
// required address in watch only wallet
//  - account client: `wallet_address`
//  - account others: `wallet_multisig_address`
// this func is expected to be used by only keygen
func (w *ColdWallet) ExportAccountKey(accountType account.AccountType, addrStatus address.AddrStatus) (string, error) {
	if w.wtype != types.WalletTypeKeyGen {
		return "", errors.New("it's available on keygen wallet")
	}

	//Note: condition of data in database at keygen wallet
	// - account: client, addr_status==1, isMultisig==false then export address for `wallet_address`
	// - acccunt: others, addr_status==1, isMultisig==false then export address for `full_public_key`
	// - account: others, addr_status==3, isMultisig==true then export address for `wallet_multisig_address`

	// exptected key status for update
	updateAddrStatus := getAddrStatus(addrStatus, accountType)
	if updateAddrStatus == "" {
		return "", errors.Errorf("addrStatus would be out of range to export: %s", addrStatus.String())
	}

	// get account key
	accountKeyTable, err := w.repo.AccountKey().GetAllAddrStatus(accountType, addrStatus)
	if err != nil {
		return "", errors.Wrap(err, "fail to call repo.AccountKey().GetAllAddrStatus()")
	}
	if len(accountKeyTable) == 0 {
		w.logger.Info("no records in account_key table")
		return "", nil
	}

	//export csv file
	fileName, err := w.exportAccountKey(accountKeyTable, accountType,
		address.AddrStatusValue[addrStatus])
	if err != nil {
		return "", errors.Wrap(err, "fail to call w.exportAccountKeyTable()")
	}

	//update table
	wifs := make([]string, len(accountKeyTable))
	for idx, record := range accountKeyTable {
		wifs[idx] = record.WalletImportFormat
	}
	_, err = w.repo.AccountKey().UpdateAddrStatus(accountType, updateAddrStatus, wifs)
	if err != nil {
		return "", errors.Wrap(err, "fail to call repo.AccountKey().UpdateAddrStatus()")
	}

	w.logger.Debug(
		"address type of account",
		zap.String("accountType", accountType.String()),
		zap.Bool("isMultisig", account.AccountTypeMultisig[accountType]))

	return fileName, nil
}

func getAddrStatus(currentKey address.AddrStatus, accountType account.AccountType) address.AddrStatus {
	//TODO: Though file is already exported, allow to export again?? Yes
	// if you wanna export file again, update keystatus in database manually
	if !account.AccountTypeMultisig[accountType] {
		// not multisig account
		//TODO: current key status should be checked as well
		return address.AddrStatusAddressExported //4
	}
	// multisig account
	if currentKey == address.AddrStatusPrivKeyImported { //1
		return address.AddrStatusPubkeyExported //2
	} else if currentKey == address.AddrStatusMultiAddressImported { //3
		return address.AddrStatusAddressExported //4
	}
	return ""
}

// exportAccountKey export account_key_table as csv file
// TODO: export logic could be defined as address.Storager
func (w *ColdWallet) exportAccountKey(accountKeyTable []*models.AccountKey, accountType account.AccountType, addrStatusVal uint8) (string, error) {
	//create fileName
	fileName := w.addrFileRepo.CreateFilePath(accountType, addrStatusVal)

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
			record.P2SHSegwitAddress,
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
