package keygensrv

import (
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/repository/coldrepo"
	"github.com/hiromaily/go-bitcoin/pkg/wallet"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api/ethgrp"
)

// PrivKey type
type PrivKey struct {
	eth            ethgrp.Ethereumer
	logger         *zap.Logger
	accountKeyRepo coldrepo.AccountKeyRepositorier
	wtype          wallet.WalletType
}

// NewPrivKey returns privKey object
func NewPrivKey(
	eth ethgrp.Ethereumer,
	logger *zap.Logger,
	accountKeyRepo coldrepo.AccountKeyRepositorier,
	wtype wallet.WalletType) *PrivKey {

	return &PrivKey{
		eth:            eth,
		logger:         logger,
		accountKeyRepo: accountKeyRepo,
		wtype:          wtype,
	}
}

// Import imports privKey for accountKey
// TODO: implement
func (p *PrivKey) Import(accountType account.AccountType) error {
	//ImportRawKey(hexKey, passPhrase string) (string, error)
	//p.eth.ImportRawKey()
	return nil
}
