package keygensrv

import (
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/address"
	models "github.com/hiromaily/go-bitcoin/pkg/models/rdb"
	"github.com/hiromaily/go-bitcoin/pkg/repository/coldrepo"
	"github.com/hiromaily/go-bitcoin/pkg/wallet"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api"
)

// FullPubKeyImporter is FullPubkeyImport service
type FullPubKeyImporter interface {
	ImportFullPubKey(fileName string) error
}

// FullPubkeyImport type
type FullPubkeyImport struct {
	btc                api.Bitcoiner
	logger             *zap.Logger
	authFullPubKeyRepo coldrepo.AuthFullPubkeyRepositorier
	pubkeyFileRepo     address.FileStorager
	wtype              wallet.WalletType
}

// NewPubkeyImport returns pubkeyImport object
func NewFullPubkeyImport(
	btc api.Bitcoiner,
	logger *zap.Logger,
	authFullPubKeyRepo coldrepo.AuthFullPubkeyRepositorier,
	pubkeyFileRepo address.FileStorager,
	wtype wallet.WalletType) *FullPubkeyImport {

	return &FullPubkeyImport{
		btc:                btc,
		logger:             logger,
		authFullPubKeyRepo: authFullPubKeyRepo,
		pubkeyFileRepo:     pubkeyFileRepo,
		wtype:              wtype,
	}
}

// ImportPubKey imports auth pubKey from csv in keygen wallet
func (p *FullPubkeyImport) ImportFullPubKey(fileName string) error {
	//validate file name
	//if err := p.pubkeyFileRepo.ValidateFilePath(fileName, account.AccountTypeAuthorization); err != nil {
	//	return err
	//}

	// read file for full public key
	pubKeys, err := p.pubkeyFileRepo.ImportAddress(fileName)
	if err != nil {
		return errors.Wrapf(err, "fail to call fileStorager.ImportPubKey() fileName: %s", fileName)
	}

	// insert full pubKey into auth_fullpubkey_table
	fullPubKeys := make([]*models.AuthFullpubkey, len(pubKeys))
	for i, key := range pubKeys {
		inner := strings.Split(key, ",")

		fullPubKeys[i] = &models.AuthFullpubkey{
			Coin:          inner[0],
			AuthAccount:   inner[1],
			FullPublicKey: inner[2],
		}
	}
	//TODO:Upsert would be better to prevent error which occur when data is already inserted
	err = p.authFullPubKeyRepo.InsertBulk(fullPubKeys)
	if err != nil {
		if strings.Contains(err.Error(), "1062: Duplicate entry") {
			p.logger.Info("full-pubkey is already imported")
		} else {
			return errors.Wrap(err, "fail to call authFullPubKeyRepo.InsertBulk()")
		}
	}

	return nil
}
