package btc

import (
	"context"
	"fmt"

	"github.com/btcsuite/btcd/btcutil"

	apibtc "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/btc"
	repocold "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository/cold"
	signusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/sign"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainWallet "github.com/hiromaily/go-crypto-wallet/internal/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

type importPrivateKeyUseCase struct {
	btc         apibtc.PrivateKeyImportClient
	authKeyRepo repocold.AuthAccountKeyRepositorier
	authType    domainAccount.AuthType
	wtype       domainWallet.WalletType
}

// NewImportPrivateKeyUseCase creates a new ImportPrivateKeyUseCase for sign wallet
func NewImportPrivateKeyUseCase(
	btc apibtc.PrivateKeyImportClient,
	authKeyRepo repocold.AuthAccountKeyRepositorier,
	authType domainAccount.AuthType,
	wtype domainWallet.WalletType,
) signusecase.ImportPrivateKeyUseCase {
	return &importPrivateKeyUseCase{
		btc:         btc,
		authKeyRepo: authKeyRepo,
		authType:    authType,
		wtype:       wtype,
	}
}

func (u *importPrivateKeyUseCase) Import(ctx context.Context, input signusecase.ImportPrivateKeyInput) error {
	// 1. retrieve records(private key) from account_key table
	authKeyItem, err := u.authKeyRepo.GetOne(ctx, u.authType)
	if err != nil {
		return fmt.Errorf("fail to call authKeyRepo.GetOne(): %w", err)
	}
	if authKeyItem.AddrStatus != domainAddress.AddrStatusHDKeyGenerated {
		logger.Info("no unimported private key")
		return nil
	}

	logger.Debug(
		"target records",
		"auth_type", u.authType.String(),
		"P2PKH_address", authKeyItem.P2pkhAddress,
		"P2SH_segwit_address", authKeyItem.P2shSegwitAddress,
		"wif", authKeyItem.WalletImportFormat)

	// decode wif
	wif, err := btcutil.DecodeWIF(authKeyItem.WalletImportFormat)
	if err != nil {
		return fmt.Errorf(
			"fail to call btcutil.DecodeWIF(%s). WIF is invalid format: %w",
			authKeyItem.WalletImportFormat, err)
	}

	// import private key by wif without rescan
	err = u.btc.ImportPrivKeyWithoutReScan(wif, u.authType.String())
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
	_, err = u.authKeyRepo.UpdateAddrStatus(domainAddress.AddrStatusPrivKeyImported, authKeyItem.WalletImportFormat)
	if err != nil {
		logger.Error(
			"fail to call authKeyRepo.UpdateAddrStatus()",
			"target_table", "auth_account_key",
			"auth_type", u.authType.String(),
			"record.WalletImportFormat", authKeyItem.WalletImportFormat,
			"error", err)
	}

	// check address was stored in bitcoin core by importing private key
	u.checkImportedAddress(authKeyItem.P2pkhAddress, authKeyItem.P2shSegwitAddress, authKeyItem.FullPublicKey)

	return nil
}

// checkImportedAddress check address was stored in bitcoin core by importing private key
// debug use
func (u *importPrivateKeyUseCase) checkImportedAddress(walletAddress, p2shSegwitAddress, fullPublicKey string) {
	// Note,
	// GetAccount() calls GetAddressInfo() internally

	var (
		targetAddr string
		addrType   domainAddress.AddrType
	)

	switch u.btc.CoinTypeCode() {
	case domainCoin.BTC:
		targetAddr = p2shSegwitAddress
		addrType = domainAddress.AddrTypeP2shSegwit
	case domainCoin.BCH:
		targetAddr = walletAddress
		addrType = domainAddress.AddrTypeBCHCashAddr
	case domainCoin.LTC, domainCoin.ETH, domainCoin.XRP, domainCoin.HYT:
		logger.Warn("this coin type is not implemented in checkImportedAddress()",
			"coin_type_code", u.btc.CoinTypeCode().String())
		return
	default:
		logger.Warn("this coin type is not implemented in checkImportedAddress()",
			"coin_type_code", u.btc.CoinTypeCode().String())
		return
	}

	// 1.call `getaccount` by target_address
	// FIXME: error occurred in BCH
	acnt, err := u.btc.GetAccount(targetAddr)
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
	addrInfo, err := u.btc.GetAddressInfo(targetAddr)
	if err != nil {
		logger.Warn(
			"fail to call btc.GetAddressInfo()",
			addrType.String(), targetAddr,
			"error", err)
	} else if addrInfo.PubKey != fullPublicKey {
		logger.Warn(
			"pubkey is not matched",
			"in_bitcoin_core", addrInfo.PubKey,
			"in_database", fullPublicKey)
	}
}
