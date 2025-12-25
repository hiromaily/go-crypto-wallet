package shared

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/cold"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file"
	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/address"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
)

type exportAddressUseCase struct {
	accountKeyRepo  cold.AccountKeyRepositorier
	addrFileRepo    file.AddressFileRepositorier
	multisigAccount account.MultisigAccounter
	coinTypeCode    domainCoin.CoinTypeCode
}

// NewExportAddressUseCase creates a new ExportAddressUseCase
func NewExportAddressUseCase(
	accountKeyRepo cold.AccountKeyRepositorier,
	addrFileRepo file.AddressFileRepositorier,
	multisigAccount account.MultisigAccounter,
	coinTypeCode domainCoin.CoinTypeCode,
) keygenusecase.ExportAddressUseCase {
	return &exportAddressUseCase{
		accountKeyRepo:  accountKeyRepo,
		addrFileRepo:    addrFileRepo,
		multisigAccount: multisigAccount,
		coinTypeCode:    coinTypeCode,
	}
}

func (u *exportAddressUseCase) Export(
	ctx context.Context,
	input keygenusecase.ExportAddressInput,
) (keygenusecase.ExportAddressOutput, error) {
	// Get target status for account based on coin type
	var targetAddrStatus address.AddrStatus
	switch u.coinTypeCode {
	case domainCoin.BTC, domainCoin.BCH:
		if !u.multisigAccount.IsMultisigAccount(input.AccountType) {
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
		return keygenusecase.ExportAddressOutput{}, fmt.Errorf("coinType[%s] is not implemented yet", u.coinTypeCode)
	default:
		return keygenusecase.ExportAddressOutput{}, fmt.Errorf("coinType[%s] is not implemented yet", u.coinTypeCode)
	}

	// Get account key
	accountKeyTable, err := u.accountKeyRepo.GetAllAddrStatus(input.AccountType, targetAddrStatus)
	if err != nil {
		return keygenusecase.ExportAddressOutput{},
			fmt.Errorf("fail to call accountKeyRepo.GetAllAddrStatus(): %w", err)
	}
	if len(accountKeyTable) == 0 {
		logger.Info("no records to export in account_key table")
		return keygenusecase.ExportAddressOutput{
			FileName: "",
		}, nil
	}

	// Export csv file
	fileName, err := u.exportAccountKey(accountKeyTable, input.AccountType)
	if err != nil {
		return keygenusecase.ExportAddressOutput{}, err
	}

	// Update addrStatus in account_key
	updatedItems := make([]string, len(accountKeyTable))
	for idx, record := range accountKeyTable {
		updatedItems[idx] = record.WalletImportFormat
	}
	_, err = u.accountKeyRepo.UpdateAddrStatus(input.AccountType, address.AddrStatusAddressExported, updatedItems)
	if err != nil {
		return keygenusecase.ExportAddressOutput{},
			fmt.Errorf("fail to call accountKeyRepo.UpdateAddrStatus(): %w", err)
	}

	return keygenusecase.ExportAddressOutput{
		FileName: fileName,
	}, nil
}

// exportAccountKey exports account_key_table as csv file
func (u *exportAddressUseCase) exportAccountKey(
	accountKeyTable []*models.AccountKey, accountType domainAccount.AccountType,
) (string, error) {
	// Create fileName
	fileName := u.addrFileRepo.CreateFilePath(accountType)

	file, err := os.Create(fileName) //nolint:gosec
	if err != nil {
		return "", fmt.Errorf("fail to call os.Create(%s): %w", fileName, err)
	}

	defer func() {
		if cerr := file.Close(); cerr != nil {
			err = fmt.Errorf("failed to close file: %w", cerr)
		}
	}()

	writer := bufio.NewWriter(file)

	// Export any address, wallet side chooses proper address
	for _, record := range accountKeyTable {
		// Each line of csv data
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
