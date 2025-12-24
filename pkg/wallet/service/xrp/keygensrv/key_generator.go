package keygensrv

import (
	"context"
	"database/sql"
	"fmt"

	domainAccount "github.com/hiromaily/go-crypto-wallet/pkg/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/pkg/domain/coin"
	domainKey "github.com/hiromaily/go-crypto-wallet/pkg/domain/key"
	domainWallet "github.com/hiromaily/go-crypto-wallet/pkg/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/ripple/xrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/repository/cold"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/ripple"
)

// XRPKeyGenerator is XRP key generation service
type XRPKeyGenerator interface {
	Generate(accountType domainAccount.AccountType, isKeyPair bool, keys []domainKey.WalletKey) error
}

// XRPKeyGenerate type
type XRPKeyGenerate struct {
	xrp               ripple.Rippler
	dbConn            *sql.DB
	coinTypeCode      domainCoin.CoinTypeCode
	wtype             domainWallet.WalletType
	accountKeyRepo    cold.AccountKeyRepositorier
	xrpAccountKeyRepo cold.XRPAccountKeyRepositorier
}

// NewXRPKeyGenerate returns XRPKeyGenerate object
func NewXRPKeyGenerate(
	xrpAPI ripple.Rippler,
	dbConn *sql.DB,
	coinTypeCode domainCoin.CoinTypeCode,
	wtype domainWallet.WalletType,
	accountKeyRepo cold.AccountKeyRepositorier,
	xrpAccountKeyRepo cold.XRPAccountKeyRepositorier,
) *XRPKeyGenerate {
	return &XRPKeyGenerate{
		xrp:               xrpAPI,
		dbConn:            dbConn,
		coinTypeCode:      coinTypeCode,
		wtype:             wtype,
		accountKeyRepo:    accountKeyRepo,
		xrpAccountKeyRepo: xrpAccountKeyRepo,
	}
}

// Generate generate xrp keys for account
func (k *XRPKeyGenerate) Generate(
	accountType domainAccount.AccountType,
	isKeyPair bool,
	keys []domainKey.WalletKey,
) error {
	logger.Debug("generate keys for XRP",
		"account_type", accountType.String(),
		"len(keys)", len(keys),
	)

	// transaction
	dtx, err := k.dbConn.Begin()
	if err != nil {
		return fmt.Errorf("failed to call db.Begin(): %w", err)
	}
	defer func() {
		if err != nil {
			dtx.Rollback()
		} else {
			dtx.Commit()
		}
	}()

	// *xrp.ResponseWalletPropose
	items := make([]*models.XRPAccountKey, 0, len(keys))
	for _, v := range keys {
		// TODO:
		// - WIF => badSeed
		// - P2PKHAddr => badSeed
		var generatedKey *xrp.ResponseWalletPropose
		generatedKey, err = k.xrp.WalletPropose(context.TODO(), v.P2SHSegWitAddr)
		if err != nil {
			return fmt.Errorf("fail to call xrp.WalletPropose(): %w", err)
		}
		if generatedKey.Status == xrp.StatusCodeError.String() {
			return fmt.Errorf("fail to call xrp.WalletPropose() %s", generatedKey.Error)
		}
		// TODO: passphrase or related ID should be stored in table??
		items = append(items, &models.XRPAccountKey{
			Coin:             k.coinTypeCode.String(),
			Account:          accountType.String(),
			AccountID:        generatedKey.Result.AccountID,
			KeyType:          xrp.GetXRPKeyTypeValue(generatedKey.Result.KeyType),
			MasterKey:        generatedKey.Result.MasterKey,
			MasterSeed:       generatedKey.Result.MasterSeed,
			MasterSeedHex:    generatedKey.Result.MasterSeedHex,
			PublicKey:        generatedKey.Result.PublicKey,
			PublicKeyHex:     generatedKey.Result.PublicKeyHex,
			IsRegularKeyPair: isKeyPair,
			AllocatedID:      0,
		})

		// update account_key table for address as ripple address
		_, err = k.accountKeyRepo.UpdateAddr(accountType, generatedKey.Result.AccountID, v.P2SHSegWitAddr)
		if err != nil {
			return fmt.Errorf("fail to call accountKeyRepo.UpdateAddr(): %w", err)
		}
	}

	// insert keys to DB
	err = k.xrpAccountKeyRepo.InsertBulk(items)
	if err != nil {
		return fmt.Errorf("fail to call accountKeyRepo.InsertBulk() for XRP: %w", err)
	}

	return nil
}
