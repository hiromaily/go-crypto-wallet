package xrp

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	dtoxrp "github.com/hiromaily/go-crypto-wallet/internal/application/dto/xrp"
	apixrp "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/xrp"
	repocold "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository/cold"
	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	domainXRP "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/xrp"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
	apixrpimpl "github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/xrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// xrpKeyGenClient defines the interface for XRP operations needed by generateKeyUseCase.
// This follows the Interface Segregation Principle - depend only on what you need.
type xrpKeyGenClient interface {
	apixrp.KeyGenerator
}

type generateKeyUseCase struct {
	xrp               xrpKeyGenClient
	dbConn            *sql.DB
	coinTypeCode      domainCoin.CoinTypeCode
	xrpAccountKeyRepo repocold.XRPAccountKeyRepositorier
}

// NewGenerateKeyUseCase creates a new GenerateKeyUseCase.
// The xrp parameter accepts any type that implements xrpKeyGenClient (KeyGenerator).
// Typically, apixrp.XRPer is passed which implements all required methods.
func NewGenerateKeyUseCase(
	xrp xrpKeyGenClient,
	dbConn *sql.DB,
	coinTypeCode domainCoin.CoinTypeCode,
	xrpAccountKeyRepo repocold.XRPAccountKeyRepositorier,
) keygenusecase.GenerateKeyUseCase {
	return &generateKeyUseCase{
		xrp:               xrp,
		dbConn:            dbConn,
		coinTypeCode:      coinTypeCode,
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
	items := make([]*domainXRP.XRPAccountKey, 0, len(walletKeys))
	for _, v := range walletKeys {
		// TODO:
		// - WIF => badSeed
		// - P2PKHAddr => badSeed
		var generatedKey *dtoxrp.ResponseWalletPropose
		generatedKey, err = u.xrp.WalletPropose(ctx, v.P2SHSegWitAddr)
		if err != nil {
			return fmt.Errorf("fail to call xrp.WalletPropose(): %w", err)
		}
		if generatedKey.Warning != "" {
			return fmt.Errorf("fail to call xrp.WalletPropose(): %s", generatedKey.Warning)
		}

		// TODO: passphrase or related ID should be stored in table??
		xrpKey, err := domainXRP.NewXRPAccountKey(
			u.coinTypeCode,
			input.AccountType,
			generatedKey.AccountID,
			domainXRP.XRPKeyType(apixrpimpl.GetXRPKeyTypeValue(generatedKey.KeyType)),
			generatedKey.MasterSeed,
			generatedKey.MasterSeedHex,
			generatedKey.PublicKey,
			generatedKey.PublicKeyHex,
			input.IsKeyPair,
			0,
		)
		if err != nil {
			return fmt.Errorf("fail to create XrpAccountKey: %w", err)
		}

		// Set deprecated MasterKey field if present
		if generatedKey.MasterKey != "" {
			xrpKey.MasterKey = generatedKey.MasterKey
		}

		items = append(items, xrpKey)

		// TODO: Legacy cross-coin table update - removed as XRP should only use XRP repository
		// Previously this was updating the BTC account_key table with XRP address
		// This functionality may need to be reimplemented if it was intentional
	}

	// Insert keys to DB
	err = u.xrpAccountKeyRepo.InsertBulk(ctx, items)
	if err != nil {
		return fmt.Errorf("fail to call xrpAccountKeyRepo.InsertBulk() for XRP: %w", err)
	}

	return nil
}
