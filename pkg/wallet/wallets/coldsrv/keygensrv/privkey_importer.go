package keygensrv

import (
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/address"
	"github.com/hiromaily/go-bitcoin/pkg/repository/coldrepo"
	"github.com/hiromaily/go-bitcoin/pkg/wallet"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api"
)

// PrivKeyer is PrivKeyer service
type PrivKeyer interface {
	Import(accountType account.AccountType) error
}

// PrivKey type
type PrivKey struct {
	btc            api.Bitcoiner
	logger         *zap.Logger
	accountKeyRepo coldrepo.AccountKeyRepositorier
	wtype          wallet.WalletType
}

// NewPrivKey returns privKey object
func NewPrivKey(
	btc api.Bitcoiner,
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
				"fail to call repo.AccountKey().UpdateAddrStatus()",
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
func (p *PrivKey) checkImportedAddress(walletAddress, p2shSegwitAddress, fullPublicKey string) {
	//Note,
	//GetAccount() calls GetAddressInfo() internally

	//1.call `getaccount` by wallet_address
	acnt, err := p.btc.GetAccount(walletAddress)
	if err != nil {
		p.logger.Warn(
			"fail to call btc.GetAccount()",
			zap.String("walletAddress", walletAddress),
			zap.Error(err))
	} else {
		p.logger.Debug(
			"account is found",
			zap.String("account", acnt),
			zap.String("walletAddress", walletAddress))
	}

	//2.call `getaccount` by p2sh_segwit_address
	acnt, err = p.btc.GetAccount(p2shSegwitAddress)
	if err != nil {
		p.logger.Warn(
			"fail to call btc.GetAccount()",
			zap.String("p2shSegwitAddress", p2shSegwitAddress),
			zap.Error(err))
		return
	}
	p.logger.Debug(
		"account is found by p2sh_segwit_address",
		zap.String("account", acnt),
		zap.String("p2shSegwitAddress", p2shSegwitAddress))

	//3.call `getaddressinfo` by p2sh_segwit_address
	addrInfo, err := p.btc.GetAddressInfo(p2shSegwitAddress)
	if err != nil {
		p.logger.Warn(
			"fail to call btc.GetAddressInfo()",
			zap.String("p2shSegwitAddress", p2shSegwitAddress),
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
