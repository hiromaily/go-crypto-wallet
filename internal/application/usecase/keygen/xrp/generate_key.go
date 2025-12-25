package xrp

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/ripple"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/ripple/xrp"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/cold"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
)

type generateKeyUseCase struct {
	xrp               ripple.Rippler
	dbConn            *sql.DB
	coinTypeCode      domainCoin.CoinTypeCode
	accountKeyRepo    cold.AccountKeyRepositorier
	xrpAccountKeyRepo cold.XRPAccountKeyRepositorier
}

// NewGenerateKeyUseCase creates a new GenerateKeyUseCase
func NewGenerateKeyUseCase(
	xrp ripple.Rippler,
	dbConn *sql.DB,
	coinTypeCode domainCoin.CoinTypeCode,
	accountKeyRepo cold.AccountKeyRepositorier,
	xrpAccountKeyRepo cold.XRPAccountKeyRepositorier,
) keygenusecase.GenerateKeyUseCase {
	return &generateKeyUseCase{
		xrp:               xrp,
		dbConn:            dbConn,
		coinTypeCode:      coinTypeCode,
		accountKeyRepo:    accountKeyRepo,
		xrpAccountKeyRepo: xrpAccountKeyRepo,
	}
}

func (u *generateKeyUseCase) Generate(ctx context.Context, input keygenusecase.GenerateKeyInput) error {
	// Convert interface{} to []domainKey.WalletKey
	walletKeys, ok := input.WalletKeys.([]domainKey.WalletKey)
	if !ok {
		return errors.New("invalid wallet keys type")
	}

	logger.Debug("generate keys for XRP",
		"account_type", input.AccountType.String(),
		"len(keys)", len(walletKeys),
	)

	// Start transaction
	dtx, err := u.dbConn.Begin()
	if err != nil {
		return fmt.Errorf("failed to call db.Begin(): %w", err)
	}
	defer func() {
		if err != nil {
			_ = dtx.Rollback() // Error already being handled
		} else {
			_ = dtx.Commit() // Error already being handled
		}
	}()

	// Generate XRP keys
	items := make([]*models.XRPAccountKey, 0, len(walletKeys))
	for _, v := range walletKeys {
		// TODO:
		// - WIF => badSeed
		// - P2PKHAddr => badSeed
		var generatedKey *xrp.ResponseWalletPropose
		generatedKey, err = u.xrp.WalletPropose(ctx, v.P2SHSegWitAddr)
		if err != nil {
			return fmt.Errorf("fail to call xrp.WalletPropose(): %w", err)
		}
		if generatedKey.Status == xrp.StatusCodeError.String() {
			return fmt.Errorf("fail to call xrp.WalletPropose() %s", generatedKey.Error)
		}

		// TODO: passphrase or related ID should be stored in table??
		items = append(items, &models.XRPAccountKey{
			Coin:             u.coinTypeCode.String(),
			Account:          input.AccountType.String(),
			AccountID:        generatedKey.Result.AccountID,
			KeyType:          xrp.GetXRPKeyTypeValue(generatedKey.Result.KeyType),
			MasterKey:        generatedKey.Result.MasterKey,
			MasterSeed:       generatedKey.Result.MasterSeed,
			MasterSeedHex:    generatedKey.Result.MasterSeedHex,
			PublicKey:        generatedKey.Result.PublicKey,
			PublicKeyHex:     generatedKey.Result.PublicKeyHex,
			IsRegularKeyPair: input.IsKeyPair,
			AllocatedID:      0,
		})

		// Update account_key table for address as ripple address
		_, err = u.accountKeyRepo.UpdateAddr(input.AccountType, generatedKey.Result.AccountID, v.P2SHSegWitAddr)
		if err != nil {
			return fmt.Errorf("fail to call accountKeyRepo.UpdateAddr(): %w", err)
		}
	}

	// Insert keys to DB
	err = u.xrpAccountKeyRepo.InsertBulk(items)
	if err != nil {
		return fmt.Errorf("fail to call xrpAccountKeyRepo.InsertBulk() for XRP: %w", err)
	}

	return nil
}
