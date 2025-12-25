package btc

import (
	"context"
	"fmt"
	"strings"

	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/bitcoin"
	models "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/database/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/cold"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/fullpubkey"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

type importFullPubkeyUseCase struct {
	btc                bitcoin.Bitcoiner
	authFullPubKeyRepo cold.AuthFullPubkeyRepositorier
	pubkeyFileRepo     file.AddressFileRepositorier
}

// NewImportFullPubkeyUseCase creates a new ImportFullPubkeyUseCase
func NewImportFullPubkeyUseCase(
	btc bitcoin.Bitcoiner,
	authFullPubKeyRepo cold.AuthFullPubkeyRepositorier,
	pubkeyFileRepo file.AddressFileRepositorier,
) keygenusecase.ImportFullPubkeyUseCase {
	return &importFullPubkeyUseCase{
		btc:                btc,
		authFullPubKeyRepo: authFullPubKeyRepo,
		pubkeyFileRepo:     pubkeyFileRepo,
	}
}

func (u *importFullPubkeyUseCase) Import(
	ctx context.Context,
	input keygenusecase.ImportFullPubkeyInput,
) error {
	// Read file for full public key
	pubKeys, err := u.pubkeyFileRepo.ImportAddress(input.FileName)
	if err != nil {
		return fmt.Errorf("fail to call pubkeyFileRepo.ImportAddress() fileName: %s: %w", input.FileName, err)
	}

	// Insert full pubKey into auth_fullpubkey_table
	fullPubKeys := make([]*models.AuthFullpubkey, len(pubKeys))
	for i, key := range pubKeys {
		inner := strings.Split(key, ",")

		fpk, err := fullpubkey.ConvertLine(u.btc.CoinTypeCode(), inner)
		if err != nil {
			return err
		}

		fullPubKeys[i] = &models.AuthFullpubkey{
			Coin:          fpk.CoinTypeCode.String(),
			AuthAccount:   fpk.AuthType.String(),
			FullPublicKey: fpk.FullPubKey,
		}
	}

	// TODO: Upsert would be better to prevent error which occur when data is already inserted
	err = u.authFullPubKeyRepo.InsertBulk(fullPubKeys)
	if err != nil {
		if strings.Contains(err.Error(), "1062: Duplicate entry") {
			logger.Info("full-pubkey is already imported")
		} else {
			return fmt.Errorf("fail to call authFullPubKeyRepo.InsertBulk(): %w", err)
		}
	}

	return nil
}
