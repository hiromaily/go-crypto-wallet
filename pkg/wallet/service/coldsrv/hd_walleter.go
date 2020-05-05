package coldsrv

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	models "github.com/hiromaily/go-bitcoin/pkg/models/rdb"
	"github.com/hiromaily/go-bitcoin/pkg/repository/coldrepo"
	"github.com/hiromaily/go-bitcoin/pkg/wallet"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/coin"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/key"
)

// HDWalleter is HD wallet key generation service
type HDWalleter interface {
	Generate(accountType account.AccountType, seed []byte, count uint32) ([]key.WalletKey, error)
}

// HDWallet type
type HDWallet struct {
	logger       *zap.Logger
	repo         HDWalletRepo
	keygen       key.Generator
	coinTypeCode coin.CoinTypeCode
	wtype        wallet.WalletType
}

// NewHDWallet returns hdWallet object
func NewHDWallet(
	logger *zap.Logger,
	repo HDWalletRepo,
	keygen key.Generator,
	coinTypeCode coin.CoinTypeCode,
	wtype wallet.WalletType) *HDWallet {

	return &HDWallet{
		logger:       logger,
		repo:         repo,
		keygen:       keygen,
		coinTypeCode: coinTypeCode,
		wtype:        wtype,
	}
}

// Generate generate hd wallet keys for account
func (h *HDWallet) Generate(
	accountType account.AccountType,
	seed []byte, count uint32) ([]key.WalletKey, error) {

	h.logger.Debug("generate HDWallet", zap.String("account_type", accountType.String()))

	//get latest index
	idxFrom, err := h.repo.GetMaxIndex(accountType)
	if err != nil {
		h.logger.Info(err.Error())
		return nil, nil
	}
	h.logger.Debug("max_index",
		zap.String("account_type", accountType.String()),
		zap.Int64("current_index", idxFrom),
	)

	// generate hd wallet key
	walletKeys, err := h.generateHDKey(accountType, seed, uint32(idxFrom), count)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call key.generateAccountKeyData()")
	}

	// insert key information to account_key_table / auth_account_key_table
	err = h.repo.Insert(walletKeys, idxFrom, h.coinTypeCode, accountType)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call repo.Insert()")
	}

	return walletKeys, err
}

func (h *HDWallet) generateHDKey(
	accountType account.AccountType,
	seed []byte,
	idxFrom,
	count uint32) ([]key.WalletKey, error) {

	// generate key
	walletKeys, err := h.keygen.CreateKey(seed, accountType, idxFrom, count)
	if err != nil {
		return nil, errors.Wrap(err, "fail to call keyData.CreateKey()")
	}
	return walletKeys, nil
}

//-----------------------------------------------------------------------------
// HDWalletRepo interface
//-----------------------------------------------------------------------------

// HDWalletRepo is HDWalletRepo interface
type HDWalletRepo interface {
	GetMaxIndex(accountType account.AccountType) (int64, error)
	Insert(keys []key.WalletKey, idx int64, coinTypeCode coin.CoinTypeCode, accountType account.AccountType) error
}

//-----------------------------------------------------------------------------
// AuthHDWalletRepo
//-----------------------------------------------------------------------------

// AuthHDWalletRepo is AuthHDWalletRepo object
type AuthHDWalletRepo struct {
	authKeyRepo coldrepo.AuthAccountKeyRepositorier
	authType    account.AuthType
}

// NewAuthHDWalletRepo returns AuthHDWalletRepo
func NewAuthHDWalletRepo(authKeyRepo coldrepo.AuthAccountKeyRepositorier, authType account.AuthType) HDWalletRepo {
	return &AuthHDWalletRepo{
		authKeyRepo: authKeyRepo,
		authType:    authType,
	}
}

// GetMaxIndex returns index
func (w *AuthHDWalletRepo) GetMaxIndex(accountType account.AccountType) (int64, error) {
	_, err := w.authKeyRepo.GetOne(w.authType)
	if err != nil {
		return 0, nil
	}
	return 0, errors.New("auth key has already been created. only one record is allowed")
}

// Insert inserts key to auth_account_key table
func (w *AuthHDWalletRepo) Insert(keys []key.WalletKey, idx int64, coinTypeCode coin.CoinTypeCode, accountType account.AccountType) error {
	if len(keys) != 1 {
		return errors.New("only one key is allowed")
	}
	key := keys[0]
	item := &models.AuthAccountKey{
		Coin:               coinTypeCode.String(),
		AuthAccount:        w.authType.String(),
		P2PKHAddress:       key.P2PKHAddr,
		P2SHSegwitAddress:  key.P2SHSegWitAddr,
		FullPublicKey:      key.FullPubKey,
		MultisigAddress:    "",
		RedeemScript:       key.RedeemScript,
		WalletImportFormat: key.WIF,
		Idx:                idx,
	}

	return w.authKeyRepo.Insert(item)
}

//-----------------------------------------------------------------------------
// AccountHDWalletRepo
//-----------------------------------------------------------------------------

// AccountHDWalletRepo is AccountHDWalletRepo interface
type AccountHDWalletRepo struct {
	accountKeyRepo coldrepo.AccountKeyRepositorier
}

// NewAccountHDWalletRepo returns HDWalletRepo
func NewAccountHDWalletRepo(accountKeyRepo coldrepo.AccountKeyRepositorier) HDWalletRepo {
	return &AccountHDWalletRepo{
		accountKeyRepo: accountKeyRepo,
	}
}

// GetMaxIndex returns index
func (w *AccountHDWalletRepo) GetMaxIndex(accountType account.AccountType) (int64, error) {
	idx, err := w.accountKeyRepo.GetMaxIndex(accountType)
	if err != nil {
		return 0, nil
	}
	idx++
	return idx, nil
}

// Insert inserts key to account_key_table
func (w *AccountHDWalletRepo) Insert(keys []key.WalletKey, idxFrom int64, coinTypeCode coin.CoinTypeCode, accountType account.AccountType) error {
	// insert key information to account_key_table
	accountKeyItems := make([]*models.AccountKey, len(keys))
	for idx, key := range keys {
		accountKeyItems[idx] = &models.AccountKey{
			Coin:               coinTypeCode.String(),
			Account:            accountType.String(),
			P2PKHAddress:       key.P2PKHAddr,
			P2SHSegwitAddress:  key.P2SHSegWitAddr,
			FullPublicKey:      key.FullPubKey,
			MultisigAddress:    "",
			RedeemScript:       key.RedeemScript,
			WalletImportFormat: key.WIF,
			Idx:                idxFrom,
		}
		idxFrom++
	}
	return w.accountKeyRepo.InsertBulk(accountKeyItems)

}
