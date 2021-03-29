package coldsrv

import (
	"bufio"
	"os"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/address"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/repository/coldrepo"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

// AddressExport type
type AddressExport struct {
	logger          *zap.Logger
	accountKeyRepo  coldrepo.AccountKeyRepositorier
	addrFileRepo    address.FileRepositorier
	multisigAccount account.MultisigAccounter
	coinTypeCode    coin.CoinTypeCode
	wtype           wallet.WalletType
}

// NewAddressExport returns addressExport
func NewAddressExport(
	logger *zap.Logger,
	accountKeyRepo coldrepo.AccountKeyRepositorier,
	addrFileRepo address.FileRepositorier,
	multisigAccount account.MultisigAccounter,
	coinTypeCode coin.CoinTypeCode,
	wtype wallet.WalletType) *AddressExport {
	return &AddressExport{
		logger:          logger,
		accountKeyRepo:  accountKeyRepo,
		addrFileRepo:    addrFileRepo,
		multisigAccount: multisigAccount,
		coinTypeCode:    coinTypeCode,
		wtype:           wtype,
	}
}

// ExportAddress exports addresses in account_key_table as csv file
func (a *AddressExport) ExportAddress(accountType account.AccountType) (string, error) {
	// get target status for account
	var targetAddrStatus address.AddrStatus
	switch a.coinTypeCode {
	case coin.BTC, coin.BCH:
		if !a.multisigAccount.IsMultisigAccount(accountType) {
			// non-multisig account
			targetAddrStatus = address.AddrStatusPrivKeyImported
		} else {
			targetAddrStatus = address.AddrStatusMultisigAddressGenerated
		}
	case coin.ETH:
		targetAddrStatus = address.AddrStatusPrivKeyImported
	case coin.XRP:
		targetAddrStatus = address.AddrStatusHDKeyGenerated
	default:
		return "", errors.Errorf("coinType[%s] is not implemented yet.", a.coinTypeCode)
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
	// create fileName
	fileName := a.addrFileRepo.CreateFilePath(accountType)

	file, err := os.Create(fileName)
	if err != nil {
		return "", errors.Wrapf(err, "fail to call os.Create(%s)", fileName)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// export any address, wallet side chooses proper address/
	for _, record := range accountKeyTable {
		// each line of csv data
		tmpData := address.CreateLine(record)
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
