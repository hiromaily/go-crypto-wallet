package keygensrv

import (
	"go.uber.org/zap"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/repository/coldrepo"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/xrpgrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

// XRPKeyGenerator is XRP key generation service
type XRPKeyGenerator interface {
	Generate(accountType account.AccountType, isKeyPair bool) error
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
func (k *XRPKeyGenerate) Generate(accountType account.AccountType, isKeyPair bool) error {
	k.logger.Debug("generate keys for XRP", zap.String("account_type", accountType.String()))
	return nil
}
