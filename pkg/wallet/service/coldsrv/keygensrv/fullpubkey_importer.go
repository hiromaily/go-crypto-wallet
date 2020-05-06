package keygensrv

import (
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/address"
	"github.com/hiromaily/go-bitcoin/pkg/fullpubkey"
	models "github.com/hiromaily/go-bitcoin/pkg/models/rdb"
	"github.com/hiromaily/go-bitcoin/pkg/repository/coldrepo"
	"github.com/hiromaily/go-bitcoin/pkg/wallet"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api/btcgrp"
)

// FullPubKeyImporter is FullPubkeyImport service
type FullPubKeyImporter interface {
	ImportFullPubKey(fileName string) error
}

// FullPubkeyImport type
type FullPubkeyImport struct {
	btc                btcgrp.Bitcoiner
	logger             *zap.Logger
	authFullPubKeyRepo coldrepo.AuthFullPubkeyRepositorier
	pubkeyFileRepo     address.FileRepositorier
	wtype              wallet.WalletType
}

// NewFullPubkeyImport returns FullPubkeyImport object
func NewFullPubkeyImport(
	btc btcgrp.Bitcoiner,
	logger *zap.Logger,
	authFullPubKeyRepo coldrepo.AuthFullPubkeyRepositorier,
	pubkeyFileRepo address.FileRepositorier,
	wtype wallet.WalletType) *FullPubkeyImport {

	return &FullPubkeyImport{
		btc:                btc,
		logger:             logger,
		authFullPubKeyRepo: authFullPubKeyRepo,
		pubkeyFileRepo:     pubkeyFileRepo,
		wtype:              wtype,
	}
}

// ImportFullPubKey imports auth fullpubKey from csv
func (p *FullPubkeyImport) ImportFullPubKey(fileName string) error {

	// read file for full public key
	pubKeys, err := p.pubkeyFileRepo.ImportAddress(fileName)
	if err != nil {
		return errors.Wrapf(err, "fail to call fileStorager.ImportPubKey() fileName: %s", fileName)
	}

	// insert full pubKey into auth_fullpubkey_table
	fullPubKeys := make([]*models.AuthFullpubkey, len(pubKeys))
	for i, key := range pubKeys {
		inner := strings.Split(key, ",")

		fpk, err := fullpubkey.ConvertLine(p.btc.CoinTypeCode(), inner)
		if err != nil {
			return err
		}

		fullPubKeys[i] = &models.AuthFullpubkey{
			Coin:          fpk.CoinTypeCode.String(),
			AuthAccount:   fpk.AuthType.String(),
			FullPublicKey: fpk.FullPubKey,
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
