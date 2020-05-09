package keygensrv

import (
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/address"
	"github.com/hiromaily/go-bitcoin/pkg/repository/coldrepo"
	"github.com/hiromaily/go-bitcoin/pkg/wallet"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api/btcgrp"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/coin"
)

// PrivKey type
type PrivKey struct {
	btc            btcgrp.Bitcoiner
	logger         *zap.Logger
	accountKeyRepo coldrepo.AccountKeyRepositorier
	wtype          wallet.WalletType
}

// NewPrivKey returns privKey object
func NewPrivKey(
	btc btcgrp.Bitcoiner,
	logger *zap.Logger,
	accountKeyRepo coldrepo.AccountKeyRepositorier,
	wtype wallet.WalletType) *PrivKey {

	return &PrivKey{
		btc:            btc,
		logger:         logger,
		accountKeyRepo: accountKeyRepo,
		wtype:          wtype,
	}
}

// Import imports privKey for accountKey
//  - get WIF whose `is_imported_priv_key` is false
//  - then call ImportPrivKey(wif) without rescan
func (p *PrivKey) Import(accountType account.AccountType) error {

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
			zap.String("P2PKH_address", record.P2PKHAddress),
			zap.String("P2SH_segwit_address", record.P2SHSegwitAddress),
			zap.String("wif", record.WalletImportFormat))
		// decode wif
		wif, err := btcutil.DecodeWIF(record.WalletImportFormat)
		if err != nil {
			return errors.Wrapf(err, "fail to call btcutil.DecodeWIF(%s). WIF is invalid format", record.WalletImportFormat)
		}

		// import private key by wif without rescan
		err = p.btc.ImportPrivKeyWithoutReScan(wif, accountType.String())
		if err != nil {
			//error would be returned sometimes according to condition of bitcoin core
			//for now, it continues even if error occurred
			p.logger.Warn(
				"fail to call btc.ImportPrivKeyWithoutReScan()",
				zap.String("wif", record.WalletImportFormat),
				zap.Error(err))
			continue
		}

		//update DB
		_, err = p.accountKeyRepo.UpdateAddrStatus(accountType, address.AddrStatusPrivKeyImported, []string{record.WalletImportFormat})
		if err != nil {
			p.logger.Error(
				"fail to call accountKeyRepo.UpdateAddrStatus(), but privKey import is done",
				zap.String("target_table", "account_key_account"),
				zap.String("account_type", accountType.String()),
				zap.String("record.WalletImportFormat", record.WalletImportFormat),
				zap.Error(err))
		}

		// check address was stored in bitcoin core by importing private key
		p.checkImportedAddress(record.P2PKHAddress, record.P2SHSegwitAddress, record.FullPublicKey)
	}

	return nil
}

// checkImportedAddress check address was stored in bitcoin core by importing private key
// debug usage
// FIXME: this code is same to signsrv/privkey_importer.go
func (p *PrivKey) checkImportedAddress(walletAddress, p2shSegwitAddress, fullPublicKey string) {
	//Note,
	//GetAccount() calls GetAddressInfo() internally

	var (
		targetAddr string
		addrType   address.AddrType
	)

	switch p.btc.CoinTypeCode() {
	case coin.BTC:
		targetAddr = p2shSegwitAddress
		addrType = address.AddrTypeP2shSegwit
	case coin.BCH:
		targetAddr = walletAddress
		addrType = address.AddrTypeBCHCashAddr
	default:
		p.logger.Warn("this coin type is not implemented in checkImportedAddress()",
			zap.String("coin_type_code", p.btc.CoinTypeCode().String()))
		return
	}

	// 1.call `getaccount` by target_address
	acnt, err := p.btc.GetAccount(targetAddr)
	if err != nil {
		p.logger.Warn(
			"fail to call btc.GetAccount()",
			zap.String(addrType.String(), targetAddr),
			zap.Error(err))
		return
	}
	p.logger.Debug(
		"account is found",
		zap.String("account", acnt),
		zap.String(addrType.String(), targetAddr))

	// 2.call `getaddressinfo` by target_address
	addrInfo, err := p.btc.GetAddressInfo(targetAddr)
	if err != nil {
		p.logger.Warn(
			"fail to call btc.GetAddressInfo()",
			zap.String(addrType.String(), targetAddr),
			zap.Error(err))
	} else {
		if addrInfo.Pubkey != fullPublicKey {
			p.logger.Warn(
				"pubkey is not matched",
				zap.String("in_bitcoin_core", addrInfo.Pubkey),
				zap.String("in_database", fullPublicKey))
		}
	}
}
