package coldwallet

import (
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/address"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/types"
)

// ImportPrivateKey imports privKey
//  - get WIF whose is_imported_priv_key is false for given account from database
//  - then call ImportPrivKey(wif) without rescan
func (w *ColdWallet) ImportPrivateKey(accountType account.AccountType) error {
	if w.wtype == types.WalletTypeWatchOnly {
		return errors.New("it's available on only coldwallet")
	}

	//1. retrieve records(private key) from account_key table
	accountKeyTable, err := w.repo.GetAllAccountKeyByAddrStatus(accountType, address.AddrStatusHDKeyGenerated) //addr_status=0
	if err != nil {
		return errors.Wrap(err, "fail to call repo.GetAllAccountKeyByAddrStatus()")
	}
	if len(accountKeyTable) == 0 {
		w.logger.Info("no unimported private key")
		return nil
	}

	for _, record := range accountKeyTable {
		w.logger.Debug(
			"target records",
			zap.String("account_type", accountType.String()),
			zap.String("wallet_address", record.WalletAddress),
			zap.String("wif", record.WalletImportFormat))
		// decode wif
		wif, err := btcutil.DecodeWIF(record.WalletImportFormat)
		if err != nil {
			return errors.Wrapf(err, "fail to call btcutil.DecodeWIF(%s). WIF is invalid format", record.WalletImportFormat)
		}

		// import private key by wif without rescan
		err = w.btc.ImportPrivKeyWithoutReScan(wif, accountType.String())
		if err != nil {
			//error would be returned sometimes according to condition of bitcoin core
			//for now, it continues even if error occurred
			w.logger.Warn(
				"fail to call btc.ImportPrivKeyWithoutReScan()",
				zap.String("wif", record.WalletImportFormat),
				zap.Error(err))
			continue
		}

		//update DB
		_, err = w.repo.UpdateAddrStatusByWIF(accountType, address.AddrStatusPrivKeyImported, record.WalletImportFormat, nil, true)
		if err != nil {
			w.logger.Error(
				"fail to update table by calling btc.UpdateAddrStatusByWIF()",
				zap.String("target_table", "account_key_account"),
				zap.String("account_type", accountType.String()),
				zap.String("record.WalletImportFormat", record.WalletImportFormat),
				zap.Error(err))
		}

		// check address was stored in bitcoin core by importing private key
		w.checkImportedAddress(record.WalletAddress, record.P2shSegwitAddress, record.FullPublicKey)
	}

	return nil
}

// checkImportedAddress check address was stored in bitcoin core by importing private key
// debug usage
func (w *ColdWallet) checkImportedAddress(walletAddress, p2shSegwitAddress, fullPublicKey string) {
	//Note,
	//GetAccount() calls GetAddressInfo() internally

	//1.call `getaccount` by wallet_address
	acnt, err := w.btc.GetAccount(walletAddress)
	if err != nil {
		w.logger.Warn(
			"fail to call btc.GetAccount()",
			zap.String("walletAddress", walletAddress),
			zap.Error(err))
	} else {
		w.logger.Debug(
			"account is found",
			zap.String("account", acnt),
			zap.String("walletAddress", walletAddress))
	}

	//2.call `getaccount` by p2sh_segwit_address
	acnt, err = w.btc.GetAccount(p2shSegwitAddress)
	if err != nil {
		w.logger.Warn(
			"fail to call btc.GetAccount()",
			zap.String("p2shSegwitAddress", p2shSegwitAddress),
			zap.Error(err))
		return
	}
	w.logger.Debug(
		"account is found by p2sh_segwit_address",
		zap.String("account", acnt),
		zap.String("p2shSegwitAddress", p2shSegwitAddress))

	//3.call `getaddressinfo` by p2sh_segwit_address
	addrInfo, err := w.btc.GetAddressInfo(p2shSegwitAddress)
	if err != nil {
		w.logger.Warn(
			"fail to call btc.GetAddressInfo()",
			zap.String("p2shSegwitAddress", p2shSegwitAddress),
			zap.Error(err))
	} else {
		if addrInfo.Pubkey != fullPublicKey {
			w.logger.Warn(
				"pubkey is not matched",
				zap.String("in_bitcoin_core", addrInfo.Pubkey),
				zap.String("in_database", fullPublicKey))
		}
	}
}
