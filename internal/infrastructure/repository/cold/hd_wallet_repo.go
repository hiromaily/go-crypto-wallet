package cold

import (
	"errors"
	"time"

	"github.com/hiromaily/go-crypto-wallet/internal/application/ports/persistence"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainAuth "github.com/hiromaily/go-crypto-wallet/internal/domain/auth"
	domainBitcoin "github.com/hiromaily/go-crypto-wallet/internal/domain/bitcoin"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
)

//-----------------------------------------------------------------------------
// AuthHDWalletRepo
//-----------------------------------------------------------------------------

// AuthHDWalletRepo implements HDWalletRepo for auth account keys
type AuthHDWalletRepo struct {
	authKeyRepo persistence.AuthAccountKeyRepositorier
	authType    domainAccount.AuthType
}

// NewAuthHDWalletRepo creates a new AuthHDWalletRepo
func NewAuthHDWalletRepo(
	authKeyRepo persistence.AuthAccountKeyRepositorier,
	authType domainAccount.AuthType,
) persistence.HDWalletRepo {
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
	accountType domainAccount.AccountType,
	keyType domainKey.KeyType,
) error {
	if len(keys) != 1 {
		return errors.New("only one key is allowed")
	}
	keyItem := keys[0]

	item, err := domainAuth.NewAuthAccountKey(
		coinTypeCode,
		keyType.String(),
		w.authType,
		accountType,
		keyItem.P2PKHAddr,
		keyItem.P2SHSegWitAddr,
		keyItem.Bech32Addr,
		keyItem.FullPubKey,
		keyItem.WIF,
		idx,
	)
	if err != nil {
		return err
	}

	if keyItem.TaprootAddr != "" {
		item.SetTaprootAddress(keyItem.TaprootAddr)
	}
	if keyItem.RedeemScript != "" {
		item.RedeemScript = keyItem.RedeemScript
	}

	return w.authKeyRepo.Insert(item)
}

//-----------------------------------------------------------------------------
// AccountHDWalletRepo
//-----------------------------------------------------------------------------

// AccountHDWalletRepo implements HDWalletRepo for account keys
type AccountHDWalletRepo struct {
	accountKeyRepo persistence.BTCAccountKeyRepositorier
}

// NewAccountHDWalletRepo creates a new AccountHDWalletRepo
func NewAccountHDWalletRepo(accountKeyRepo persistence.BTCAccountKeyRepositorier) persistence.HDWalletRepo {
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

// Insert inserts keys to btc_account_key table
func (w *AccountHDWalletRepo) Insert(
	keys []domainKey.WalletKey,
	idxFrom int64,
	coinTypeCode domainCoin.CoinTypeCode,
	accountType domainAccount.AccountType,
	keyType domainKey.KeyType,
) error {
	// insert key information to btc_account_key table
	accountKeyItems := make([]*domainBitcoin.BtcAccountKey, len(keys))
	now := time.Now()
	for idx, keyItem := range keys {
		var taprootAddr *string
		if keyItem.TaprootAddr != "" {
			taprootAddr = &keyItem.TaprootAddr
		}

		accountKeyItems[idx] = &domainBitcoin.BtcAccountKey{
			CoinTypeCode:       coinTypeCode,
			KeyType:            keyType.String(),
			Account:            accountType,
			P2pkhAddress:       keyItem.P2PKHAddr,
			P2shSegwitAddress:  keyItem.P2SHSegWitAddr,
			Bech32Address:      keyItem.Bech32Addr,
			TaprootAddress:     taprootAddr,
			FullPublicKey:      keyItem.FullPubKey,
			MultisigAddress:    "",
			RedeemScript:       keyItem.RedeemScript,
			WalletImportFormat: keyItem.WIF,
			Idx:                idxFrom,
			AddrStatus:         domainAddress.AddrStatusHDKeyGenerated,
			UpdatedAt:          &now,
		}
		idxFrom++
	}
	return w.accountKeyRepo.InsertBulk(accountKeyItems)
}
