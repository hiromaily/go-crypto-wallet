package signsrv

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

// PrivKeyer is PrivKeyer service
type PrivKeyer interface {
	Import() error
}

// PrivKey type
type PrivKey struct {
	btc         bitcoin.Bitcoiner
	authKeyRepo cold.AuthAccountKeyRepositorier
	authType    domainAccount.AuthType
	wtype       domainWallet.WalletType
}

// NewPrivKey returns privKey object
func NewPrivKey(
	btc bitcoin.Bitcoiner,
	authKeyRepo cold.AuthAccountKeyRepositorier,
	authType domainAccount.AuthType,
	wtype domainWallet.WalletType,
) *PrivKey {
	return &PrivKey{
		btc:         btc,
		authKeyRepo: authKeyRepo,
		authType:    authType,
		wtype:       wtype,
	}
}

// Import imports privKey for authKey
//   - get WIF whose `is_imported_priv_key` is false
//   - then call ImportPrivKey(wif) without rescan
func (p *PrivKey) Import() error {
	// 1. retrieve records(private key) from account_key table
	authKeyItem, err := p.authKeyRepo.GetOne(p.authType)
	if err != nil {
		return fmt.Errorf("fail to call authKeyRepo.GetOne(): %w", err)
	}
	if authKeyItem.AddrStatus != address.AddrStatusHDKeyGenerated.Int8() {
		logger.Info("no unimported private key")
		return nil
	}

	logger.Debug(
		"target records",
		"auth_type", p.authType.String(),
		"P2PKH_address", authKeyItem.P2PKHAddress,
		"P2SH_segwit_address", authKeyItem.P2SHSegwitAddress,
		"wif", authKeyItem.WalletImportFormat)
	// decode wif
	wif, err := btcutil.DecodeWIF(authKeyItem.WalletImportFormat)
	if err != nil {
		return fmt.Errorf(
			"fail to call btcutil.DecodeWIF(%s). WIF is invalid format: %w",
			authKeyItem.WalletImportFormat, err)
	}

	// import private key by wif without rescan
	err = p.btc.ImportPrivKeyWithoutReScan(wif, p.authType.String())
	if err != nil {
		// error would be returned sometimes according to condition of bitcoin core
		// for now, it continues even if error occurred
		logger.Warn(
			"fail to call btc.ImportPrivKeyWithoutReScan()",
			"wif", authKeyItem.WalletImportFormat,
			"error", err)
		return fmt.Errorf("fail to call btc.ImportPrivKeyWithoutReScan(): %w", err)
	}

	// update DB
	_, err = p.authKeyRepo.UpdateAddrStatus(address.AddrStatusPrivKeyImported, authKeyItem.WalletImportFormat)
	if err != nil {
		logger.Error(
			"fail to call repo.AccountKey().UpdateAddrStatus()",
			"target_table", "auth_account_key",
			"auth_type", p.authType.String(),
			"record.WalletImportFormat", authKeyItem.WalletImportFormat,
			"error", err)
	}

	// check address was stored in bitcoin core by importing private key
	p.checkImportedAddress(authKeyItem.P2PKHAddress, authKeyItem.P2SHSegwitAddress, authKeyItem.FullPublicKey)

	return nil
}

// checkImportedAddress check address was stored in bitcoin core by importing private key
// debug use
// FIXME: this code is same to keygensrv/privkey_importer.go
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
	// FIXME: error occurred in BCH
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
