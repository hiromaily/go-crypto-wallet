package xrp

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	file "github.com/hiromaily/go-crypto-wallet/internal/application/ports/file"
	repocold "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository/cold"
	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainXRP "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/xrp"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

type exportAddressUseCase struct {
	xrpAccountKeyRepo repocold.XRPAccountKeyRepositorier
	addrFileRepo      file.AddressFileRepositorier
	coinTypeCode      domainCoin.CoinTypeCode
}

// NewExportAddressUseCase creates a new ExportAddressUseCase for XRP.
func NewExportAddressUseCase(
	xrpAccountKeyRepo repocold.XRPAccountKeyRepositorier,
	addrFileRepo file.AddressFileRepositorier,
	coinTypeCode domainCoin.CoinTypeCode,
) keygenusecase.ExportAddressUseCase {
	return &exportAddressUseCase{
		xrpAccountKeyRepo: xrpAccountKeyRepo,
		addrFileRepo:      addrFileRepo,
		coinTypeCode:      coinTypeCode,
	}
}

func (u *exportAddressUseCase) Export(
	ctx context.Context,
	input keygenusecase.ExportAddressInput,
) (keygenusecase.ExportAddressOutput, error) {
	// XRP keys are stored with AddrStatusHDKeyGenerated after HD wallet generation.
	// There is no separate "import private key" step for XRP.
	keys, err := u.xrpAccountKeyRepo.GetAllAddrStatus(ctx, input.AccountType, domainAddress.AddrStatusHDKeyGenerated)
	if err != nil {
		return keygenusecase.ExportAddressOutput{},
			fmt.Errorf("fail to call xrpAccountKeyRepo.GetAllAddrStatus(): %w", err)
	}
	if len(keys) == 0 {
		logger.Info("no records to export in xrp_account_key table")
		return keygenusecase.ExportAddressOutput{FileName: ""}, nil
	}

	fileName, err := u.exportXRPAccountKey(keys, input.AccountType)
	if err != nil {
		return keygenusecase.ExportAddressOutput{}, err
	}

	// Update status to AddrStatusAddressExported using AccountIDs as identifiers
	accountIDs := make([]string, len(keys))
	for i, key := range keys {
		accountIDs[i] = key.AccountID
	}
	if _, err = u.xrpAccountKeyRepo.UpdateAddrStatus(
		ctx, input.AccountType, domainAddress.AddrStatusAddressExported, accountIDs,
	); err != nil {
		return keygenusecase.ExportAddressOutput{},
			fmt.Errorf("fail to call xrpAccountKeyRepo.UpdateAddrStatus(): %w", err)
	}

	return keygenusecase.ExportAddressOutput{FileName: fileName}, nil
}

// exportXRPAccountKey exports XRP account keys as a CSV file.
// CSV format (10 fields, compatible with ConvertAddressLine):
//
//	coin, account, classicAddress, "", "", "", publicKey, "", "", idx
func (u *exportAddressUseCase) exportXRPAccountKey(
	keys []*domainXRP.XRPAccountKey,
	accountType domainAccount.AccountType,
) (fileName string, err error) {
	fileName = u.addrFileRepo.CreateFilePath(accountType)

	f, err := os.Create(fileName) //nolint:gosec
	if err != nil {
		return "", fmt.Errorf("fail to call os.Create(%s): %w", fileName, err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("failed to close file: %w", cerr)
		}
	}()

	writer := bufio.NewWriter(f)
	for _, key := range keys {
		line := createXRPAddressLine(u.coinTypeCode, key)
		if _, err = writer.WriteString(strings.Join(line, ",") + "\n"); err != nil {
			return "", fmt.Errorf("fail to call writer.WriteString(%s): %w", fileName, err)
		}
	}
	if err = writer.Flush(); err != nil {
		return "", fmt.Errorf("fail to call writer.Flush(%s): %w", fileName, err)
	}

	return fileName, nil
}

// createXRPAddressLine creates a 10-field CSV line for an XRP account key.
// The classic address goes in the P2PKHAddress field (index 2).
func createXRPAddressLine(coinTypeCode domainCoin.CoinTypeCode, key *domainXRP.XRPAccountKey) []string {
	return []string{
		coinTypeCode.String(),
		key.Account.String(),
		key.AccountID,    // XRP classic address (r...) — stored in P2PKHAddress position
		"",               // P2SHSegwit — unused for XRP
		"",               // Bech32 — unused for XRP
		"",               // Taproot — unused for XRP
		key.PublicKeyHex, // Public key hex
		"",               // MultisigAddress — unused for XRP
		"",               // RedeemScript — unused for XRP
		strconv.FormatInt(key.AllocatedID, 10),
	}
}
