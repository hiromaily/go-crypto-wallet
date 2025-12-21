package keygensrv

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/pkg/errors"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/address"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	"github.com/hiromaily/go-crypto-wallet/pkg/repository/coldrepo"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/btcgrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
)

// PrivKey type
type PrivKey struct {
	btc            btcgrp.Bitcoiner
	logger         logger.Logger
	accountKeyRepo coldrepo.AccountKeyRepositorier
	wtype          wallet.WalletType
}

// NewPrivKey returns privKey object
func NewPrivKey(
	btc btcgrp.Bitcoiner,
	logger logger.Logger,
	accountKeyRepo coldrepo.AccountKeyRepositorier,
	wtype wallet.WalletType,
) *PrivKey {
	return &PrivKey{
		btc:            btc,
		logger:         logger,
		accountKeyRepo: accountKeyRepo,
		wtype:          wtype,
	}
}

// Import imports privKey for accountKey
//   - get WIF whose `is_imported_priv_key` is false
//   - then call ImportPrivKey(wif) without rescan
func (p *PrivKey) Import(accountType account.AccountType) error {
	// 1. retrieve records(private key) from account_key table
	// addr_status=0
	accountKeyTable, err := p.accountKeyRepo.GetAllAddrStatus(accountType, address.AddrStatusHDKeyGenerated)
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
			"account_type", accountType.String(),
			"P2PKH_address", record.P2PKHAddress,
			"P2SH_segwit_address", record.P2SHSegwitAddress,
			"wif", record.WalletImportFormat)
		// decode wif
		var wif *btcutil.WIF
		wif, err = btcutil.DecodeWIF(record.WalletImportFormat)
		if err != nil {
			return errors.Wrapf(
				err, "fail to call btcutil.DecodeWIF(%s). WIF is invalid format",
				record.WalletImportFormat)
		}

		// import private key by wif without rescan
		err = p.btc.ImportPrivKeyWithoutReScan(wif, accountType.String())
		if err != nil {
			// error would be returned sometimes according to condition of bitcoin core
			// for now, it continues even if error occurred
			p.logger.Warn(
				"fail to call btc.ImportPrivKeyWithoutReScan()",
				"wif", record.WalletImportFormat,
				"error", err)
			// continue
			return err
		}

		// update DB
		_, err = p.accountKeyRepo.UpdateAddrStatus(
			accountType, address.AddrStatusPrivKeyImported, []string{record.WalletImportFormat})
		if err != nil {
			p.logger.Error(
				"fail to call accountKeyRepo.UpdateAddrStatus(), but privKey import is done",
				"target_table", "account_key_account",
				"account_type", accountType.String(),
				"record.WalletImportFormat", record.WalletImportFormat,
				"error", err)
			return err
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
	// Note,
	// GetAccount() calls GetAddressInfo() internally

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
	case coin.LTC, coin.ETH, coin.XRP, coin.ERC20, coin.HYC:
		p.logger.Warn("this coin type is not implemented in checkImportedAddress()",
			"coin_type_code", p.btc.CoinTypeCode().String())
		return
	default:
		p.logger.Warn("this coin type is not implemented in checkImportedAddress()",
			"coin_type_code", p.btc.CoinTypeCode().String())
		return
	}

	// 1.call `getaccount` by target_address
	acnt, err := p.btc.GetAccount(targetAddr)
	if err != nil {
		p.logger.Warn(
			"fail to call btc.GetAccount()",
			addrType.String(), targetAddr,
			"error", err)
		return
	}
	p.logger.Debug(
		"account is found",
		"account", acnt,
		addrType.String(), targetAddr)

	// 2.call `getaddressinfo` by target_address
	addrInfo, err := p.btc.GetAddressInfo(targetAddr)
	if err != nil {
		p.logger.Warn(
			"fail to call btc.GetAddressInfo()",
			addrType.String(), targetAddr,
			"error", err)
	} else if addrInfo.Pubkey != fullPublicKey {
		p.logger.Warn(
			"pubkey is not matched",
			"in_bitcoin_core", addrInfo.Pubkey,
			"in_database", fullPublicKey)
	}
}
