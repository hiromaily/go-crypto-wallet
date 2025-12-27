package key

import (
	"github.com/btcsuite/btcd/chaincfg"

	domainAccount "github.com/hiromaily/go-crypto-wallet/internal/domain/account"
	domainCoin "github.com/hiromaily/go-crypto-wallet/internal/domain/coin"
	domainKey "github.com/hiromaily/go-crypto-wallet/internal/domain/key"
	"github.com/hiromaily/go-crypto-wallet/internal/infrastructure/storage/file/address"
)

// BIP84Generator implements Generator interface for BIP84 (Native SegWit Bech32 addresses)
type BIP84Generator struct {
	hdKey *HDKey
}

// NewBIP84Generator returns BIP84Generator
func NewBIP84Generator(coinTypeCode domainCoin.CoinTypeCode, conf *chaincfg.Params) *BIP84Generator {
	return &BIP84Generator{
		hdKey: NewHDKey(PurposeTypeBIP84, coinTypeCode, conf),
	}
}

// KeyType returns the key type this generator supports
func (*BIP84Generator) KeyType() domainKey.KeyType {
	return domainKey.KeyTypeBIP84
}

// CreateKey creates keys based on BIP84 standard
func (g *BIP84Generator) CreateKey(
	seed []byte,
	accountType domainAccount.AccountType,
	idxFrom, count uint32,
) ([]domainKey.WalletKey, error) {
	return g.hdKey.CreateKey(seed, accountType, idxFrom, count)
}

// SupportsAddressType checks if this generator supports the given address type
func (*BIP84Generator) SupportsAddressType(addrType address.AddrType) bool {
	return addrType == address.AddrTypeBech32
}

// GetDerivationPath returns the BIP84 derivation path
func (g *BIP84Generator) GetDerivationPath(accountType domainAccount.AccountType, index uint32) string {
	return g.hdKey.GetDerivationPath(accountType, index)
}
