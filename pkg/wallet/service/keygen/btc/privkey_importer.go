package btc

import (
	"fmt"

	"github.com/btcsuite/btcd/btcutil"

	"github.com/hiromaily/go-crypto-wallet/pkg/address"
	domainAccount "github.com/hiromaily/go-crypto-wallet/pkg/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/pkg/domain/coin"
	domainWallet "github.com/hiromaily/go-crypto-wallet/pkg/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/bitcoin"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/repository/cold"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// PrivKey type
type PrivKey struct {
	btc            bitcoin.Bitcoiner
	accountKeyRepo cold.AccountKeyRepositorier
	wtype          domainWallet.WalletType
}

// NewPrivKey returns privKey object
func NewPrivKey(
	btc bitcoin.Bitcoiner,
	accountKeyRepo cold.AccountKeyRepositorier,
	wtype domainWallet.WalletType,
) *PrivKey {
	return &PrivKey{
		btc:            btc,
		accountKeyRepo: accountKeyRepo,
		wtype:          wtype,
	}
}

// Import imports privKey for accountKey
//   - get WIF whose `is_imported_priv_key` is false
//   - then call ImportPrivKey(wif) without rescan
func (p *PrivKey) Import(accountType domainAccount.AccountType) error {
	// 1. retrieve records(private key) from account_key table
	// addr_status=0
	accountKeyTable, err := p.accountKeyRepo.GetAllAddrStatus(accountType, address.AddrStatusHDKeyGenerated)
	if err != nil {
		return fmt.Errorf("fail to call repo.GetAllAccountKeyByAddrStatus(): %w", err)
	}
	if len(accountKeyTable) == 0 {
		logger.Info("no unimported private key")
		return nil
	}

	for _, record := range accountKeyTable {
		logger.Debug(
			"target records",
			"account_type", accountType.String(),
			"P2PKH_address", record.P2PKHAddress,
			"P2SH_segwit_address", record.P2SHSegwitAddress,
			"wif", record.WalletImportFormat)
		// decode wif
		var wif *btcutil.WIF
		wif, err = btcutil.DecodeWIF(record.WalletImportFormat)
		if err != nil {
			return fmt.Errorf(
				"fail to call btcutil.DecodeWIF(%s). WIF is invalid format: %w",
				record.WalletImportFormat, err)
		}

		// import private key by wif without rescan
		err = p.btc.ImportPrivKeyWithoutReScan(wif, accountType.String())
		if err != nil {
			// error would be returned sometimes according to condition of bitcoin core
			// for now, it continues even if error occurred
			logger.Warn(
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
			logger.Error(
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
	case domainCoin.BTC:
		targetAddr = p2shSegwitAddress
		addrType = address.AddrTypeP2shSegwit
	case domainCoin.BCH:
		targetAddr = walletAddress
		addrType = address.AddrTypeBCHCashAddr
	case domainCoin.LTC, domainCoin.ETH, domainCoin.XRP, domainCoin.ERC20, domainCoin.HYT:
		logger.Warn("this coin type is not implemented in checkImportedAddress()",
			"coin_type_code", p.btc.CoinTypeCode().String())
		return
	default:
		logger.Warn("this coin type is not implemented in checkImportedAddress()",
			"coin_type_code", p.btc.CoinTypeCode().String())
		return
	}

	// 1.call `getaccount` by target_address
	acnt, err := p.btc.GetAccount(targetAddr)
	if err != nil {
		logger.Warn(
			"fail to call btc.GetAccount()",
			addrType.String(), targetAddr,
			"error", err)
		return
	}
	logger.Debug(
		"account is found",
		"account", acnt,
		addrType.String(), targetAddr)

	// 2.call `getaddressinfo` by target_address
	addrInfo, err := p.btc.GetAddressInfo(targetAddr)
	if err != nil {
		logger.Warn(
			"fail to call btc.GetAddressInfo()",
			addrType.String(), targetAddr,
			"error", err)
	} else if addrInfo.Pubkey != fullPublicKey {
		logger.Warn(
			"pubkey is not matched",
			"in_bitcoin_core", addrInfo.Pubkey,
			"in_database", fullPublicKey)
	}
}
