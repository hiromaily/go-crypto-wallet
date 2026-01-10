package btc

import (
	"context"
	"fmt"
	"strings"

	portsBtc "github.com/hiromaily/go-crypto-wallet/internal/application/ports/btc"
	"github.com/hiromaily/go-crypto-wallet/internal/application/ports/persistence"
	portsStorage "github.com/hiromaily/go-crypto-wallet/internal/application/ports/storage"
	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	domainAuth "github.com/hiromaily/go-crypto-wallet/internal/domain/auth"
	domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/fullpubkey"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

type importFullPubkeyUseCase struct {
	btc                portsBtc.Bitcoiner
	authFullPubKeyRepo persistence.AuthFullPubkeyRepositorier
	pubkeyFileRepo     portsStorage.AddressFileRepositorier
}

// NewImportFullPubkeyUseCase creates a new ImportFullPubkeyUseCase
func NewImportFullPubkeyUseCase(
	btc portsBtc.Bitcoiner,
	authFullPubKeyRepo persistence.AuthFullPubkeyRepositorier,
	pubkeyFileRepo portsStorage.AddressFileRepositorier,
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
	fullPubKeys := make([]*domainAuth.AuthFullPubkey, 0, len(pubKeys))
	for _, key := range pubKeys {
		inner := strings.Split(key, ",")

		fpk, err := fullpubkey.ConvertLine(u.btc.CoinTypeCode(), inner)
		if err != nil {
			return err
		}

		// Create base entity (supports both legacy and extended formats)
		// Legacy format: only FullPubKey is set
		// Extended format: ExtendedPubKey, Fingerprint, DerivationPath are also set
		authFullPubkey, err := domainAuth.NewAuthFullPubkey(
			fpk.CoinTypeCode,
			fpk.AuthType,
			fpk.FullPubKey,
		)
		if err != nil {
			return fmt.Errorf("fail to create AuthFullPubkey: %w", err)
		}

		// If extended format is detected, set additional fields
		if fpk.ExtendedPubKey != "" {
			authFullPubkey.SetExtendedPubKey(fpk.ExtendedPubKey)
		}
		if fpk.Fingerprint != "" {
			// Convert string fingerprint to domainKey.Fingerprint
			fingerprint := domainKey.Fingerprint(fpk.Fingerprint)
			authFullPubkey.SetFingerprint(fingerprint)
		}
		if fpk.DerivationPath != "" {
			authFullPubkey.SetDerivationPath(fpk.DerivationPath)
		}

		fullPubKeys = append(fullPubKeys, authFullPubkey)
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
