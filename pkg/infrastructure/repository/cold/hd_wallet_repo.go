package cold

import (
	"errors"

	domainAccount "github.com/hiromaily/go-crypto-wallet/pkg/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/pkg/domain/coin"
	domainKey "github.com/hiromaily/go-crypto-wallet/pkg/domain/key"
	models "github.com/hiromaily/go-crypto-wallet/pkg/models/rdb"
)

//-----------------------------------------------------------------------------
// HDWalletRepo interface
//-----------------------------------------------------------------------------

// HDWalletRepo is an interface for HD wallet key storage operations.
// It abstracts over both AccountKeyRepository and AuthAccountKeyRepository
// to allow the same use case code to work with either account or auth keys.
type HDWalletRepo interface {
	GetMaxIndex(accountType domainAccount.AccountType) (int64, error)
	Insert(
		keys []domainKey.WalletKey,
		idx int64,
		coinTypeCode domainCoin.CoinTypeCode,
		accountType domainAccount.AccountType,
	) error
}

//-----------------------------------------------------------------------------
// AuthHDWalletRepo
//-----------------------------------------------------------------------------

// AuthHDWalletRepo implements HDWalletRepo for auth account keys
type AuthHDWalletRepo struct {
	authKeyRepo AuthAccountKeyRepositorier
	authType    domainAccount.AuthType
}

// NewAuthHDWalletRepo creates a new AuthHDWalletRepo
func NewAuthHDWalletRepo(
	authKeyRepo AuthAccountKeyRepositorier,
	authType domainAccount.AuthType,
) HDWalletRepo {
	return &AuthHDWalletRepo{
		authKeyRepo: authKeyRepo,
		authType:    authType,
	}
}

// GetMaxIndex returns index for auth keys (always 0 since only one auth key is allowed)
func (w *AuthHDWalletRepo) GetMaxIndex(_ domainAccount.AccountType) (int64, error) {
	_, err := w.authKeyRepo.GetOne(w.authType)
	if err != nil {
		return 0, nil
	}
	return 0, errors.New("auth key has already been created. only one record is allowed")
}

// Insert inserts key to auth_account_key table
func (w *AuthHDWalletRepo) Insert(
	keys []domainKey.WalletKey,
	idx int64,
	coinTypeCode domainCoin.CoinTypeCode,
	_ domainAccount.AccountType,
) error {
	if len(keys) != 1 {
		return errors.New("only one key is allowed")
	}
	keyItem := keys[0]
	item := &models.AuthAccountKey{
		Coin:               coinTypeCode.String(),
		AuthAccount:        w.authType.String(),
		P2PKHAddress:       keyItem.P2PKHAddr,
		P2SHSegwitAddress:  keyItem.P2SHSegWitAddr,
		Bech32Address:      keyItem.Bech32Addr,
		FullPublicKey:      keyItem.FullPubKey,
		MultisigAddress:    "",
		RedeemScript:       keyItem.RedeemScript,
		WalletImportFormat: keyItem.WIF,
		Idx:                idx,
	}

	return w.authKeyRepo.Insert(item)
}

//-----------------------------------------------------------------------------
// AccountHDWalletRepo
//-----------------------------------------------------------------------------

// AccountHDWalletRepo implements HDWalletRepo for account keys
type AccountHDWalletRepo struct {
	accountKeyRepo AccountKeyRepositorier
}

// NewAccountHDWalletRepo creates a new AccountHDWalletRepo
func NewAccountHDWalletRepo(accountKeyRepo AccountKeyRepositorier) HDWalletRepo {
	return &AccountHDWalletRepo{
		accountKeyRepo: accountKeyRepo,
	}
}

// GetMaxIndex returns the next available index for account keys
func (w *AccountHDWalletRepo) GetMaxIndex(accountType domainAccount.AccountType) (int64, error) {
	idx, err := w.accountKeyRepo.GetMaxIndex(accountType)
	if err != nil {
		return 0, nil
	}
	idx++
	return idx, nil
}

// Insert inserts keys to account_key_table
func (w *AccountHDWalletRepo) Insert(
	keys []domainKey.WalletKey,
	idxFrom int64,
	coinTypeCode domainCoin.CoinTypeCode,
	accountType domainAccount.AccountType,
) error {
	// insert key information to account_key_table
	accountKeyItems := make([]*models.AccountKey, len(keys))
	for idx, keyItem := range keys {
		accountKeyItems[idx] = &models.AccountKey{
			Coin:               coinTypeCode.String(),
			Account:            accountType.String(),
			P2PKHAddress:       keyItem.P2PKHAddr,
			P2SHSegwitAddress:  keyItem.P2SHSegWitAddr,
			Bech32Address:      keyItem.Bech32Addr,
			FullPublicKey:      keyItem.FullPubKey,
			MultisigAddress:    "",
			RedeemScript:       keyItem.RedeemScript,
			WalletImportFormat: keyItem.WIF,
			Idx:                idxFrom,
		}
		idxFrom++
	}
	return w.accountKeyRepo.InsertBulk(accountKeyItems)
}
