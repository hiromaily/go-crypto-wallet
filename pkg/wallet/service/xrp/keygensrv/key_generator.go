package keygensrv

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/repository/coldrepo"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/xrpgrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/xrpgrp/xrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/key"
)

// XRPKeyGenerator is XRP key generation service
type XRPKeyGenerator interface {
	Generate(accountType account.AccountType, isKeyPair bool, keys []key.WalletKey) error
}

// XRPKeyGenerate type
type XRPKeyGenerate struct {
	xrp            xrpgrp.Rippler
	logger         *zap.Logger
	coinTypeCode   coin.CoinTypeCode
	wtype          wallet.WalletType
	accountKeyRepo coldrepo.XRPAccountKeyRepositorier
}

// NewXRPKeyGenerate returns XRPKeyGenerate object
func NewXRPKeyGenerate(
	xrp xrpgrp.Rippler,
	logger *zap.Logger,
	coinTypeCode coin.CoinTypeCode,
	wtype wallet.WalletType,
	accountKeyRepo coldrepo.XRPAccountKeyRepositorier) *XRPKeyGenerate {

	return &XRPKeyGenerate{
		xrp:            xrp,
		logger:         logger,
		coinTypeCode:   coinTypeCode,
		wtype:          wtype,
		accountKeyRepo: accountKeyRepo,
	}
}

// Generate generate xrp keys for account
func (k *XRPKeyGenerate) Generate(accountType account.AccountType, isKeyPair bool, keys []key.WalletKey) error {
	k.logger.Debug("generate keys for XRP",
		zap.String("account_type", accountType.String()),
		zap.Int("len(keys)", len(keys)),
	)

	//*xrp.ResponseWalletPropose
	items := make([]*models.XRPAccountKey, 0, len(keys))
	for _, v := range keys {
		//TODO:
		// - WIF => badSeed
		// - P2PKHAddr => badSeed
		generatedKey, err := k.xrp.WalletPropose(v.P2SHSegWitAddr)
		if err != nil {
			return errors.Wrap(err, "fail to call xrp.WalletPropose()")
		}
		if generatedKey.Status == xrp.StatusCodeError.String() {
			return errors.Errorf("fail to call xrp.WalletPropose() %s", generatedKey.Error)
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
	}

	// insert keys to DB
	if err := k.accountKeyRepo.InsertBulk(items); err != nil {
		return errors.Wrap(err, "fail to call accountKeyRepo.InsertBulk() for XRP")
	}

	return nil
}
