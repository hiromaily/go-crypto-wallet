package key

import (
	"github.com/btcsuite/btcd/chaincfg"

	portsWallet "github.com/hiromaily/go-crypto-wallet/internal/application/ports/wallet"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainAddress "github.com/hiromaily/go-crypto-wallet/internal/domain/address"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
)

// BIP86Generator implements Generator interface for BIP86 (Taproot addresses)
type BIP86Generator struct {
	hdKey *HDKey
}

// Compile-time check that generator implements Generator interface
var _ portsWallet.Generator = (*BIP86Generator)(nil)

// NewBIP86Generator returns BIP86Generator
func NewBIP86Generator(coinTypeCode domainCoin.CoinTypeCode, conf *chaincfg.Params) *BIP86Generator {
	return &BIP86Generator{
		hdKey: NewHDKey(PurposeTypeBIP86, coinTypeCode, conf),
	}
}

// KeyType returns the key type this generator supports
func (*BIP86Generator) KeyType() domainKey.KeyType {
	return domainKey.KeyTypeBIP86
}

// CreateKey creates keys based on BIP86 standard
func (g *BIP86Generator) CreateKey(
	seed []byte,
	accountType domainAccount.AccountType,
	idxFrom, count uint32,
) ([]domainKey.WalletKey, error) {
	return g.hdKey.CreateKey(seed, accountType, idxFrom, count)
}

// CreateKeyWithAccountXpriv creates HD keys and returns the account-level extended private key.
// This method delegates to the underlying HDKey implementation.
func (g *BIP86Generator) CreateKeyWithAccountXpriv(
	seed []byte,
	accountType domainAccount.AccountType,
	idxFrom, count uint32,
) ([]domainKey.WalletKey, string, error) {
	return g.hdKey.CreateKeyWithAccountXpriv(seed, accountType, idxFrom, count)
}

// SupportsAddressType checks if this generator supports the given address type
func (*BIP86Generator) SupportsAddressType(addrType domainAddress.AddrType) bool {
	return addrType == domainAddress.AddrTypeTaproot
}

// GetDerivationPath returns the BIP86 derivation path
func (g *BIP86Generator) GetDerivationPath(accountType domainAccount.AccountType, index uint32) string {
	return g.hdKey.GetDerivationPath(accountType, index)
}
