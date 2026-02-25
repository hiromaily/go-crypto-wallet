package eth

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"

	filerepo "github.com/hiromaily/go-crypto-wallet/internal/application/ports/file"
	repocold "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository/cold"
	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

type exportFullPubkeyUseCase struct {
	accountKeyRepo repocold.ETHAccountKeyRepositorier
	pubkeyFileRepo filerepo.AddressFileRepositorier
	coinTypeCode   domainCoin.CoinTypeCode
}

// NewExportFullPubkeyUseCase creates a new ExportFullPubkeyUseCase for ETH keygen wallet.
// It exports the account-level extended public key (accountXpub) to a file so the
// Watch wallet can derive and verify child addresses without holding private keys.
func NewExportFullPubkeyUseCase(
	accountKeyRepo repocold.ETHAccountKeyRepositorier,
	pubkeyFileRepo filerepo.AddressFileRepositorier,
	coinTypeCode domainCoin.CoinTypeCode,
) keygenusecase.ExportFullPubkeyUseCase {
	return &exportFullPubkeyUseCase{
		accountKeyRepo: accountKeyRepo,
		pubkeyFileRepo: pubkeyFileRepo,
		coinTypeCode:   coinTypeCode,
	}
}

// Export derives the accountXpub from the stored accountXpriv and writes it to a CSV file.
// The file uses a 5-field ETH-specific CSV format (no multisig auth type):
//
//	coinTypeCode,accountType,purpose,extendedPubKey,derivationPath
func (u *exportFullPubkeyUseCase) Export(
	ctx context.Context,
	input keygenusecase.ExportFullPubkeyInput,
) (keygenusecase.ExportFullPubkeyOutput, error) {
	logger.Debug("export ETH full pubkey", "account_type", input.AccountType.String())

	accountXpub, derivationPath, err := u.deriveAccountXpub(input.AccountType)
	if err != nil {
		return keygenusecase.ExportFullPubkeyOutput{}, err
	}

	fileName, err := u.writeXpubFile(input.AccountType, accountXpub, derivationPath)
	if err != nil {
		return keygenusecase.ExportFullPubkeyOutput{}, err
	}

	logger.Info("exported ETH account xpub",
		"account_type", input.AccountType.String(),
		"file", fileName,
	)

	return keygenusecase.ExportFullPubkeyOutput{FileName: fileName}, nil
}

// deriveAccountXpub retrieves the accountXpriv from the database, derives the accountXpub,
// and returns it along with the BIP-44 derivation path.
func (u *exportFullPubkeyUseCase) deriveAccountXpub(
	accountType domainAccount.AccountType,
) (xpub string, derivationPath string, err error) {
	key, err := u.accountKeyRepo.GetOneMaxID(accountType)
	if err != nil {
		return "", "", fmt.Errorf("failed to get ETH account key: %w", err)
	}
	if key.AccountExtendedPrivkey == "" {
		return "", "", errors.New(
			"accountXpriv not found for this account; run key generation first",
		)
	}

	extKey, err := hdkeychain.NewKeyFromString(key.AccountExtendedPrivkey)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse accountXpriv: %w", err)
	}

	pubKey, err := extKey.Neuter()
	if err != nil {
		return "", "", fmt.Errorf("failed to derive accountXpub from accountXpriv: %w", err)
	}

	// BIP-44 derivation path at account level: m/44'/coin'/account'
	// ETH BIP44 coin type is 60
	const ethCoinType = 60
	path := fmt.Sprintf("m/44'/%d'/%d'", ethCoinType, accountType.BIP44AccountIndex())

	return pubKey.String(), path, nil
}

// writeXpubFile creates the output CSV file and writes the accountXpub line.
func (u *exportFullPubkeyUseCase) writeXpubFile(
	accountType domainAccount.AccountType,
	accountXpub string,
	derivationPath string,
) (string, error) {
	fileName := u.pubkeyFileRepo.CreateFilePath(accountType)

	f, err := os.Create(fileName) //nolint:gosec
	if err != nil {
		return "", fmt.Errorf("failed to create pubkey file %s: %w", fileName, err)
	}
	defer f.Close() //nolint:errcheck

	writer := bufio.NewWriter(f)

	// ETH single-sig xpub file format (CSV):
	// coinTypeCode,accountType,purpose,extendedPubKey,derivationPath
	// This format is ETH-specific (no multisig auth type) and follows the same
	// structural pattern as the BTC sign wallet fullpubkey export.
	line := fmt.Sprintf("%s,%s,44,%s,%s\n",
		u.coinTypeCode.String(),
		accountType.String(),
		accountXpub,
		derivationPath,
	)

	if _, err := writer.WriteString(line); err != nil {
		return "", fmt.Errorf("failed to write pubkey file: %w", err)
	}
	if err := writer.Flush(); err != nil {
		return "", fmt.Errorf("failed to flush pubkey file: %w", err)
	}

	return fileName, nil
}
