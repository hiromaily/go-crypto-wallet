package keygensrv

import (
	"fmt"
	"strings"

	"github.com/hiromaily/go-crypto-wallet/pkg/address"
	"github.com/hiromaily/go-crypto-wallet/pkg/fullpubkey"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/repository/coldrepo"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/btcgrp"
)

// FullPubkeyImport type
type FullPubkeyImport struct {
	btc                btcgrp.Bitcoiner
	logger             logger.Logger
	authFullPubKeyRepo coldrepo.AuthFullPubkeyRepositorier
	pubkeyFileRepo     address.FileRepositorier
	wtype              wallet.WalletType
}

// NewFullPubkeyImport returns FullPubkeyImport object
func NewFullPubkeyImport(
	btc btcgrp.Bitcoiner,
	logger logger.Logger,
	authFullPubKeyRepo coldrepo.AuthFullPubkeyRepositorier,
	pubkeyFileRepo address.FileRepositorier,
	wtype wallet.WalletType,
) *FullPubkeyImport {
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
		return fmt.Errorf("fail to call fileStorager.ImportPubKey() fileName: %s: %w", fileName, err)
	}

	// insert full pubKey into auth_fullpubkey_table
	fullPubKeys := make([]*models.AuthFullpubkey, len(pubKeys))
	for i, key := range pubKeys {
		inner := strings.Split(key, ",")

		var fpk *fullpubkey.FullPubKeyFormat
		fpk, err = fullpubkey.ConvertLine(p.btc.CoinTypeCode(), inner)
		if err != nil {
			return err
		}

		fullPubKeys[i] = &models.AuthFullpubkey{
			Coin:          fpk.CoinTypeCode.String(),
			AuthAccount:   fpk.AuthType.String(),
			FullPublicKey: fpk.FullPubKey,
		}
	}
	// TODO:Upsert would be better to prevent error which occur when data is already inserted
	err = p.authFullPubKeyRepo.InsertBulk(fullPubKeys)
	if err != nil {
		if strings.Contains(err.Error(), "1062: Duplicate entry") {
			p.logger.Info("full-pubkey is already imported")
		} else {
			return fmt.Errorf("fail to call authFullPubKeyRepo.InsertBulk(): %w", err)
		}
	}

	return nil
}
