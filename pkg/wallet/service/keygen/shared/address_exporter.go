package shared

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/address"
	domainAccount "github.com/hiromaily/go-crypto-wallet/pkg/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/pkg/domain/coin"
	domainWallet "github.com/hiromaily/go-crypto-wallet/pkg/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/repository/cold"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/storage/file"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
)

// AddressExport type
type AddressExport struct {
	accountKeyRepo  cold.AccountKeyRepositorier
	addrFileRepo    file.AddressFileRepositorier
	multisigAccount account.MultisigAccounter
	coinTypeCode    domainCoin.CoinTypeCode
	wtype           domainWallet.WalletType
}

// NewAddressExport returns addressExport
func NewAddressExport(
	accountKeyRepo cold.AccountKeyRepositorier,
	addrFileRepo file.AddressFileRepositorier,
	multisigAccount account.MultisigAccounter,
	coinTypeCode domainCoin.CoinTypeCode,
	wtype domainWallet.WalletType,
) *AddressExport {
	return &AddressExport{
		accountKeyRepo:  accountKeyRepo,
		addrFileRepo:    addrFileRepo,
		multisigAccount: multisigAccount,
		coinTypeCode:    coinTypeCode,
		wtype:           wtype,
	}
}

// ExportAddress exports addresses in account_key_table as csv file
func (a *AddressExport) ExportAddress(accountType domainAccount.AccountType) (string, error) {
	// get target status for account
	var targetAddrStatus address.AddrStatus
	switch a.coinTypeCode {
	case domainCoin.BTC, domainCoin.BCH:
		if !a.multisigAccount.IsMultisigAccount(accountType) {
			// non-multisig account
			targetAddrStatus = address.AddrStatusPrivKeyImported
		} else {
			targetAddrStatus = address.AddrStatusMultisigAddressGenerated
		}
	case domainCoin.ETH:
		targetAddrStatus = address.AddrStatusPrivKeyImported
	case domainCoin.XRP:
		targetAddrStatus = address.AddrStatusHDKeyGenerated
	case domainCoin.LTC, domainCoin.ERC20, domainCoin.HYT:
		return "", fmt.Errorf("coinType[%s] is not implemented yet", a.coinTypeCode)
	default:
		return "", fmt.Errorf("coinType[%s] is not implemented yet", a.coinTypeCode)
	}

	// get account key
	accountKeyTable, err := a.accountKeyRepo.GetAllAddrStatus(accountType, targetAddrStatus)
	if err != nil {
		return "", fmt.Errorf("fail to call accountKeyRepo.GetAllAddrStatus(): %w", err)
	}
	if len(accountKeyTable) == 0 {
		logger.Info("no records to export in account_key table")
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
		return "", fmt.Errorf("fail to call a.accountKeyRepo.UpdateAddrStatus(): %w", err)
	}

	return fileName, nil
}

// exportAccountKey export account_key_table as csv file
func (a *AddressExport) exportAccountKey(
	accountKeyTable []*models.AccountKey, accountType domainAccount.AccountType,
) (string, error) {
	// create fileName
	fileName := a.addrFileRepo.CreateFilePath(accountType)

	file, err := os.Create(fileName) //nolint:gosec
	if err != nil {
		return "", fmt.Errorf("fail to call os.Create(%s): %w", fileName, err)
	}

	defer file.Close()

	writer := bufio.NewWriter(file)

	// export any address, wallet side chooses proper address/
	for _, record := range accountKeyTable {
		// each line of csv data
		tmpData := address.CreateLine(record)
		_, err = writer.WriteString(strings.Join(tmpData, ",") + "\n")
		if err != nil {
			return "", fmt.Errorf("fail to call writer.WriteString(%s): %w", fileName, err)
		}
	}
	err = writer.Flush()
	if err != nil {
		return "", fmt.Errorf("fail to call writer.Flush(%s): %w", fileName, err)
	}

	return fileName, nil
}
