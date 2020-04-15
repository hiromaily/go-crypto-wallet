package coldwallet

import (
	"bufio"
	"os"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/address"
	"github.com/hiromaily/go-bitcoin/pkg/model/rdb/coldrepo"
	"github.com/hiromaily/go-bitcoin/pkg/wallets/types"
)

//ExportAddedPubkeyHistory export data in added_pubkey_history table as csv file
func (w *ColdWallet) ExportAddedPubkeyHistory(accountType account.AccountType) (string, error) {
	if w.wtype != types.WalletTypeSignature {
		return "", errors.New("it's available on sign wallet")
	}

	// get record in added_pybkey_history table
	// condition: is_exported==false and multisig_address is already created
	addedPubkeyHistoryTable, err := w.storager.GetAddedPubkeyHistoryTableByNotExported(accountType)
	if err != nil {
		return "", errors.Wrap(err, "storager.GetAddedPubkeyHistoryTableByNotExported()")
	}

	if len(addedPubkeyHistoryTable) == 0 {
		w.logger.Info(
			"no records in added_pubkey_history table",
			zap.String("account", accountType.String()))
		return "", nil
	}

	// export data in added_pubkey_history table as csv file
	fileName, err := w.exportAddedPubkeyHistoryTable(addedPubkeyHistoryTable, accountType,
		address.AddrStatusValue[address.AddrStatusPubkeyExported])
	if err != nil {
		return "", errors.Wrap(err, "fail to call exportAddedPubkeyHistoryTable()")
	}

	// update current status
	ids := make([]int64, len(addedPubkeyHistoryTable))
	for idx, record := range addedPubkeyHistoryTable {
		ids[idx] = record.ID
	}
	_, err = w.storager.UpdateIsExportedOnAddedPubkeyHistoryTable(accountType, ids, nil, true)
	if err != nil {
		return "", errors.Wrap(err, "fail to call storager.UpdateIsExportedOnAddedPubkeyHistoryTable()")
	}

	return fileName, nil
}

// TODO: export logic could be defined as address.Storager
func (w *ColdWallet) exportAddedPubkeyHistoryTable(addedPubkeyHistoryTable []coldrepo.AddedPubkeyHistoryTable, accountType account.AccountType, keyStatus uint8) (string, error) {
	//fileName
	fileName := w.addrFileRepo.CreateFilePath(accountType, keyStatus)

	file, err := os.Create(fileName)
	if err != nil {
		return "", errors.Wrapf(err, "fail to call os.Create(%s)", fileName)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for _, record := range addedPubkeyHistoryTable {
		//each line of csv data
		tmpData := []string{
			record.FullPublicKey,
			record.AuthAddress1,
			record.AuthAddress2,
			record.WalletMultisigAddress,
			record.RedeemScript,
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
