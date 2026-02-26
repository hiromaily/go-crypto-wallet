package btc

import (
	"errors"
	"time"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
)

// BTCAccountKey represents a Bitcoin account key with associated addresses in the domain layer.
//
// SECURITY NOTE: The WalletImportFormat field contains the private key in WIF format.
// This must NEVER be logged or exposed in error messages. Handle with extreme care.
type BTCAccountKey struct {
	ID                     int64
	CoinTypeCode           domainCoin.CoinTypeCode
	KeyType                string // BIP type: bip44, bip49, bip84, bip86, musig2
	Account                domainAccount.AccountType
	P2pkhAddress           string  // Pay To PubKey Hash address (legacy)
	P2shSegwitAddress      string  // P2SH-SegWit address (wrapped SegWit)
	Bech32Address          string  // Native SegWit address (bech32)
	TaprootAddress         *string // Taproot address (BIP86) - nullable
	FullPublicKey          string  // Full public key
	MultisigAddress        string  // Multisig address
	RedeemScript           string  // Redeem script for multisig
	WalletImportFormat     string  // WIF - NEVER log this field
	AccountExtendedPrivkey *string // Account xpriv for BIP32 derivation - NEVER log this field
	Idx                    int64   // HD wallet index
	AddrStatus             domainAddress.AddrStatus
	UpdatedAt              *time.Time
}

// NewBtcAccountKey creates a new BtcAccountKey entity for key generation.
//
// SECURITY: The walletImportFormat parameter contains the private key.
// Ensure this value is handled securely and never logged.
func NewBtcAccountKey(
	coinTypeCode domainCoin.CoinTypeCode,
	keyType string,
	account domainAccount.AccountType,
	p2pkhAddress string,
	p2shSegwitAddress string,
	bech32Address string,
	fullPublicKey string,
	walletImportFormat string,
	idx int64,
) (*BTCAccountKey, error) {
	if coinTypeCode == "" {
		return nil, errors.New("coin type code cannot be empty")
	}
	if account == "" {
		return nil, errors.New("account cannot be empty")
	}
	if fullPublicKey == "" {
		return nil, errors.New("full public key cannot be empty")
	}
	if walletImportFormat == "" {
		return nil, errors.New("wallet import format cannot be empty")
	}
	if keyType == "" {
		return nil, errors.New("key type cannot be empty")
	}

	now := time.Now()
	return &BTCAccountKey{
		CoinTypeCode:       coinTypeCode,
		KeyType:            keyType,
		Account:            account,
		P2pkhAddress:       p2pkhAddress,
		P2shSegwitAddress:  p2shSegwitAddress,
		Bech32Address:      bech32Address,
		FullPublicKey:      fullPublicKey,
		WalletImportFormat: walletImportFormat,
		Idx:                idx,
		AddrStatus:         domainAddress.AddrStatusHDKeyGenerated,
		UpdatedAt:          &now,
	}, nil
}

// SetTaprootAddress sets the taproot address (BIP86).
func (k *BTCAccountKey) SetTaprootAddress(address string) {
	k.TaprootAddress = &address
	k.updateTimestamp()
}

// SetMultisigInfo sets multisig address and redeem script.
func (k *BTCAccountKey) SetMultisigInfo(multisigAddress, redeemScript string) {
	k.MultisigAddress = multisigAddress
	k.RedeemScript = redeemScript
	k.updateTimestamp()
}

// UpdateAddrStatus updates the address status.
func (k *BTCAccountKey) UpdateAddrStatus(status domainAddress.AddrStatus) {
	k.AddrStatus = status
	k.updateTimestamp()
}

// UpdateAddress updates the P2PKH address.
func (k *BTCAccountKey) UpdateAddress(p2pkhAddress string) {
	k.P2pkhAddress = p2pkhAddress
	k.updateTimestamp()
}

// GetWIF returns the wallet import format (private key).
// SECURITY: This returns sensitive private key data. Never log the return value.
func (k *BTCAccountKey) GetWIF() string {
	return k.WalletImportFormat
}

// SetAccountExtendedPrivkey sets the account-level extended private key for BIP32 derivation.
// SECURITY: This contains sensitive private key data. Never log the parameter or field value.
func (k *BTCAccountKey) SetAccountExtendedPrivkey(xpriv string) {
	k.AccountExtendedPrivkey = &xpriv
	k.updateTimestamp()
}

// GetAccountExtendedPrivkey returns the account-level extended private key.
// SECURITY: This returns sensitive private key data. Never log the return value.
func (k *BTCAccountKey) GetAccountExtendedPrivkey() *string {
	return k.AccountExtendedPrivkey
}

func (k *BTCAccountKey) updateTimestamp() {
	now := time.Now()
	k.UpdatedAt = &now
}
