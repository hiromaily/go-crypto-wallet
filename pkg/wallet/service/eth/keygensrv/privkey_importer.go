package keygensrv

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/address"
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
// FIXME: this code is similar to service/btc/coldsrv/privkey_importer.go
//  - it may be better to integrate
func (p *PrivKey) Import(accountType account.AccountType) error {
	//ImportRawKey(hexKey, passPhrase string) (string, error)
	//p.eth.ImportRawKey()
	//1. retrieve records(private key) from account_key table
	accountKeyTable, err := p.accountKeyRepo.GetAllAddrStatus(accountType, address.AddrStatusHDKeyGenerated) //addr_status=0
	if err != nil {
		return errors.Wrap(err, "fail to call repo.GetAllAccountKeyByAddrStatus()")
	}
	if len(accountKeyTable) == 0 {
		p.logger.Info("no unimported private key")
		return nil
	}

	for _, record := range accountKeyTable {
		p.logger.Debug(
			"target records",
			zap.String("account_type", accountType.String()),
			zap.String("address", record.P2PKHAddress),
			zap.String("private key", record.WalletImportFormat))

		// format privkey
		// FIXME: how to link imported key to specific accountName like client, deposit (grouping)
		// TODO: where password should come from
		generatedAddr, err := p.eth.ImportRawKey(record.WalletImportFormat, "password")
		if err != nil {
			// it continues even if error occurred
			// because database stores status, import run again by same command for this key
			p.logger.Warn(
				"fail to call eth.ImportRawKey()",
				zap.String("private key", record.WalletImportFormat),
				zap.Error(err))
			continue
		}
		// check generated address
		if generatedAddr != record.P2PKHAddress {
			p.logger.Warn("inconsistency between generated address",
				zap.String("old_address", record.P2PKHAddress),
				zap.String("new_address", generatedAddr),
			)
		}

		//update DB
		_, err = p.accountKeyRepo.UpdateAddrStatus(accountType, address.AddrStatusPrivKeyImported, []string{record.WalletImportFormat})
		if err != nil {
			p.logger.Error(
				"fail to call repo.AccountKey().UpdateAddrStatus()",
				zap.String("target_table", "account_key_account"),
				zap.String("account_type", accountType.String()),
				zap.String("private key", record.WalletImportFormat),
				zap.Error(err))
		}
	}

	return nil
}
