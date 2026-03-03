package eth

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"

	apieth "github.com/hiromaily/go-crypto-wallet/internal/application/ports/api/eth"
	repocold "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository/cold"
	keygenusecase "github.com/hiromaily/go-crypto-wallet/internal/application/usecase/keygen"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	pkgeth "github.com/hiromaily/go-crypto-wallet/pkg/chains/eth"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

type importPrivateKeyUseCase struct {
	eth              apieth.ETHKeyAccessor
	accountKeyRepo   repocold.ETHAccountKeyRepositorier
	keystorePassword string
}

// NewImportPrivateKeyUseCase creates a new ImportPrivateKeyUseCase
func NewImportPrivateKeyUseCase(
	eth apieth.ETHKeyAccessor,
	accountKeyRepo repocold.ETHAccountKeyRepositorier,
	keystorePassword string,
) keygenusecase.ImportPrivateKeyUseCase {
	return &importPrivateKeyUseCase{
		eth:              eth,
		accountKeyRepo:   accountKeyRepo,
		keystorePassword: keystorePassword,
	}
}

// Import imports private keys from database to local keystore
// This implementation is Anvil-compatible as it uses local keystore.ImportECDSA()
// instead of the personal_importRawKey RPC (which Anvil doesn't support).
// Private keys are written directly to the local filesystem keystore directory.
func (u *importPrivateKeyUseCase) Import(
	ctx context.Context,
	input keygenusecase.ImportPrivateKeyInput,
) error {
	// Retrieve records (private key) from account_key table with addr_status=0
	accountKeyTable, err := u.accountKeyRepo.GetAllAddrStatus(input.AccountType, domainAddress.AddrStatusHDKeyGenerated)
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
			"address", record.Address,
			"private key", record.PrivateKey)

		// Convert private key to ECDSA
		ecdsaKey, convertErr := pkgeth.ToECDSA(record.PrivateKey)
		if convertErr != nil {
			logger.Warn(
				"fail to call pkgeth.ToECDSA()",
				"private key", record.PrivateKey,
				"error", convertErr)
			return fmt.Errorf("fail to call pkgeth.ToECDSA(): %w", convertErr)
		}

		// Import ECDSA key into keystore
		// FIXME: how to link imported key to specific accountName like client, deposit (grouping)
		var acct accounts.Account
		acct, err = ks.ImportECDSA(ecdsaKey, u.keystorePassword)
		if err != nil {
			// It continues even if error occurred
			// Because database stores status, import run again by same command for this key
			logger.Warn(
				"fail to call ks.ImportECDSA()",
				"private key", record.PrivateKey,
				"error", err)
			return fmt.Errorf("fail to call ks.ImportECDSA(): %w", err)
		}

		logger.Debug("key account is generated",
			"account.Address.Hex()", acct.Address.Hex(),
			"account.Address.String()", acct.Address.String(),
			"account.URL.String()", acct.URL.String(),
		)

		// Check generated address
		if acct.Address.Hex() != record.Address {
			logger.Warn("inconsistency between generated address",
				"old_address", record.Address,
				"new_address", acct.Address.Hex(),
			)
		}

		// Update DB
		_, err = u.accountKeyRepo.UpdateAddrStatus(
			input.AccountType, domainAddress.AddrStatusPrivKeyImported, []string{record.PrivateKey})
		if err != nil {
			logger.Error(
				"fail to call accountKeyRepo.UpdateAddrStatus(), but privKey import is done",
				"target_table", "eth_account_key",
				"account_type", input.AccountType.String(),
				"private key", record.PrivateKey,
				"error", err)
		}
	}

	return nil
}
