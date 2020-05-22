package signsrv

import (
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/address"
	"github.com/hiromaily/go-crypto-wallet/pkg/repository/coldrepo"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/btcgrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

// PrivKeyer is PrivKeyer service
type PrivKeyer interface {
	Import() error
}

// PrivKey type
type PrivKey struct {
	btc         btcgrp.Bitcoiner
	logger      *zap.Logger
	authKeyRepo coldrepo.AuthAccountKeyRepositorier
	authType    account.AuthType
	wtype       wallet.WalletType
}

// NewPrivKey returns privKey object
func NewPrivKey(
	btc btcgrp.Bitcoiner,
	logger *zap.Logger,
	authKeyRepo coldrepo.AuthAccountKeyRepositorier,
	authType account.AuthType,
	wtype wallet.WalletType) *PrivKey {

	return &PrivKey{
		btc:         btc,
		logger:      logger,
		authKeyRepo: authKeyRepo,
		authType:    authType,
		wtype:       wtype,
	}
}

// Import imports privKey for authKey
//  - get WIF whose `is_imported_priv_key` is false
//  - then call ImportPrivKey(wif) without rescan
func (p *PrivKey) Import() error {

	//1. retrieve records(private key) from account_key table
	authKeyItem, err := p.authKeyRepo.GetOne(p.authType)
	if err != nil {
		return errors.Wrap(err, "fail to call authKeyRepo.GetOne()")
	}
	if authKeyItem.AddrStatus != address.AddrStatusHDKeyGenerated.Int8() {
		p.logger.Info("no unimported private key")
		return nil
	}

	p.logger.Debug(
		"target records",
		zap.String("auth_type", p.authType.String()),
		zap.String("P2PKH_address", authKeyItem.P2PKHAddress),
		zap.String("P2SH_segwit_address", authKeyItem.P2SHSegwitAddress),
		zap.String("wif", authKeyItem.WalletImportFormat))
	// decode wif
	wif, err := btcutil.DecodeWIF(authKeyItem.WalletImportFormat)
	if err != nil {
		return errors.Wrapf(err, "fail to call btcutil.DecodeWIF(%s). WIF is invalid format", authKeyItem.WalletImportFormat)
	}

	// import private key by wif without rescan
	err = p.btc.ImportPrivKeyWithoutReScan(wif, p.authType.String())
	if err != nil {
		//error would be returned sometimes according to condition of bitcoin core
		//for now, it continues even if error occurred
		p.logger.Warn(
			"fail to call btc.ImportPrivKeyWithoutReScan()",
			zap.String("wif", authKeyItem.WalletImportFormat),
			zap.Error(err))
		return errors.Wrapf(err, "fail to call btc.ImportPrivKeyWithoutReScan()")
	}

	//update DB
	_, err = p.authKeyRepo.UpdateAddrStatus(address.AddrStatusPrivKeyImported, authKeyItem.WalletImportFormat)
	if err != nil {
		p.logger.Error(
			"fail to call repo.AccountKey().UpdateAddrStatus()",
			zap.String("target_table", "auth_account_key"),
			zap.String("auth_type", p.authType.String()),
			zap.String("record.WalletImportFormat", authKeyItem.WalletImportFormat),
			zap.Error(err))
	}

	// check address was stored in bitcoin core by importing private key
	p.checkImportedAddress(authKeyItem.P2PKHAddress, authKeyItem.P2SHSegwitAddress, authKeyItem.FullPublicKey)

	return nil
}

// checkImportedAddress check address was stored in bitcoin core by importing private key
// debug use
// FIXME: this code is same to keygensrv/privkey_importer.go
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
	//FIXME: error occurred in BCH
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
