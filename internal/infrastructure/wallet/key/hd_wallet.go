package key

import (
	"fmt"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"

	portsWallet "github.com/hiromaily/go-crypto-wallet/internal/application/ports/wallet"
	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
)

// TODO: except client address, same address is used only once due to security
// - address could be traced easily
// - so any internal addresses should be disposable

// BIP44
// https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki#Purpose
// m / purpose' / coin_type' / account' / change / address_index

// e.g.  BIP44, Bitcoin Mainnet
// Client account  => m/44/0/0/0/0~xxxxx
// Receipt account => m/44/0/1/0/0~xxxxx
// Payment account => m/44/0/2/0/0~xxxxx

// PurposeType BIP44/BIP49, for now 44 is used as fixed value
type PurposeType uint32

// Uint32 converter
func (t PurposeType) Uint32() uint32 {
	return uint32(t)
}

// purpose depends on BIP, BIP44  is a constant set to `44`
const (
	PurposeTypeBIP44 PurposeType = 44 // BIP44
	PurposeTypeBIP49 PurposeType = 49 // BIP49
	PurposeTypeBIP84 PurposeType = 84 // BIP84
	PurposeTypeBIP86 PurposeType = 86 // BIP86
)

// CoinType creates a separate subtree for every cryptocoin
//
//	which come from `CoinType` in go-crypto-wallet/pkg/wallet/coin/types.go
type CoinType uint32

// Uint32 converter
func (t CoinType) Uint32() uint32 {
	return uint32(t)
}

// coin_type
// https://github.com/satoshilabs/slips/blob/master/slip-0044.md

// Account
// account come from `AccountType` in go-crypto-wallet/pkg/account/public_account.go

// ChangeType  external or internal use
type ChangeType uint32

// Uint32 converter
func (t ChangeType) Uint32() uint32 {
	return uint32(t)
}

// change_type
const (
	ChangeTypeExternal ChangeType = 0 // constant 0 is used for external chain
	ChangeTypeInternal ChangeType = 1 // constant 1 for internal chain (also known as change addresses)
)

// maxMultisigAccountIndex is the maximum account index for multisig accounts.
// Accounts with indices 0, 1, and 2 (deposit, payment, stored) are used for multisig
// and require non-hardened derivation to match xpub-derived keys.
const maxMultisigAccountIndex = 2

// HDKey HD Wallet Key object
type HDKey struct {
	purpose      PurposeType
	coinType     domainCoin.CoinType
	coinTypeCode domainCoin.CoinTypeCode
	conf         *chaincfg.Params
	strategy     portsWallet.CoinKeyStrategy // Strategy for coin-specific key generation
}

// NewHDKey returns Key
func NewHDKey(
	purpose PurposeType,
	coinTypeCode domainCoin.CoinTypeCode,
	conf *chaincfg.Params,
	strategy portsWallet.CoinKeyStrategy,
) *HDKey {
	// Validate that strategy is not nil
	if strategy == nil {
		panic(fmt.Sprintf(
			"strategy cannot be nil for coin %s (purpose: %d)",
			coinTypeCode.String(),
			purpose,
		))
	}

	keyData := HDKey{
		purpose:      purpose,
		coinType:     domainCoin.GetCoinType(coinTypeCode, conf),
		coinTypeCode: coinTypeCode,
		conf:         conf,
		strategy:     strategy,
	}

	return &keyData
}

// CreateKey create hd key
func (k *HDKey) CreateKey(
	seed []byte,
	accountType domainAccount.AccountType,
	idxFrom, count uint32,
) ([]domainKey.WalletKey, error) {
	// create privateKey, publicKey by account level
	privKey, _, err := k.createKeyByAccount(seed, accountType)
	if err != nil {
		return nil, fmt.Errorf("fail to call createKeyByAccount(): %w", err)
	}
	// create keys by index and count
	return k.createKeysWithIndex(privKey, idxFrom, count)
}

// CreateKeyWithAccountXpriv creates HD keys and returns the account-level extended private key.
// This extended private key can be used for BIP32 derivation at different indices.
func (k *HDKey) CreateKeyWithAccountXpriv(
	seed []byte,
	accountType domainAccount.AccountType,
	idxFrom, count uint32,
) (keys []domainKey.WalletKey, accountXpriv string, err error) {
	// create privateKey, publicKey by account level
	privKey, _, err := k.createKeyByAccount(seed, accountType)
	if err != nil {
		return nil, "", fmt.Errorf("fail to call createKeyByAccount(): %w", err)
	}

	// Get account-level extended private key string
	accountXpriv = privKey.String()

	// create keys by index and count
	keys, err = k.createKeysWithIndex(privKey, idxFrom, count)
	if err != nil {
		return nil, "", err
	}

	return keys, accountXpriv, nil
}

// KeyType returns the key type this generator supports (implements Generator interface)
func (k *HDKey) KeyType() domainKey.KeyType {
	switch k.purpose {
	case PurposeTypeBIP44:
		return domainKey.KeyTypeBIP44
	case PurposeTypeBIP49:
		return domainKey.KeyTypeBIP49
	case PurposeTypeBIP84:
		return domainKey.KeyTypeBIP84
	case PurposeTypeBIP86:
		return domainKey.KeyTypeBIP86
	default:
		return domainKey.KeyTypeBIP44
	}
}

// SupportsAddressType checks if this generator supports the given address type (implements Generator interface)
func (k *HDKey) SupportsAddressType(addrType domainAddress.AddrType) bool {
	switch k.purpose {
	case PurposeTypeBIP44:
		return addrType == domainAddress.AddrTypeLegacy
	case PurposeTypeBIP49:
		return addrType == domainAddress.AddrTypeP2shSegwit
	case PurposeTypeBIP84:
		return addrType == domainAddress.AddrTypeBech32
	case PurposeTypeBIP86:
		return addrType == domainAddress.AddrTypeTaproot
	default:
		return false
	}
}

// buildDerivationPath builds a BIP32 derivation path string
// Format: m/purpose'/coinType'/account'/0/index
// This helper function ensures consistency across all BIP generators
func buildDerivationPath(purpose, coinType, account, index uint32) string {
	return fmt.Sprintf("m/%d'/%d'/%d'/0/%d", purpose, coinType, account, index)
}

// GetDerivationPath returns the derivation path for the given account and index (implements Generator interface)
func (k *HDKey) GetDerivationPath(accountType domainAccount.AccountType, index uint32) string {
	return buildDerivationPath(
		k.purpose.Uint32(),
		k.coinType.Uint32(),
		accountType.Uint32(),
		index,
	)
}

// createKeyByAccount create privateKey, publicKey by account level
func (k *HDKey) createKeyByAccount(
	seed []byte, accountType domainAccount.AccountType,
) (*hdkeychain.ExtendedKey, *hdkeychain.ExtendedKey, error) {
	// Master
	masterKey, err := hdkeychain.NewMaster(seed, k.conf)
	if err != nil {
		return nil, nil, err
	}
	// Purpose
	purpose, err := masterKey.Derive(hdkeychain.HardenedKeyStart + k.purpose.Uint32())
	if err != nil {
		return nil, nil, err
	}
	// CoinType
	coinType, err := purpose.Derive(hdkeychain.HardenedKeyStart + k.coinType.Uint32())
	if err != nil {
		return nil, nil, err
	}
	// Account
	// For multisig accounts (deposit=0, payment=1, stored=2), use non-hardened derivation
	// because watch wallet derives from coin-level xpubs which can only derive non-hardened children.
	// For other accounts (auth1=11, auth2=12, etc.), use hardened derivation for enhanced security.
	accountIndex := accountType.BIP44AccountIndex()
	isMultisigAccount := accountIndex <= maxMultisigAccountIndex // deposit=0, payment=1, stored=2

	logger.Debug(
		"create_key_by_account",
		"account_type", accountType.String(),
		"account_value", accountIndex,
		"is_multisig_account", isMultisigAccount,
	)

	var accountPrivKey *hdkeychain.ExtendedKey
	if isMultisigAccount {
		// Non-hardened derivation for multisig accounts to match xpub-derived keys
		accountPrivKey, err = coinType.Derive(accountIndex)
	} else {
		// Hardened derivation for non-multisig accounts (enhanced security)
		accountPrivKey, err = coinType.Derive(hdkeychain.HardenedKeyStart + accountIndex)
	}
	if err != nil {
		return nil, nil, err
	}
	// Change
	// Index

	// get pubKey
	publicKey, err := accountPrivKey.Neuter()
	if err != nil {
		return nil, nil, err
	}

	// strPrivateKey := account.String()
	// strPublicKey := publicKey.String()
	return accountPrivKey, publicKey, nil
}

// createKeysWithIndex create keys by index and count
// e.g. - idxFrom:0,  count 10 => 0-9
//   - idxFrom:10, count 10 => 10-19
func (k *HDKey) createKeysWithIndex(
	accountPrivKey *hdkeychain.ExtendedKey, idxFrom, count uint32,
) ([]domainKey.WalletKey, error) {
	// accountPrivKey, err := hdkeychain.NewKeyFromString(accountPrivKey)

	// Change
	change, err := accountPrivKey.Derive(ChangeTypeExternal.Uint32())
	if err != nil {
		return nil, err
	}

	// Index
	walletKeys := make([]domainKey.WalletKey, count)
	for i := uint32(0); i < count; i++ {
		// Derive child key
		child, err := change.Derive(idxFrom + i)
		if err != nil {
			return nil, fmt.Errorf("failed to derive child key at index %d: %w", idxFrom+i, err)
		}

		// Get private key
		privateKey, err := child.ECPrivKey()
		if err != nil {
			return nil, fmt.Errorf("failed to get EC private key at index %d: %w", idxFrom+i, err)
		}

		// Use strategy to generate wallet key
		walletKey, err := k.strategy.GenerateWalletKey(privateKey)
		if err != nil {
			return nil, fmt.Errorf("failed to generate wallet key at index %d: %w", idxFrom+i, err)
		}

		walletKeys[i] = *walletKey
	}

	return walletKeys, nil
}

// DeriveAccountKey derives an account-level extended key from a coin-level extended public key.
//
// The input extended public key should be at m/purpose'/coin' level (exported from sign wallets).
// This function derives to m/purpose'/coin'/account level and returns the account extended key.
//
// Parameters:
//   - coinLevelExtendedKey: Extended public key at m/purpose'/coin' level (xpub/tpub format)
//   - accountType: Account type (deposit=0, payment=1, stored=2, etc.)
//
// Returns:
//   - Derived account-level extended key (hdkeychain.ExtendedKey)
//   - Error if derivation fails
//
// Note: Since we're deriving from an extended public key (xpub), we can only derive non-hardened keys.
// The coin level (m/purpose'/coin') is already hardened.
func DeriveAccountKey(
	coinLevelExtendedKey string,
	accountType domainAccount.AccountType,
) (*hdkeychain.ExtendedKey, error) {
	// Parse extended public key
	coinLevelKey, err := hdkeychain.NewKeyFromString(coinLevelExtendedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse coin-level extended public key: %w", err)
	}

	// Derive account-specific key: m/purpose'/coin/account (non-hardened account index)
	// Note: Since we're deriving from an extended public key (xpub), we can only
	// derive non-hardened keys. The coin level is already hardened.
	// Use BIP44AccountIndex (deposit=0, payment=1, stored=2) not Uint32() (deposit=1, payment=2, stored=3)
	accountKey, err := coinLevelKey.Derive(accountType.BIP44AccountIndex())
	if err != nil {
		return nil, fmt.Errorf("failed to derive account key: %w", err)
	}

	return accountKey, nil
}
