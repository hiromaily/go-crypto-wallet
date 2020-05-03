package keygensrv

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
	"github.com/hiromaily/go-bitcoin/pkg/repository/coldrepo"
	"github.com/hiromaily/go-bitcoin/pkg/wallet"
)

// AddressExporter is AddressExporter service
type AddressExporter interface {
	ExportAddress(accountType account.AccountType) (string, error)
}

// AddressExport type
type AddressExport struct {
	logger         *zap.Logger
	accountKeyRepo coldrepo.AccountKeyRepositorier
	addrFileRepo   address.FileRepositorier
	wtype          wallet.WalletType
}

// NewAddressExport returns addressExport
func NewAddressExport(
	logger *zap.Logger,
	accountKeyRepo coldrepo.AccountKeyRepositorier,
	addrFileRepo address.FileRepositorier,
	wtype wallet.WalletType) *AddressExport {

	return &AddressExport{
		logger:         logger,
		accountKeyRepo: accountKeyRepo,
		addrFileRepo:   addrFileRepo,
		wtype:          wtype,
	}
}

// ExportAddress exports addresses in account_key_table as csv file
func (a *AddressExport) ExportAddress(accountType account.AccountType) (string, error) {
	// get target status for account
	var targetAddrStatus address.AddrStatus
	if !account.IsMultisigAccount(accountType) {
		// non-multisig account
		targetAddrStatus = address.AddrStatusPrivKeyImported
	} else {
		targetAddrStatus = address.AddrStatusMultisigAddressGenerated
	}

	// get account key
	accountKeyTable, err := a.accountKeyRepo.GetAllAddrStatus(accountType, targetAddrStatus)
	if err != nil {
		return "", errors.Wrap(err, "fail to call accountKeyRepo.GetAllAddrStatus()")
	}
	if len(accountKeyTable) == 0 {
		a.logger.Info("no records to export in account_key table")
		return "", nil
	}

	// export csv file
	fileName, err := a.exportAccountKey(accountKeyTable, accountType)
	if err != nil {
		return "", err
	}

	// update addrStatus in account_key
	updatedItems := make([]string, len(accountKeyTable))
	for idx, record := range accountKeyTable {
		updatedItems[idx] = record.WalletImportFormat
	}
	_, err = a.accountKeyRepo.UpdateAddrStatus(accountType, address.AddrStatusAddressExported, updatedItems)
	if err != nil {
		return "", errors.Wrap(err, "fail to call a.accountKeyRepo.UpdateAddrStatus()")
	}

	return fileName, nil
}

// exportAccountKey export account_key_table as csv file
func (a *AddressExport) exportAccountKey(accountKeyTable []*models.AccountKey, accountType account.AccountType) (string, error) {
	//create fileName
	fileName := a.addrFileRepo.CreateFilePath(accountType)

	file, err := os.Create(fileName)
	if err != nil {
		return "", errors.Wrapf(err, "fail to call os.Create(%s)", fileName)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// export any address, wallet side chooses proper address/
	for _, record := range accountKeyTable {
		//each line of csv data
		tmpData := []string{
			record.Account,
			record.P2PKHAddress,
			record.P2SHSegwitAddress,
			record.FullPublicKey,
			record.MultisigAddress,
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
