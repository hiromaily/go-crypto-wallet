package btc

import (
	"bufio"
	"context"
	"fmt"
	"os"

	signusecase "github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/sign"
	domainAccount "github.com/hiromaily/go-crypto-wallet/pkg/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/pkg/domain/coin"
	domainWallet "github.com/hiromaily/go-crypto-wallet/pkg/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/fullpubkey"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/repository/cold"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/storage/file"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
)

type exportFullPubkeyUseCase struct {
	authKeyRepo    cold.AuthAccountKeyRepositorier
	pubkeyFileRepo file.AddressFileRepositorier
	coinTypeCode   domainCoin.CoinTypeCode
	authType       domainAccount.AuthType
	wtype          domainWallet.WalletType
}

// NewExportFullPubkeyUseCase creates a new ExportFullPubkeyUseCase for sign wallet
func NewExportFullPubkeyUseCase(
	authKeyRepo cold.AuthAccountKeyRepositorier,
	pubkeyFileRepo file.AddressFileRepositorier,
	coinTypeCode domainCoin.CoinTypeCode,
	authType domainAccount.AuthType,
	wtype domainWallet.WalletType,
) signusecase.ExportFullPubkeyUseCase {
	return &exportFullPubkeyUseCase{
		authKeyRepo:    authKeyRepo,
		pubkeyFileRepo: pubkeyFileRepo,
		coinTypeCode:   coinTypeCode,
		authType:       authType,
		wtype:          wtype,
	}
}

func (u *exportFullPubkeyUseCase) Export(ctx context.Context) (signusecase.ExportFullPubkeyOutput, error) {
	// get account key
	authKeyTable, err := u.authKeyRepo.GetOne(u.authType)
	if err != nil {
		return signusecase.ExportFullPubkeyOutput{},
			fmt.Errorf("fail to call authKeyRepo.GetOne(%s): %w", u.authType.String(), err)
	}

	// export csv file
	fileName, err := u.exportAccountKey(authKeyTable, u.authType)
	if err != nil {
		return signusecase.ExportFullPubkeyOutput{}, err
	}

	return signusecase.ExportFullPubkeyOutput{
		FileName: fileName,
	}, nil
}

// exportAccountKey export account_key_table as csv file
func (u *exportFullPubkeyUseCase) exportAccountKey(
	authKeyTable *models.AuthAccountKey, authType domainAccount.AuthType,
) (string, error) {
	// create fileName
	fileName := u.pubkeyFileRepo.CreateFilePath(u.authType.AccountType())

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

	// output: coinType, authType, fullPubkey
	_, err = writer.WriteString(fullpubkey.CreateLine(u.coinTypeCode, authType, authKeyTable.FullPublicKey))
	if err != nil {
		return "", fmt.Errorf("fail to call writer.WriteString(%s): %w", fileName, err)
	}
	if err = writer.Flush(); err != nil {
		return "", fmt.Errorf("fail to call writer.Flush(%s): %w", fileName, err)
	}
	return fileName, nil
}
