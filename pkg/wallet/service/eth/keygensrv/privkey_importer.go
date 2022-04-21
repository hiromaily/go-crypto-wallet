package keygensrv

import (
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/address"
	"github.com/hiromaily/go-crypto-wallet/pkg/repository/coldrepo"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/ethgrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/ethgrp/eth"
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
	wtype wallet.WalletType,
) *PrivKey {
	return &PrivKey{
		eth:            eth,
		logger:         logger,
		accountKeyRepo: accountKeyRepo,
		wtype:          wtype,
	}
}

// Import imports privKey for accountKey for ETH
func (p *PrivKey) Import(accountType account.AccountType) error {
	// 1. retrieve records(private key) from account_key table
	accountKeyTable, err := p.accountKeyRepo.GetAllAddrStatus(accountType, address.AddrStatusHDKeyGenerated) // addr_status=0
	if err != nil {
		return errors.Wrap(err, "fail to call repo.GetAllAccountKeyByAddrStatus()")
	}
	if len(accountKeyTable) == 0 {
		p.logger.Info("no unimported private key")
		return nil
	}

	// keystore directory is linked to any apis to get accounts
	// so multiple directories are not good idea
	p.logger.Debug("NewKeyStore", zap.String("key_dir", p.eth.GetKeyDir()))
	// keyDir := fmt.Sprintf("%s/%s", p.keyDir, accountType.String())
	ks := keystore.NewKeyStore(p.eth.GetKeyDir(), keystore.StandardScryptN, keystore.StandardScryptP)

	for _, record := range accountKeyTable {
		p.logger.Debug(
			"target records",
			zap.String("account_type", accountType.String()),
			zap.String("address", record.P2PKHAddress),
			zap.String("private key", record.WalletImportFormat))

		// generatedAddr, err := p.eth.ImportRawKey(record.WalletImportFormat, "password")
		ecdsaKey, err := p.eth.ToECDSA(record.WalletImportFormat)
		if err != nil {
			p.logger.Warn(
				"fail to call key.ToECDSA()",
				zap.String("private key", record.WalletImportFormat),
				zap.Error(err))
			// continue
			return errors.Wrap(err, "fail to call key.ToECDSA()")
		}
		// FIXME: how to link imported key to specific accountName like client, deposit (grouping)
		// TODO: where password should come from // ImportRawKey(hexKey, passPhrase string) (string, error)
		account, err := ks.ImportECDSA(ecdsaKey, eth.Password)
		if err != nil {
			// it continues even if error occurred
			// because database stores status, import run again by same command for this key
			p.logger.Warn(
				"fail to call eth.ImportECDSA()",
				zap.String("private key", record.WalletImportFormat),
				zap.Error(err))
			// continue
			return errors.Wrap(err, "fail to call eth.ImportECDSA()")
		}
		p.logger.Debug("key account is generated",
			zap.String("account.Address.Hex()", account.Address.Hex()),
			zap.String("account.Address.String()", account.Address.String()),
			zap.String("account.URL.String()", account.URL.String()),
		)

		// check generated address
		if account.Address.Hex() != record.P2PKHAddress {
			p.logger.Warn("inconsistency between generated address",
				zap.String("old_address", record.P2PKHAddress),
				zap.String("new_address", account.Address.Hex()),
			)
		}

		// update DB
		_, err = p.accountKeyRepo.UpdateAddrStatus(accountType, address.AddrStatusPrivKeyImported, []string{record.WalletImportFormat})
		if err != nil {
			p.logger.Error(
				"fail to call accountKeyRepo.UpdateAddrStatus(), but privKey import is done",
				zap.String("target_table", "account_key_account"),
				zap.String("account_type", accountType.String()),
				zap.String("private key", record.WalletImportFormat),
				zap.Error(err))
		}
	}

	return nil
}
