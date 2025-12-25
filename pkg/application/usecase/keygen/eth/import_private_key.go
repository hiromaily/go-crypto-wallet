package eth

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"

	"github.com/hiromaily/go-crypto-wallet/pkg/address"
	keygenusecase "github.com/hiromaily/go-crypto-wallet/pkg/application/usecase/keygen"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/ethereum"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/ethereum/eth"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/repository/cold"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

type importPrivateKeyUseCase struct {
	eth            ethereum.Ethereumer
	accountKeyRepo cold.AccountKeyRepositorier
}

// NewImportPrivateKeyUseCase creates a new ImportPrivateKeyUseCase
func NewImportPrivateKeyUseCase(
	eth ethereum.Ethereumer,
	accountKeyRepo cold.AccountKeyRepositorier,
) keygenusecase.ImportPrivateKeyUseCase {
	return &importPrivateKeyUseCase{
		eth:            eth,
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

	// Keystore directory is linked to any APIs to get accounts
	// So multiple directories are not good idea
	logger.Debug("NewKeyStore", "key_dir", u.eth.GetKeyDir())
	ks := keystore.NewKeyStore(u.eth.GetKeyDir(), keystore.StandardScryptN, keystore.StandardScryptP)

	for _, record := range accountKeyTable {
		logger.Debug(
			"target records",
			"account_type", input.AccountType.String(),
			"address", record.P2PKHAddress,
			"private key", record.WalletImportFormat)

		// Convert private key to ECDSA
		ecdsaKey, convertErr := u.eth.ToECDSA(record.WalletImportFormat)
		if convertErr != nil {
			logger.Warn(
				"fail to call eth.ToECDSA()",
				"private key", record.WalletImportFormat,
				"error", convertErr)
			return fmt.Errorf("fail to call eth.ToECDSA(): %w", convertErr)
		}

		// Import ECDSA key into keystore
		// FIXME: how to link imported key to specific accountName like client, deposit (grouping)
		// TODO: where password should come from
		var acct accounts.Account
		acct, err = ks.ImportECDSA(ecdsaKey, eth.Password)
		if err != nil {
			// It continues even if error occurred
			// Because database stores status, import run again by same command for this key
			logger.Warn(
				"fail to call ks.ImportECDSA()",
				"private key", record.WalletImportFormat,
				"error", err)
			return fmt.Errorf("fail to call ks.ImportECDSA(): %w", err)
		}

		logger.Debug("key account is generated",
			"account.Address.Hex()", acct.Address.Hex(),
			"account.Address.String()", acct.Address.String(),
			"account.URL.String()", acct.URL.String(),
		)

		// Check generated address
		if acct.Address.Hex() != record.P2PKHAddress {
			logger.Warn("inconsistency between generated address",
				"old_address", record.P2PKHAddress,
				"new_address", acct.Address.Hex(),
			)
		}

		// Update DB
		_, err = u.accountKeyRepo.UpdateAddrStatus(
			input.AccountType, address.AddrStatusPrivKeyImported, []string{record.WalletImportFormat})
		if err != nil {
			logger.Error(
				"fail to call accountKeyRepo.UpdateAddrStatus(), but privKey import is done",
				"target_table", "account_key_account",
				"account_type", input.AccountType.String(),
				"private key", record.WalletImportFormat,
				"error", err)
		}
	}

	return nil
}
