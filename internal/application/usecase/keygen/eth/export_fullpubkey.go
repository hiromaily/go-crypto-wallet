package eth

import (
	"context"
	"fmt"

	filerepo "github.com/hiromaily/go-crypto-wallet/internal/application/ports/file"
	repocold "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository/cold"
	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
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

// Export retrieves the accountXpriv from storage, derives the accountXpub via the domain
// entity, and writes it to a CSV file via the file port.
// The file uses a 5-field ETH-specific CSV format (no multisig auth type):
//
//	coinTypeCode,accountType,purpose,extendedPubKey,derivationPath
func (u *exportFullPubkeyUseCase) Export(
	ctx context.Context,
	input keygenusecase.ExportFullPubkeyInput,
) (keygenusecase.ExportFullPubkeyOutput, error) {
	logger.Debug("export ETH full pubkey", "account_type", input.AccountType.String())

	key, err := u.accountKeyRepo.GetOneMaxID(input.AccountType)
	if err != nil {
		return keygenusecase.ExportFullPubkeyOutput{}, fmt.Errorf("failed to get ETH account key: %w", err)
	}

	xpub, err := key.DeriveAccountXpub()
	if err != nil {
		return keygenusecase.ExportFullPubkeyOutput{}, err
	}

	derivationPath := domainCoin.BIP44AccountPath(u.coinTypeCode, input.AccountType.BIP44AccountIndex())

	fileName, err := u.pubkeyFileRepo.WriteXpubLine(
		input.AccountType,
		u.coinTypeCode.String(),
		xpub,
		derivationPath,
	)
	if err != nil {
		return keygenusecase.ExportFullPubkeyOutput{}, fmt.Errorf("failed to write xpub file: %w", err)
	}

	logger.Info("exported ETH account xpub",
		"account_type", input.AccountType.String(),
		"file", fileName,
	)

	return keygenusecase.ExportFullPubkeyOutput{FileName: fileName}, nil
}
