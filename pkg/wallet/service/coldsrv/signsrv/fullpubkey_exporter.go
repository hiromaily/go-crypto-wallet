package signsrv

import (
	"bufio"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/address"
	models "github.com/hiromaily/go-bitcoin/pkg/models/rdb"
	"github.com/hiromaily/go-bitcoin/pkg/repository/coldrepo"
	"github.com/hiromaily/go-bitcoin/pkg/wallet"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/coin"
)

// FullPubkeyExporter is FullPubkeyExporter service
type FullPubkeyExporter interface {
	ExportFullPubkey() (string, error)
}

// FullPubkeyExport type
type FullPubkeyExport struct {
	logger         *zap.Logger
	authKeyRepo    coldrepo.AuthAccountKeyRepositorier
	pubkeyFileRepo address.FileRepositorier
	coinTypeCode   coin.CoinTypeCode
	authType       account.AuthType
	wtype          wallet.WalletType
}

// NewFullPubkeyExport returns fullPubkeyExport
func NewFullPubkeyExport(
	logger *zap.Logger,
	authKeyRepo coldrepo.AuthAccountKeyRepositorier,
	pubkeyFileRepo address.FileRepositorier,
	coinTypeCode coin.CoinTypeCode,
	authType account.AuthType,
	wtype wallet.WalletType) *FullPubkeyExport {

	return &FullPubkeyExport{
		logger:         logger,
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
		return "", errors.Wrapf(err, "fail to call authKeyRepo.GetOne(%s)", f.authType.String())
	}

	// export csv file
	fileName, err := f.exportAccountKey(authKeyTable, f.authType)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

// exportAccountKey export account_key_table as csv file
func (f *FullPubkeyExport) exportAccountKey(authKeyTable *models.AuthAccountKey, authType account.AuthType) (string, error) {
	//create fileName
	fileName := f.pubkeyFileRepo.CreateFilePath(f.authType.AccountType())

	file, err := os.Create(fileName)
	if err != nil {
		return "", errors.Wrapf(err, "fail to call os.Create(%s)", fileName)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	//output: coinType, authType, fullPubkey
	if _, err = writer.WriteString(
		fmt.Sprintf("%s,%s,%s\n", f.coinTypeCode.String(), authType.String(), authKeyTable.FullPublicKey)); err != nil {
		return "", errors.Wrapf(err, "fail to call writer.WriteString(%s)", fileName)
	}
	if err = writer.Flush(); err != nil {
		return "", errors.Wrapf(err, "fail to call writer.Flush(%s)", fileName)
	}
	return fileName, nil
}
