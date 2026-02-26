package btc

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	appdto "github.com/hiromaily/go-crypto-wallet/internal/application/dto"
	file "github.com/hiromaily/go-crypto-wallet/internal/application/ports/file"
	repocold "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository/cold"
	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainBTC "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/btc"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

type exportAddressUseCase struct {
	accountKeyRepo  repocold.BTCAccountKeyRepositorier
	addrFileRepo    file.AddressFileRepositorier
	multisigAccount *domainAccount.MultisigConfig
	coinTypeCode    domainCoin.CoinTypeCode
}

// NewExportAddressUseCase creates a new ExportAddressUseCase for BTC/BCH.
// Note: This use case is specific to BTC/BCH. ETH and XRP require separate implementations
// due to different repository interfaces and domain types.
func NewExportAddressUseCase(
	accountKeyRepo repocold.BTCAccountKeyRepositorier,
	addrFileRepo file.AddressFileRepositorier,
	multisigAccount *domainAccount.MultisigConfig,
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
	// Validate coin type - this use case only supports BTC/BCH
	if u.coinTypeCode != domainCoin.BTC && u.coinTypeCode != domainCoin.BCH {
		return keygenusecase.ExportAddressOutput{},
			fmt.Errorf("unsupported coinType[%s]: this use case only supports BTC/BCH", u.coinTypeCode)
	}

	// Get target status for account based on multisig configuration
	var targetAddrStatus domainAddress.AddrStatus
	if !u.multisigAccount.IsMultisigAccount(input.AccountType) {
		// non-multisig account
		targetAddrStatus = domainAddress.AddrStatusPrivKeyImported
	} else {
		targetAddrStatus = domainAddress.AddrStatusMultisigAddressGenerated
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
	_, err = u.accountKeyRepo.UpdateAddrStatus(input.AccountType, domainAddress.AddrStatusAddressExported, updatedItems)
	if err != nil {
		return keygenusecase.ExportAddressOutput{},
			fmt.Errorf("fail to call accountKeyRepo.UpdateAddrStatus(): %w", err)
	}

	return keygenusecase.ExportAddressOutput{
		FileName: fileName,
	}, nil
}

// exportAccountKey exports btc_account_key table as csv file
func (u *exportAddressUseCase) exportAccountKey(
	accountKeyTable []*domainBTC.BTCAccountKey, accountType domainAccount.AccountType,
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
		tmpData := appdto.CreateAddressLine(record)
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
