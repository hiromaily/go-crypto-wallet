package btc

import (
	"context"
	"fmt"

	"github.com/btcsuite/btcd/btcutil"

	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/api/bitcoin"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/repository/cold"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/address"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

type importPrivateKeyUseCase struct {
	btc            bitcoin.Bitcoiner
	accountKeyRepo cold.AccountKeyRepositorier
}

// NewImportPrivateKeyUseCase creates a new ImportPrivateKeyUseCase
func NewImportPrivateKeyUseCase(
	btc bitcoin.Bitcoiner,
	accountKeyRepo cold.AccountKeyRepositorier,
) keygenusecase.ImportPrivateKeyUseCase {
	return &importPrivateKeyUseCase{
		btc:            btc,
		accountKeyRepo: accountKeyRepo,
	}
}

func (u *importPrivateKeyUseCase) Import(
	ctx context.Context,
	input keygenusecase.ImportPrivateKeyInput,
) error {
	// Retrieve records (private key) from account_key table with addr_status=0
	accountKeyTable, err := u.accountKeyRepo.GetAllAddrStatus(input.AccountType, address.AddrStatusHDKeyGenerated)
	if err != nil {
		return fmt.Errorf("fail to call accountKeyRepo.GetAllAddrStatus(): %w", err)
	}
	if len(accountKeyTable) == 0 {
		logger.Info("no unimported private key")
		return nil
	}

	for _, record := range accountKeyTable {
		logger.Debug(
			"target records",
			"account_type", input.AccountType.String(),
			"P2PKH_address", record.P2PKHAddress,
			"P2SH_segwit_address", record.P2SHSegwitAddress,
			"wif", record.WalletImportFormat)

		// Decode WIF
		wif, err := btcutil.DecodeWIF(record.WalletImportFormat)
		if err != nil {
			return fmt.Errorf(
				"fail to call btcutil.DecodeWIF(%s). WIF is invalid format: %w",
				record.WalletImportFormat, err)
		}

		// Import private key by WIF without rescan
		err = u.btc.ImportPrivKeyWithoutReScan(wif, input.AccountType.String())
		if err != nil {
			// Error would be returned sometimes according to condition of bitcoin core
			// For now, it continues even if error occurred
			logger.Warn(
				"fail to call btc.ImportPrivKeyWithoutReScan()",
				"wif", record.WalletImportFormat,
				"error", err)
			return err
		}

		// Update DB
		_, err = u.accountKeyRepo.UpdateAddrStatus(
			input.AccountType, address.AddrStatusPrivKeyImported, []string{record.WalletImportFormat})
		if err != nil {
			logger.Error(
				"fail to call accountKeyRepo.UpdateAddrStatus(), but privKey import is done",
				"target_table", "account_key_account",
				"account_type", input.AccountType.String(),
				"record.WalletImportFormat", record.WalletImportFormat,
				"error", err)
			return err
		}

		// Check address was stored in bitcoin core by importing private key
		u.checkImportedAddress(record.P2PKHAddress, record.P2SHSegwitAddress, record.FullPublicKey)
	}

	return nil
}

// checkImportedAddress checks if address was stored in bitcoin core by importing private key
// Debug usage
func (u *importPrivateKeyUseCase) checkImportedAddress(walletAddress, p2shSegwitAddress, fullPublicKey string) {
	// Note: GetAccount() calls GetAddressInfo() internally

	var (
		targetAddr string
		addrType   address.AddrType
	)

	switch u.btc.CoinTypeCode() {
	case domainCoin.BTC:
		targetAddr = p2shSegwitAddress
		addrType = address.AddrTypeP2shSegwit
	case domainCoin.BCH:
		targetAddr = walletAddress
		addrType = address.AddrTypeBCHCashAddr
	case domainCoin.LTC, domainCoin.ETH, domainCoin.XRP, domainCoin.ERC20, domainCoin.HYT:
		logger.Warn("this coin type is not implemented in checkImportedAddress()",
			"coin_type_code", u.btc.CoinTypeCode().String())
		return
	default:
		logger.Warn("this coin type is not implemented in checkImportedAddress()",
			"coin_type_code", u.btc.CoinTypeCode().String())
		return
	}

	// Call `getaccount` by target_address
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

	// Call `getaddressinfo` by target_address
	addrInfo, err := u.btc.GetAddressInfo(targetAddr)
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
