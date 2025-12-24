package btc

import (
	"fmt"
	"strings"

	domainWallet "github.com/hiromaily/go-crypto-wallet/pkg/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/fullpubkey"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/bitcoin"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/repository/cold"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/storage/file"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
)

// FullPubkeyImport type
type FullPubkeyImport struct {
	btc                bitcoin.Bitcoiner
	authFullPubKeyRepo cold.AuthFullPubkeyRepositorier
	pubkeyFileRepo     file.AddressFileRepositorier
	wtype              domainWallet.WalletType
}

// NewFullPubkeyImport returns FullPubkeyImport object
func NewFullPubkeyImport(
	btc bitcoin.Bitcoiner,
	authFullPubKeyRepo cold.AuthFullPubkeyRepositorier,
	pubkeyFileRepo file.AddressFileRepositorier,
	wtype domainWallet.WalletType,
) *FullPubkeyImport {
	return &FullPubkeyImport{
		btc:                btc,
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
			logger.Info("full-pubkey is already imported")
		} else {
			return fmt.Errorf("fail to call authFullPubKeyRepo.InsertBulk(): %w", err)
		}
	}

	return nil
}
