package eth

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"

	"github.com/hiromaily/go-crypto-wallet/pkg/address"
	domainAccount "github.com/hiromaily/go-crypto-wallet/pkg/domain/account"
	domainWallet "github.com/hiromaily/go-crypto-wallet/pkg/domain/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/ethereum"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/api/ethereum/eth"
	"github.com/hiromaily/go-crypto-wallet/pkg/infrastructure/repository/cold"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// PrivKey type
type PrivKey struct {
	eth            ethereum.Ethereumer
	accountKeyRepo cold.AccountKeyRepositorier
	wtype          domainWallet.WalletType
}

// NewPrivKey returns privKey object
func NewPrivKey(
	ethAPI ethereum.Ethereumer,
	accountKeyRepo cold.AccountKeyRepositorier,
	wtype domainWallet.WalletType,
) *PrivKey {
	return &PrivKey{
		eth:            ethAPI,
		accountKeyRepo: accountKeyRepo,
		wtype:          wtype,
	}
}

// Import imports privKey for accountKey for ETH
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

	// keystore directory is linked to any apis to get accounts
	// so multiple directories are not good idea
	logger.Debug("NewKeyStore", "key_dir", p.eth.GetKeyDir())
	// keyDir := fmt.Sprintf("%s/%s", p.keyDir, accountType.String())
	ks := keystore.NewKeyStore(p.eth.GetKeyDir(), keystore.StandardScryptN, keystore.StandardScryptP)

	for _, record := range accountKeyTable {
		logger.Debug(
			"target records",
			"account_type", accountType.String(),
			"address", record.P2PKHAddress,
			"private key", record.WalletImportFormat)

		// generatedAddr, err := p.eth.ImportRawKey(record.WalletImportFormat, "password")
		ecdsaKey, convertErr := p.eth.ToECDSA(record.WalletImportFormat)
		if convertErr != nil {
			logger.Warn(
				"fail to call key.ToECDSA()",
				"private key", record.WalletImportFormat,
				"error", convertErr)
			// continue
			return fmt.Errorf("fail to call key.ToECDSA(): %w", convertErr)
		}
		// FIXME: how to link imported key to specific accountName like client, deposit (grouping)
		// TODO: where password should come from // ImportRawKey(hexKey, passPhrase string) (string, error)
		var acct accounts.Account
		acct, err = ks.ImportECDSA(ecdsaKey, eth.Password)
		if err != nil {
			// it continues even if error occurred
			// because database stores status, import run again by same command for this key
			logger.Warn(
				"fail to call eth.ImportECDSA()",
				"private key", record.WalletImportFormat,
				"error", err)
			// continue
			return fmt.Errorf("fail to call eth.ImportECDSA(): %w", err)
		}
		logger.Debug("key account is generated",
			"account.Address.Hex()", acct.Address.Hex(),
			"account.Address.String()", acct.Address.String(),
			"account.URL.String()", acct.URL.String(),
		)

		// check generated address
		if acct.Address.Hex() != record.P2PKHAddress {
			logger.Warn("inconsistency between generated address",
				"old_address", record.P2PKHAddress,
				"new_address", acct.Address.Hex(),
			)
		}

		// update DB
		_, err = p.accountKeyRepo.UpdateAddrStatus(
			accountType, address.AddrStatusPrivKeyImported, []string{record.WalletImportFormat})
		if err != nil {
			logger.Error(
				"fail to call accountKeyRepo.UpdateAddrStatus(), but privKey import is done",
				"target_table", "account_key_account",
				"account_type", accountType.String(),
				"private key", record.WalletImportFormat,
				"error", err)
		}
	}

	return nil
}
