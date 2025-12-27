package key

import (
	"github.com/btcsuite/btcd/chaincfg"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/address"
)

// BIP44Generator implements Generator interface for BIP44 (Legacy P2PKH addresses)
type BIP44Generator struct {
	hdKey *HDKey
}

// NewBIP44Generator returns BIP44Generator
func NewBIP44Generator(coinTypeCode domainCoin.CoinTypeCode, conf *chaincfg.Params) *BIP44Generator {
	return &BIP44Generator{
		hdKey: NewHDKey(PurposeTypeBIP44, coinTypeCode, conf),
	}
}

// KeyType returns the key type this generator supports
func (*BIP44Generator) KeyType() domainKey.KeyType {
	return domainKey.KeyTypeBIP44
}

// CreateKey creates keys based on BIP44 standard
func (g *BIP44Generator) CreateKey(
	seed []byte,
	accountType domainAccount.AccountType,
	idxFrom, count uint32,
) ([]domainKey.WalletKey, error) {
	return g.hdKey.CreateKey(seed, accountType, idxFrom, count)
}

// SupportsAddressType checks if this generator supports the given address type
func (*BIP44Generator) SupportsAddressType(addrType address.AddrType) bool {
	return addrType == address.AddrTypeLegacy
}

// GetDerivationPath returns the BIP44 derivation path
func (g *BIP44Generator) GetDerivationPath(accountType domainAccount.AccountType, index uint32) string {
	return g.hdKey.GetDerivationPath(accountType, index)
}
