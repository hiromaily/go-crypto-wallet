package signsrv

import (
	"bufio"
	"fmt"
	"os"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/address"
	"github.com/hiromaily/go-crypto-wallet/pkg/fullpubkey"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/repository/coldrepo"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

// FullPubkeyExport type
type FullPubkeyExport struct {
	authKeyRepo    coldrepo.AuthAccountKeyRepositorier
	pubkeyFileRepo address.FileRepositorier
	coinTypeCode   coin.CoinTypeCode
	authType       account.AuthType
	wtype          wallet.WalletType
}

// NewFullPubkeyExport returns fullPubkeyExport
func NewFullPubkeyExport(
	authKeyRepo coldrepo.AuthAccountKeyRepositorier,
	pubkeyFileRepo address.FileRepositorier,
	coinTypeCode coin.CoinTypeCode,
	authType account.AuthType,
	wtype wallet.WalletType,
) *FullPubkeyExport {
	return &FullPubkeyExport{
		authKeyRepo:    authKeyRepo,
		pubkeyFileRepo: pubkeyFileRepo,
		coinTypeCode:   coinTypeCode,
		authType:       authType,
		wtype:          wtype,
	}
}

// ExportFullPubkey exports full-pubkey in auth_account_key_table as csv file
func (f *FullPubkeyExport) ExportFullPubkey() (string, error) {
	// get account key
	authKeyTable, err := f.authKeyRepo.GetOne(f.authType)
	if err != nil {
		return "", fmt.Errorf("fail to call authKeyRepo.GetOne(%s): %w", f.authType.String(), err)
	}

	// export csv file
	fileName, err := f.exportAccountKey(authKeyTable, f.authType)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

// exportAccountKey export account_key_table as csv file
func (f *FullPubkeyExport) exportAccountKey(
	authKeyTable *models.AuthAccountKey, authType account.AuthType,
) (string, error) {
	// create fileName
	fileName := f.pubkeyFileRepo.CreateFilePath(f.authType.AccountType())

	file, err := os.Create(fileName) //nolint:gosec
	if err != nil {
		return "", fmt.Errorf("fail to call os.Create(%s): %w", fileName, err)
	}

	defer file.Close()

	writer := bufio.NewWriter(file)

	// output: coinType, authType, fullPubkey
	_, err = writer.WriteString(fullpubkey.CreateLine(f.coinTypeCode, authType, authKeyTable.FullPublicKey))
	if err != nil {
		return "", fmt.Errorf("fail to call writer.WriteString(%s): %w", fileName, err)
	}
	if err = writer.Flush(); err != nil {
		return "", fmt.Errorf("fail to call writer.Flush(%s): %w", fileName, err)
	}
	return fileName, nil
}
