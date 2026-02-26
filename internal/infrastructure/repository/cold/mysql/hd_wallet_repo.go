package mysql

import (
	"context"
	"errors"
	"fmt"
	"time"

	repocold "github.com/hiromaily/go-crypto-wallet/internal/application/ports/repository/cold"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainAuth "github.com/hiromaily/go-crypto-wallet/internal/domain/auth"
	domainBTC "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/btc"
	domainETH "github.com/hiromaily/go-crypto-wallet/internal/domain/chains/eth"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
)

//-----------------------------------------------------------------------------
// AuthHDWalletRepo
//-----------------------------------------------------------------------------

// AuthHDWalletRepo implements HDWalletRepo for auth account keys
type AuthHDWalletRepo struct {
	authKeyRepo repocold.AuthAccountKeyRepositorier
	authType    domainAccount.AuthType
}

// NewAuthHDWalletRepo creates a new AuthHDWalletRepo
func NewAuthHDWalletRepo(
	authKeyRepo repocold.AuthAccountKeyRepositorier,
	authType domainAccount.AuthType,
) repocold.HDWalletRepo {
	return &AuthHDWalletRepo{
		authKeyRepo: authKeyRepo,
		authType:    authType,
	}
}

// GetMaxIndex returns index for auth keys (always 0 since only one auth key is allowed per account)
func (w *AuthHDWalletRepo) GetMaxIndex(ctx context.Context, accountType domainAccount.AccountType) (int64, error) {
	_, err := w.authKeyRepo.GetByAccount(w.authType, accountType)
	if err != nil {
		return 0, nil
	}
	return 0, errors.New("auth key has already been created for this account. only one record is allowed per account")
}

// Insert inserts key to auth_account_key table
func (w *AuthHDWalletRepo) Insert(
	keys []domainKey.WalletKey,
	accountXpriv string,
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

	// Store account-level extended private key for BIP32 derivation
	if accountXpriv != "" {
		item.SetAccountExtendedPrivkey(accountXpriv)
	}

	return w.authKeyRepo.Insert(item)
}

//-----------------------------------------------------------------------------
// AccountHDWalletRepo
//-----------------------------------------------------------------------------

// AccountHDWalletRepo implements HDWalletRepo for account keys
type AccountHDWalletRepo struct {
	accountKeyRepo repocold.BTCAccountKeyRepositorier
}

// NewAccountHDWalletRepo creates a new AccountHDWalletRepo
func NewAccountHDWalletRepo(accountKeyRepo repocold.BTCAccountKeyRepositorier) repocold.HDWalletRepo {
	return &AccountHDWalletRepo{
		accountKeyRepo: accountKeyRepo,
	}
}

// GetMaxIndex returns the next available index for account keys
func (w *AccountHDWalletRepo) GetMaxIndex(ctx context.Context, accountType domainAccount.AccountType) (int64, error) {
	idx, err := w.accountKeyRepo.GetMaxIndex(ctx, accountType)
	if err != nil {
		return 0, nil
	}
	idx++
	return idx, nil
}

// Insert inserts keys to btc_account_key table
func (w *AccountHDWalletRepo) Insert(
	keys []domainKey.WalletKey,
	accountXpriv string,
	idxFrom int64,
	coinTypeCode domainCoin.CoinTypeCode,
	accountType domainAccount.AccountType,
	keyType domainKey.KeyType,
) error {
	// insert key information to btc_account_key table
	accountKeyItems := make([]*domainBTC.BTCAccountKey, len(keys))
	now := time.Now()
	for idx, keyItem := range keys {
		var taprootAddr *string
		if keyItem.TaprootAddr != "" {
			taprootAddr = &keyItem.TaprootAddr
		}

		accountKeyItems[idx] = &domainBTC.BTCAccountKey{
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

		// Store account-level extended private key for BIP32 derivation
		if accountXpriv != "" {
			accountKeyItems[idx].SetAccountExtendedPrivkey(accountXpriv)
		}

		idxFrom++
	}
	return w.accountKeyRepo.InsertBulk(accountKeyItems)
}

//-----------------------------------------------------------------------------
// ETHHDWalletRepo
//-----------------------------------------------------------------------------

// ETHHDWalletRepo implements repocold.HDWalletRepo for ETH using ETHAccountKeyRepositorier.
// This allows the shared GenerateHDWalletUseCase to work with ETH key generation.
type ETHHDWalletRepo struct {
	ethAccountKeyRepo repocold.ETHAccountKeyRepositorier
}

// NewETHHDWalletRepo creates a new ETHHDWalletRepo
func NewETHHDWalletRepo(ethAccountKeyRepo repocold.ETHAccountKeyRepositorier) *ETHHDWalletRepo {
	return &ETHHDWalletRepo{
		ethAccountKeyRepo: ethAccountKeyRepo,
	}
}

// GetMaxIndex returns the maximum key index for the given account type.
func (r *ETHHDWalletRepo) GetMaxIndex(ctx context.Context, accountType domainAccount.AccountType) (int64, error) {
	return r.ethAccountKeyRepo.GetMaxIndex(ctx, accountType)
}

// Insert stores newly generated ETH wallet keys with the account-level extended private key.
func (r *ETHHDWalletRepo) Insert(
	keys []domainKey.WalletKey,
	accountXpriv string,
	idx int64,
	coinTypeCode domainCoin.CoinTypeCode,
	accountType domainAccount.AccountType,
	keyType domainKey.KeyType,
) error {
	items := make([]*domainETH.ETHAccountKey, 0, len(keys))
	for i, key := range keys {
		ethKey, err := domainETH.NewETHAccountKey(
			accountType,
			key.P2PKHAddr, // Ethereum address stored in P2PKHAddr field
			key.FullPubKey,
			key.WIF, // Private key stored in WIF field for ETH
			idx+int64(i),
		)
		if err != nil {
			return fmt.Errorf("failed to create ETH account key at index %d: %w", idx+int64(i), err)
		}
		if accountXpriv != "" {
			ethKey.SetAccountExtendedPrivkey(accountXpriv)
		}
		items = append(items, ethKey)
	}
	if err := r.ethAccountKeyRepo.InsertBulk(items); err != nil {
		return fmt.Errorf("failed to insert ETH account keys: %w", err)
	}
	return nil
}
