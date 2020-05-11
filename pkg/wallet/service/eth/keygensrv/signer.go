package keygensrv

import (
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/repository/coldrepo"
	"github.com/hiromaily/go-bitcoin/pkg/tx"
	"github.com/hiromaily/go-bitcoin/pkg/wallet"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api/ethgrp"
)

// Sign type
type Sign struct {
	eth            ethgrp.Ethereumer
	logger         *zap.Logger
	accountKeyRepo coldrepo.AccountKeyRepositorier
	authKeyRepo    coldrepo.AuthAccountKeyRepositorier
	txFileRepo     tx.FileRepositorier
	wtype          wallet.WalletType
}

// NewSign returns sign object
func NewSign(
	eth ethgrp.Ethereumer,
	logger *zap.Logger,
	accountKeyRepo coldrepo.AccountKeyRepositorier,
	authKeyRepo coldrepo.AuthAccountKeyRepositorier,
	txFileRepo tx.FileRepositorier,
	wtype wallet.WalletType) *Sign {

	return &Sign{
		eth:            eth,
		logger:         logger,
		accountKeyRepo: accountKeyRepo,
		authKeyRepo:    authKeyRepo,
		txFileRepo:     txFileRepo,
		wtype:          wtype,
	}
}

// SignTx sign on tx in csv file
// TODO: implementation
func (s *Sign) SignTx(filePath string) (string, bool, string, error) {
	return "", false, "", nil
}
